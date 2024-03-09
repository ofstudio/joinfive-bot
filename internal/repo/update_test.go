package repo

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"joinfive-bot/internal/config"
	"joinfive-bot/internal/models"
)

func TestUpdate(t *testing.T) {
	suite.Run(t, new(TestUpdateSuite))
}

type TestUpdateSuite struct {
	suite.Suite
	repo *SQLiteRepo
}

func (suite *TestUpdateSuite) SetupSubTest() {
	var err error
	suite.repo, err = NewSQLiteRepo(config.DB{
		Filepath:        filepath.Join(suite.T().TempDir(), "test.db"),
		RequiredVersion: 1,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(suite.repo)
}

func (suite *TestUpdateSuite) TearDownSubTest() {
	suite.repo.Close()
}

func (suite *TestUpdateSuite) TestUpdateCreate() {
	suite.Run("success", func() {
		update := &models.Update{
			Status: "test",
			Chat: models.Chat{
				Id:       -1,
				Type:     "test",
				Title:    "test",
				Username: "",
			},
			Member: models.Member{
				Id:        1,
				FirstName: "test",
				LastName:  "",
				Username:  "",
				IsBot:     false,
			},
		}
		err := suite.repo.UpdateCreate(context.Background(), update)
		suite.NoError(err)
		suite.NotEmpty(update.Id)
		suite.NotEmpty(update.CreatedAt)
	})

}

func (suite *TestUpdateSuite) TestUpdateGetById() {
	suite.Run("success", func() {
		update := &models.Update{
			Status: "test",
			Chat: models.Chat{
				Id:    -1,
				Type:  "test",
				Title: "test",
			},
			Member: models.Member{
				Id:        1,
				FirstName: "test",
				IsBot:     false,
			},
		}
		err := suite.repo.UpdateCreate(context.Background(), update)
		suite.NoError(err)
		suite.NotEmpty(update.Id)

		got, err := suite.repo.UpdateGetById(context.Background(), update.Id)
		suite.NoError(err)
		suite.Equal(update, got)
	})

	suite.Run("not found", func() {
		got, err := suite.repo.UpdateGetById(context.Background(), -1)
		suite.Error(err)
		suite.ErrorIs(err, ErrNotFound)
		suite.Nil(got)
	})

}
