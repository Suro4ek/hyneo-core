package command

import (
	"hyneo/internal/auth/services"
)

type Command struct {
	Name    string `json:"name"`
	Payload int    `json:"payload"`
	Exec    func(message interface{}, service services.Service)
}
