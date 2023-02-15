package vk

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/go-redis/redis/v9"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/password"
	"hyneo/internal/social/services"
	"hyneo/internal/user"
	"hyneo/pkg/logging"
	"hyneo/pkg/mysql"
)

type Service struct {
	Vk              *api.VK
	Code            *code.Service
	Redis           *redis.Client
	ServiceID       int
	log             *logging.Logger
	PasswordService password.Service
	userService     user.Service
	client          *mysql.Client
}

func NewVkService(
	VK *api.VK,
	Code *code.Service,
	redis *redis.Client,
	ServiceID int,
	log *logging.Logger,
	passwordService password.Service,
	userService user.Service,
	client *mysql.Client,
) services.Service {
	return &Service{
		Vk:              VK,
		Code:            Code,
		ServiceID:       ServiceID,
		Redis:           redis,
		log:             log,
		PasswordService: passwordService,
		userService:     userService,
		client:          client,
	}
}

func (s *Service) GetMessage(messageObject interface{}) services.Message {
	message := messageObject.(events.MessageNewObject)
	return services.Message{
		Text:   message.Message.Text,
		ChatID: int64(message.Message.PeerID),
	}
}

func (s *Service) GetService() *services.GetService {
	return &services.GetService{
		ServiceID: s.ServiceID,
		Code:      s.Code,
		Redis:     s.Redis,
		Password:  s.PasswordService,
		Client:    s.client,
	}
}

func (s *Service) GetUser(ID int64) (user1 []user.LinkUser, err error) {
	users, err := s.userService.GetLinkUserByServiceIdAndServiceUserID(ID, s.ServiceID)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return *users, nil
}

func (s *Service) GetUserID(userId int64) (user1 *user.LinkUser, err error) {
	u, err := s.userService.GetLinkUserByUserID(userId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return u, nil
}

func (s *Service) GetMCUser(username string) (*user.User, error) {
	u, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return u, nil
}

func (s *Service) SendMessage(message string, chadID int64) {
	m := params.NewMessagesSendBuilder()
	m.Message(message)
	m.PeerID(int(chadID))
	m.RandomID(0)
	_, err := s.Vk.MessagesSend(m.Params)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Service) ClearKeyboard(message string, chadID int64) {
	m := params.NewMessagesSendBuilder()
	m.Message(message)
	m.PeerID(int(chadID))
	k := &object.MessagesKeyboard{
		Buttons: [][]object.MessagesKeyboardButton{},
		OneTime: true,
	}
	m.Keyboard(k)
	m.RandomID(0)
	_, err := s.Vk.MessagesSend(m.Params)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Service) SoloUserKeyBoard(user user.LinkUser) *object.MessagesKeyboard {
	buttons := object.NewMessagesKeyboard(false)
	buttons.AddRow()
	buttons.AddTextButton("Информация о аккаунте", fmt.Sprintf("status %d", user.User.ID), "primary")
	buttons.AddRow()
	if user.Notificated {
		buttons.AddTextButton("Отключить уведомления", fmt.Sprintf("notify %d", user.User.ID), "negative")
	} else {
		buttons.AddTextButton("Включить уведомления", fmt.Sprintf("notify %d", user.User.ID), "positive")
	}
	if user.Banned {
		buttons.AddTextButton("Разбанить", fmt.Sprintf("ban %d", user.User.ID), "positive")
	} else {
		buttons.AddTextButton("Забанить", fmt.Sprintf("ban %d", user.User.ID), "negative")
	}
	buttons.AddRow()
	buttons.AddTextButton("Кикнуть", fmt.Sprintf("kick %d", user.User.ID), "negative")
	buttons.AddTextButton("Восставновить", fmt.Sprintf("restore %d", user.User.ID), "positive")
	buttons.AddTextButton("Отвязать", fmt.Sprintf("unlink %d", user.User.ID), "negative")
	//for _, keyboardConfig := range keyboard.Keyboard {
	//	buttons.AddRow()
	//	for _, button := range keyboardConfig.KeyboardButtons {
	//		buttons.AddTextButton(button.Name, fmt.Sprintf("%s %d", button.Payload, userID), button.Color)
	//	}
	//}
	//keyboard.AddRow()
	//keyboard.AddTextButton("Статус", fmt.Sprintf("status %d", userID), "primary")
	//keyboard.AddRow()
	//keyboard.AddTextButton("Восстановить", fmt.Sprintf("restore %d", userID), "positive")
	//keyboard.AddRow().AddTextButton("Уведомления", fmt.Sprintf("notify %d", userID), "positive").
	//	AddTextButton("Кикнуть", fmt.Sprintf("kick %d", userID), "negative").
	//	AddTextButton("Заблокировать", fmt.Sprintf("ban %d", userID), "negative")
	//keyboard.AddRow()
	//keyboard.AddTextButton("Отвязать", fmt.Sprintf("unlink %d", userID), "negative")
	return buttons
}

func (s *Service) AccountKeyboard(message string, chatID int64, user user.LinkUser) {
	m := params.NewMessagesSendBuilder()
	soloUserKeyBoard := s.SoloUserKeyBoard(user)
	soloUserKeyBoard.AddTextButton("Назад", "accounts", "secondary")
	m.Message(message)
	m.PeerID(int(chatID))
	m.Keyboard(soloUserKeyBoard)
	m.RandomID(0)
	_, err := s.Vk.MessagesSend(m.Params)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Service) SendKeyboard(message string, ChatID int64) {
	m := params.NewMessagesSendBuilder()
	users, err := s.userService.GetLinkUserByServiceIdAndServiceUserID(ChatID, s.ServiceID)
	if err != nil {
		s.log.Error(err)
		return
	}
	messagesKeyboard := object.NewMessagesKeyboard(false)
	if len(*users) == 1 {
		u := (*users)[0]
		messagesKeyboard = s.SoloUserKeyBoard(u)
	} else {
		for _, u := range *users {
			messagesKeyboard.AddRow()
			messagesKeyboard.AddTextButton(u.User.Username, fmt.Sprintf("user %d", u.UserID), "primary")
		}
	}
	m.Message(message)
	m.PeerID(int(ChatID))
	m.Keyboard(messagesKeyboard)
	m.RandomID(0)
	_, err = s.Vk.MessagesSend(m.Params)
	if err != nil {
		s.log.Error(err)
	}
}
