package app

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"

	tele "gopkg.in/telebot.v3"

	"joinfive-bot/internal/helpers"
)

// allowedUpdates is a list of the allowed update types for the bot
var allowedUpdates = []string{
	"message",
	"chat_member",
}

// newBot creates a new bot instance
func (a *App) newBot() (*tele.Bot, error) {
	var err error
	var poller tele.Poller
	var pollerAttr slog.Attr

	if a.cfg.Token == "" {
		return nil, fmt.Errorf("%w: %w", ErrBotCreate, ErrBotTokenNotSet)
	}

	// Create poller
	if a.cfg.Bot.UseWebhook {
		pollerAttr = slog.Group("poller",
			slog.String("type", "webhook"),
			slog.String("public_url", a.cfg.Bot.WebhookPublicURL),
			slog.String("listen", a.cfg.Bot.WebhookListen),
		)
		poller, err = a.newWebhook()
	} else {
		pollerAttr = slog.Group("poller",
			slog.String("type", "long_poller"),
		)
		poller, err = a.newLongPoller()
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBotCreate, err)
	}

	// Create bot
	bot, err := tele.NewBot(tele.Settings{
		Token: a.cfg.Bot.Token,
		OnError: func(err error, c tele.Context) {
			slog.Error(err.Error(), helpers.TeleContextAttrs(c)...)
		},
		Poller: poller,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBotCreate, err)
	}

	// Remove webhook if it was set before
	if err = bot.RemoveWebhook(false); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBotCreate, err)
	}

	slog.Info("bot: created",
		slog.String("bot_username", bot.Me.Username),
		pollerAttr,
	)

	return bot, nil
}

// newWebhook creates a new webhook poller
func (a *App) newWebhook() (tele.Poller, error) {
	secretToken, err := a.newSecretToken()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBotSecretToken, err)
	}

	return &tele.Webhook{
		Listen:         a.cfg.Bot.WebhookListen,
		AllowedUpdates: allowedUpdates,
		SecretToken:    secretToken,
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: a.cfg.Bot.WebhookPublicURL,
		},
	}, nil
}

// newLongPoller creates a new long poller
func (a *App) newLongPoller() (tele.Poller, error) {
	return &tele.LongPoller{
		Timeout:        a.cfg.Bot.LongPollTimeout,
		AllowedUpdates: allowedUpdates,
	}, nil
}

// newSecretToken creates a new secret token
// https://core.telegram.org/bots/api#setwebhook
func (a *App) newSecretToken() (string, error) {
	b := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
