package service

import (
	"hyneo/internal/auth"
	"hyneo/internal/auth/mc"
	"hyneo/internal/auth/password"
	"hyneo/pkg/logging"
	"hyneo/pkg/mysql"
	"time"

	"github.com/mergestat/timediff"
)

type service struct {
	client  *mysql.Client
	service password.Service
	log     *logging.Logger
}

func NewMCService(client *mysql.Client, pservice password.Service, log *logging.Logger) mc.Service {
	return &service{
		client:  client,
		service: pservice,
		log:     log,
	}
}

/*
	Функции авторизации принимающие в себя
	username, password и возращает пользователя,
    если есть такой пользователь и пароль коректный,
    либо ошибки UserNotFound, IncorrectPassword, Fault - ошибка с бд,
    создает сессию на 24 часа и делает пользователя авторизированным
*/
func (s *service) Login(username string, password string) (*auth.User, error) {
	var user auth.User
	err := s.client.DB.Model(&auth.User{}).Where(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.UserNotFound
	}
	if !s.service.ComparePassword(user.PasswordHash, password) {
		s.log.Error(err)
		return nil, mc.IncorrectPassword
	}
	user.Authorized = true
	user.LastJoin = time.Now()
	user.Session = time.Now().Add(time.Hour * 24)
	err = s.client.DB.Save(&user).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	return &user, nil
}

/*
	Функция регистрации принимает в себя пользователя и возращяет пользователя
	берет passwordHash и создает хешированный пароль в sha256 с salt,
    создает сессию на 24 часа и делает пользователя авторизированным
    и возращяет пользователя
*/
func (s *service) Register(user *auth.User) (*auth.User, error) {
	user.PasswordHash = s.service.CreatePassword(user.PasswordHash)
	user.Session = time.Now().Add(time.Hour * 24)
	user.Authorized = true
	user.LastJoin = time.Now()
	err := s.client.DB.Create(user).Where(&auth.User{Username: user.Username}).First(user).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	return user, nil
}

/*
	Функции изменены пароля
	вводные имя пользователя, старый пароль и новый пароль,
	если пароль старый не правильный возращяет ошибку,
	либо если пользователь не найден и если ошибка с базой данных,
	так же хеширует пароль в sha256
*/
func (s *service) ChangePassword(username string, oldPassword string, newPassword string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{}).Where(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		s.log.Error(err)
		return mc.UserNotFound
	}
	if !s.service.ComparePassword(user.PasswordHash, oldPassword) {
		s.log.Error(err)
		return mc.IncorrectPassword
	}
	user.PasswordHash = s.service.CreatePassword(newPassword)
	err = s.client.DB.Save(&user).Error
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}

/*
	Фукнция выхода пользователя
	делает пользователя не авторизованным и сохроняет в бд
	Возращяет ошибку пользователь не найден, либо ошибку с бд
*/
func (s *service) Logout(username string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{}).Where(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		s.log.Error(err)
		return mc.UserNotFound
	}
	user.Authorized = false
	err = s.client.DB.Save(&user).Error
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}

/*
	Возращяет сколько дней назад заходил игрок
	Либо возращяет ошибку, что пользователь не найден
*/
func (s *service) LastLogin(username string) (string, error) {
	var user auth.User
	err := s.client.DB.Model(auth.User{}).Where(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		s.log.Error(err)
		return "", mc.UserNotFound
	}
	return s.LeftTime(user.LastJoin), nil
}

/*
	Функция возращяет пользователя по имя пользователя
	Возращяет так же ошибку, если пользователь не найден
	UserNotFound или ошибку с бд Fault,
	так же проверяет сессию игрока
*/
func (s *service) GetUser(username string) (*auth.User, error) {
	var user auth.User
	err := s.client.DB.Model(&auth.User{}).Where(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.UserNotFound
	}
	if user.Session.Sub(time.Now()) < 0 {
		user.Authorized = false
	}
	err = s.client.DB.Save(&user).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	return &user, nil
}

/*
	Функция удаление пользователя
	Удаляет пользователя по имя пользователя
	и может вернут ошибку UserNotFound - если игрок не найден,
 	или Fault - ошибка бд
*/
func (s *service) UnRegister(username string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{}).Where(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		s.log.Error(err)
		return mc.UserNotFound
	}
	err = s.client.DB.Delete(&user).Error
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}

/*
	Функция превращяющая time.Time в строку и сколько прошло с данного момента
*/
func (s *service) LeftTime(t time.Time) string {
	return timediff.TimeDiff(t, timediff.WithLocale("ru-RU"))
}

/*
	Функция обновления пользователя
	Вводные данные пользователь
	ищет пользователя в бд и обновляет его данные в бд
*/
func (s *service) UpdateUser(user *auth.User) (*auth.User, error) {
	user1 := &auth.User{}
	err := s.client.DB.Model(&auth.User{}).Where(&auth.User{ID: user.ID}).First(user1).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	user.LastJoin = time.Now()
	err = s.client.DB.Model(user1).Where(&auth.User{ID: user.ID}).Omit("id").Updates(*user).First(user1).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	return user1, nil
}

/*
	Функция обновления последнего сервера
*/
func (s *service) UpdateLastServer(userId int64, server string) error {
	user := &auth.User{}
	err := s.client.DB.Model(&auth.User{}).Where(&auth.User{ID: uint32(userId)}).First(user).Error
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	user.LastServer = server
	err = s.client.DB.Save(user).Error
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}
