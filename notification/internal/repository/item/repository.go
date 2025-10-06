package item

import (
	"github.com/jackc/pgx/v5"
)

type repository struct {
	db *pgx.Conn
}

func NewItemRepository(db *pgx.Conn) *repository {
	return &repository{
		db: db,
	}
}
