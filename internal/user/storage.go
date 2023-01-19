package user

type Service interface {
	CreateUser(user User) (*User, error)

	GetUserByID(id int64) (*User, error)
	GetUserByName(username string) (*User, error)

	UpdateUser(id int64, user User) (*User, error)

	RemoveUser(id int64) error

	CreateLinkUser(user LinkUser) (*LinkUser, error)

	GetLinkUserByID(id uint32) (*LinkUser, error)
	GetLinkUserByUserID(id int64) (*LinkUser, error)

	GetLinkUserByServiceIdAndServiceUserID(id int64, serviceID int) (*[]LinkUser, error)

	UpdateLinkUser(id uint32, user LinkUser) (*LinkUser, error)

	RemoveLinkUser(id uint32) error

	CountAccounts(ip string) (int, error)

	GetLinkedUsers(userId int64) ([]LinkUser, error)

	AddIgnore(userId int64, ignoreUserId int64) error
	RemoveIgnore(userId int64, ignoreUserId int64) error
	GetIgnore(userId int64) (*[]IgnoreUser, error)
}
