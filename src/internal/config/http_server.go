package config

const (
	httpPortEnvName = "HTTP_PORT"
)

type HTTPConfig interface {
	Port() string
}

type httpConfig struct {
	port string
}

func NewHTTPConfig() (HTTPConfig, error) {
	port, err := getEnvVariable(httpPortEnvName)
	if err != nil {
		return nil, err
	}

	return &httpConfig{
		port: port,
	}, nil
}

func (cfg *httpConfig) Port() string {
	return cfg.port
}
