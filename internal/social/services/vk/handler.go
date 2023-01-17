package vk

import (
	"context"
	"github.com/go-redis/redis/v9"
	"hyneo/internal/social/services"
	command2 "hyneo/internal/social/services/command"
	"hyneo/internal/user"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type handler struct {
	lp      *longpoll.LongPoll
	service *services.Service
}

func NewVKHandler(lp *longpoll.LongPoll, service *services.Service) *handler {
	return &handler{
		lp:      lp,
		service: service,
	}
}

// Message
// функция обработки сообщений из вконтакте
func (h *handler) Message() {
	h.lp.MessageNew(func(_ context.Context, m events.MessageNewObject) {
		mstr := strings.TrimSpace(m.Message.Text)
		if mstr == "" {
			return
		}
		marray := strings.Fields(mstr)
		cmd, userId := h.GetCommandByPayload(strings.ReplaceAll(m.Message.Payload, "\"", ""))
		if cmd != nil {
			if cmd.WithoutUser {
				go cmd.Exec(m, &user.LinkUser{}, *h.service)
				return
			}
			userIdInt, err := strconv.ParseInt(userId, 10, 64)
			if err != nil {
				(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", int64(m.Message.FromID))
				return
			}
			u := &user.LinkUser{}
			ser := (*h.service).GetService()
			err = ser.Redis.HGetAll(context.Background(), "link:"+userId).Scan(&u)
			if err != nil {
				u, err = (*h.service).GetUserID(userIdInt)
				if err != nil {
					(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", int64(m.Message.FromID))
					return
				} else {
					ctx := context.TODO()
					if _, err := ser.Redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
						rdb.HSet(ctx, "link:"+userId, "id", u.ID)
						rdb.HSet(ctx, "link:"+userId, "service_id", u.ServiceId)
						rdb.HSet(ctx, "link:"+userId, "service_user_id", u.ServiceUserID)
						rdb.HSet(ctx, "link:"+userId, "notificated", u.Notificated)
						rdb.HSet(ctx, "link:"+userId, "banned", u.Banned)
						rdb.HSet(ctx, "link:"+userId, "double_auth", u.DoubleAuth)
						rdb.HSet(ctx, "link:"+userId, "user_id", u.UserID)
						return nil
					}); err != nil {
						(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", int64(m.Message.FromID))
						return
					}
					ser.Redis.Expire(ctx, "link:"+userId, time.Second*60)
				}
			}
			go cmd.Exec(m, u, *h.service)
		} else {
			if cmd, ok := command2.GetCommands()[strings.ToLower(marray[0])]; ok {
				if cmd.Payload == "-1" {
					go cmd.Exec(m, &user.LinkUser{}, *h.service)
				}
			}
		}
	})
	log.Println("Start Long Poll")
	if err := h.lp.Run(); err != nil {
		log.Fatal(err)
	}
}

/*
	Функция нахождение команды по payload
	Возращяет указатель на команду может быть nil и userId из auth.LinkUser
*/
func (h *handler) GetCommandByPayload(payload string) (cmd *command2.Command, userId string) {
	for _, cmd := range command2.GetCommands() {
		if strings.HasPrefix(payload, cmd.Payload) {
			return cmd, strings.TrimSpace(payload[len(cmd.Payload):])
		}
	}
	return nil, ""
}
