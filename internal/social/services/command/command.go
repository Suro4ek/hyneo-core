package command

import (
	"hyneo/internal/auth"
	"hyneo/internal/social/services"
)

type Command struct {
	Name        string   `json:"name"`
	Payload     string   `json:"payload"`
	WithoutUser bool     `json:"without_user"`
	Alias       []string `json:"alias"`
	Exec        func(message interface{}, userId *auth.LinkUser, service services.Service)
}
