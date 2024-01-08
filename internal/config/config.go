package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Env      string `env:"APP_ENV" env-default:"prod"` // default prod in more secure
		HTTPPort string `env:"HTTP_PORT" env-default:"8080"`
	}
	Postgres struct {
		Host string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port string `env:"POSTGRES_PORT" env-default:"5432"`
		User string `env:"POSTGRES_USER" env-default:"postgres"`
		DB   string `env:"POSTGRES_DB" env-default:"postgres"`
	}
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatalf("cannot read .env file: %s", err)
	}

	return &cfg
}
