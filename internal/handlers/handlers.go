package handlers

import (
	"context"
	"log/slog"

	tele "gopkg.in/telebot.v3"

	"joinfive-bot/internal/helpers"
	"joinfive-bot/internal/models"
	"joinfive-bot/internal/usecases"
)

// Handlers - telegram bot handlers.
type Handlers struct {
	ctx      context.Context
	usecases *usecases.UseCases
}

// NewHandlers - Handlers constructor.
func NewHandlers(ctx context.Context, u *usecases.UseCases) *Handlers {
	slog.Info("handlers: created")
	return &Handlers{
		ctx:      ctx,
		usecases: u,
	}
}

// Start - start command handler.
func (h *Handlers) Start(c tele.Context) error {
	slog.With(helpers.TeleContextAttrs(c)...).
		Info("handlers: start")
	return c.Send("👋")
}

// ChatMember - chat member event handler.
func (h *Handlers) ChatMember(c tele.Context) error {
	slog.With(helpers.TeleContextAttrs(c)...).
		Info("handlers: chat_member",
			slog.String("status", string(c.ChatMember().NewChatMember.Role)),
		)

	upd := models.NewUpdate(c.ChatMember())
	return h.usecases.UpdateCreate(h.ctx, upd)
}
