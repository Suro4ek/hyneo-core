package storage

import (
	"hyneo/internal/auth/mc"
	"hyneo/internal/user"
	"hyneo/pkg/mysql"
	"strings"
)

type storageUser struct {
	client *mysql.Client
}

func CreateStorageUser(client *mysql.Client) user.Service {
	return &storageUser{
		client: client,
	}
}

func (s storageUser) CreateUser(user user.User) (*user.User, error) {
	err := s.client.DB.Create(&user).Error
	return &user, err
}

func (s storageUser) GetUserByID(id uint32) (*user.User, error) {
	var getUser *user.User
	err := s.client.DB.
		Model(&user.User{}).
		Where(&user.User{ID: id}).
		First(getUser).Error
	return getUser, err
}

func (s storageUser) GetUserByName(username string) (*user.User, error) {
	var getUser *user.User
	err := s.client.DB.
		Model(&user.User{}).
		Where(&user.User{Username: strings.ToLower(username)}).
		First(getUser).Error
	return getUser, err
}

func (s storageUser) UpdateUser(id uint32, updateUser user.User) (*user.User, error) {
	err := s.client.DB.Model(&user.User{}).Where(&user.User{ID: id}).Updates(updateUser).First(&updateUser).Error
	return &updateUser, err
}

func (s storageUser) RemoveUser(id uint32) error {
	err := s.client.DB.Delete(&user.User{}, id).Error
	return err
}

// TODO Надо чекнуть
func (s storageUser) CountAccounts(ip string) (int, error) {
	var size = 0
	var users []user.User
	err := s.client.DB.Model(&user.User{}).Where(&user.User{RegisteredIP: ip}).Find(&users).Error
	if err != nil {
		return 0, mc.Fault
	}
	size += len(users)
	err = s.client.DB.Model(&user.User{}).Where(&user.User{IP: ip}).Find(&users).Error
	if err != nil {
		return 0, mc.Fault
	}
	size += len(users)
	return size / 2, nil
}

func (s storageUser) CreateLinkUser(user user.LinkUser) (*user.LinkUser, error) {
	err := s.client.DB.Create(&user).Error
	return &user, err
}

func (s storageUser) GetLinkUserByID(id uint32) (*user.LinkUser, error) {
	var getUserLink *user.LinkUser
	err := s.client.DB.
		Model(&user.LinkUser{}).
		Where(&user.LinkUser{ID: id}).
		First(getUserLink).Error
	return getUserLink, err
}

func (s storageUser) GetLinkUserByUserID(id int64) (*user.LinkUser, error) {
	var getUserLink *user.LinkUser
	err := s.client.DB.
		Model(&user.LinkUser{}).
		Joins("User").
		Where(&user.LinkUser{UserID: id}).
		First(getUserLink).Error
	return getUserLink, err
}

func (s storageUser) UpdateLinkUser(id uint32, updateUser user.LinkUser) (*user.LinkUser, error) {
	err := s.client.DB.Model(&user.LinkUser{}).Where(&user.LinkUser{ID: id}).Updates(updateUser).First(&updateUser).Error
	return &updateUser, err
}

func (s storageUser) RemoveLinkUser(id uint32) error {
	err := s.client.DB.Delete(&user.LinkUser{}, id).Error
	return err
}

func (s storageUser) GetLinkUserByServiceIdAndServiceUserID(id int64, serviceID int) (*[]user.LinkUser, error) {
	var users *[]user.LinkUser
	err := s.client.DB.Model(&user.LinkUser{}).Where(&user.LinkUser{
		ServiceId:     serviceID,
		ServiceUserID: id,
	}).First(&users).Error
	return users, err
}

func (s storageUser) GetLinkedUsers(userId int64) ([]user.LinkUser, error) {
	var users []user.LinkUser
	err := s.client.DB.Model(&user.LinkUser{}).Where(&user.LinkUser{UserID: userId}).Find(&users).Error
	if err != nil {
		return nil, mc.Fault
	}
	return users, nil
}
