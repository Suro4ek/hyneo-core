package command

import (
	"hyneo/internal/auth/services"
)

var UnBind = &Command{
	Name:    "отвязать",
	Payload: "unlink",
	Alias:   []string{"unlink"},
	Exec: func(message interface{}, userId int64, service services.Service) {
		msg := service.GetMessage(message)
		user, err := service.GetUserID(userId)
		if err != nil {
			service.SendMessage("Не удалось отвязать аккаунт", msg.ChatID)
			return
		}
		s := service.GetService()
		err = s.Client.DB.Delete(&user).Error
		if err != nil {
			service.SendMessage("Не удалось отвязать аккаунт", msg.ChatID)
			return
		}
		users, err := service.GetUser(msg.ChatID)
		if users != nil {
			service.SendKeyboard("Аккаунт отвязан "+user.User.Username, msg.ChatID)
		} else {
			service.ClearKeyboard("Аккаунт отвязан "+user.User.Username, msg.ChatID)
		}
	},
}
