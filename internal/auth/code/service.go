package code

import (
	"context"
	"github.com/go-redis/redis/v9"
	"math/rand"
	"time"
)

type Service struct {
	Client *redis.Client
}

type User struct {
	VKId int    `redis:"vk_id"`
	Code string `redis:"code"`
}

//TODO error handler all function

func (c *Service) CreateCode(username string, userId int) string {
	code := c.RandStringRunes(6)
	ctx := context.Background()
	if _, err := c.Client.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "code:"+username, "vk_id", userId)
		rdb.HSet(ctx, "code:"+username, "code", code)
		return nil
	}); err != nil {
		//log error
	}
	c.Client.Expire(ctx, "code:"+username, time.Minute*15)
	return code
}

func (c *Service) CompareCode(username string, code string) bool {
	var user User
	err := c.Client.HGetAll(context.Background(), "code:"+username).Scan(&user)
	if err != nil {
		return false
	}
	if user.Code == code {
		return true
	}
	return false
}

func (c *Service) GetCode(username string) User {
	var user User
	err := c.Client.HGetAll(context.Background(), "code:"+username).Scan(&user)
	if err != nil {
		return User{}
	}
	return user
}

func (c *Service) RemoveCode(username string) {
	ctx := context.Background()
	c.Client.Del(ctx, "code:"+username)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (c *Service) RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
