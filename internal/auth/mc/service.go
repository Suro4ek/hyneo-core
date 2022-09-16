package mc

import (
	"hyneo/internal/auth"
	"time"
)

type Service interface {
	/*
			Login авторизации принимающие в себя
			username, password и возращает пользователя,
		    если есть такой пользователь и пароль коректный,
		    либо ошибки UserNotFound, IncorrectPassword, Fault - ошибка с бд,
		    создает сессию на 24 часа и делает пользователя авторизированным
	*/
	Login(username string, password string) (*auth.User, error)
	/*
		   	Register регистрации принимает в себя пользователя и возращяет пользователя
			берет passwordHash и создает хешированный пароль в sha256 с salt,
			создает сессию на 24 часа и делает пользователя авторизированным
			и возращяет пользователя
	*/
	Register(*auth.User) (*auth.User, error)
	/*
		ChangePassword изменение пароля
		вводные имя пользователя, старый пароль и новый пароль,
		если пароль старый не правильный возращяет ошибку,
		либо если пользователь не найден и если ошибка с базой данных,
		так же хеширует пароль в sha256
	*/
	ChangePassword(username string, oldPassword string, newPassword string) error
	/*
		Logout выхода пользователя
		делает пользователя не авторизованным и сохроняет в бд
		Возращяет ошибку пользователь не найден, либо ошибку с бд
	*/
	Logout(username string) error
	/*
		LastLogin Возращяет сколько дней назад заходил игрок
		Либо возращяет ошибку, что пользователь не найден
	*/
	LastLogin(username string) (string, error)
	/*
		GetUser возращяет пользователя по имя пользователя
		Возращяет так же ошибку, если пользователь не найден
		UserNotFound или ошибку с бд Fault,
		так же проверяет сессию игрока
	*/
	GetUser(id string) (*auth.User, error)
	/*
	   	UnRegister удаление пользователя
	   	Удаляет пользователя по имя пользователя
	   	и может вернут ошибку UserNotFound - если игрок не найден,
	    или Fault - ошибка бд
	*/
	UnRegister(username string) error
	/*
		LeftTime превращяющая time.Time в строку и сколько прошло с данного момента
	*/
	LeftTime(t time.Time) string
	/*
		UpdateUser обновления пользователя
		Вводные данные пользователь
		ищет пользователя в бд и обновляет его данные в бд
	*/
	UpdateUser(user *auth.User) (*auth.User, error)
	/*
		UpdateLastServer обновления последнего сервера
	*/
	UpdateLastServer(userId int64, server string) error
}
