package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func getEnvVariable(envVariable string) (string, error) {
	if value, exists := os.LookupEnv(envVariable); exists {
		return value, nil
	}
	return "", fmt.Errorf("env variable %s not found", envVariable)
}
