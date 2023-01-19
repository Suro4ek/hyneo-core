package storage

import (
	"context"
	"github.com/go-redis/redis/v9"
	"hyneo/internal/user"
	"hyneo/pkg/mysql"
	"strconv"
	"strings"
	"time"
)

type storageUser struct {
	client *mysql.Client
	redis  *redis.Client
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
	var getUser user.User
	err := s.client.DB.
		Model(&user.User{}).
		Where(&user.User{ID: id}).
		First(&getUser).
		Error
	return &getUser, err
}

func (s storageUser) GetUserByName(username string) (*user.User, error) {
	var getUser user.User
	err := s.client.DB.
		Model(&user.User{}).
		Where(&user.User{Username: strings.ToLower(username)}).
		First(&getUser).Error
	return &getUser, err
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
		return 0, user.Fault
	}
	size += len(users)
	err = s.client.DB.Model(&user.User{}).Where(&user.User{IP: ip}).Find(&users).Error
	if err != nil {
		return 0, user.Fault
	}
	size += len(users)
	return size / 2, nil
}

func (s storageUser) CreateLinkUser(user user.LinkUser) (*user.LinkUser, error) {
	err := s.client.DB.Create(&user).Error
	return &user, err
}

func (s storageUser) GetLinkUserByID(id uint32) (*user.LinkUser, error) {
	var getUserLink user.LinkUser
	err := s.client.DB.
		Model(&user.LinkUser{}).
		Where(&user.LinkUser{ID: id}).
		First(&getUserLink).Error
	return &getUserLink, err
}

func (s storageUser) GetLinkUserByUserID(id int64) (*user.LinkUser, error) {
	var getUserLink user.LinkUser
	err := s.client.DB.
		Model(&user.LinkUser{}).
		Joins("User").
		Where(&user.LinkUser{UserID: id}).
		First(&getUserLink).Error
	return &getUserLink, err
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
	var users []user.LinkUser
	err := s.client.DB.Model(&user.LinkUser{}).Where(&user.LinkUser{
		ServiceId:     serviceID,
		ServiceUserID: id,
	}).First(&users).Error
	return &users, err
}

func (s storageUser) GetLinkedUsers(userId int64) ([]user.LinkUser, error) {
	var users []user.LinkUser
	err := s.client.DB.Model(&user.LinkUser{}).Where(&user.LinkUser{UserID: userId}).Find(&users).Error
	if err != nil {
		return nil, user.Fault
	}
	return users, nil
}

func (s storageUser) AddIgnore(userId uint32, ignoreUserId int32) error {
	if ignoreUserId == -1 {
		err := s.client.DB.
			Model(&user.IgnoreUser{}).
			Where(&user.IgnoreUser{
				UserID: userId}).
			Delete(&user.IgnoreUser{}).Error
		if err != nil {
			return err
		}
		err = s.redis.Del(context.TODO(), "ignore:"+strconv.Itoa(int(userId))).Err()
		if err != nil {
			return err
		}
		err = s.redis.HSet(context.TODO(), "ignore:"+strconv.Itoa(int(userId)), strconv.Itoa(int(-1))).Err()
		if err != nil {
			return err
		}
		err = s.redis.Expire(context.TODO(), "ignore:"+strconv.Itoa(int(userId)), time.Second*60).Err()
		return err
	} else {
		err := s.client.DB.Create(&user.IgnoreUser{UserID: userId, IgnoreID: ignoreUserId}).Error
		if err != nil {
			return err
		}
		err = s.redis.HSet(context.TODO(), "ignore:"+strconv.Itoa(int(userId)), strconv.Itoa(int(ignoreUserId))).Err()
		if err != nil {
			return err
		}
		err = s.redis.Expire(context.TODO(), "ignore:"+strconv.Itoa(int(userId)), time.Second*60).Err()
		return err
	}
}

func (s storageUser) RemoveIgnore(userId uint32, ignoreUserId int32) error {
	err := s.client.DB.
		Model(&user.IgnoreUser{}).
		Where(&user.IgnoreUser{
			UserID:   userId,
			IgnoreID: ignoreUserId}).
		Delete(&user.IgnoreUser{}).Error
	if err != nil {
		return err
	}
	err = s.redis.HDel(context.TODO(), "ignore:"+strconv.Itoa(int(userId)), strconv.Itoa(int(ignoreUserId))).Err()
	if err != nil {
		return err
	}
	return err
}

func (s storageUser) GetIgnore(userId uint32) (*[]user.IgnoreUser, error) {
	var users []user.IgnoreUser
	err := s.redis.HGetAll(context.TODO(), "ignore:"+strconv.Itoa(int(userId))).Scan(users)
	if err != nil {
		err := s.client.DB.
			Model(&user.IgnoreUser{}).
			Where(&user.IgnoreUser{UserID: userId}).
			Find(&users).Error
		if err != nil {
			return nil, err
		}
		ctx := context.TODO()
		if _, err := s.redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
			for _, u := range users {
				rdb.HSet(ctx, "ignore:"+strconv.Itoa(int(userId)), strconv.Itoa(int(u.IgnoreID)))
			}
			return nil
		}); err != nil {
			return nil, err
		}
		err = s.redis.Expire(context.TODO(), "ignore:"+strconv.Itoa(int(userId)), time.Second*60).Err()
		if err != nil {
			return nil, err
		}
	}
	return &users, err
}
