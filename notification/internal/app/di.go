package app

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	httpClient "github.com/dfg007star/avito_informer/notification/internal/client/http"
	telegramClient "github.com/dfg007star/avito_informer/notification/internal/client/http/telegram"
	"github.com/dfg007star/avito_informer/notification/internal/config"
	"github.com/dfg007star/avito_informer/notification/internal/repository"
	itemRepository "github.com/dfg007star/avito_informer/notification/internal/repository/item"
	"github.com/dfg007star/avito_informer/notification/internal/service"
	itemService "github.com/dfg007star/avito_informer/notification/internal/service/item"
	telegramService "github.com/dfg007star/avito_informer/notification/internal/service/telegram"
	"github.com/jackc/pgx/v5"
)

type diContainer struct {
	telegramService service.TelegramService
	telegramClient  httpClient.TelegramClient
	telegramBot     *bot.Bot

	itemService    service.ItemService
	itemRepository repository.ItemRepository

	postgresClient *pgx.Conn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(d.TelegramClient(ctx))
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().Telegram.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

func (d *diContainer) ItemService(ctx context.Context) service.ItemService {
	if d.itemService == nil {
		d.itemService = itemService.NewItemService(d.ItemRepository(ctx))
	}

	return d.itemService
}

func (d *diContainer) ItemRepository(ctx context.Context) repository.ItemRepository {
	if d.itemRepository == nil {
		d.itemRepository = itemRepository.NewItemRepository(d.PostgresClient(ctx))
	}

	return d.itemRepository
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
