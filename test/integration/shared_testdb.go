package integration

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/pkg/migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	sharedDB     *SharedTestDatabase
	sharedDBOnce sync.Once
	sharedDBMux  sync.RWMutex
)

// SharedTestDatabase represents a shared ScyllaDB instance for all tests
type SharedTestDatabase struct {
	Container testcontainers.Container
	Session   *gocql.Session
	Host      string
	Port      string
	refCount  int
	mu        sync.Mutex
}

// GetSharedTestDatabase returns a shared test database instance
func GetSharedTestDatabase(ctx context.Context) (*SharedTestDatabase, error) {
	var err error
	sharedDBOnce.Do(func() {
		sharedDB, err = newSharedTestDatabase(ctx)
	})

	if err != nil {
		return nil, err
	}

	sharedDBMux.Lock()
	sharedDB.refCount++
	sharedDBMux.Unlock()

	return sharedDB, nil
}

// newSharedTestDatabase creates a new shared ScyllaDB test container
func newSharedTestDatabase(ctx context.Context) (*SharedTestDatabase, error) {
	req := testcontainers.ContainerRequest{
		Image:        "scylladb/scylla:6.1",
		ExposedPorts: []string{"9042/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Starting listening for CQL clients"),
			wait.ForListeningPort("9042/tcp"),
		).WithDeadline(5 * time.Minute),
		Cmd: []string{
			"--smp", "1",
			"--memory", "256M",
			"--overprovisioned", "1",
			"--developer-mode", "1",
			"--disable-version-check",
			"--skip-wait-for-gossip-to-settle 0",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start ScyllaDB container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "9042")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	port := mappedPort.Port()

	// Create session with optimized settings for testing
	session, err := createOptimizedSession(host, port)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	testDB := &SharedTestDatabase{
		Container: container,
		Session:   session,
		Host:      host,
		Port:      port,
		refCount:  0,
	}

	// Apply production migrations
	if err := testDB.ApplyMigrations(ctx); err != nil {
		testDB.closeInternal(ctx)
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return testDB, nil
}

// createOptimizedSession creates a new CQL session optimized for testing
func createOptimizedSession(host, port string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(fmt.Sprintf("%s:%s", host, port))
	cluster.Consistency = gocql.One
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.NumConns = 1
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 5}
	cluster.DisableInitialHostLookup = true

	var session *gocql.Session
	var err error

	// Retry with exponential backoff (max 20 attempts for slower systems)
	for i := 0; i < 20; i++ {
		session, err = cluster.CreateSession()
		if err == nil {
			break
		}

		backoff := time.Duration(i+1) * 1 * time.Second
		if backoff > 10*time.Second {
			backoff = 10 * time.Second
		}

		log.Printf("Attempt %d: Failed to connect to ScyllaDB, retrying in %v: %v", i+1, backoff, err)
		time.Sleep(backoff)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create session after retries: %w", err)
	}

	return session, nil
}

// ApplyMigrations applies production migrations to the test database
func (sdb *SharedTestDatabase) ApplyMigrations(ctx context.Context) error {
	// Get the path to migrations directory relative to the project root
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}

	// Go up from test/integration/shared_testdb.go to project root
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	migrationsDir := filepath.Join(projectRoot, "migrations")

	log.Printf("Applying migrations from: %s", migrationsDir)

	// Use the production migration system
	if err := migrate.ApplyMigrations(sdb.Session, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply production migrations: %w", err)
	}

	log.Println("Production migrations applied successfully to test database")
	return nil
}

// Clean truncates all tables for test cleanup efficiently
func (sdb *SharedTestDatabase) Clean(ctx context.Context) error {
	sdb.mu.Lock()
	defer sdb.mu.Unlock()

	tables := []string{
		"mingle.profiles",
		"mingle.posts",
		"mingle.posts_by_author",
		"mingle.comments",
		"mingle.comments_by_post",
		"mingle.reactions",
		"mingle.reactions_by_target",
		"mingle.subscriptions",
		"mingle.followers",
		"mingle.user_feed",
		"mingle.user_activity",
	}

	// Use individual truncates for better reliability
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE %s", table)
		if err := sdb.Session.Query(query).Exec(); err != nil {
			log.Printf("Warning: failed to truncate table %s: %v", table, err)
			// Continue with other tables even if one fails
		}
	}

	return nil
}

// Release decrements the reference count
func (sdb *SharedTestDatabase) Release() {
	sharedDBMux.Lock()
	defer sharedDBMux.Unlock()

	sdb.refCount--
	if sdb.refCount <= 0 {
		// Don't actually close during tests - let it be cleaned up at the end
		log.Printf("Shared test database released (refCount: %d)", sdb.refCount)
	}
}

// closeInternal closes the session and terminates the container
func (sdb *SharedTestDatabase) closeInternal(ctx context.Context) error {
	if sdb.Session != nil {
		sdb.Session.Close()
	}
	if sdb.Container != nil {
		return sdb.Container.Terminate(ctx)
	}
	return nil
}

// ForceClose forcefully closes the shared database (use only in cleanup)
func ForceCloseSharedDatabase(ctx context.Context) error {
	sharedDBMux.Lock()
	defer sharedDBMux.Unlock()

	if sharedDB != nil {
		err := sharedDB.closeInternal(ctx)
		sharedDB = nil
		return err
	}
	return nil
}

// ConnectionString returns the connection string for the test database
func (sdb *SharedTestDatabase) ConnectionString() string {
	return fmt.Sprintf("%s:%s", sdb.Host, sdb.Port)
}

// TestDatabase wraps SharedTestDatabase to maintain compatibility
type TestDatabase struct {
	*SharedTestDatabase
}

// NewTestDatabase creates a new test database (now uses shared instance)
func NewTestDatabase(ctx context.Context) (*TestDatabase, error) {
	sharedDb, err := GetSharedTestDatabase(ctx)
	if err != nil {
		return nil, err
	}

	return &TestDatabase{SharedTestDatabase: sharedDb}, nil
}

// Close releases the shared database reference
func (tdb *TestDatabase) Close(ctx context.Context) error {
	if tdb.SharedTestDatabase != nil {
		tdb.SharedTestDatabase.Release()
	}
	return nil
}
