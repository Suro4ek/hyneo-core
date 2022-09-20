package command

import (
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
)

var Notify = &Command{
	Name:        "notify",
	Payload:     "notify",
	WithoutUser: false,
	Exec: func(message interface{}, user *auth.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		ser := service.GetService()
		err := ser.Client.DB.Model(&auth.LinkUser{}).Where("user_id = ?", user.UserID).
			Update("notificated", !user.Notificated).Error
		if err != nil {
			service.SendMessage("Не удалось выполнить команду", msg.ChatID)
			return
		}
		user.Notificated = !user.Notificated
		if user.Notificated {
			service.AccountKeyboard("Вы отписались от уведомлений", msg.ChatID, *user)
		} else {
			service.AccountKeyboard("Вы подписались на уведомления", msg.ChatID, *user)
		}
	},
}
