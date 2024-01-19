package config

import "fmt"

const (
	postgresHostEnvName     = "POSTGRES_HOST"
	postgresPortEnvName     = "POSTGRES_PORT"
	postgresUserEnvName     = "POSTGRES_USER"
	postgresDatabaseEnvName = "POSTGRES_DB"
)

type PostgresConfig interface {
	DSN() string
}

type postgresconfig struct {
	host string
	port string
	user string
	db   string
}

func NewPostgresConfig() (PostgresConfig, error) {
	host, err := getEnvVariable(postgresHostEnvName)
	if err != nil {
		return nil, err
	}

	port, err := getEnvVariable(postgresPortEnvName)
	if err != nil {
		return nil, err
	}

	user, err := getEnvVariable(postgresUserEnvName)
	if err != nil {
		return nil, err
	}

	db, err := getEnvVariable(postgresDatabaseEnvName)
	if err != nil {
		return nil, err
	}

	return &postgresconfig{
		host: host,
		port: port,
		user: user,
		db:   db,
	}, nil
}

func (cfg *postgresconfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		cfg.host,
		cfg.port,
		cfg.user,
		cfg.db,
	)
}
