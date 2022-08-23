package command

import (
	"context"
	"hyneo/internal/auth/services"
)

var Kick = &Command{
	Name:    "kick",
	Payload: "kick",
	Exec: func(message interface{}, userId int64, service services.Service) {
		msg := service.GetMessage(message)
		ser := service.GetService()
		ser.Redis.Publish(context.Background(), "kick", userId)
		service.SendMessage("Вы кикнули ", msg.ChatID)
	},
}
