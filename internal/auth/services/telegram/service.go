package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/services"
	"hyneo/pkg/mysql"
)

type telegramService struct {
	bot       *tgbotapi.BotAPI
	Client    *mysql.Client
	Code      *code.Service
	ServiceID int
}

func NewTelegramService(client *mysql.Client, bot *tgbotapi.BotAPI, Code *code.Service, ServiceID int) services.Service {
	return &telegramService{
		Client:    client,
		bot:       bot,
		Code:      Code,
		ServiceID: ServiceID,
	}
}

func (s *telegramService) SendMessage(message string, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, message)
	s.bot.Send(msg)
}

func (s *telegramService) GetService() *services.GetService {
	return &services.GetService{
		ServiceID: s.ServiceID,
		Client:    s.Client,
		Code:      s.Code,
	}
}

func (s *telegramService) GetUser(ID int64) (user1 []auth.LinkUser, err error) {
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

func (s *telegramService) GetMCUser(username string) (*auth.User, error) {
	var user auth.User
	err := s.Client.DB.Model(&auth.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *telegramService) GetUserID(userId int64) (user1 *auth.LinkUser, err error) {
	var user auth.LinkUser
	err = s.Client.DB.Model(&auth.LinkUser{}).Joins("User").Where(auth.LinkUser{
		UserID: userId,
	}).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *telegramService) GetMessage(messageObject interface{}) services.Message {
	messageTG := messageObject.(*tgbotapi.Message)
	return services.Message{
		Text:   messageTG.Text,
		ChatID: messageTG.From.ID,
	}
}

func (s *telegramService) ClearKeyboard(message string, chatID int64) {
}

func (s *telegramService) AccountKeyboard(message string, chatID int64, userID int64) {
	msg := tgbotapi.NewMessage(chatID, message)
	keyboard := s.SoloUserKeyBoard(userID)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "accounts"),
		))
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *telegramService) SoloUserKeyBoard(userId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статус", fmt.Sprintf("status %d", userId)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Восстановить", fmt.Sprintf("restore %d", userId)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Уведомления", fmt.Sprintf("notify %d", userId)),
			tgbotapi.NewInlineKeyboardButtonData("Кикнуть", fmt.Sprintf("kick %d", userId)),
			tgbotapi.NewInlineKeyboardButtonData("Заблокировать", fmt.Sprintf("ban %d", userId)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отвязать аккаунт", fmt.Sprintf("unlink %d", userId)),
		),
	)
}

func (s *telegramService) SendKeyboard(message string, chatId int64) {
	var users []auth.LinkUser
	s.Client.DB.Model(&auth.LinkUser{}).Joins("User").Where(auth.LinkUser{
		ServiceId:     s.ServiceID,
		ServiceUserID: chatId}).Find(&users)
	var numericKeyboard tgbotapi.InlineKeyboardMarkup
	if len(users) == 1 {
		user := users[0].User
		numericKeyboard = s.SoloUserKeyBoard(user.ID)
	} else {
		numericKeyboard = tgbotapi.NewInlineKeyboardMarkup()
		for _, user := range users {
			numericKeyboard.InlineKeyboard = append(numericKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(user.User.Username, fmt.Sprintf("user %d", user.UserID)),
			))
		}
	}
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ReplyMarkup = numericKeyboard
	s.bot.Send(msg)
}
