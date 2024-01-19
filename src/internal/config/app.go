package config

const (
	appEnvironmentEnvName = "APP_ENV"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type ApplicationConfig interface {
	Env() string
}

type applicationConfig struct {
	env string
}

func NewApplicationConfig() (ApplicationConfig, error) {
	env, err := getEnvVariable(appEnvironmentEnvName)
	if err != nil {
		return nil, err
	}

	return &applicationConfig{
		env: env,
	}, nil
}

func (cfg *applicationConfig) Env() string {
	return cfg.env
}
