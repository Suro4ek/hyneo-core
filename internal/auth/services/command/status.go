package command

import (
	"hyneo/internal/auth"
	"hyneo/internal/auth/services"
)

var Status = &Command{
	Name:        "status",
	Payload:     "status",
	WithoutUser: false,
	Exec: func(message interface{}, userId *auth.LinkUser, service services.Service) {
		//msg := service.GetMessage(message)
		//ser := service.GetService()

	},
}
