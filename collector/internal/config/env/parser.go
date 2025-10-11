package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type parserEnvConfig struct {
	DelayBetweenLinks  time.Duration `env:"DELAY_BETWEEN_LINKS,required"`
	GetCookieTimeout   time.Duration `env:"GET_COOKIE_TIMEOUT,required"`
	GetItemsRetryDelay time.Duration `env:"GET_ITEMS_RETRY_DELAY,required"`
	GetItemsMaxRetry   int           `env:"GET_ITEMS_MAX_RETRY,required"`
}

type parserConfig struct {
	raw parserEnvConfig
}

func NewParserConfig() (*parserConfig, error) {
	var raw parserEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &parserConfig{raw: raw}, nil
}

func (cfg *parserConfig) DelayBetweenLinks() time.Duration {
	return cfg.raw.DelayBetweenLinks
}

func (cfg *parserConfig) GetCookieTimeout() time.Duration {
	return cfg.raw.GetCookieTimeout
}

func (cfg *parserConfig) GetItemsRetryDelay() time.Duration {
	return cfg.raw.GetItemsRetryDelay
}

func (cfg *parserConfig) GetItemsMaxRetry() int {
	return cfg.raw.GetItemsMaxRetry
}
