package auth

import "time"

type User struct {
	ID           uint32    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	RegisteredIP string    `json:"registered_ip"`
	IP           string    `json:"ip"`
	Session      time.Time `json:"session"`
	LastJoin     time.Time `json:"last_join"`
	LastServer   string    `json:"last_server"`
	Authorized   bool      `json:"authorized"`
}

//TODO подумать насчет DoubleAuth оставить его или нет
type LinkUser struct {
	ID            int64 `json:"id" redis:"id"`
	ServiceId     int   `json:"service_id" redis:"service_id"`
	ServiceUserID int64 `json:"service_user_id" redis:"service_user_id"`
	Notificated   bool  `json:"notificated" redis:"notificated"`
	Banned        bool  `json:"banned" redis:"banned"`
	DoubleAuth    bool  `json:"double_auth" redis:"double_auth"`
	UserID        int64 `json:"user_id" redis:"user_id"`
	User          User
}
