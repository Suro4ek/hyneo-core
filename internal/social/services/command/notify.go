package command

import (
	"context"
	"fmt"
	"hyneo/internal/social/services"
	"hyneo/internal/user"
)

var Notify = &Command{
	Name:        "notify",
	Payload:     "notify",
	WithoutUser: false,
	Exec: func(message interface{}, user *user.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		ser := service.GetService()
		user.Notificated = !user.Notificated
		_, err := ser.User.UpdateLinkUser(user.ID, *user)
		if err != nil {
			service.SendMessage("Не удалось выполнить команду", msg.ChatID)
			return
		}
		user.Notificated = !user.Notificated
		ser.Redis.HSet(context.Background(), fmt.Sprintf("link:%d", user.UserID), "notificated", user.Banned)
		if user.Notificated {
			service.AccountKeyboard("Вы подписались на уведомления", msg.ChatID, *user)
		} else {
			service.AccountKeyboard("Вы отписались от уведомлений", msg.ChatID, *user)
		}
	},
}
