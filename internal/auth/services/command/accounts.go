package command

import (
	"hyneo/internal/auth"
	"hyneo/internal/auth/services"
)

var Accounts = &Command{
	Name:        "аккаунты",
	Payload:     "accounts",
	WithoutUser: true,
	Exec: func(message interface{}, user *auth.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		service.SendKeyboard("Выберите аккаунт", msg.ChatID)
	},
}
