package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DB_HEALTH_MSG   string        = "database system is ready to accept connections"
	SSL_MODE_PARAM  string        = "sslmode=disable"
	POSTGESQL_IMAGE string        = "postgres:16-alpine"
	DB_NAME         string        = "mingle-db"
	DB_USER         string        = "postgres"
	DB_PASSWORD     string        = "example"
	STARTUP_TIMEOUT time.Duration = 6 * time.Second
	STRATEGY_OCC    int           = 2
)

var TestStore IStore

func TestMain(m *testing.M) {
	var (
		username = "neo4j"
		password = "exam"
	)
	ctx := context.Background()

	container, connURL, err := startContainer(ctx, username, password)
	if err != nil {
		log.Fatal("can not create container:", err)
	}
	defer container.Terminate(ctx)

	driver, err := neo4j.NewDriverWithContext(connURL, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatal("can not create neo4j driver:", err)
	}

	TestStore = NewVstore(driver)

	os.Exit(m.Run())
}

func startContainer(ctx context.Context, username, password string) (testcontainers.Container, string, error) {
	request := testcontainers.ContainerRequest{
		Image:        "neo4j",
		ExposedPorts: []string{"7687/tcp"},
		Env:          map[string]string{"NEO4J_AUTH": fmt.Sprintf("%s/%s", username, password)},
		WaitingFor:   wait.ForLog("Bolt enabled"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	port, err := container.MappedPort(ctx, "7687")
	if err != nil {
		return nil, "", err
	}

	address := fmt.Sprintf("bolt://localhost:%d", port.Int())

	return container, address, nil
}

func SetupTest(_ testing.TB) func(tb testing.TB) {
	log.Println("setup test")

	return func(_ testing.TB) {
		log.Println("teardown test")
	}
}
