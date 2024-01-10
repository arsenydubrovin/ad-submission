package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	app struct {
		env      string `env:"APP_ENV" env-default:"prod"` // default prod in more secure
		httpport string `env:"HTTP_PORT" env-default:"8080"`
	}
	postgres struct {
		host string `env:"POSTGRES_HOST" env-default:"localhost"`
		port string `env:"POSTGRES_PORT" env-default:"5432"`
		user string `env:"POSTGRES_USER" env-default:"postgres"`
		db   string `env:"POSTGRES_DB" env-default:"postgres"`
	}
}

func loadConfig() *config {
	var cfg config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatalf("cannot read .env file: %s", err)
	}

	return &cfg
}
