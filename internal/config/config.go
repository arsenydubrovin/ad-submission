package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPPort     string `env:"HTTP_PORT" env-default:"8080"`
	PostgresHost string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresDB   string `env:"POSTGRES_DB" env-default:"postgres"`
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatalf("cannot read .env file: %s", err)
	}

	return &cfg
}
