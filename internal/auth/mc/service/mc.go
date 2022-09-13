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

func (s *service) Login(username string, password string) (*auth.User, error) {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
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

func (s *service) Register(user *auth.User) (*auth.User, error) {
	user.PasswordHash = s.service.CreatePassword(user.PasswordHash)
	err := s.client.DB.Create(user).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	return user, nil
}

func (s *service) ChangePassword(username string, oldPassword string, newPassword string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
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

func (s *service) Logout(username string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
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

func (s *service) LastLogin(username string) (string, error) {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		s.log.Error(err)
		return "", mc.UserNotFound
	}
	return s.LeftTime(user.LastJoin), nil
}

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

func (s *service) UnRegister(username string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
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

func (s *service) LeftTime(t time.Time) string {
	return timediff.TimeDiff(t, timediff.WithLocale("ru-RU"))
}

func (s *service) UpdateUser(user *auth.User) (*auth.User, error) {
	user1 := &auth.User{}
	err := s.client.DB.Find(user).Scan(user1).Error
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	s.log.Info(*user)
	err = s.client.DB.Model(user1).Updates(*user).Scan(user1).Error
	s.log.Info(*user1)
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	return user1, nil
}

func (s *service) UpdateLastServer(userId int64, server string) error {
	user := &auth.User{}
	err := s.client.DB.Model(&auth.User{ID: uint32(userId)}).First(user).Error
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
