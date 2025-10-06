package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dfg007star/avito_informer/notification/internal/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	itemService := a.diContainer.ItemService(ctx)
	telegramService := a.diContainer.TelegramService(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			items, err := itemService.GetNotNotifiedItems(ctx)
			if err != nil {
				fmt.Printf("failed to get not notified items: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}

			for _, item := range items {
				event := model.ItemEvent{
					Title:       item.Title,
					Description: item.Description,
					Price:       item.Price,
					Url:         item.Url,
					PreviewUrl:  item.PreviewUrl,
					CreatedAt:   item.CreatedAt,
				}
				err := telegramService.SendItemNotification(ctx, event)
				if err != nil {
					fmt.Printf("failed to send telegram message for item %d: %v\n", item.ID, err)
					continue
				}

				item.NeedNotify = true
				err = itemService.UpdateItem(ctx, item)
				if err != nil {
					fmt.Printf("failed to update item %d need_notify status: %v\n", item.ID, err)
					continue
				}
			}

			time.Sleep(10 * time.Second)
		}
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initTelegramBot,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	telegramBot := a.diContainer.TelegramBot(ctx)

	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "🚀 avito bomb activate",
		})
		if err != nil {
			fmt.Printf("failed to send activation message: %v\n", err)
		}
	})

	go func() {
		fmt.Println("🤖 telegram bot started...")
		telegramBot.Start(ctx)
	}()

	return nil
}
