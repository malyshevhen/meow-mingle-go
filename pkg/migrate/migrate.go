package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gocql/gocql"
)

// ApplyMigrations applies all pending migrations to the database.
func ApplyMigrations(session *gocql.Session, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.up.cql"))
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Applying migration %s\n", file)
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if err := session.Query(string(content)).Exec(); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}
	}

	return nil
}
