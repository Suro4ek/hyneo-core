package command

import (
	"hyneo/internal/auth/services"
)

var Account = &Command{
	Name:    "аккаунт",
	Payload: "user",
	Exec: func(message interface{}, userId int64, service services.Service) {
		msg := service.GetMessage(message)
		users, err := service.GetUser(msg.ChatID)
		user, err := service.GetUserID(userId)
		if err != nil || user == nil {
			if users != nil && len(users) == 1 {
				service.AccountKeyboard("Этот аккаунт не привязан к вам", msg.ChatID, users[0].UserID)
			} else {
				service.ClearKeyboard("Этот аккаунт не привязан к вам", msg.ChatID)
			}
			return
		}
		service.AccountKeyboard("Настройки аккаунта "+user.User.Username, msg.ChatID, userId)
	},
}
