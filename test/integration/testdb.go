package integration

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDatabase represents a test ScyllaDB instance
type TestDatabase struct {
	Container testcontainers.Container
	Session   *gocql.Session
	Host      string
	Port      string
}

// NewTestDatabase creates a new ScyllaDB test container
func NewTestDatabase(ctx context.Context) (*TestDatabase, error) {
	req := testcontainers.ContainerRequest{
		Image:        "scylladb/scylla:6.1",
		ExposedPorts: []string{"9042/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Starting listening for CQL clients"),
			wait.ForListeningPort("9042/tcp"),
		).WithDeadline(5 * time.Minute),
		Cmd: []string{
			"--smp", "1",
			"--memory", "1G",
			"--overprovisioned", "1",
			"--api-address", "0.0.0.0",
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

	// Create session
	session, err := createSession(host, port)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	testDB := &TestDatabase{
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

	return testDB, nil
}

// createSession creates a new CQL session with retries
func createSession(host, port string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(fmt.Sprintf("%s:%s", host, port))
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	// Retry connection with exponential backoff
	var session *gocql.Session
	var err error
	for i := range 30 {
		session, err = cluster.CreateSession()
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Failed to connect to ScyllaDB: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create session after retries: %w", err)
	}

	return session, nil
}

// ApplyMigrations applies test migrations to the database
func (tdb *TestDatabase) ApplyMigrations(ctx context.Context) error {
	// Create keyspace
	keyspaceQuery := `CREATE KEYSPACE IF NOT EXISTS mingle WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`
	if err := tdb.Session.Query(keyspaceQuery).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %w", err)
	}

	// Apply all table creations
	migrations := []string{
		// Profiles table
		`CREATE TABLE IF NOT EXISTS mingle.profiles (
			user_id text PRIMARY KEY,
			email text,
			first_name text,
			last_name text,
			bio text,
			avatar_url text,
			created_at timestamp,
			updated_at timestamp
		)`,

		// Profiles email index
		`CREATE INDEX IF NOT EXISTS profiles_email_idx ON mingle.profiles (email)`,

		// Posts table
		`CREATE TABLE IF NOT EXISTS mingle.posts (
			id uuid PRIMARY KEY,
			author_id text,
			content text,
			image_urls list<text>,
			created_at timestamp,
			updated_at timestamp
		)`,

		// Posts by author table
		`CREATE TABLE IF NOT EXISTS mingle.posts_by_author (
			author_id text,
			created_at timestamp,
			post_id uuid,
			content text,
			image_urls list<text>,
			updated_at timestamp,
			PRIMARY KEY (author_id, created_at, post_id)
		) WITH CLUSTERING ORDER BY (created_at DESC, post_id ASC)`,

		// Comments table
		`CREATE TABLE IF NOT EXISTS mingle.comments (
			id uuid PRIMARY KEY,
			post_id uuid,
			author_id text,
			content text,
			created_at timestamp,
			updated_at timestamp
		)`,

		// Comments by post table
		`CREATE TABLE IF NOT EXISTS mingle.comments_by_post (
			post_id uuid,
			created_at timestamp,
			comment_id uuid,
			author_id text,
			content text,
			updated_at timestamp,
			PRIMARY KEY (post_id, created_at, comment_id)
		) WITH CLUSTERING ORDER BY (created_at DESC, comment_id ASC)`,

		// Reactions table
		`CREATE TABLE IF NOT EXISTS mingle.reactions (
			target_id uuid,
			target_type text,
			author_id text,
			reaction_type text,
			created_at timestamp,
			PRIMARY KEY ((target_id, target_type), author_id)
		)`,

		// Reactions by target table
		`CREATE TABLE IF NOT EXISTS mingle.reactions_by_target (
			target_id uuid,
			target_type text,
			reaction_type text,
			author_id text,
			created_at timestamp,
			PRIMARY KEY ((target_id, target_type), reaction_type, author_id)
		)`,

		// Subscriptions table
		`CREATE TABLE IF NOT EXISTS mingle.subscriptions (
			follower_id text,
			following_id text,
			created_at timestamp,
			PRIMARY KEY (follower_id, following_id)
		)`,

		// Followers table
		`CREATE TABLE IF NOT EXISTS mingle.followers (
			following_id text,
			follower_id text,
			created_at timestamp,
			PRIMARY KEY (following_id, follower_id)
		)`,

		// User feed table
		`CREATE TABLE IF NOT EXISTS mingle.user_feed (
			user_id text,
			created_at timestamp,
			post_id uuid,
			author_id text,
			content text,
			image_urls list<text>,
			PRIMARY KEY (user_id, created_at, post_id)
		) WITH CLUSTERING ORDER BY (created_at DESC, post_id ASC)`,

		// User activity table
		`CREATE TABLE IF NOT EXISTS mingle.user_activity (
			user_id text,
			activity_id uuid,
			activity_type text,
			target_id uuid,
			metadata map<text, text>,
			created_at timestamp,
			PRIMARY KEY (user_id, created_at, activity_id)
		) WITH CLUSTERING ORDER BY (created_at DESC, activity_id ASC)`,
	}

	for _, migration := range migrations {
		if err := tdb.Session.Query(migration).Exec(); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

// Clean truncates all tables for test cleanup
func (tdb *TestDatabase) Clean(ctx context.Context) error {
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

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE %s", table)
		if err := tdb.Session.Query(query).Exec(); err != nil {
			log.Printf("Warning: failed to truncate table %s: %v", table, err)
		}
	}

	return nil
}

// Close closes the session and terminates the container
func (tdb *TestDatabase) Close(ctx context.Context) error {
	if tdb.Session != nil {
		tdb.Session.Close()
	}
	if tdb.Container != nil {
		return tdb.Container.Terminate(ctx)
	}
	return nil
}

// ConnectionString returns the connection string for the test database
func (tdb *TestDatabase) ConnectionString() string {
	return fmt.Sprintf("%s:%s", tdb.Host, tdb.Port)
}
