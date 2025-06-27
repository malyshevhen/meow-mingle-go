package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
		if err := applyMigrationFile(session, file); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}
		fmt.Printf("Successfully applied migration %s\n", file)
	}

	return nil
}

// applyMigrationFile reads and executes all CQL statements in a migration file
func applyMigrationFile(session *gocql.Session, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	statements := splitCQLStatements(string(content))

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		fmt.Printf("  Executing statement %d: %s\n", i+1, truncateStatement(stmt))
		if err := session.Query(stmt).Exec(); err != nil {
			return fmt.Errorf("failed to execute statement %d (%s): %w", i+1, truncateStatement(stmt), err)
		}
	}

	return nil
}

// splitCQLStatements splits a string containing multiple CQL statements into individual statements
func splitCQLStatements(content string) []string {
	// Remove comments and split by semicolons
	lines := strings.Split(content, "\n")
	var cleanedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "//") {
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}

	// Join lines and split by semicolons
	cleanContent := strings.Join(cleanedLines, " ")
	statements := strings.Split(cleanContent, ";")

	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}

// truncateStatement truncates a CQL statement for logging purposes
func truncateStatement(stmt string) string {
	const maxLength = 100
	if len(stmt) <= maxLength {
		return stmt
	}
	return stmt[:maxLength] + "..."
}
