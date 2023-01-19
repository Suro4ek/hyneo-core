package storage

import (
	"context"
	"fmt"
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

func CreateStorageUser(client *mysql.Client, redis *redis.Client) user.Service {
	return &storageUser{
		client: client,
		redis:  redis,
	}
}

func (s storageUser) CreateUser(user user.User) (*user.User, error) {
	err := s.client.DB.Create(&user).Error
	return &user, err
}

func (s storageUser) GetUserByID(id int64) (*user.User, error) {
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

func (s storageUser) UpdateUser(id int64, updateUser user.User) (*user.User, error) {
	err := s.client.DB.Model(&user.User{}).Where(&user.User{ID: id}).Updates(updateUser).First(&updateUser).Error
	return &updateUser, err
}

func (s storageUser) RemoveUser(id int64) error {
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

func (s storageUser) AddIgnore(userId int64, ignoreUserId int64) error {
	if ignoreUserId == -1 {
		err := s.client.DB.
			Model(&user.IgnoreUser{}).
			Where(&user.IgnoreUser{
				UserID: userId}).
			Delete(&user.IgnoreUser{}).Error
		if err != nil {
			return err
		}
		err = s.redis.Del(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores").Err()
		if err != nil {
			return err
		}
		err = s.redis.HSet(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores", strconv.Itoa(-1), "null").Err()
		if err != nil {
			return err
		}
		err = s.redis.Expire(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores", time.Second*60).Err()
		return err
	} else {
		u := &user.IgnoreUser{}
		err := s.client.DB.Create(&user.IgnoreUser{UserID: userId, IgnoreID: ignoreUserId}).Error
		if err != nil {
			return err
		}
		err = s.client.DB.Model(&user.IgnoreUser{}).Where(&user.IgnoreUser{UserID: userId, IgnoreID: ignoreUserId}).Joins("User").First(u).Error
		if err != nil {
			return err
		}
		err = s.redis.HSet(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores", strconv.Itoa(int(ignoreUserId)), u.IgnoreUser.Username).Err()
		if err != nil {
			return err
		}
		err = s.redis.Expire(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores", time.Second*60).Err()
		return err
	}
}

func (s storageUser) RemoveIgnore(userId int64, ignoreUserId int64) error {
	err := s.client.DB.
		Model(&user.IgnoreUser{}).
		Where(&user.IgnoreUser{
			UserID:   userId,
			IgnoreID: ignoreUserId}).
		Delete(&user.IgnoreUser{}).Error
	if err != nil {
		return err
	}
	err = s.redis.HDel(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores", strconv.Itoa(int(ignoreUserId))).Err()
	if err != nil {
		return err
	}
	return err
}

func (s storageUser) GetIgnore(userId int64) (*[]user.IgnoreUser, error) {
	var users []user.IgnoreUser
	keys, err := s.redis.HKeys(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores").Result()
	if err != nil {
		err := s.client.DB.
			Model(&user.IgnoreUser{}).
			Where(&user.IgnoreUser{UserID: userId}).
			Joins("User").
			Find(&users).Error
		if err != nil {
			return nil, err
		}
		ctx := context.TODO()
		if _, err := s.redis.Pipelined(ctx, func(rdb redis.Pipeliner) error {
			for _, u := range users {
				rdb.HSet(ctx, "user:"+strconv.Itoa(int(userId))+":ignores", strconv.Itoa(int(u.IgnoreID)), u.IgnoreUser.Username)
			}
			return nil
		}); err != nil {
			return nil, err
		}
		err = s.redis.Expire(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores", time.Second*60).Err()
		if err != nil {
			return nil, err
		}
	} else {
		for _, key := range keys {
			useranem, _ := s.redis.HGet(context.TODO(), "user:"+strconv.Itoa(int(userId))+":ignores", key).Result()
			fmt.Println(useranem)
			keyInt, _ := strconv.ParseInt(key, 10, 64)
			users = append(users, user.IgnoreUser{IgnoreID: int64(keyInt), IgnoreUser: user.User{Username: useranem}})
		}
	}
	return &users, err
}
