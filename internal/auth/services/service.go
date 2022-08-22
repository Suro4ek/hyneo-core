package services

import (
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/pkg/mysql"
)

type Service interface {
	SendMessage(message string, chatID int64)
	SendKeyboard(message string, chatID int64)
	AccountKeyboard(message string, chatID int64, userID int64)
	ClearKeyboard(message string, chatID int64)

	GetUser(ID int64) (user []auth.LinkUser, err error)
	GetMessage(messageObject interface{}) (message Message)
	GetMCUser(username string) (*auth.User, error)
	GetUserID(userId int64) (user *auth.LinkUser, err error)
	GetService() (service *GetService)
}

type Message struct {
	Text   string `json:"text"`
	ChatID int64  `json:"chat_id"`
}

type GetService struct {
	ServiceID int
	Client    *mysql.Client
	Code      *code.Service
}
