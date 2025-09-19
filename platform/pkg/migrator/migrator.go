package migrator

import "database/sql"

type migrator struct {
	db            *sql.DB
	migrationsDir string
}

type Migrator interface {
	GetDB() *sql.DB
	GetMigrationsDir() string
}

func Init(db *sql.DB, migrationsDir string) Migrator {
	return &migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *migrator) GetDB() *sql.DB {
	return m.db
}

func (m *migrator) GetMigrationsDir() string {
	return m.migrationsDir
}
