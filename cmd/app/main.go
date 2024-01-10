package main

import (
	"log/slog"
	"os"

	"github.com/arsenydubrovin/ad-submission/internal/models"
)

func main() {
	cfg := loadConfig()

	log := setupLogger(cfg.app.env)
	slog.SetDefault(log)

	db, err := openDB(cfg.postgres.host, cfg.postgres.port, cfg.postgres.user, cfg.postgres.db)
	if err != nil {
		log.Error("failed to open database", wrapErr(err))
		os.Exit(1)
	}
	log.Info("connection to the database is set")

	md := models.New(db)

	advert := &models.Advert{
		Title:       "Фотик",
		Price:       456,
		Description: "Зачётный пленочный фотки «Зенит». Достался от деда»",
		PhotoLinks: []string{
			"фотка 1",
			"фотка 2",
			"фотка 3",
		},
	}

	err = md.Adverts.Insert(advert)
	if err != nil {
		log.Error("failed to insert advert", wrapErr(err))
		os.Exit(1)
	}
	log.Info("advert is inserted", slog.Int("id", advert.Id))
}
