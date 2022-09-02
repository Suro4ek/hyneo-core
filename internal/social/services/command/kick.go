package command

import (
	"context"
	"encoding/json"
	"fmt"
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
)

var Kick = &Command{
	Name:        "kick",
	Payload:     "kick",
	WithoutUser: false,
	Exec: func(message interface{}, user *auth.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		ser := service.GetService()
		out, _ := json.Marshal(services.RedisSend{
			Channel: "kick",
			UserId:  fmt.Sprintf("%d", user.UserID),
			Message: "§cВы кикнуты из игры с помощью бота ВК",
		})
		ser.Redis.Publish(context.Background(), "messenger.bungee", string(out))
		service.SendMessage("Вы кикнули ", msg.ChatID)
	},
}
