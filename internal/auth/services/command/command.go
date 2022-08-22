package command

import (
	"hyneo/internal/auth/services"
)

type Command struct {
	Name    string   `json:"name"`
	Payload string   `json:"payload"`
	Alias   []string `json:"alias"`
	Exec    func(message interface{}, userId int64, service services.Service)
}
