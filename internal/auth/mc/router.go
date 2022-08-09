package mc

import (
	"context"
	"hyneo/eu.suro/hyneo/protos/auth"
	auth2 "hyneo/internal/auth"
	"hyneo/pkg/mysql"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type routerService struct {
	client  *mysql.Client
	service Service
	auth.UnimplementedAuthServer
}

func NewAuthRouter(client *mysql.Client, service Service) auth.AuthServer {
	return &routerService{
		client:  client,
		service: service,
	}
}

func (r *routerService) Login(ctx context.Context, res *auth.LoginRequest) (*auth.Empty, error) {
	err := r.service.Login(res.User.Username, res.User.Password)
	if err != nil {
		return nil, err
	}
	return &auth.Empty{}, nil
}

func (r *routerService) Register(ctx context.Context, res *auth.RegisterRequest) (*auth.Empty, error) {
	err := r.service.Register(&auth2.User{
		Username:     res.User.Username,
		PasswordHash: res.User.Password,
		LastJoin:     time.Now(),
		Lisense:      res.User.Lisense,
		Authorized:   true,
		Session:      time.Now().Add(24 * time.Hour),
		IP:           res.User.Ip,
		RegisteredIP: res.User.RegisteredIp,
		LastServer:   res.User.LastServer,
	})
	if err != nil {
		return nil, err
	}
	return &auth.Empty{}, nil
}

func (r *routerService) ChangePassword(ctx context.Context, res *auth.ChangePasswordRequest) (*auth.Empty, error) {
	err := r.service.ChangePassword(res.Username, res.OldPassword, res.NewPassword)
	if err != nil {
		return nil, err
	}
	return &auth.Empty{}, nil
}

func (r *routerService) Logout(ctx context.Context, res *auth.LogoutRequest) (*auth.Empty, error) {
	err := r.service.Logout(res.Username)
	if err != nil {
		return nil, err
	}
	return &auth.Empty{}, nil
}

func (r *routerService) LastLogin(ctx context.Context, res *auth.LastLoginRequest) (*auth.LastLoginResponse, error) {
	lastLogin, err := r.service.LastLogin(res.Username)
	if err != nil {
		return nil, err
	}
	return &auth.LastLoginResponse{
		LastLogin: lastLogin,
	}, nil
}

func (r *routerService) GetUser(ctx context.Context, res *auth.GetUserRequest) (*auth.GetUserResponse, error) {
	user, err := r.service.GetUser(res.Username)
	if err != nil {
		return nil, err
	}
	return &auth.GetUserResponse{
		User: &auth.User{
			Username:     user.Username,
			LastLogin:    r.service.LetfTime(user.LastJoin),
			Lisense:      user.Lisense,
			Ip:           user.IP,
			RegisteredIp: user.RegisteredIP,
			LastServer:   user.LastServer,
			Session:      timestamppb.New(user.Session),
			Auth:         user.Authorized,
		},
	}, nil
}

func (r *routerService) UnRegister(ctx context.Context, res *auth.UnRegisterRequest) (*auth.Empty, error) {
	err := r.service.UnRegister(res.Username)
	if err != nil {
		return nil, err
	}
	return &auth.Empty{}, nil
}
