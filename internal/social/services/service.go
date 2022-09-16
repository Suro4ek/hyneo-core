package services

import (
	"github.com/go-redis/redis/v9"
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/password"
	"hyneo/pkg/mysql"
)

type Service interface {
	// SendMessage Отправка сообщения пользователю
	SendMessage(message string, chatID int64)
	//SendKeyboard Отправка клавиатуры может быть и с несколькими пользователями
	SendKeyboard(message string, chatID int64)
	//AccountKeyboard Отправка клавиатуры конкретного пользователя
	AccountKeyboard(message string, chatID int64, userID int64)
	//ClearKeyboard Очистка клавиатуры
	ClearKeyboard(message string, chatID int64)

	// GetUser Получение пользователей []auth.LinkUser по ID пользователя сети
	GetUser(ID int64) (user []auth.LinkUser, err error)
	//GetMessage получить Message по messageObject
	GetMessage(messageObject interface{}) (message Message)
	//GetMCUser получить *auth.User по никнейму из игры
	GetMCUser(username string) (*auth.User, error)
	//GetUserID получить *auth.LinkUser по id строки в бд
	GetUserID(userId int64) (user *auth.LinkUser, err error)
	//GetService получить *GetService
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
	Redis     *redis.Client
	Password  password.Service
}

type RedisSend struct {
	Channel string `json:"channel"`
	UserId  string `json:"userId"`
	Message string `json:"message"`
}
