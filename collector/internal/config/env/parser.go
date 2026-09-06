package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type parserEnvConfig struct {
	// DelayBetweenLinks is the fallback pause between links when the jitter
	// range below is not configured. Also used as the between-cycles pause.
	DelayBetweenLinks time.Duration `env:"DELAY_BETWEEN_LINKS,required"`
	// LinkDelayMin/Max bound the randomised pause between links. If Max <= Min
	// the code falls back to DelayBetweenLinks.
	LinkDelayMin time.Duration `env:"LINK_DELAY_MIN" envDefault:"0s"`
	LinkDelayMax time.Duration `env:"LINK_DELAY_MAX" envDefault:"0s"`

	GetCookieTimeout   time.Duration `env:"GET_COOKIE_TIMEOUT,required"`
	GetItemsRetryDelay time.Duration `env:"GET_ITEMS_RETRY_DELAY,required"`
	GetItemsMaxRetry   int           `env:"GET_ITEMS_MAX_RETRY,required"`
	Proxy              string        `env:"PROXY"`

	// ProfileDir is the Chromium --user-data-dir, kept on a Docker volume so the
	// logged-in Avito session survives restarts.
	ProfileDir string `env:"PROFILE_DIR" envDefault:"/data/avito-profile"`
	// NavTimeout bounds a single page navigation + wait-for-results.
	NavTimeout time.Duration `env:"NAV_TIMEOUT" envDefault:"45s"`
	// Headless runs Chromium without a display. Leave false in the container
	// (Xvfb provides the display); useful to set true for local debugging.
	Headless bool `env:"PARSER_HEADLESS" envDefault:"false"`
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

func (cfg *parserConfig) LinkDelayRange() (min, max time.Duration) {
	return cfg.raw.LinkDelayMin, cfg.raw.LinkDelayMax
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

func (cfg *parserConfig) Proxy() string {
	return cfg.raw.Proxy
}

func (cfg *parserConfig) ProfileDir() string {
	return cfg.raw.ProfileDir
}

func (cfg *parserConfig) NavTimeout() time.Duration {
	return cfg.raw.NavTimeout
}

func (cfg *parserConfig) Headless() bool {
	return cfg.raw.Headless
}
