package command

import (
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
)

var Account = &Command{
	Name:        "аккаунт",
	Payload:     "user",
	WithoutUser: false,
	Exec: func(message interface{}, user *auth.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		users, err := service.GetUser(msg.ChatID)
		if err != nil || user == nil {
			if users != nil && len(users) == 1 {
				service.AccountKeyboard("Этот аккаунт не привязан к вам", msg.ChatID, users[0])
			} else {
				service.ClearKeyboard("Этот аккаунт не привязан к вам", msg.ChatID)
			}
			return
		}
		service.AccountKeyboard("Настройки аккаунта "+user.User.Username, msg.ChatID, *user)
	},
}
