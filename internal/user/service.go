package user

import (
	"hyneo/pkg/logging"
	"time"
)

type UserService struct {
	userService Service
	log         *logging.Logger
}

func CreateUserService(userService Service, log *logging.Logger) *UserService {
	return &UserService{
		userService: userService,
		log:         log,
	}
}

func (s *UserService) GetUser(username string) (*User, error) {
	userByName, err := s.userService.GetUserByName(username)
	if err != nil {
		s.log.Error(err)
		return nil, UserNotFound
	}
	if userByName.Session.Sub(time.Now()) < 0 {
		userByName.Authorized = false
	}
	_, err = s.userService.UpdateUser(userByName.ID, *userByName)
	if err != nil {
		s.log.Error(err)
		return nil, Fault
	}
	return userByName, nil
}

func (s *UserService) UpdateUser(user *User) (*User, error) {
	if user.ID == 0 {
		return nil, UserNotFound
	}
	user.LastJoin = time.Now()
	u, err := s.userService.UpdateUser(user.ID, *user)
	if err != nil {
		s.log.Error(err)
		return nil, Fault
	}
	return u, nil
}

func (s *UserService) GetLinkedUsers(userId int64) ([]LinkUser, error) {
	return s.userService.GetLinkedUsers(userId)
}

func (s *UserService) AddIgnore(userId uint32, ignoreId int32) error {
	return s.userService.AddIgnore(userId, ignoreId)
}

func (s *UserService) RemoveIgnore(userId uint32, ignoreId int32) error {
	return s.userService.RemoveIgnore(userId, ignoreId)
}

func (s *UserService) IgnoreList(userId uint32) (*[]IgnoreUser, error) {
	return s.userService.GetIgnore(userId)
}
