package config

import (
	"github.com/dfg007star/avito_informer/http/internal/config/env"
	"github.com/joho/godotenv"
	"os"
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
