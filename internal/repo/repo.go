package repo

import (
	"fmt"
	"log/slog"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/jmoiron/sqlx"

	"joinfive-bot/internal/config"
)

// SQLiteRepo is a repository implementation for SQLite database
type SQLiteRepo struct {
	db *sqlx.DB
}

// NewSQLiteRepo - SQLiteRepo constructor
func NewSQLiteRepo(c config.DB) (*SQLiteRepo, error) {

	if c.Filepath == "" {
		return nil, ErrNoFilePath
	}

	// migrate db
	version, err := migrateDB(c.Filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDBMigration, err)
	}

	// check db version
	if version != c.RequiredVersion {
		return nil, fmt.Errorf(
			"%w: got '%d', want '%d'",
			ErrDBVersion,
			version,
			c.RequiredVersion,
		)
	}

	// connect to db
	db, err := sqlx.Open("sqlite", c.Filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDBOpen, err)
	}

	slog.Info("repo: created",
		slog.String("driver", "sqlite"),
		slog.String("path", c.Filepath),
		slog.Int("version", int(version)),
	)

	return &SQLiteRepo{
		db: db,
	}, nil
}

// Close closes the repository
func (r *SQLiteRepo) Close() {
	if err := r.db.Close(); err != nil {
		slog.Error(fmt.Sprintf("repo: failed to close: %v", err))
	} else {
		slog.Info("repo: closed")
	}
}
