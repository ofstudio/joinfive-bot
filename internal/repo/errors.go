package repo

import (
	"errors"
)

var (
	ErrDBMigration = errors.New("repo: failed to migrate db")
	ErrDBVersion   = errors.New("repo: db version mismatch")
	ErrDBDirty     = errors.New("repo: db is dirty")
	ErrDBOpen      = errors.New("repo: failed to open db")
	ErrStmtPrepare = errors.New("repo: failed to prepare statement")
	ErrNoFilePath  = errors.New("repo: db file path not set")

	ErrNotFound      = errors.New("repo: not found")
	ErrUpdateCreate  = errors.New("repo: failed to create update")
	ErrUpdateGetById = errors.New("repo: failed to get update by id")
)
