package command

import (
	"errors"
	"gorm.io/gorm"
	"hyneo/internal/auth"
	"hyneo/internal/auth/services"
	"strings"
)

var Bind = &Command{
	Name:        "привязать",
	Payload:     "-1",
	WithoutUser: true,
	Alias:       []string{"link"},
	Exec: func(message interface{}, user1 *auth.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		length := strings.Fields(msg.Text)
		if len(length) != 2 {
			service.SendMessage("VK - привязать [ник]\n TG - /link [ник]", msg.ChatID)
			return
		}
		users, err := service.GetUser(msg.ChatID)
		if err != nil {
			if !errors.As(err, &gorm.ErrRecordNotFound) {
				service.SendMessage("Не удалось выполнить команду", msg.ChatID)
				return
			}
		}
		if err == nil {
			if len(users) > 2 {
				service.SendMessage("Вы не можете привязать больше двух аккаунтов", msg.ChatID)
				return
			}
		}
		mcuser, err := service.GetMCUser(length[1])
		if err != nil {
			service.SendMessage("Не удалось выполнить команду", msg.ChatID)
			return
		}
		user, err := service.GetUserID(mcuser.ID)
		if err != nil {
			if !errors.As(err, &gorm.ErrRecordNotFound) {
				service.SendMessage("Не удалось выполнить команду", msg.ChatID)
				return
			}
		}
		if err == nil {
			if user.ID != 0 {
				service.SendMessage("Аккаунт уже привязан", msg.ChatID)
				return
			}
		}
		s := service.GetService()
		createCode := s.Code.CreateCode(mcuser.Username, msg.ChatID, s.ServiceID)
		service.SendMessage("Зайдите в игру и введите код: /code "+createCode, msg.ChatID)
	},
}
