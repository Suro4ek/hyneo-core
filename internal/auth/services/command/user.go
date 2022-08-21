package command

import (
	"hyneo/internal/auth/services"
	"strconv"
)

var Account = &Command{
	Name:    "аккаунт",
	Payload: "user",
	Exec: func(message interface{}, userId string, service services.Service) {
		msg := service.GetMessage(message)
		userIdInt, _ := strconv.Atoi(userId)
		user, err := service.GetUserID(int64(userIdInt))
		if err != nil {
			service.SendMessage("Не удалось отвязать аккаунт", msg.ChatID)
			return
		}
	},
}
