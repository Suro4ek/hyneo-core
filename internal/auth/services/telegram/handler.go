package telegram

import (
	"hyneo/internal/auth/services"
	"hyneo/internal/auth/services/command"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler struct {
	bot     *tgbotapi.BotAPI
	service *services.Service
}

func NewTelegramHandler(bot *tgbotapi.BotAPI, service *services.Service) *handler {
	return &handler{
		bot:     bot,
		service: service,
	}
}

func (h *handler) Message() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		if cmd, ok := command.GetCommands()[strings.ToLower(update.Message.Command())]; ok {
			go cmd.Exec(update.Message, *h.service)
		}
	}
}
