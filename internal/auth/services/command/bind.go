package command

import (
	"errors"
	"hyneo/internal/auth/services"
)

var Bind = &Command{
	Name:    "привязать",
	Payload: -1,
	Exec: func(message interface{}, service services.Service) {
		err := service.BindAccount(message)
		if err != nil {
			if errors.As(err, &services.HelpError) {
				service.SendMessage("Помощь по командам", message)
			} else {
				service.SendMessage("Не удалось привязать аккаунт", message)
			}
		}
	},
}
