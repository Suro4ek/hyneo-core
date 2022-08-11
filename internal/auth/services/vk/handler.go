package vk

import (
	"context"
	"hyneo/internal/auth/services/command"
	"strings"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type handler struct {
	lp      *longpoll.LongPoll
	service *VKService
}

func NewVKHandler(lp *longpoll.LongPoll, service *VKService) *handler {
	return &handler{
		lp:      lp,
		service: service,
	}
}

func (h *handler) Message() {
	h.lp.MessageNew(func(_ context.Context, m events.MessageNewObject) {
		mstr := strings.TrimSpace(m.Message.Text)
		marray := strings.Fields(mstr)
		if cmd, ok := command.GetCommands()[strings.ToLower(marray[0])]; ok {
			if cmd.Payload == -1 {
				go cmd.Exec(m, h.service)
			} else {
				if cmd.Payload == m.Message.PeerID {
					go cmd.Exec(m, h.service)
				}
			}
		}
	})
}
