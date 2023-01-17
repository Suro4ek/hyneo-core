package command

import (
	"fmt"
	"hyneo/internal/social/services"
	"hyneo/internal/user"
)

var Restore = &Command{
	Name:        "restore",
	Payload:     "restore",
	WithoutUser: false,
	Exec: func(message interface{}, user *user.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		ser := service.GetService()
		usermc, err := service.GetMCUser(user.User.Username)
		if err != nil {
			service.SendMessage("Произошла ошибка", msg.ChatID)
			return
		}
		newPassword := ser.Code.RandStringBytesMaskImprSrcUnsafe(10)
		usermc.PasswordHash = ser.Password.CreatePassword(newPassword)
		usermc.Authorized = false
		err = ser.Client.DB.Save(usermc).Error
		if err != nil {
			service.SendMessage("Произошла ошибка", msg.ChatID)
			return
		}
		service.SendMessage(fmt.Sprintf("Ваш новый пароль: %s", newPassword), msg.ChatID)
	},
}
