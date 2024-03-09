package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"joinfive-bot/internal/models"
)

var updateCreate = addNamedStmt(`
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
`)

// UpdateCreate creates a new update record.
func (r *SQLiteRepo) UpdateCreate(ctx context.Context, update *models.Update) error {
	if err := r.namedStmts[updateCreate].
		QueryRowxContext(ctx, update).
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

var updateGetById = addStmt(`
SELECT * FROM updates
WHERE id = ?1
`)

// UpdateGetById retrieves an update by id.
func (r *SQLiteRepo) UpdateGetById(ctx context.Context, id int64) (*models.Update, error) {
	var update models.Update
	if err := r.stmts[updateGetById].
		QueryRowxContext(ctx, id).
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
