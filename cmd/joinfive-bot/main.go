package main

import (
	"log/slog"
	"os"

	"joinfive-bot/internal/app"
	"joinfive-bot/internal/config"
)

func main() {
	slog.Info("starting")

	cfg, err := config.NewConfig()
	if err != nil {
		fatal(err)
	}

	if err = app.NewApp(cfg).Start(); err != nil {
		fatal(err)
	}

	slog.Info("exiting")
}

func fatal(err error) {
	slog.Error("fatal: " + err.Error())
	os.Exit(1)
}
