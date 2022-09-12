package mc

import (
	"hyneo/internal/auth"
	"time"
)

type Service interface {
	Login(username string, password string) (*auth.User, error)
	Register(*auth.User) (*auth.User, error)
	ChangePassword(username string, oldPassword string, newPassword string) error
	Logout(username string) error
	LastLogin(username string) (string, error)
	GetUser(id string) (*auth.User, error)
	UnRegister(username string) error
	LeftTime(t time.Time) string
	UpdateUser(user *auth.User) (*auth.User, error)
}
