package vk

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/services"
	"hyneo/pkg/mysql"
	"log"
)

type Service struct {
	Client    *mysql.Client
	Vk        *api.VK
	Code      *code.Service
	ServiceID int
}

func NewVkService(Client *mysql.Client, VK *api.VK, Code *code.Service, ServiceID int) services.Service {
	return &Service{
		Client:    Client,
		Vk:        VK,
		Code:      Code,
		ServiceID: ServiceID,
	}
}

func (s *Service) GetMessage(messageObject interface{}) services.Message {
	message := messageObject.(*object.MessagesMessage)
	return services.Message{
		Text:   message.Text,
		ChatID: int64(message.PeerID),
	}
}

func (s *Service) GetService() *services.GetService {
	return &services.GetService{
		ServiceID: s.ServiceID,
		Client:    s.Client,
		Code:      s.Code,
	}
}

func (s *Service) GetUser(ID int64) (user1 []auth.LinkUser, err error) {
	var users []auth.LinkUser
	err = s.Client.DB.Model(&auth.LinkUser{}).Where(&auth.LinkUser{
		ServiceId:     s.ServiceID,
		ServiceUserID: ID,
	}).First(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) GetUserID(userId int64) (user1 *auth.LinkUser, err error) {
	var user auth.LinkUser
	err = s.Client.DB.Model(&auth.LinkUser{}).Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetMCUser(username string) (*auth.User, error) {
	var user auth.User
	err := s.Client.DB.Model(&auth.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) SendMessage(message string, chadID int64) {
	m := params.NewMessagesSendBuilder()
	m.Message(message)
	m.PeerID(int(chadID))
	m.RandomID(0)
	s.Vk.MessagesSend(m.Params)
}

func (s *Service) ClearKeyboard(message string, chadID int64) {
	m := params.NewMessagesSendBuilder()
	m.Message(message)
	m.PeerID(int(chadID))
	keyboard := &object.MessagesKeyboard{
		Buttons: [][]object.MessagesKeyboardButton{},
		OneTime: true,
	}
	m.Keyboard(keyboard)
	m.RandomID(0)
	s.Vk.MessagesSend(m.Params)
}

func (s *Service) SoloUserKeyBoard(userID int64) *object.MessagesKeyboard {
	keyboard := object.NewMessagesKeyboard(false)
	keyboard.AddRow()
	keyboard.AddTextButton("Убрать клавиатуру", "clear_keyboard", "secondary")
	keyboard.AddRow()
	keyboard.AddTextButton("Статус", fmt.Sprintf("status %d", userID), "primary")
	keyboard.AddRow()
	keyboard.AddTextButton("Восстановить", fmt.Sprintf("restore %d", userID), "positive")
	keyboard.AddRow().AddTextButton("Уведомления", fmt.Sprintf("notify %d", userID), "positive").
		AddTextButton("Кикнуть", fmt.Sprintf("kick %d", userID), "negative").
		AddTextButton("Заблокировать", fmt.Sprintf("ban %d", userID), "negative")
	keyboard.AddRow()
	keyboard.AddTextButton("Отвязать", fmt.Sprintf("unlink %d", userID), "negative")
	return keyboard
}

func (s *Service) AccountKeyboard(message string, chatID int64, userID int64) {
	m := params.NewMessagesSendBuilder()
	keyboard := s.SoloUserKeyBoard(userID)
	keyboard.AddTextButton("Назад", "accounts", "secondary")
	m.Message(message)
	m.PeerID(int(chatID))
	m.Keyboard(keyboard)
	m.RandomID(0)
	s.Vk.MessagesSend(m.Params)
}

func (s *Service) SendKeyboard(message string, ChatID int64) {
	m := params.NewMessagesSendBuilder()
	var users []auth.LinkUser
	s.Client.DB.Model(&auth.LinkUser{}).Joins("User").Where(auth.LinkUser{
		ServiceId:     s.ServiceID,
		ServiceUserID: ChatID,
	}).Find(&users)
	keyboard := object.NewMessagesKeyboard(false)
	if len(users) == 1 {
		user := users[0].User
		keyboard = s.SoloUserKeyBoard(user.ID)
	} else {
		for _, user := range users {
			keyboard.AddRow()
			keyboard.AddTextButton(user.User.Username, fmt.Sprintf("user %d", user.UserID), "primary")
		}
	}
	m.Message(message)
	m.PeerID(int(ChatID))
	m.Keyboard(keyboard)
	m.RandomID(0)
	_, err := s.Vk.MessagesSend(m.Params)
	if err != nil {
		log.Println(err)
	}
}
