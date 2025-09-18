package internal

import (
	"context"
	"fmt"
	repository "github.com/dfg007star/avito_informer/http/internal/repository"
	linkRepository "github.com/dfg007star/avito_informer/http/internal/repository/link"
	"github.com/jackc/pgx/v5"
)

type diContainer struct {
	linkRepository repository.LinkRepository

	postgresClient *pgx.Conn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) LinkRepository(ctx context.Context) repository.LinkRepository {
	if d.linkRepository == nil {
		d.linkRepository = linkRepository.NewRepository(d.PostgresClient(ctx))
	}

	return d.linkRepository
}

func (d *diContainer) PostgresClient(ctx context.Context) *pgx.Conn {
	if d.postgresClient == nil {
		con, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to database: %w", err))
		}

		err = con.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("database is unavailable: %w", err))
		}

		d.postgresClient = con
	}

	return d.postgresClient
}
