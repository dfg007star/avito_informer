package pg

import (
	"database/sql"

	"github.com/pressly/goose/v3"

	"github.com/dfg007star/avito_informer/platform/pkg/migrator"
)

type PostgresMigrator struct {
	migrator.Migrator
}

func New(db *sql.DB, migrationsDir string) *PostgresMigrator {
	return &PostgresMigrator{
		Migrator: migrator.Init(db, migrationsDir),
	}
}

func (m *PostgresMigrator) Up() error {
	err := goose.Up(m.GetDB(), m.GetMigrationsDir())
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresMigrator) Down() error {
	return goose.Down(m.GetDB(), m.GetMigrationsDir())
}
