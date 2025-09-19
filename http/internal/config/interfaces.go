package config

import (
	"time"
)

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}

type HTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}
