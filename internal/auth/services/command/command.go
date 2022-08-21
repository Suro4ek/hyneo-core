package command

import (
	"hyneo/internal/auth/services"
)

type Command struct {
	Name    string   `json:"name"`
	Payload string   `json:"payload"`
	Alias   []string `json:"alias"`
	Exec    func(message interface{}, userid string, service services.Service)
}
