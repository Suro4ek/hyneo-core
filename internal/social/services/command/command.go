package command

import (
	"hyneo/internal/social/services"
	"hyneo/internal/user"
)

type Command struct {
	Name        string   `json:"name"`
	Payload     string   `json:"payload"`
	WithoutUser bool     `json:"without_user"`
	Alias       []string `json:"alias"`
	Exec        func(message interface{}, userId *user.LinkUser, service services.Service)
}
