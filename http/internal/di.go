package internal

import (
	"github.com/jackc/pgx/v5"
)

type diContainer struct {
	postgresClient *pgx.Conn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}
