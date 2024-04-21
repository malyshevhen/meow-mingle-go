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
		"file://./../../db/migration",
		connURL)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func runPostgresContainer(
	ctx context.Context,
) (pgContainer *postgres.PostgresContainer, connStr string, err error) {
	pgContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("mingle-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("example"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(6*time.Second)),
	)
	if err != nil {
		return nil, "", err
	}

	connStr, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", err
	}
	return
}
