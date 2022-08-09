package service

import (
	"fmt"
	"hyneo/internal/auth"
	"hyneo/internal/auth/mc"
	"hyneo/internal/auth/password"
	"hyneo/pkg/mysql"
	"time"

	"github.com/mergestat/timediff"
)

type service struct {
	client   *mysql.Client
	pservice password.PasswordService
}

func NewMCService(client *mysql.Client, pservice password.PasswordService) mc.Service {
	return &service{
		client:   client,
		pservice: pservice,
	}
}

func (s *service) Login(username string, password string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		return err
	}
	if !s.pservice.ComparePassword(user.PasswordHash, password) {
		return fmt.Errorf("password is not correct")
	}
	user.Authorized = true
	err = s.client.DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Register(user *auth.User) error {
	user.PasswordHash = s.pservice.CreatePassword(user.PasswordHash)
	err := s.client.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ChangePassword(username string, old_password string, new_password string) error {
	var user auth.User
	err := s.client.DB.Model(&auth.User{Username: username}).First(&user).Error
	if err != nil {
		return err
	}
	if !s.pservice.ComparePassword(user.PasswordHash, old_password) {
		return fmt.Errorf("password is not correct")
	}
	user.PasswordHash = s.pservice.CreatePassword(new_password)
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
	return s.LetfTime(user.LastJoin), nil
}

func (s *service) GetUser(id string) (*auth.User, error) {
	var user auth.User
	err := s.client.DB.Model(&auth.User{}).Where("id = ?", id).First(&user).Error
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

func (s *service) LetfTime(t time.Time) string {
	return timediff.TimeDiff(t, timediff.WithLocale("ru-RU"))
}
