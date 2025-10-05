package config

import (
	"os"

	"github.com/dfg007star/avito_informer/collector/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Postgres PostgresConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	PostgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Postgres: PostgresCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
