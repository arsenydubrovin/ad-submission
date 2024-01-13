package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	config struct {
		app      app
		postgres postgres
	}

	app struct {
		env      string
		httpPort string
	}

	postgres struct {
		host string
		port string
		user string
		db   string
	}
)

func loadConfig() *config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file, %s", err)
	}

	return &config{
		app: app{
			env:      getEnv("APP_ENV", "prod"), // default prod in more secure
			httpPort: getEnv("HTTP_PORT", "8080"),
		},
		postgres: postgres{
			host: getEnv("POSTGRES_HOST", "localhost"),
			port: getEnv("POSTGRES_PORT", "5432"),
			user: getEnv("POSTGRES_USER", "postgres"),
			db:   getEnv("POSTGRES_DB", "postgres"),
		},
	}
}

func getEnv(envVariable, defaultValue string) string {
	if value, exists := os.LookupEnv(envVariable); exists {
		return value
	}
	return defaultValue
}
