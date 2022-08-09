package command

import (
	"hyneo/internal/auth/vk"

	"github.com/SevereCloud/vksdk/v2/events"
)

var UnBind = &Command{
	Name:    "отвязать",
	Payload: 5,
	Exec: func(message events.MessageNewObject, service *vk.VKService) {
		err := service.UnBindAccount(message.Message.FromID)
		if err != nil {
			service.SendMessage("Не удалось отвязать аккаунт", message.Message.PeerID, message.Message.RandomID, nil)
		}
	},
}
