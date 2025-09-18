package link

import (
	"github.com/jackc/pgx/v5"
)

type repository struct {
	db *pgx.Conn
}

func NewRepository(clientPostgres *pgx.Conn) *repository {
	return &repository{
		db: clientPostgres,
	}
}
