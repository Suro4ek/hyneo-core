package service

import (
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hyneo/internal/auth"
	"hyneo/internal/auth/mc"
	"hyneo/internal/auth/password"
	"hyneo/pkg/mysql"
	"time"

	"github.com/mergestat/timediff"
)

type service struct {
	client  *mysql.Client
	service password.Service
}

var (
	IncorrectPassword = status.New(codes.Unauthenticated, "incorrect password").Err()
)

func NewMCService(client *mysql.Client, pservice password.Service) mc.Service {
	return &service{
		client:  client,
		service: pservice,
	}
}

func (s *service) Login(username string, password string) (*auth.User, error) {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}
	if !s.service.ComparePassword(user.PasswordHash, password) {
		return nil, IncorrectPassword
	}
	user.Authorized = true
	user.LastJoin = time.Now()
	user.Session = time.Now().Add(time.Hour * 24)
	err = s.client.DB.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) Register(user *auth.User) (*auth.User, error) {
	user.PasswordHash = s.service.CreatePassword(user.PasswordHash)
	err := s.client.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) ChangePassword(username string, oldPassword string, newPassword string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		return err
	}
	if !s.service.ComparePassword(user.PasswordHash, oldPassword) {
		return IncorrectPassword
	}
	user.PasswordHash = s.service.CreatePassword(newPassword)
	err = s.client.DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Logout(username string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		return err
	}
	user.Authorized = false
	err = s.client.DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *service) LastLogin(username string) (string, error) {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		return "", err
	}
	return s.LeftTime(user.LastJoin), nil
}

func (s *service) GetUser(id string) (*auth.User, error) {
	var user auth.User
	err := s.client.DB.Model(&auth.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	if user.Session.Sub(time.Now()) < 0 {
		user.Authorized = false
	}
	err = s.client.DB.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) UnRegister(username string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		return err
	}
	err = s.client.DB.Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *service) LeftTime(t time.Time) string {
	return timediff.TimeDiff(t, timediff.WithLocale("ru-RU"))
}
