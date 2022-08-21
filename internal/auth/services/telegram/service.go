package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/services"
	"hyneo/pkg/mysql"
	"strconv"
)

type telegramService struct {
	bot       *tgbotapi.BotAPI
	Client    mysql.Client
	Code      code.Service
	ServiceID int
}

func NewTelegramService(client *mysql.Client, bot *tgbotapi.BotAPI, Code code.Service, ServiceID int) services.Service {
	return &telegramService{
		Client:    *client,
		bot:       bot,
		Code:      Code,
		ServiceID: ServiceID,
	}
}

func (s *telegramService) SendMessage(message string, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, message)
	s.bot.Send(msg)
}

func (s *telegramService) clearKeyBoard(message string, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.bot.Send(msg)
}

func (s *telegramService) GetUser(ID int) (user1 interface{}, err error) {
	var users []auth.LinkUser
	err = s.Client.DB.Model(&auth.LinkUser{}).Where(&auth.LinkUser{
		ServiceId:     s.ServiceID,
		ServiceUserID: int64(ID),
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

func (s *telegramService) CheckCode(username string, code string) error {
	user, err := s.GetMCUser(username)
	if err != nil {
		return err
	}
	TGiD := s.Code.GetCode(username)
	if TGiD == nil {
		return nil
	}
	if TGiD.Service != s.ServiceID {
		return nil
	}
	if !s.Code.CompareCode(username, code) {
		return services.InvalidCode
	}
	vkUser := &auth.LinkUser{
		ServiceUserID: TGiD.UserID,
		User:          *user,
		ServiceId:     s.ServiceID,
	}
	err = s.Client.DB.Save(vkUser).Error
	if err != nil {
		return err
	}
	s.SendKeyboard("Вы успешно привязали аккаунт", TGiD.UserID)
	s.Code.RemoveCode(username)
	return nil
}

func (s *telegramService) UnBindAccount(messageObject interface{}, userId string) error {
	messageTG := messageObject.(*tgbotapi.Message)
	tgID := messageTG.Chat.ID
	userIdInt, _ := strconv.Atoi(userId)
	user, err := s.GetUserID(int64(userIdInt))
	if err != nil {
		return err
	}
	luser := user.(*auth.LinkUser)
	err = s.Client.DB.Delete(luser).Error
	if err != nil {
		return err
	}
	users, err := s.GetUser(int(tgID))
	if users != nil {
		s.SendKeyboard("Вы успешно отвязали аккаунт", tgID)
	} else {
		s.clearKeyBoard("Вы успешно отвязали аккаунт", tgID)
	}
	return nil
}

func (s *telegramService) ClearKeyboard(messageObject interface{}) error {
	return nil
}

func (s *telegramService) Status(messageObject interface{}, userId string) error {
	return nil
}

func (s *telegramService) Restore(messageObject interface{}, userId string) error {
	return nil
}

func (s *telegramService) Notify(messageObject interface{}, userId string) error {
	return nil
}

func (s *telegramService) Kick(messageObject interface{}, userId string) error {
	return nil
}

func (s *telegramService) Ban(messageObject interface{}, userId string) error {
	return nil
}

func (s *telegramService) Account(messageObject interface{}, userId string) error {
	messageTG := messageObject.(*tgbotapi.Message)
	userIdInt, _ := strconv.Atoi(userId)
	luser, err := s.GetUserID(int64(userIdInt))
	if err != nil || luser == nil {
		s.clearKeyBoard("Вы не привязаны к аккаунту ", messageTG.Chat.ID)
		return err
	}
	user := luser.(*auth.LinkUser)
	s.SendAccountKeyBoard("Настройки аккаунта "+user.User.Username, userIdInt, messageTG.Chat.ID)
	return nil
}

func (s *telegramService) Accounts(messageObject interface{}) error {
	s.SendKeyboard("Выберите аккаунт", messageObject.(*tgbotapi.Message).Chat.ID)
	return nil
}

func (s *telegramService) NotifyServer(userId string, server string) error {
	var tgUser auth.LinkUser
	err := s.Client.DB.Joins("User", s.Client.DB.Where("id = ?", userId)).First(&tgUser).Error
	if err != nil {
		return err
	}
	s.SendMessage("Вы подключились к серверу "+server, tgUser.UserID)
	return nil
}

func (s *telegramService) Join(userId string, ip string) error {
	var tgUser auth.LinkUser
	err := s.Client.DB.Joins("User", s.Client.DB.Where("id = ?", userId)).First(&tgUser).Error
	if err != nil {
		return err
	}
	s.SendMessage("Вы подключились к серверу с "+ip, tgUser.UserID)
	return nil
}

func (s *telegramService) SendAccountKeyBoard(message string, userId int, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, message)
	keyboard := s.SoloUserKeyBoard(userId)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "accounts"),
		))
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *telegramService) SoloUserKeyBoard(userId int) tgbotapi.InlineKeyboardMarkup {
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
		numericKeyboard = s.SoloUserKeyBoard(int(user.ID))
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
