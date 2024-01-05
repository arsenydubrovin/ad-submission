package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port        string `env:"HTTP_PORT" env-default:"8080"`
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatalf("cannot read .env file: %s", err)
	}

	return &cfg
}
