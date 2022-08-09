package vk

import (
	"context"
	"log"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type handler struct {
	lp *longpoll.LongPoll
}

func NewVKHandler(lp *longpoll.LongPoll) *handler {
	return &handler{
		lp: lp,
	}
}

func (h *handler) Message() {
	h.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
	})
}
