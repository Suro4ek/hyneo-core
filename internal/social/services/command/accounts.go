package command

import (
	"hyneo/internal/social/services"
	"hyneo/internal/user"
)

var Accounts = &Command{
	Name:        "аккаунты",
	Payload:     "accounts",
	WithoutUser: true,
	Exec: func(message interface{}, user *user.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		service.SendKeyboard("Выберите аккаунт", msg.ChatID)
	},
}
