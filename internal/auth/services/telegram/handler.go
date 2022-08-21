package telegram

import (
	"fmt"
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
		if update.Message != nil { // ignore any non-Message updates
			mstr := strings.TrimSpace(update.Message.Text)
			if mstr == "" {
				return
			}
			if update.Message.IsCommand() {
				if cmd, ok := command.GetCommands()[strings.ToLower(update.Message.Command())]; ok {
					go cmd.Exec(update.Message, "", *h.service)
				} else {
					cmd := h.GetCommand(update.Message.Command())
					if cmd != nil {
						go cmd.Exec(update.Message, "", *h.service)
					}
				}
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			_, err := h.bot.Request(callback)
			if err != nil {
				return
			}
			cmd, user_id := h.GetCommandByPayload(update.CallbackQuery.Data)
			fmt.Println(user_id)
			if cmd != nil {
				go cmd.Exec(update.CallbackQuery.Message, user_id, *h.service)
			}
		}
	}
}

//get command by command payload prefix
func (h *handler) GetCommandByPayload(payload string) (cmd *command.Command, userId string) {
	for _, cmd := range command.GetCommands() {
		if strings.HasPrefix(payload, cmd.Payload) {
			return cmd, strings.TrimSpace(payload[len(cmd.Payload):])
		}
	}
	return nil, ""
}

func (h *handler) GetCommand(cmd1 string) *command.Command {
	for _, cmd := range command.GetCommands() {
		for _, alias := range cmd.Alias {
			if alias == cmd1 {
				return cmd
			}
		}
	}
	return nil
}
