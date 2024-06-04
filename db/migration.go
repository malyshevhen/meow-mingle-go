package migrations

import (
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed all:migration
var MigrationsFS embed.FS

func Init() (sd source.Driver, sourceName string, initErr error) {
	sourceName = "iofs"
	sd, err := iofs.New(MigrationsFS, "migration")
	if err != nil {
		initErr = fmt.Errorf("failed to create migration source driver: %s", err.Error())
		return
	}

	embededFiles, err := fs.Glob(MigrationsFS, "migration/*.sql")
	if err != nil {
		initErr = fmt.Errorf("unable to read migrations: %s", err.Error())
		return
	}

	for _, file := range embededFiles {
		log.Printf("%-15s ==> Migration file: %s\n", "Migrations", file)
	}

	return
}
