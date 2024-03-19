package repo

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrateFS embed.FS

// migrateDB migrates the database.
// It returns the current version of the database.
func migrateDB(filepath string) (uint, error) {
	data, err := iofs.New(migrateFS, "migrations")
	if err != nil {
		return 0, err
	}

	m, err := migrate.NewWithSourceInstance("iofs", data, "sqlite://"+filepath)
	if err != nil {
		return 0, err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return 0, err
	}

	version, isDirty, err := m.Version()
	if err != nil {
		return 0, err
	}

	if isDirty {
		return 0, ErrDBDirty
	}

	return version, nil
}
