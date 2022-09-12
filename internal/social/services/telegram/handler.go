package telegram

import (
	"context"
	"github.com/go-redis/redis/v9"
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
	command2 "hyneo/internal/social/services/command"
	"strconv"
	"strings"
	"time"

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
				if cmd, ok := command2.GetCommands()[strings.ToLower(update.Message.Command())]; ok {
					go cmd.Exec(update.Message, &auth.LinkUser{}, *h.service)
				} else {
					cmd := h.GetCommand(update.Message.Command())
					if cmd != nil {
						go cmd.Exec(update.Message, &auth.LinkUser{}, *h.service)
					}
				}
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			_, err := h.bot.Request(callback)
			if err != nil {
				return
			}
			cmd, userId := h.GetCommandByPayload(update.CallbackQuery.Data)
			if cmd == nil {
				return
			}
			if cmd.WithoutUser {
				go cmd.Exec(update.CallbackQuery.Message, &auth.LinkUser{}, *h.service)
				return
			}
			userIdInt, err := strconv.ParseInt(userId, 10, 64)
			if err != nil {
				(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", update.Message.From.ID)
				return
			}
			user := &auth.LinkUser{}
			ser := (*h.service).GetService()
			err = ser.Redis.HGetAll(context.Background(), "link:"+userId).Scan(&user)
			if err != nil {
				user, err = (*h.service).GetUserID(userIdInt)
				if err != nil {
					(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", update.Message.From.ID)
					return
				} else {
					ctx := context.TODO()
					if _, err := ser.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
						rdb.HSet(ctx, "link:"+userId, "id", user.ID)
						rdb.HSet(ctx, "link:"+userId, "service_id", user.ServiceId)
						rdb.HSet(ctx, "link:"+userId, "service_user_id", user.ServiceUserID)
						rdb.HSet(ctx, "link:"+userId, "notificated", user.Notificated)
						rdb.HSet(ctx, "link:"+userId, "banned", user.Banned)
						rdb.HSet(ctx, "link:"+userId, "double_auth", user.DoubleAuth)
						rdb.HSet(ctx, "link:"+userId, "user_id", user.UserID)
						return nil
					}); err != nil {
						(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", update.Message.From.ID)
						return
					}
					ser.Redis.Expire(ctx, "link:"+userId, time.Minute*5)
				}
			}
			if cmd != nil {
				go cmd.Exec(update.CallbackQuery.Message, user, *h.service)
			}
		}
	}
}

func (h *handler) GetCommandByPayload(payload string) (cmd *command2.Command, userId string) {
	for _, cmd := range command2.GetCommands() {
		if strings.HasPrefix(payload, cmd.Payload) {
			return cmd, strings.TrimSpace(payload[len(cmd.Payload):])
		}
	}
	return nil, ""
}

func (h *handler) GetCommand(cmd1 string) *command2.Command {
	for _, cmd := range command2.GetCommands() {
		for _, alias := range cmd.Alias {
			if alias == cmd1 {
				return cmd
			}
		}
	}
	return nil
}
