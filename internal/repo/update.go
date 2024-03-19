package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"joinfive-bot/internal/models"
)

// UpdateCreate creates a new update record.
func (r *SQLiteRepo) UpdateCreate(ctx context.Context, update *models.Update) error {
	const query = `
INSERT INTO updates (chat_id,
                     chat_type,
                     chat_title,
                     chat_username,
                     member_id,
                     member_first_name,
                     member_last_name,
                     member_username,
                     member_is_bot,
                     status)
VALUES (:chat_id,
        :chat_type,
        :chat_title,
        :chat_username,
        :member_id,
        :member_first_name,
        :member_last_name,
        :member_username,
        :member_is_bot,
        :status) 
RETURNING 
	id, 
	created_at
`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUpdateCreate, err)
	}
	//goland:noinspection ALL
	defer stmt.Close()

	if err = stmt.QueryRowxContext(ctx, update).
		Scan(
			&update.Id,
			&update.CreatedAt,
		); err != nil {
		return fmt.Errorf("%w: %w", ErrUpdateCreate, err)
	}
	slog.Info(
		"repo: update created",
		slog.Int64("id", update.Id),
	)
	return nil
}

// UpdateGetById retrieves an update by id.
func (r *SQLiteRepo) UpdateGetById(ctx context.Context, id int64) (*models.Update, error) {
	const query = `SELECT * FROM updates WHERE id = ?1`
	var update models.Update
	if err := r.db.QueryRowxContext(ctx, query, id).
		StructScan(&update); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %w", ErrUpdateGetById, err)
	}
	slog.Info(
		"repo: update retrieved",
		slog.Int64("id", update.Id),
	)
	return &update, nil
}
