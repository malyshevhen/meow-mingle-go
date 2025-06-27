package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/pkg/migrate"
)

func main() {
	var (
		host          = flag.String("host", "127.0.0.1:9042", "Database host")
		user          = flag.String("user", "scylla", "Database user")
		password      = flag.String("password", "scyllapassword", "Database password")
		migrationsDir = flag.String("dir", "./migrations", "Migrations directory")
		help          = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		fmt.Println("Migration tool for Meow Mingle")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  migrate [options]")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("Environment variables:")
		fmt.Println("  SCYLLA_URL      - Database host (overrides -host)")
		fmt.Println("  SCYLLA_USER     - Database user (overrides -user)")
		fmt.Println("  SCYLLA_PASS     - Database password (overrides -password)")
		fmt.Println("  MIGRATIONS_DIR  - Migrations directory (overrides -dir)")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  migrate")
		fmt.Println("  migrate -host 192.168.1.100:9042 -user admin -password secret")
		fmt.Println("  SCYLLA_URL=localhost:9042 migrate")
		return
	}

	// Override with environment variables if set
	if envHost := os.Getenv("SCYLLA_URL"); envHost != "" {
		*host = envHost
	}
	if envUser := os.Getenv("SCYLLA_USER"); envUser != "" {
		*user = envUser
	}
	if envPassword := os.Getenv("SCYLLA_PASS"); envPassword != "" {
		*password = envPassword
	}
	if envDir := os.Getenv("MIGRATIONS_DIR"); envDir != "" {
		*migrationsDir = envDir
	}

	fmt.Printf("🗄️  Meow Mingle Database Migration Tool\n")
	fmt.Printf("=====================================\n\n")
	fmt.Printf("Database host: %s\n", *host)
	fmt.Printf("Database user: %s\n", *user)
	fmt.Printf("Migrations directory: %s\n", *migrationsDir)
	fmt.Printf("\n")

	// Check if migrations directory exists
	if _, err := os.Stat(*migrationsDir); os.IsNotExist(err) {
		log.Fatalf("❌ Migrations directory does not exist: %s", *migrationsDir)
	}

	// Create cluster configuration
	cluster := gocql.NewCluster(*host)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: *user,
		Password: *password,
	}
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4

	fmt.Printf("🔌 Connecting to database...\n")
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer session.Close()

	fmt.Printf("✅ Connected to database successfully\n\n")

	// Apply migrations
	fmt.Printf("🚀 Starting migration process...\n")
	if err := migrate.ApplyMigrations(session, *migrationsDir); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Printf("\n✅ All migrations applied successfully!\n")
	fmt.Printf("🎉 Database schema is up to date.\n")
}
