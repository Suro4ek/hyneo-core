package command

import (
	"hyneo/internal/auth/services"
)

var Account = &Command{
	Name:    "аккаунт",
	Payload: "user",
	Exec: func(message interface{}, userId int64, service services.Service) {
		msg := service.GetMessage(message)
		user, err := service.GetUserID(userId)
		if err != nil || user == nil {
			service.ClearKeyboard("Вы не привязаны к аккаунт", msg.ChatID)
			return
		}
		service.AccountKeyboard("Настройки аккаунта "+user.User.Username, msg.ChatID, userId)
	},
}
