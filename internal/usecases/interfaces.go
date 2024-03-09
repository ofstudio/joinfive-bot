package usecases

import (
	"context"

	"joinfive-bot/internal/models"
)

// Notifier - notifier interface.
type Notifier interface {
	Notify(update *models.Update)
}

// Repo - repository interface.
type Repo interface {
	UpdateCreate(context.Context, *models.Update) error
}
