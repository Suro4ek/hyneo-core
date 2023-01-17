package user

type Service interface {
	CreateUser(user User) (*User, error)

	GetUserByID(id uint32) (*User, error)
	GetUserByName(username string) (*User, error)

	UpdateUser(id uint32, user User) (*User, error)

	RemoveUser(id uint32) error

	CreateLinkUser(user LinkUser) (*LinkUser, error)

	GetLinkUserByID(id uint32) (*LinkUser, error)
	GetLinkUserByUserID(id int64) (*LinkUser, error)

	UpdateLinkUser(id uint32, user LinkUser) (*LinkUser, error)

	RemoveLinkUser(id uint32) error

	CountAccounts(ip string) (int, error)

	GetLinkedUsers(userId int64) ([]LinkUser, error)
}
