package integration

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/pkg/migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SimpleTestDB represents a simple test database instance
type SimpleTestDB struct {
	Container testcontainers.Container
	Session   *gocql.Session
	Host      string
	Port      string
}

// NewSimpleTestDatabase creates a new ScyllaDB test container
func NewSimpleTestDatabase(ctx context.Context) (*SimpleTestDB, error) {
	log.Println("Starting ScyllaDB test container...")

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

	// Create session with retries
	session, err := createSessionWithRetries(host, port)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	testDB := &SimpleTestDB{
		Container: container,
		Session:   session,
		Host:      host,
		Port:      port,
	}

	// Apply migrations
	if err := testDB.ApplyMigrations(ctx); err != nil {
		testDB.Close(ctx)
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("ScyllaDB test container ready")
	return testDB, nil
}

// createSessionWithRetries creates a new CQL session with retries
func createSessionWithRetries(host, port string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(fmt.Sprintf("%s:%s", host, port))
	cluster.Consistency = gocql.One
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.NumConns = 1
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}
	cluster.DisableInitialHostLookup = true

	var session *gocql.Session
	var err error

	// Retry connection with exponential backoff
	for i := range 15 {
		session, err = cluster.CreateSession()
		if err == nil {
			break
		}

		backoff := min(time.Duration(i+1)*time.Second, 10*time.Second)

		log.Printf("Attempt %d: Failed to connect to ScyllaDB, retrying in %v: %v", i+1, backoff, err)
		time.Sleep(backoff)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create session after retries: %w", err)
	}

	return session, nil
}

// ApplyMigrations applies production migrations to the test database
func (db *SimpleTestDB) ApplyMigrations(ctx context.Context) error {
	// Get the path to migrations directory relative to the project root
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}

	// Go up from test/integration/simple_testdb.go to project root
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	migrationsDir := filepath.Join(projectRoot, "migrations")

	log.Printf("Applying migrations from: %s", migrationsDir)

	// Use the production migration system
	if err := migrate.ApplyMigrations(db.Session, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply production migrations: %w", err)
	}

	log.Println("Production migrations applied successfully to test database")
	return nil
}

// Clean truncates all tables for test cleanup
func (db *SimpleTestDB) Clean(ctx context.Context) error {
	if err := db.Session.Query("DROP KEYSPACE IF EXISTS mingle").Exec(); err != nil {
		log.Printf("Warning: failed to drop keyspace: %v", err)
	}
	return db.ApplyMigrations(ctx)
}

// Close closes the session and terminates the container
func (db *SimpleTestDB) Close(ctx context.Context) error {
	if db.Session != nil {
		db.Session.Close()
	}
	if db.Container != nil {
		return db.Container.Terminate(ctx)
	}
	return nil
}
