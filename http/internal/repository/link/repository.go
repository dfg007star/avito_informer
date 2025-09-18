package link

import (
	"github.com/jackc/pgx/v5"

	def "github.com/dfg007star/avito_informer/http/"
)

var _ def.LinkRepository = (*repository)(nil)

type repository struct {
	db *pgx.Conn
}

func NewRepository(clientPostgres *pgx.Conn) *repository {
	return &repository{
		db: clientPostgres,
	}
}
