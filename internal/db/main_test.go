package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DB_HEALTH_MSG        string        = "database system is ready to accept connections"
	SSL_MODE_PARAM       string        = "sslmode=disable"
	POSTGESQL_IMAGE      string        = "postgres:16-alpine"
	DB_NAME              string        = "mingle-db"
	DB_USER              string        = "postgres"
	DB_PASSWORD          string        = "example"
	MIGRATION_SOURCE_URL string        = "file://./../../db/migration"
	STARTUP_TIMEOUT      time.Duration = 6 * time.Second
	STRATEGY_OCC         int           = 2
)

var (
	TestStore IStore
	Migration *migrate.Migrate
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, connURL, err := runPostgresContainer(ctx)
	if err != nil {
		log.Fatal("can not create container:", err)
	}
	defer container.Terminate(ctx)

	conn, err := sql.Open("postgres", connURL)
	if err != nil {
		log.Fatal("can not connect to the DB:", err)
	}

	TestStore = NewSQLStore(conn)

	Migration, err = migrate.New(
		MIGRATION_SOURCE_URL,
		connURL)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func runPostgresContainer(ctx context.Context) (pgContainer *postgres.PostgresContainer, connStr string, err error) {
	withImage := testcontainers.WithImage(POSTGESQL_IMAGE)
	withDB := postgres.WithDatabase(DB_NAME)
	withUsername := postgres.WithUsername(DB_USER)
	withPassword := postgres.WithPassword(DB_PASSWORD)
	withStrategy := testcontainers.WithWaitStrategy(
		wait.
			ForLog(DB_HEALTH_MSG).
			WithOccurrence(STRATEGY_OCC).
			WithStartupTimeout(STARTUP_TIMEOUT),
	)

	pgContainer, err = postgres.RunContainer(
		ctx,
		withImage,
		withDB,
		withUsername,
		withPassword,
		withStrategy,
	)
	if err != nil {
		return nil, "", err
	}

	connStr, err = pgContainer.ConnectionString(ctx, SSL_MODE_PARAM)
	if err != nil {
		return nil, "", err
	}
	return
}
