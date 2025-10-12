package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type httpEnvConfig struct {
	Host        string        `env:"HTTP_HOST,required"`
	Port        string        `env:"HTTP_PORT,required"`
	ReadTimeout time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	Password    string        `env:"HTTP_PASSWORD,required"`
}

type httpConfig struct {
	raw httpEnvConfig
}

func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &httpConfig{raw: raw}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *httpConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}

func (cfg *httpConfig) Password() string {
	return cfg.raw.Password
}
