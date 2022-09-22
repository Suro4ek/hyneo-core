package command

import (
	"context"
	"encoding/json"
	"fmt"
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
)

var Ban = &Command{
	Name:        "ban",
	Payload:     "ban",
	WithoutUser: false,
	Exec: func(message interface{}, user *auth.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		ser := service.GetService()
		err := ser.Client.DB.Model(&auth.LinkUser{}).Where("user_id = ?", user.UserID).
			Update("banned", !user.Banned).Error
		if err != nil {
			service.SendMessage("Не удалось выполнить команду", msg.ChatID)
			return
		}
		user.Banned = !user.Banned
		if !user.Banned {
			out, _ := json.Marshal(services.RedisSend{
				Channel: "unban",
				UserId:  fmt.Sprintf("%d", user.UserID),
				Message: "§cВы разбанены на сервере с помощью бота",
			})
			ser.Redis.Publish(context.Background(), "messenger.bungee", string(out))
			service.AccountKeyboard("Вы разбанили "+user.User.Username, msg.ChatID, *user)
		} else {
			out, _ := json.Marshal(services.RedisSend{
				Channel: "ban",
				UserId:  fmt.Sprintf("%d", user.UserID),
				Message: "§cВы забанены на сервере с помощью бота",
			})
			ser.Redis.Publish(context.Background(), "messenger.bungee", string(out))
			service.AccountKeyboard("Вы забанили "+user.User.Username, msg.ChatID, *user)
		}

	},
}
