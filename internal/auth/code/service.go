package code

import (
	"math/rand"
)

type CodeService struct {
}

type CodeUser struct {
	VKId int
	Code string
}

var codes = make(map[string]CodeUser, 0)

func (c *CodeService) CreateCode(username string, userId int) string {
	code := c.RandStringRunes(6)
	codes[username] = CodeUser{
		VKId: userId,
		Code: code,
	}
	return code
}

func (c *CodeService) CompareCode(username string, code string) bool {
	if _, ok := codes[username]; !ok {
		return false
	}
	if codes[username].Code == code {
		return true
	}
	return false
}

func (c *CodeService) GetCode(username string) CodeUser {
	return codes[username]
}

func (c *CodeService) RemoveCode(username string) {
	delete(codes, username)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (c *CodeService) RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
