package services

type Service interface {
	SendMessage(message string, messageObject interface{})
	GetUser(ID int) (user interface{}, err error)
	BindAccount(messageObject interface{}) error
	CheckCode(username string, code string) error
	UnBindAccount(ID interface{}) error
	NotifyServer(user_id string, server string) error
	Join(user_id string, ip string) error
}
