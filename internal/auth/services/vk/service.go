package vk

import (
	"errors"
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/services"
	"hyneo/pkg/mysql"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
)

type VKService struct {
	Client mysql.Client
	Vk     *api.VK
	Code   code.CodeService
}

func NewVkService(Client *mysql.Client, VK *api.VK, Code code.CodeService) services.Service {
	return &VKService{
		Client: *Client,
		Vk:     VK,
		Code:   Code,
	}
}

func (s *VKService) GetUser(vkID int) (user1 interface{}, err error) {
	var user auth.VkUser
	err = s.Client.DB.Model(&auth.VkUser{}).Where("vk_id = ?", vkID).First(&user).Error
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
		return errors.New("Помощь")
	}
	vkID := messageVK.Message.FromID
	_, err := s.GetUser(vkID)
	if err != nil {
		return err
	}
	mcuser, err := s.GetMCUser(length[1])
	if err != nil {
		return err
	}
	code := s.Code.CreateCode(mcuser.Username, vkID)
	s.SendMessage("Зайдите в игру и введите код: /code "+code, message)
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

func (s *VKService) UnBindAccount(message interface{}) error {
	messageVK := message.(events.MessageNewObject)
	vkId := messageVK.Message.FromID
	user, err := s.GetUser(vkId)
	if err != nil {
		return err
	}
	err = s.Client.DB.Delete(user).Error
	if err != nil {
		return err
	}
	s.clearKeyBoard("Вы успешно отвязали аккаунт", vkId, 0)
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
