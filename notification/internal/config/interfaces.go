package config

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}

type TelegramBotConfig interface {
	Token() string
	ChatID() int64
}
