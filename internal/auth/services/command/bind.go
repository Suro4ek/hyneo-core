package command

import (
	"hyneo/internal/auth/services"
)

var Bind = &Command{
	Name:    "привязать",
	Payload: -1,
	Exec: func(message interface{}, service services.Service) {
		err := service.BindAccount(message)
		if err != nil {
			service.SendMessage("Не удалось привязать аккаунт", message)
		}
	},
}
