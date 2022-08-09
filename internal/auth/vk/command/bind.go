package command

import (
	"hyneo/internal/auth/vk"

	"github.com/SevereCloud/vksdk/v2/events"
)

var Bind = &Command{
	Name:    "привязать",
	Payload: -1,
	Exec: func(message events.MessageNewObject, service *vk.VKService) {
		err := service.BindAccount(message)
		if err != nil {
			service.SendMessage("Не удалось привязать аккаунт", message.Message.PeerID, message.Message.RandomID, nil)
		}
	},
}
