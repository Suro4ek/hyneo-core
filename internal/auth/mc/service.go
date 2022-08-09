package mc

import (
	"hyneo/internal/auth"
	"time"
)

type Service interface {
	Login(username string, password string) error
	Register(*auth.User) error
	ChangePassword(username string, old_password string, new_password string) error
	Logout(username string) error
	LastLogin(username string) (string, error)
	GetUser(id string) (*auth.User, error)
	UnRegister(username string) error
	LetfTime(t time.Time) string
}
