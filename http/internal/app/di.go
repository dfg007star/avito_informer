package app

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dfg007star/avito_informer/http/internal/config"
	"github.com/dfg007star/avito_informer/http/internal/handler"
	auth "github.com/dfg007star/avito_informer/http/internal/handler/auth"
	linkHandler "github.com/dfg007star/avito_informer/http/internal/handler/link"
	log "github.com/dfg007star/avito_informer/http/internal/handler/log"
	repository "github.com/dfg007star/avito_informer/http/internal/repository"
	linkRepository "github.com/dfg007star/avito_informer/http/internal/repository/link"
	"github.com/dfg007star/avito_informer/http/internal/service"
	linkService "github.com/dfg007star/avito_informer/http/internal/service/link"
	"github.com/jackc/pgx/v5"
)

type diContainer struct {
	linkRouter     *http.ServeMux
	linkHandler    handler.LinkHandler
	authHandler    handler.AuthHandler
	logHandler     handler.LogHandler
	linkService    service.LinkService
	linkRepository repository.LinkRepository
	templates      *template.Template

	postgresClient *pgx.Conn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Templates() *template.Template {
	if d.templates == nil {
		tmpl, err := template.ParseGlob("./html/*.html")
		if err != nil {
			panic(fmt.Errorf("failed to parse templates: %w", err))
		}
		d.templates = tmpl
	}

	return d.templates
}

func (d *diContainer) LinkRouter(
	ctx context.Context,
	linkHandler handler.LinkHandler,
	authHandler handler.AuthHandler,
	logHandler handler.LogHandler,
) *http.ServeMux {
	if d.linkRouter == nil {
		r := http.NewServeMux()

		fileServer := http.FileServer(http.Dir("./static/"))
		r.Handle("/static/", http.StripPrefix("/static/", fileServer))

		r.Handle("GET /{$}", auth.Middleware(http.HandlerFunc(linkHandler.IndexHandler)))
		r.Handle("POST /links", auth.Middleware(http.HandlerFunc(linkHandler.CreateLinkHandler)))
		r.Handle("DELETE /links/{id}", auth.Middleware(http.HandlerFunc(linkHandler.DeleteLinkHandler)))
		r.Handle("GET /links/{id}", auth.Middleware(http.HandlerFunc(linkHandler.ShowLinkHandler)))
		r.HandleFunc("/login", authHandler.LoginHandler)
		r.Handle("GET /log", auth.Middleware(http.HandlerFunc(logHandler.IndexHandler)))

		d.linkRouter = r
	}

	return d.linkRouter
}

func (d *diContainer) LinkHandler(ctx context.Context) handler.LinkHandler {
	if d.linkHandler == nil {
		d.linkHandler = linkHandler.NewHandler(d.LinkService(ctx), d.Templates())
	}

	return d.linkHandler
}

func (d *diContainer) LogHandler(ctx context.Context) handler.LogHandler {
	if d.logHandler == nil {
		d.logHandler = log.NewHandler(d.Templates())
	}

	return d.logHandler
}

func (d *diContainer) AuthHandler(ctx context.Context) handler.AuthHandler {
	if d.authHandler == nil {
		d.authHandler = auth.NewHandler(d.Templates())
	}

	return d.authHandler
}

func (d *diContainer) LinkService(ctx context.Context) service.LinkService {
	if d.linkService == nil {
		d.linkService = linkService.NewLinkService(d.LinkRepository(ctx))
	}

	return d.linkService
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
