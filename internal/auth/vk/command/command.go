package command

import (
	"hyneo/internal/auth/vk"

	"github.com/SevereCloud/vksdk/v2/events"
)

type Command struct {
	Name    string `json:"name"`
	Payload int    `json:"payload"`
	Exec    func(message events.MessageNewObject, service *vk.VKService)
}
