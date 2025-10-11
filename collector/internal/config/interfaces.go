package config

import "time"

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}

type ParserConfig interface {
	DelayBetweenLinks() time.Duration
	GetCookieTimeout() time.Duration
	GetItemsRetryDelay() time.Duration
	GetItemsMaxRetry() int
}
