package vk

import (
	"context"
	"fmt"
	"hyneo/internal/auth/services"
	"hyneo/internal/auth/services/command"
	"log"
	"strings"

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
		fmt.Println(strings.ReplaceAll(m.Message.Payload, "\"", ""))
		cmd, userId := h.GetCommandByPayload(strings.ReplaceAll(m.Message.Payload, "\"", ""))
		if cmd != nil {
			go cmd.Exec(m, userId, *h.service)
		} else {
			if cmd, ok := command.GetCommands()[strings.ToLower(marray[0])]; ok {
				if cmd.Payload == "-1" {
					go cmd.Exec(m, "", *h.service)
				}
			}
		}
	})
	log.Println("Start Long Poll")
	if err := h.lp.Run(); err != nil {
		log.Fatal(err)
	}
}

func (h *handler) GetCommandByPayload(payload string) (cmd *command.Command, userId string) {
	for _, cmd := range command.GetCommands() {
		if strings.HasPrefix(payload, cmd.Payload) {
			return cmd, strings.TrimSpace(payload[len(cmd.Payload):])
		}
	}
	return nil, ""
}
