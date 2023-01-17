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
)

type telegramService struct {
	bot             *tgbotapi.BotAPI
	Code            *code.Service
	Redis           *redis.Client
	ServiceID       int
	log             *logging.Logger
	PasswordService password.Service
	userService     user.Service
}

func NewTelegramService(
	bot *tgbotapi.BotAPI,
	Code *code.Service,
	redis *redis.Client,
	ServiceID int,
	log *logging.Logger,
	passwordService password.Service,
	userService user.Service,
) services.Service {
	return &telegramService{
		bot:             bot,
		Code:            Code,
		ServiceID:       ServiceID,
		Redis:           redis,
		log:             log,
		PasswordService: passwordService,
		userService:     userService,
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
		Code:      s.Code,
		Redis:     s.Redis,
		Password:  s.PasswordService,
	}
}

func (s *telegramService) GetUser(ID int64) (user1 []user.LinkUser, err error) {
	users, err := s.userService.GetLinkUserByServiceIdAndServiceUserID(ID, s.ServiceID)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return *users, nil
}

/*
Получение *user.User по никнейму в игре
*/
func (s *telegramService) GetMCUser(username string) (*user.User, error) {
	u, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return u, nil
}

/*
Получение
*/
func (s *telegramService) GetUserID(userId int64) (user1 *user.LinkUser, err error) {
	u, err := s.userService.GetLinkUserByUserID(userId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return u, nil
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
	users, err := s.userService.GetLinkUserByServiceIdAndServiceUserID(chatId, s.ServiceID)
	if err != nil {
		s.log.Error(err)
		return
	}
	var numericKeyboard tgbotapi.InlineKeyboardMarkup
	if len(*users) == 1 {
		u := (*users)[0]
		numericKeyboard = s.SoloUserKeyBoard(u)
	} else {
		numericKeyboard = tgbotapi.NewInlineKeyboardMarkup()
		for _, u := range *users {
			numericKeyboard.InlineKeyboard = append(numericKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(u.User.Username, fmt.Sprintf("user %d", u.UserID)),
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
