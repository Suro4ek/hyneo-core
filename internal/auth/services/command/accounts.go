package command

import (
	"hyneo/internal/auth/services"
)

var Accounts = &Command{
	Name:    "аккаунты",
	Payload: "accounts",
	Exec: func(message interface{}, userId string, service services.Service) {
		msg := service.GetMessage(message)
		service.SendKeyboard("Выберите аккаунт", msg.ChatID)
	},
}
