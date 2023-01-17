package telegram

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/password"
	"hyneo/internal/social/services"
	"hyneo/internal/user"
	"hyneo/pkg/logging"
	"hyneo/pkg/mysql"
)

type telegramService struct {
	bot             *tgbotapi.BotAPI
	Client          *mysql.Client
	Code            *code.Service
	Redis           *redis.Client
	ServiceID       int
	log             *logging.Logger
	PasswordService password.Service
}

func NewTelegramService(client *mysql.Client,
	bot *tgbotapi.BotAPI,
	Code *code.Service,
	redis *redis.Client,
	ServiceID int,
	log *logging.Logger,
	passwordService password.Service) services.Service {
	return &telegramService{
		Client:          client,
		bot:             bot,
		Code:            Code,
		ServiceID:       ServiceID,
		Redis:           redis,
		log:             log,
		PasswordService: passwordService,
	}
}

func (s *telegramService) SendMessage(message string, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := s.bot.Send(msg)
	if err != nil {
		return
	}
}

/*
Получение данных сервиса
*/
func (s *telegramService) GetService() *services.GetService {
	return &services.GetService{
		ServiceID: s.ServiceID,
		Client:    s.Client,
		Code:      s.Code,
		Redis:     s.Redis,
		Password:  s.PasswordService,
	}
}

func (s *telegramService) GetUser(ID int64) (user1 []user.LinkUser, err error) {
	var users []user.LinkUser
	err = s.Client.DB.Model(&user.LinkUser{}).Where(&user.LinkUser{
		ServiceId:     s.ServiceID,
		ServiceUserID: ID,
	}).First(&users).Error
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return users, nil
}

/*
Получение *user.User по никнейму в игре
*/
func (s *telegramService) GetMCUser(username string) (*user.User, error) {
	var user user.User
	err := s.Client.DB.Model(&user.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return &user, nil
}

/*
Получение
*/
func (s *telegramService) GetUserID(userId int64) (user1 *user.LinkUser, err error) {
	var user user.LinkUser
	err = s.Client.DB.Model(&user.LinkUser{}).Joins("User").Where(user.LinkUser{
		UserID: userId,
	}).Find(&user).Error
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return &user, nil
}

func (s *telegramService) GetMessage(messageObject interface{}) services.Message {
	messageTG := messageObject.(*tgbotapi.Message)
	return services.Message{
		Text:   messageTG.Text,
		ChatID: messageTG.Chat.ID,
	}
}

func (s *telegramService) ClearKeyboard(message string, chatID int64) {
	s.SendMessage(message, chatID)
}

func (s *telegramService) AccountKeyboard(message string, chatID int64, user user.LinkUser) {
	msg := tgbotapi.NewMessage(chatID, message)
	keyBoard := s.SoloUserKeyBoard(user)
	keyBoard.InlineKeyboard = append(keyBoard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "accounts"),
		))
	msg.ReplyMarkup = keyBoard
	_, err := s.bot.Send(msg)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *telegramService) SoloUserKeyBoard(user user.LinkUser) tgbotapi.InlineKeyboardMarkup {
	buttons := tgbotapi.NewInlineKeyboardMarkup()
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Информация о аккаунте", fmt.Sprintf("status %d", user.User.ID))))
	rowChange := tgbotapi.NewInlineKeyboardRow()
	if user.Notificated {
		rowChange = append(rowChange, tgbotapi.NewInlineKeyboardButtonData("Отключить уведомления", fmt.Sprintf("notify %d", user.User.ID)))
	} else {
		rowChange = append(rowChange, tgbotapi.NewInlineKeyboardButtonData("Включить уведомления", fmt.Sprintf("notify %d", user.User.ID)))
	}
	if user.Banned {
		rowChange = append(rowChange, tgbotapi.NewInlineKeyboardButtonData("Разбанить", fmt.Sprintf("ban %d", user.User.ID)))
	} else {
		rowChange = append(rowChange, tgbotapi.NewInlineKeyboardButtonData("Забанить", fmt.Sprintf("ban %d", user.User.ID)))
	}
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, rowChange)
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Кикнуть", fmt.Sprintf("kick %d", user.User.ID)),
		tgbotapi.NewInlineKeyboardButtonData("Восставновить", fmt.Sprintf("restore %d", user.User.ID)),
		tgbotapi.NewInlineKeyboardButtonData("Отвязать", fmt.Sprintf("unlink %d", user.User.ID)),
	))
	//for _, keyboardConfig := range keyboard.Keyboard {
	//	row := tgbotapi.NewInlineKeyboardRow()
	//	for _, button := range keyboardConfig.KeyboardButtons {
	//		row = append(row, tgbotapi.NewInlineKeyboardButtonData(button.Name, fmt.Sprintf("%s %d", button.Payload, userId)))
	//	}
	//	buttons.InlineKeyboard = append(buttons.InlineKeyboard, row)
	//}
	return buttons
}

func (s *telegramService) SendKeyboard(message string, chatId int64) {
	var users []user.LinkUser
	err := s.Client.DB.Model(&user.LinkUser{}).Joins("User").Where(user.LinkUser{
		ServiceId:     s.ServiceID,
		ServiceUserID: chatId}).Find(&users).Error
	if err != nil {
		s.log.Error(err)
		return
	}
	var numericKeyboard tgbotapi.InlineKeyboardMarkup
	if len(users) == 1 {
		user := users[0]
		numericKeyboard = s.SoloUserKeyBoard(user)
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
	_, err = s.bot.Send(msg)
	if err != nil {
		s.log.Error(err)
	}
}
