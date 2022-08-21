package command

import "hyneo/internal/auth/services"

var ClearKeyboard = &Command{
	Name:    "очистить клавиатуру",
	Payload: "clear_keyboard",
	Exec: func(message interface{}, userId string, service services.Service) {
		msg := service.GetMessage(message)
		service.ClearKeyboard("Клавиатура очищена", msg.ChatID)
	},
}
