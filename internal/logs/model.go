package logs

import (
	"gorm.io/gorm"
	"hyneo/internal/user"
)

type Logs struct {
	gorm.Model
	ServerName string `json:"serverName"`
	ActionType string `json:"actionType"`
	Message    string `json:"message"`
	UserID     int64  `json:"user_id" redis:"user_id"`
	User       user.User
}
