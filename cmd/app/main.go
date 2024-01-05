package main

import (
	"log"

	"github.com/arsenydubrovin/ad-submission/internal/config"
	"github.com/arsenydubrovin/ad-submission/internal/storage/postgres"
)

func main() {
	cfg := config.Load()

	// TODO: setup logger

	_, err := postgres.New(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresDB)
	if err != nil {
		log.Fatalf("failed to initialise storage: %s", err)
	}
}
