package app

import (
	"context"
	"fmt"

	"github.com/dfg007star/avito_informer/collector/internal/config"
	"github.com/dfg007star/avito_informer/collector/internal/parser"
	"github.com/dfg007star/avito_informer/collector/internal/repository"
	"github.com/dfg007star/avito_informer/collector/internal/service"
	"github.com/jackc/pgx/v5"
)

type diContainer struct {
	service    *service.Service
	repository *repository.Repository
	parser     *parser.Parser

	postgresClient *pgx.Conn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Parser() *parser.Parser {
	if d.parser == nil {
		d.parser = parser.NewParser()
	}

	return d.parser
}

func (d *diContainer) Service(ctx context.Context) *service.Service {
	if d.service == nil {
		d.service = service.New(d.Repository(ctx))
	}

	return d.service
}

func (d *diContainer) Repository(ctx context.Context) *repository.Repository {
	if d.repository == nil {
		d.repository = repository.New(d.PostgresClient(ctx))
	}

	return d.repository
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
