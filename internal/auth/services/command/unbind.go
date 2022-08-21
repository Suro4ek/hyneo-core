package command

import (
	"hyneo/internal/auth/services"
	"strconv"
)

var UnBind = &Command{
	Name:    "отвязать",
	Payload: "unlink",
	Alias:   []string{"unlink"},
	Exec: func(message interface{}, userId string, service services.Service) {
		msg := service.GetMessage(message)
		userIdInt, _ := strconv.Atoi(userId)
		user, err := service.GetUserID(int64(userIdInt))
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
