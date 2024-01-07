package main

import (
	"log/slog"
	"os"

	"github.com/arsenydubrovin/ad-submission/internal/config"
	l "github.com/arsenydubrovin/ad-submission/internal/logger"
	"github.com/arsenydubrovin/ad-submission/internal/storage/postgres"
)

func main() {
	cfg := config.Load()

	log := l.SetupLogger(cfg.AppEnv)
	slog.SetDefault(log)

	_, err := postgres.New(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresDB)
	if err != nil {
		log.Error("failed to initialize storage", l.Err(err))
		os.Exit(1)
	}

	log.Info("storage is initialized")
}
