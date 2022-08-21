package code

import (
	"context"
	"github.com/go-redis/redis/v9"
	"math/rand"
	"time"
	"unsafe"
)

type Service struct {
	Client *redis.Client
}

type User struct {
	Service int    `redis:"service"`
	UserID  int64  `redis:"user_id"`
	Code    string `redis:"code"`
}

//TODO error handler all function

func (c *Service) CreateCode(username string, userId int64, service int) string {
	code := c.RandStringBytesMaskImprSrcUnsafe(6)
	ctx := context.Background()
	if _, err := c.Client.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, "code:"+username, "user_id", userId)
		rdb.HSet(ctx, "code:"+username, "code", code)
		rdb.HSet(ctx, "code:"+username, "service", service)
		return nil
	}); err != nil {
		//log error
	}
	c.Client.Expire(ctx, "code:"+username, time.Minute*15)
	return code
}

func (c *Service) CompareCode(username string, code string) bool {
	user := c.GetCode(username)
	if user == nil {
		return false
	}
	if user.Code == code {
		return true
	}
	return false
}

func (c *Service) GetCode(username string) *User {
	var user User
	err := c.Client.HGetAll(context.Background(), "code:"+username).Scan(&user)
	if err != nil {
		return nil
	}
	return &user
}

func (c *Service) RemoveCode(username string) {
	ctx := context.Background()
	c.Client.Del(ctx, "code:"+username)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func (c *Service) RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
