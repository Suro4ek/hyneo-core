package command

import (
	"hyneo/internal/auth/services"
)

var Status = &Command{
	Name:    "status",
	Payload: "status",
	Exec: func(message interface{}, userId int64, service services.Service) {
		//msg := service.GetMessage(message)
		//ser := service.GetService()

	},
}
