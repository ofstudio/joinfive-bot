package repo

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	"joinfive-bot/internal/config"
)

// SQLiteRepo is a repository implementation for SQLite database
type SQLiteRepo struct {
	db         *sqlx.DB
	stmts      map[stmtId]*sqlx.Stmt
	namedStmts map[namedStmtId]*sqlx.NamedStmt
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
		return nil, fmt.Errorf("%w: got '%d', want '%d'", ErrDBVersion, version, c.RequiredVersion)
	}

	// connect to db
	db, err := sqlx.Open("sqlite3", c.Filepath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDBOpen, err)
	}

	// prepare statements
	stmts, namedStmts, err := prepareStatements(db)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStmtPrepare, err)
	}

	slog.Info("repo: created",
		slog.String("driver", "sqlite3"),
		slog.String("path", c.Filepath),
		slog.Int("version", int(version)),
	)

	return &SQLiteRepo{
		db:         db,
		stmts:      stmts,
		namedStmts: namedStmts,
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

//go:embed migrations/*.sql
var fs embed.FS

// migrateDB migrates the database
func migrateDB(filepath string) (uint, error) {
	data, err := iofs.New(fs, "migrations")
	if err != nil {
		return 0, err
	}

	m, err := migrate.NewWithSourceInstance("iofs", data, "sqlite3://"+filepath)
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

type (
	stmtId      int
	namedStmtId int
)

var (
	queries      []string
	namedQueries []string
)

func addStmt(query string) stmtId {
	queries = append(queries, query)
	return stmtId(len(queries) - 1)
}

func addNamedStmt(query string) namedStmtId {
	namedQueries = append(namedQueries, query)
	return namedStmtId(len(namedQueries) - 1)
}

func prepareStatements(db *sqlx.DB) (map[stmtId]*sqlx.Stmt, map[namedStmtId]*sqlx.NamedStmt, error) {
	statements := make(map[stmtId]*sqlx.Stmt, len(queries))
	for id, query := range queries {
		s, err := db.Preparex(query)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: `%s`", err, query)
		}
		statements[stmtId(id)] = s
	}

	namedStatements := make(map[namedStmtId]*sqlx.NamedStmt, len(namedQueries))
	for id, query := range namedQueries {
		s, err := db.PrepareNamed(query)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: `%s`", err, query)
		}
		namedStatements[namedStmtId(id)] = s
	}

	return statements, namedStatements, nil
}
