package password

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

type Service struct {
}

func NewPasswordService() Service {
	return Service{}
}
func (s *Service) ComparePassword(hash string, password string) bool {
	salt := strings.Split(hash, "$")[2]
	return hash == getHash(password, salt)
}

func (s *Service) CreatePassword(password string) string {
	salt := createSalt()
	return getHash(password, salt)
}

func createSalt() string {
	b := make([]byte, 40)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}
	hasher := sha1.New()
	hasher.Write(b)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	//string substring
	sha = fmt.Sprintf("%x", sha)
	sha = sha[:16]
	return sha
}

func getHash(password string, salt string) string {
	return "$SHA$" + salt + "$" + getSha256(getSha256(password)+salt)
}

func getSha256(password string) string {
	h := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", h)
}
