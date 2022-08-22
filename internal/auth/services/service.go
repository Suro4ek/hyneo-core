package services

import (
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/pkg/mysql"
)

// TODO команды отдельно реализацию, оставить только не команды и клавиатуры и т.д для сервисов
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
	//BindAccount(messageObject interface{}) error

	//UnBindAccount(messageObject interface{}, userId string) error
	//Status(messageObject interface{}, userId string) error
	//Restore(messageObject interface{}, userId string) error
	//Notify(messageObject interface{}, userId string) error
	//Kick(messageObject interface{}, userId string) error
	//Ban(messageObject interface{}, userId string) error
	//Account(messageObject interface{}, userId string) error
	//Accounts(messageObject interface{}) error
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
