package command

import (
	"hyneo/internal/auth/services"
)

var UnBind = &Command{
	Name:    "отвязать",
	Payload: 5,
	Exec: func(message interface{}, service services.Service) {
		err := service.UnBindAccount(message)
		if err != nil {
			service.SendMessage("Не удалось отвязать аккаунт", message)
		}
	},
}
