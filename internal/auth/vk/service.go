package vk

import (
	"errors"
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/pkg/mysql"
	"strings"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/params"
	"github.com/SevereCloud/vksdk/object"
	"github.com/SevereCloud/vksdk/v2/events"
)

type VKService struct {
	Client mysql.Client
	Vk     api.VK
	Code   code.CodeService
}

func (s *VKService) GetUser(vkID int) (*auth.VkUser, error) {
	var user auth.VkUser
	err := s.Client.DB.Model(&auth.VkUser{}).Where("vk_id = ?", vkID).First(&user).Error
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

func (s *VKService) BindAccount(message events.MessageNewObject) error {
	length := strings.Split(message.Message.Text, " ")
	if len(length) != 2 {
		return errors.New("Помощь")
	}
	vkID := message.Message.FromID
	_, err := s.GetUser(vkID)
	if err != nil {
		return err
	}
	mcuser, err := s.GetMCUser(length[1])
	if err != nil {
		return err
	}
	code := s.Code.CreateCode(mcuser.Username, vkID)
	s.SendMessage("Зайдите в игру и введите код: /code "+code, message.Message.PeerID, message.Message.RandomID, nil)
	return nil
}

func (s *VKService) CheckCode(username string, code string) error {
	user, err := s.GetMCUser(username)
	if err != nil {
		return err
	}
	if !s.Code.CompareCode(username, code) {
		return errors.New("Неверный код")
	}
	VkID := s.Code.GetCode(username)
	vkUser := &auth.VkUser{
		VkID: int64(VkID.VKId),
		User: *user,
	}
	err = s.Client.DB.Save(vkUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы успешно привязали аккаунт", VkID.VKId, 0)
	return nil
}

func (s *VKService) UnBindAccount(vkId int) error {
	user, err := s.GetUser(vkId)
	if err != nil {
		return err
	}
	err = s.Client.DB.Delete(user).Error
	if err != nil {
		return err
	}
	s.SendMessage("Вы успешно отвязали аккаунт", vkId, 0, &object.MessagesKeyboard{
		Buttons: [][]object.MessagesKeyboardButton{},
		OneTime: true,
	})
	return nil
}

func (s *VKService) NotifyServer(user_id string, server string) error {
	var vkUser auth.VkUser
	err := s.Client.DB.Joins("User", s.Client.DB.Where("id = ?", user_id)).First(&vkUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы подключились к серверу "+server, int(vkUser.VkID), 0)
	return nil
}

func (s *VKService) Join(user_id string, ip string) error {
	var vkUser auth.VkUser
	err := s.Client.DB.Joins("User", s.Client.DB.Where("id = ?", user_id)).First(&vkUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы подключились к серверу с "+ip, int(vkUser.VkID), 0)
	return nil
}

func (s *VKService) SendMessage(message string, peerID int, randomID int, keyboard *object.MessagesKeyboard) {
	m := params.NewMessagesSendBuilder()
	m.Message(message)
	m.PeerID(peerID)
	if keyboard != nil {
		m.Keyboard(keyboard)
	}
	m.RandomID(randomID)
	s.Vk.MessagesSend(m.Params)
}

func (s *VKService) SendKeyboard(message string, peerID int, randomID int) {
	m := params.NewMessagesSendBuilder()
	m.Message(message)
	m.PeerID(peerID)
	keyboard := object.NewMessagesKeyboard(false)
	keyboard.AddTextButton("Убрать клавиатуру", "6", "secondary")
	keyboard.AddTextButton("Статус", "4", "primary")
	keyboard.AddTextButton("Восстановить", "3", "positive")
	keyboard.AddRow().AddTextButton("Уведомления", "2", "positive").
		AddTextButton("Кикнуть", "1", "negative").
		AddTextButton("Заблокировать", "7", "negative")
	keyboard.AddTextButton("Отвязать", "5", "negative")
	m.Keyboard(keyboard)
	m.RandomID(randomID)
	s.Vk.MessagesSend(m.Params)
}
