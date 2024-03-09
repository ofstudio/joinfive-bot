package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"

	"joinfive-bot/internal/config"
	"joinfive-bot/internal/handlers"
	"joinfive-bot/internal/notifier"
	"joinfive-bot/internal/repo"
	"joinfive-bot/internal/usecases"
)

// App is the main application.
type App struct {
	cfg config.Config
}

// NewApp - App constructor.
func NewApp(c config.Config) *App {
	return &App{cfg: c}
}

// Start starts the application.
func (a *App) Start() error {
	defer slog.Info("app: stopped")

	// app shutdown context
	Shutdown, cancel := a.ctxWithSignals()
	defer cancel()

	// create repository
	Repo, err := repo.NewSQLiteRepo(a.cfg.DB)
	if err != nil {
		return err
	}
	defer Repo.Close()

	// create bot
	Bot, err := a.newBot()
	if err != nil {
		return err
	}

	// create notifier, usecases and handlers
	Notifier := notifier.NewTelegramSingleChat(a.cfg.Settings.ReportTo, Bot)
	Usecases := usecases.NewUseCases(Repo, Notifier)
	Handlers := handlers.NewHandlers(Shutdown, Usecases)

	// bot handlers
	Bot.Use(middleware.Recover())
	Bot.Handle("/start", Handlers.Start)
	Bot.Handle(tele.OnChatMember, Handlers.ChatMember)

	// start bot
	go Bot.Start()
	slog.Info("bot: started")
	defer func() {
		Bot.Stop()
		slog.Info("bot: stopped")
	}()

	slog.Info("app: started")

	// wait for Shutdown
	<-Shutdown.Done()

	return nil
}

// ctxWithSignals returns a context that is canceled when SIGINT or SIGTERM is received.
func (a *App) ctxWithSignals() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		cancel()
	}()
	return ctx, cancel
}
