package repo

import (
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"joinfive-bot/internal/config"
)

func TestSQLiteRepo(t *testing.T) {
	suite.Run(t, new(TestSQLiteRepoSuite))
}

type TestSQLiteRepoSuite struct {
	suite.Suite
}

func (suite *TestSQLiteRepoSuite) TestNewRepo() {
	suite.Run("success", func() {
		r, err := NewSQLiteRepo(config.DB{
			Filepath:        filepath.Join(suite.T().TempDir(), "test.db"),
			RequiredVersion: 1,
		})
		suite.NoError(err)
		suite.NotNil(r)
	})

	suite.Run("invalid filepath", func() {
		r, err := NewSQLiteRepo(config.DB{
			Filepath:        "/dev/null",
			RequiredVersion: 1,
		})
		suite.ErrorIs(err, ErrDBMigration)
		suite.Nil(r)
	})

	suite.Run("invalid version", func() {
		r, err := NewSQLiteRepo(config.DB{
			Filepath:        filepath.Join(suite.T().TempDir(), "test.db"),
			RequiredVersion: 999,
		})
		suite.ErrorIs(err, ErrDBVersion)
		suite.Nil(r)
	})
}

func (suite *TestSQLiteRepoSuite) TestMigrateDB() {
	fp := filepath.Join(suite.T().TempDir(), "test.db")
	ver, err := migrateDB(fp)
	suite.NoError(err)
	suite.Equal(uint(1), ver)

	// check tables created
	db, err := sqlx.Open("sqlite", fp)
	suite.NoError(err)
	//goland:noinspection ALL
	defer db.Close()

	var tables []string
	err = db.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'")
	suite.NoError(err)
	suite.Equal([]string{"schema_migrations", "updates"}, tables)
}
