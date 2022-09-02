package vk

import (
	"context"
	"github.com/go-redis/redis/v9"
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
	command2 "hyneo/internal/social/services/command"
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
				go cmd.Exec(m, &auth.LinkUser{}, *h.service)
				return
			}
			userIdInt, err := strconv.ParseInt(userId, 10, 64)
			if err != nil {
				(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", int64(m.Message.FromID))
				return
			}
			user := &auth.LinkUser{}
			ser := (*h.service).GetService()
			err = ser.Redis.HGetAll(context.Background(), "link:"+userId).Scan(&user)
			if err != nil {
				user, err = (*h.service).GetUserID(userIdInt)
				if err != nil {
					(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", int64(m.Message.FromID))
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
						(*h.service).ClearKeyboard("Этот аккаунт не привязан к вам", int64(m.Message.FromID))
						return
					}
					ser.Redis.Expire(ctx, "link:"+userId, time.Minute*5)
				}
			}
			go cmd.Exec(m, user, *h.service)
		} else {
			if cmd, ok := command2.GetCommands()[strings.ToLower(marray[0])]; ok {
				if cmd.Payload == "-1" {
					go cmd.Exec(m, &auth.LinkUser{}, *h.service)
				}
			}
		}
	})
	log.Println("Start Long Poll")
	if err := h.lp.Run(); err != nil {
		log.Fatal(err)
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
