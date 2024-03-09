package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"joinfive-bot/internal/models"
)

// UseCases - application use cases.
type UseCases struct {
	repo     Repo
	notifier Notifier
}

// NewUseCases - UseCases constructor.
func NewUseCases(r Repo, n Notifier) *UseCases {
	slog.Info("usecases: created")
	return &UseCases{
		repo:     r,
		notifier: n,
	}
}

// UpdateCreate creates a new update record and notifies the notifier.
func (u *UseCases) UpdateCreate(ctx context.Context, update *models.Update) error {
	go u.notifier.Notify(update)
	if err := u.repo.UpdateCreate(ctx, update); err != nil {
		return fmt.Errorf("%w: %w", ErrUpdateCreate, err)
	}
	return nil
}
