package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v10"
)

// Config is application configuration
type Config struct {
	Bot      // Telegram bot configuration
	DB       // SQLite database configuration
	Settings // Application settings
}

// NewConfig loads configuration from [Default] and environment variables
func NewConfig() (Config, error) {
	c := Default
	if err := env.Parse(&c); err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrEnvParse, err)
	}
	slog.Info("config: loaded")
	return c, nil
}

// DB is SQLite database configuration
type DB struct {
	Filepath        string `env:"DB_FILEPATH,required"` // Path to database file
	RequiredVersion uint   // Required database schema version
}

// Settings - application settings
type Settings struct {
	ReportTo int64 `env:"REPORT_TO,required"` // Send reports to this telegram chat
}

// Bot is Telegram bot configuration
type Bot struct {
	Token            string        `env:"BOT_TOKEN,required,unset"`
	LongPollTimeout  time.Duration `env:"BOT_LONG_POLL_TIMEOUT"`
	UseWebhook       bool          `env:"BOT_USE_WEBHOOK"`
	WebhookListen    string        `env:"BOT_WEBHOOK_LISTEN"`
	WebhookPublicURL string        `env:"BOT_WEBHOOK_PUBLIC_URL"`
}
