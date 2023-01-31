package service

import (
	"hyneo/internal/auth/mc"
	"hyneo/internal/auth/password"
	"hyneo/internal/user"
	"hyneo/pkg/logging"
	"time"

	"github.com/mergestat/timediff"
)

type service struct {
	service     password.Service
	log         *logging.Logger
	userService user.Service
}

func NewMCService(pservice password.Service, log *logging.Logger, userService user.Service) mc.Service {
	return &service{
		service:     pservice,
		log:         log,
		userService: userService,
	}
}

func (s *service) Login(username string, password string) (*user.User, error) {
	userByName, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return nil, mc.UserNotFound
	}
	if !s.service.ComparePassword(userByName.PasswordHash, password) {
		s.log.Error(err)
		return nil, mc.IncorrectPassword
	}
	userByName.Authorized = true
	userByName.LastJoin = time.Now()
	userByName.Session = time.Now().Add(time.Hour * 24)
	_, err = s.userService.UpdateUser(userByName.ID, *userByName)
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	return userByName, nil
}

func (s *service) Register(createUser *user.User) (*user.User, error) {
	createUser.PasswordHash = s.service.CreatePassword(createUser.PasswordHash)
	createUser.Session = time.Now().Add(time.Hour * 24)
	createUser.Authorized = true
	createUser.LastJoin = time.Now()
	u, err := s.userService.CreateUser(*createUser)
	countUsers, err := s.userService.CountAccounts(createUser.RegisteredIP)
	if countUsers >= 4 {
		return nil, mc.AccountsLimit
	}
	if err != nil {
		s.log.Error(err)
		return nil, mc.Fault
	}
	return u, nil
}

func (s *service) ChangePassword(username string, oldPassword string, newPassword string) error {
	userByName, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return mc.UserNotFound
	}
	if !s.service.ComparePassword(userByName.PasswordHash, oldPassword) {
		s.log.Error(err)
		return mc.IncorrectPassword
	}
	userByName.PasswordHash = s.service.CreatePassword(newPassword)
	_, err = s.userService.UpdateUser(userByName.ID, *userByName)
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}

func (s *service) ChangePasswordConsole(username string, newPassword string) error {
	userByName, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return mc.UserNotFound
	}
	userByName.PasswordHash = s.service.CreatePassword(newPassword)
	_, err = s.userService.UpdateUser(userByName.ID, *userByName)
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}

func (s *service) Logout(username string) error {
	userByName, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return mc.UserNotFound
	}
	userByName.Authorized = false
	_, err = s.userService.UpdateUser(userByName.ID, *userByName)
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}

func (s *service) LastLogin(username string) (string, error) {
	userByName, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return "", mc.UserNotFound
	}
	return s.LeftTime(userByName.LastJoin), nil
}

func (s *service) UnRegister(username string) error {
	userByName, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return mc.UserNotFound
	}
	err = s.userService.RemoveUser(userByName.ID)
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}

func (s *service) LeftTime(t time.Time) string {
	return timediff.TimeDiff(t, timediff.WithLocale("ru-RU"))
}

func (s *service) UpdateLastServer(userId int64, server string) error {
	u, err := s.userService.GetUserByID(userId)
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	u.LastServer = server
	_, err = s.userService.UpdateUser(userId, *u)
	if err != nil {
		s.log.Error(err)
		return mc.Fault
	}
	return nil
}
