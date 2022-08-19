package telegram

import (
	"errors"
	"gorm.io/gorm"
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/services"
	"hyneo/pkg/mysql"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramService struct {
	bot    *tgbotapi.BotAPI
	Client mysql.Client
	Code   code.CodeService
}

func NewTelegramService(client *mysql.Client, bot *tgbotapi.BotAPI, Code code.CodeService) services.Service {
	return &telegramService{
		Client: *client,
		bot:    bot,
		Code:   Code,
	}
}

func (s *telegramService) SendMessage(message string, messageObject interface{}) {
	messageTG := messageObject.(*tgbotapi.Message)
	msg := tgbotapi.NewMessage(messageTG.Chat.ID, message)

	s.bot.Send(msg)
}

func (s *telegramService) clearKeyBoard(message string, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.bot.Send(msg)
}

func (s *telegramService) SendKeyboard(message string, chatId int64) {
	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Убрать клавиатуру"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Статус"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Восстановить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Уведомления"),
			tgbotapi.NewKeyboardButton("Кикнуть"),
			tgbotapi.NewKeyboardButton("Заблокировать"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отвязать"),
		),
	)
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ReplyMarkup = numericKeyboard
	s.bot.Send(msg)
}

func (s *telegramService) GetUser(ID int) (user1 interface{}, err error) {
	var user auth.LinkUser
	err = s.Client.DB.Model(&auth.LinkUser{}).Where("telegram_id = ?", ID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *telegramService) GetMCUser(username string) (*auth.User, error) {
	var user auth.User
	err := s.Client.DB.Model(&auth.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *telegramService) GetUserID(userId int64) (user1 interface{}, err error) {
	var user auth.LinkUser
	err = s.Client.DB.Model(&auth.LinkUser{}).Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *telegramService) BindAccount(message interface{}) error {
	messageTG := message.(*tgbotapi.Message)
	length := strings.Split(messageTG.Text, " ")
	if len(length) != 2 {
		return services.HelpError
	}
	tgID := messageTG.From.ID
	_, err := s.GetUser(int(tgID))
	if err != nil {
		if !errors.As(err, &gorm.ErrRecordNotFound) {
			return err
		}
	}
	mcuser, err := s.GetMCUser(length[1])
	if err != nil {
		return err
	}
	_, err = s.GetUserID(mcuser.ID)
	if err != nil {
		if !errors.As(err, &gorm.ErrRecordNotFound) {
			return err
		}
	}
	createCode := s.Code.CreateCode(mcuser.Username, int(tgID))
	s.SendMessage("Зайдите в игру и введите код: /createCode "+createCode, message)
	return nil
}

func (s *telegramService) CheckCode(username string, code string) error {
	user, err := s.GetMCUser(username)
	if err != nil {
		return err
	}
	if !s.Code.CompareCode(username, code) {
		return services.InvalidCode
	}
	TGiD := s.Code.GetCode(username)
	vkUser := &auth.LinkUser{
		TelegramID: int64(TGiD.VKId),
		User:       *user,
	}
	err = s.Client.DB.Save(vkUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы успешно привязали аккаунт", int64(TGiD.VKId))
	return nil
}
func (s *telegramService) UnBindAccount(message interface{}) error {
	messageTG := message.(*tgbotapi.Message)
	tgID := messageTG.From.ID
	user, err := s.GetUser(int(tgID))
	if err != nil {
		return err
	}
	err = s.Client.DB.Delete(user).Error
	if err != nil {
		return err
	}
	s.clearKeyBoard("Вы успешно отвязали аккаунт", tgID)
	return nil
}

func (s *telegramService) NotifyServer(user_id string, server string) error {
	var tgUser auth.LinkUser
	err := s.Client.DB.Joins("User", s.Client.DB.Where("id = ?", user_id)).First(&tgUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы подключились к серверу "+server, tgUser.TelegramID)
	return nil
}

func (s *telegramService) Join(user_id string, ip string) error {
	var tgUser auth.LinkUser
	err := s.Client.DB.Joins("User", s.Client.DB.Where("id = ?", user_id)).First(&tgUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы подключились к серверу с "+ip, tgUser.TelegramID)
	return nil
}
