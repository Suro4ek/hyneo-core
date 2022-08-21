package vk

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/services"
	"hyneo/pkg/mysql"
	"log"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
)

type VKService struct {
	Client    mysql.Client
	Vk        *api.VK
	Code      code.Service
	ServiceID int
}

func NewVkService(Client *mysql.Client, VK *api.VK, Code code.Service, ServiceID int) services.Service {
	return &VKService{
		Client:    *Client,
		Vk:        VK,
		Code:      Code,
		ServiceID: ServiceID,
	}
}

func (s *VKService) GetUser(vkID int) (user1 interface{}, err error) {
	var users []auth.LinkUser
	err = s.Client.DB.Model(&auth.LinkUser{}).Where(&auth.LinkUser{
		ServiceId:     s.ServiceID,
		ServiceUserID: int64(vkID),
	}).First(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *VKService) GetUserID(userId int64) (user1 interface{}, err error) {
	var user auth.LinkUser
	err = s.Client.DB.Model(&auth.LinkUser{}).Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *VKService) GetMCUser(username string) (*auth.User, error) {
	var user auth.User
	err := s.Client.DB.Model(&auth.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *VKService) BindAccount(message interface{}) error {
	messageVK := message.(events.MessageNewObject)
	length := strings.Split(messageVK.Message.Text, " ")
	if len(length) != 2 {
		return services.HelpError
	}
	vkID := messageVK.Message.FromID
	vkusers, err := s.GetUser(vkID)
	if err != nil {
		if !errors.As(err, &gorm.ErrRecordNotFound) {
			return err
		}
	}
	if err == nil {
		vk := vkusers.([]auth.LinkUser)
		if len(vk) == 2 {
			return services.MaxAccount
		}
	}

	mcuser, err := s.GetMCUser(length[1])
	if err != nil {
		return err
	}
	user, err := s.GetUserID(mcuser.ID)
	if err != nil {
		if !errors.As(err, &gorm.ErrRecordNotFound) {
			return err
		}
	}
	if err == nil {
		luser := user.(*auth.LinkUser)
		if luser.ID != 0 {
			return services.AlreadyBinded
		}
	}
	code := s.Code.CreateCode(mcuser.Username, vkID, s.ServiceID)
	s.SendMessage("Зайдите в игру и введите код: /code "+code, message)
	return nil
}

func (s *VKService) CheckCode(username string, code string) error {
	user, err := s.GetMCUser(username)
	if err != nil {
		return err
	}
	VkID := s.Code.GetCode(username)
	if VkID == nil {
		return services.InvalidCode
	}
	if VkID.Service != s.ServiceID {
		return nil
	}
	if !s.Code.CompareCode(username, code) {
		return services.InvalidCode
	}
	vkUser := &auth.LinkUser{
		ServiceUserID: int64(VkID.VKId),
		User:          *user,
		ServiceId:     s.ServiceID,
	}
	err = s.Client.DB.Save(vkUser).Error
	if err != nil {
		return err
	}
	s.Code.RemoveCode(username)
	s.SendKeyboard("Вы успешно привязали аккаунт", VkID.VKId, 0)
	return nil
}

func (s *VKService) UnBindAccount(message interface{}, userId string) error {
	messageVK := message.(events.MessageNewObject)
	vkId := messageVK.Message.FromID
	//userId to int
	userIdInt, _ := strconv.Atoi(userId)
	user, err := s.GetUserID(int64(userIdInt))
	if err != nil {
		return err
	}
	err = s.Client.DB.Delete(user).Error
	if err != nil {
		return err
	}
	users, err := s.GetUser(vkId)
	if users != nil {
		s.SendKeyboard("Вы успешно отвязали аккаунт", vkId, 0)
	} else {
		s.clearKeyBoard("Вы успешно отвязали аккаунт", vkId, 0)
	}

	return nil
}

func (s *VKService) NotifyServer(user_id string, server string) error {
	var vkUser auth.LinkUser
	err := s.Client.DB.Joins("User", s.Client.DB.Where("id = ?", user_id)).First(&vkUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы подключились к серверу "+server, int(vkUser.ServiceUserID), 0)
	return nil
}

func (s *VKService) Join(user_id string, ip string) error {
	var vkUser auth.LinkUser
	err := s.Client.DB.Joins("User", s.Client.DB.Where("id = ?", user_id)).First(&vkUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы подключились к серверу с "+ip, int(vkUser.ServiceUserID), 0)
	return nil
}

func (s *VKService) ClearKeyboard(messageObject interface{}) error {
	messageVK := messageObject.(events.MessageNewObject)
	s.clearKeyBoard("Клавиатура выключена", messageVK.Message.FromID, 0)
	return nil
}

func (s *VKService) SendMessage(message string, messageObject interface{}) {
	messageVK := messageObject.(events.MessageNewObject)
	m := params.NewMessagesSendBuilder()
	m.Message(message)
	m.PeerID(messageVK.Message.PeerID)
	m.RandomID(messageVK.Message.RandomID)
	s.Vk.MessagesSend(m.Params)
}

func (s *VKService) clearKeyBoard(message string, peerID int, randomID int) {
	m := params.NewMessagesSendBuilder()
	m.Message(message)
	m.PeerID(peerID)
	keyboard := &object.MessagesKeyboard{
		Buttons: [][]object.MessagesKeyboardButton{},
		OneTime: true,
	}
	m.Keyboard(keyboard)
	m.RandomID(randomID)
	s.Vk.MessagesSend(m.Params)
}

func (s *VKService) Account(messageObject interface{}, userId string) error {
	messageVK := messageObject.(events.MessageNewObject)
	s.SendAccountKeyBoard(userId, "Настройки аккаунта", messageVK.Message.FromID, 0)
	return nil
}

func (s *VKService) Accounts(messageObject interface{}) error {
	messageVK := messageObject.(events.MessageNewObject)
	s.SendKeyboard("Выберите аккаунт", messageVK.Message.FromID, 0)
	return nil
}

func (s *VKService) Ban(messageObject interface{}, userId string) error {
	_ = messageObject.(events.MessageNewObject)
	//
	return nil
}

func (s *VKService) Kick(messageObject interface{}, userId string) error {
	_ = messageObject.(events.MessageNewObject)
	//
	return nil
}

func (s *VKService) Notify(messageObject interface{}, userId string) error {
	_ = messageObject.(events.MessageNewObject)
	//
	return nil
}

func (s *VKService) Restore(messageObject interface{}, userId string) error {
	_ = messageObject.(events.MessageNewObject)
	//
	return nil
}

func (s *VKService) Status(messageObject interface{}, userId string) error {
	_ = messageObject.(events.MessageNewObject)
	//
	return nil
}

func (s *VKService) SendAccountKeyBoard(userId string, message string, peerID int, randomID int) {
	m := params.NewMessagesSendBuilder()
	var user auth.LinkUser
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = s.Client.DB.Model(&auth.LinkUser{}).Joins("User").Where(auth.LinkUser{ID: int64(userIdInt)}).Find(&user).Error
	if err != nil {
		return
	}
	keyboard := object.NewMessagesKeyboard(false)
	keyboard.AddRow()
	keyboard.AddTextButton("Убрать клавиатуру", "clear_keyboard", "secondary")
	keyboard.AddRow()
	keyboard.AddTextButton("Статус", fmt.Sprintf("status %d", user.ID), "primary")
	keyboard.AddRow()
	keyboard.AddTextButton("Восстановить", fmt.Sprintf("restore %d", user.ID), "positive")
	keyboard.AddRow().AddTextButton("Уведомления", fmt.Sprintf("notify %d", user.ID), "positive").
		AddTextButton("Кикнуть", fmt.Sprintf("kick %d", user.ID), "negative").
		AddTextButton("Заблокировать", fmt.Sprintf("ban %d", user.ID), "negative")
	keyboard.AddRow()
	keyboard.AddTextButton("Отвязать", fmt.Sprintf("unlink %d", user.ID), "negative")
	keyboard.AddTextButton("Назад", "accounts", "secondary")
	m.Message(message)
	m.PeerID(peerID)
	m.Keyboard(keyboard)
	m.RandomID(randomID)
	_, err = s.Vk.MessagesSend(m.Params)
	if err != nil {
		log.Println(err)
	}
}

func (s *VKService) SendKeyboard(message string, peerID int, randomID int) {
	m := params.NewMessagesSendBuilder()
	var users []auth.LinkUser
	s.Client.DB.Model(&auth.LinkUser{}).Joins("User").Where(auth.LinkUser{
		ServiceId:     s.ServiceID,
		ServiceUserID: int64(peerID),
	}).Find(&users)
	keyboard := object.NewMessagesKeyboard(false)
	if len(users) == 1 {
		user := users[0].User
		keyboard.AddRow()
		keyboard.AddTextButton("Убрать клавиатуру", "clear_keyboard", "secondary")
		keyboard.AddRow()
		keyboard.AddTextButton("Статус", fmt.Sprintf("status %d", user.ID), "primary")
		keyboard.AddRow()
		keyboard.AddTextButton("Восстановить", fmt.Sprintf("restore %d", user.ID), "positive")
		keyboard.AddRow().AddTextButton("Уведомления", fmt.Sprintf("notify %d", user.ID), "positive").
			AddTextButton("Кикнуть", fmt.Sprintf("kick %d", user.ID), "negative").
			AddTextButton("Заблокировать", fmt.Sprintf("ban %d", user.ID), "negative")
		keyboard.AddRow()
		keyboard.AddTextButton("Отвязать", fmt.Sprintf("unlink %d", user.ID), "negative")
	} else {
		for _, user := range users {
			keyboard.AddRow()
			keyboard.AddTextButton(user.User.Username, fmt.Sprintf("user %d", user.UserID), "primary")
		}
	}
	m.Message(message)
	m.PeerID(peerID)
	m.Keyboard(keyboard)
	m.RandomID(randomID)
	_, err := s.Vk.MessagesSend(m.Params)
	if err != nil {
		log.Println(err)
	}
}
