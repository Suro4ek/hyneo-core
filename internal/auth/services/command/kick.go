package command

import (
	"context"
	"encoding/json"
	"fmt"
	"hyneo/internal/auth/services"
)

var Kick = &Command{
	Name:    "kick",
	Payload: "kick",
	Exec: func(message interface{}, userId int64, service services.Service) {
		msg := service.GetMessage(message)
		ser := service.GetService()
		out, _ := json.Marshal(services.RedisSend{
			Channel: "kick",
			UserId:  fmt.Sprintf("%d", userId),
			Message: "§cВы кикнуты из игры с помощью бота ВК",
		})
		ser.Redis.Publish(context.Background(), "messenger.bungee", string(out))
		service.SendMessage("Вы кикнули ", msg.ChatID)
	},
}
