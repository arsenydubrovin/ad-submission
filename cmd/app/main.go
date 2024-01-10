package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := loadConfig()

	log := setupLogger(cfg.app.env)
	slog.SetDefault(log)

	_, err := openDB(cfg.postgres.host, cfg.postgres.port, cfg.postgres.user, cfg.postgres.db)
	if err != nil {
		log.Error("failed to open database", wrapErr(err))
		os.Exit(1)
	}
	log.Info("connection to the database is set")
}
