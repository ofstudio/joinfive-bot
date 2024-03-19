package repo

import (
	"errors"
)

var (
	ErrNoFilePath  = errors.New("repo: db file path not set")
	ErrDBOpen      = errors.New("repo: failed to open db")
	ErrDBMigration = errors.New("repo: failed to migrate db")
	ErrDBVersion   = errors.New("repo: db version mismatch")
	ErrDBDirty     = errors.New("repo: db is dirty")

	ErrNotFound      = errors.New("repo: not found")
	ErrUpdateCreate  = errors.New("repo: failed to create update")
	ErrUpdateGetById = errors.New("repo: failed to get update by id")
)
