package config

import (
	"os"

	"github.com/dfg007star/avito_informer/collector/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Postgres PostgresConfig
	Parser   ParserConfig
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

	ParserCfg, err := env.NewParserConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Postgres: PostgresCfg,
		Parser:   ParserCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
