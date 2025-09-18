package config

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}
