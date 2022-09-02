package command

import (
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
)

var UnBind = &Command{
	Name:        "отвязать",
	Payload:     "unlink",
	WithoutUser: false,
	Alias:       []string{"unlink"},
	Exec: func(message interface{}, user *auth.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		s := service.GetService()
		err := s.Client.DB.Delete(&user).Error
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
