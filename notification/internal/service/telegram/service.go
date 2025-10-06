package telegram

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"strings"

	"github.com/dfg007star/avito_informer/notification/internal/client/http"
	"github.com/dfg007star/avito_informer/notification/internal/config"
	"github.com/dfg007star/avito_informer/notification/internal/model"
)

//go:embed templates/item_notification.tmpl
var templateFS embed.FS

var (
	itemTemplate *template.Template
)

func init() {
	funcMap := template.FuncMap{
		"escapeHTML": escapeHTML,
	}

	itemTemplate = template.Must(template.New("item_notification.tmpl").Funcs(funcMap).ParseFS(templateFS, "templates/item_notification.tmpl"))
}

func escapeHTML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
	)
	return replacer.Replace(s)
}

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendItemNotification(ctx context.Context, item model.ItemEvent) error {
	message, err := s.buildItemMessage(item)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, config.AppConfig().Telegram.ChatID(), message)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) buildItemMessage(event model.ItemEvent) (string, error) {
	var buf bytes.Buffer
	err := itemTemplate.Execute(&buf, event)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
