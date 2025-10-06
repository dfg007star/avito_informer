package config

import (
	"github.com/dfg007star/avito_informer/notification/internal/config/env"
	"github.com/joho/godotenv"
	"os"
)

var appConfig *config

type config struct {
	Postgres PostgresConfig
	Telegram TelegramBotConfig
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

	TelegramCfg, err := env.NewTelegramBotConfig()
	if err != nil {
	}

	appConfig = &config{
		Postgres: PostgresCfg,
		Telegram: TelegramCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
