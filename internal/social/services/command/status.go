package command

import (
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
)

var Status = &Command{
	Name:        "status",
	Payload:     "status",
	WithoutUser: false,
	Exec: func(message interface{}, userId *auth.LinkUser, service services.Service) {
		msg := service.GetMessage(message)
		user, err := service.GetUserID(userId.UserID)
		if err != nil {
			service.SendMessage("Не удалось выполнить команду", msg.ChatID)
			return
		}
		str := "Аккаунт привязан к " + user.User.Username + "\n"
		if user.Notificated {
			str += "Уведомления включены"
		} else {
			str += "Уведомления выключены"
		}
		if user.Banned {
			str += "\nВы забанены"
		} else {
			str += "\nВы не забанены"
		}
		if user.DoubleAuth {
			str += "\nДвухфакторная аутентификация включена"
		} else {
			str += "\nДвухфакторная аутентификация выключена"
		}
		service.SendMessage(str, msg.ChatID)
	},
}
