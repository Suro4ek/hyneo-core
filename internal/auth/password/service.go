package password

type Service struct {
	argon2ID Argon2ID
}

func NewPasswordService() Service {
	argon2ID := NewArgon2ID()
	return Service{
		argon2ID: argon2ID,
	}
}

/*
Функция проверки хешированного пароля и обычного пароля
*/
func (s *Service) ComparePassword(hash string, password string) bool {
	compare, _ := s.argon2ID.Verify(password, hash)
	return compare
}

/*
Функция создания пароля
*/
func (s *Service) CreatePassword(password string) string {
	hash, _ := s.argon2ID.Hash(password)
	return hash
}

///*
//Функция создания соли
//*/
//func createSalt() string {
//	b := make([]byte, 40)
//	_, err := rand.Read(b)
//	// Note that err == nil only if we read len(b) bytes.
//	if err != nil {
//		panic(err)
//	}
//	hasher := sha1.New()
//	hasher.Write(b)
//	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
//	//string substring
//	sha = fmt.Sprintf("%x", sha)
//	sha = sha[:16]
//	return sha
//}
//
///*
//Функция получение хеша с помощью пароля и соли
//*/
//func getHash(password string, salt string) string {
//	return "$SHA$" + salt + "$" + getSha256(getSha256(password)+salt)
//}
//
///*
//Функция получение хеша с помощью строки
//*/
//func getSha256(password string) string {
//	h := sha256.Sum256([]byte(password))
//	return fmt.Sprintf("%x", h)
//}
