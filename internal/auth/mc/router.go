package mc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	auth2 "hyneo/internal/auth"
	"hyneo/pkg/mysql"
	"hyneo/protos/auth"
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

func (r *routerService) Login(_ context.Context, res *auth.LoginRequest) (*auth.User, error) {
	user, err := r.service.Login(res.User.Username, res.Password)
	if err != nil {
		return nil, err
	}
	return &auth.User{
		Id:           user.ID,
		Username:     user.Username,
		LastLogin:    timestamppb.New(user.LastJoin),
		Ip:           user.IP,
		RegisteredIp: user.RegisteredIP,
		LastServer:   user.LastServer,
		Auth:         user.Authorized,
		LocaleId:     0,
	}, nil
}

func (r *routerService) Register(_ context.Context, res *auth.RegisterRequest) (*auth.User, error) {
	user, err := r.service.Register(&auth2.User{
		Username:     res.User.Username,
		PasswordHash: res.Password,
		LastJoin:     time.Now(),
		Authorized:   true,
		Session:      time.Now().Add(24 * time.Hour),
		IP:           res.User.Ip,
		RegisteredIP: res.User.RegisteredIp,
		LastServer:   res.User.LastServer,
	})
	if err != nil {
		return nil, err
	}
	return &auth.User{
		Id:           user.ID,
		Username:     user.Username,
		LastLogin:    timestamppb.New(user.LastJoin),
		Ip:           user.IP,
		RegisteredIp: user.RegisteredIP,
		LastServer:   user.LastServer,
		Auth:         user.Authorized,
		LocaleId:     0,
	}, nil
}

func (r *routerService) ChangePassword(_ context.Context, res *auth.ChangePasswordRequest) (*emptypb.Empty, error) {
	err := r.service.ChangePassword(res.Username, res.OldPassword, res.NewPassword)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (r *routerService) Logout(_ context.Context, res *auth.LogoutRequest) (*emptypb.Empty, error) {
	err := r.service.Logout(res.Username)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (r *routerService) LastLogin(_ context.Context, res *auth.LastLoginRequest) (*auth.LastLoginResponse, error) {
	lastLogin, err := r.service.LastLogin(res.Username)
	if err != nil {
		return nil, err
	}
	return &auth.LastLoginResponse{
		LastLogin: lastLogin,
	}, nil
}

func (r *routerService) GetUser(_ context.Context, res *auth.GetUserRequest) (*auth.GetUserResponse, error) {
	user, err := r.service.GetUser(res.Username)
	if err != nil {
		return nil, err
	}
	return &auth.GetUserResponse{
		User: &auth.User{
			Username:     user.Username,
			LastLogin:    timestamppb.New(user.LastJoin),
			Ip:           user.IP,
			RegisteredIp: user.RegisteredIP,
			LastServer:   user.LastServer,
			Auth:         user.Authorized,
			LocaleId:     0,
		},
	}, nil
}

func (r *routerService) UnRegister(_ context.Context, res *auth.UnRegisterRequest) (*emptypb.Empty, error) {
	err := r.service.UnRegister(res.Username)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (r *routerService) UpdateUser(_ context.Context, res *auth.UpdateUserRequest) (*auth.User, error) {
	user, err := r.service.UpdateUser(&auth2.User{
		ID:         res.User.Id,
		Username:   res.User.Username,
		LastJoin:   res.User.LastLogin.AsTime(),
		Authorized: res.User.Auth,
		IP:         res.User.Ip,
		LastServer: res.User.LastServer,
	})
	if err != nil {
		return nil, err
	}
	return &auth.User{
		Id:           user.ID,
		Username:     user.Username,
		LastLogin:    timestamppb.New(user.LastJoin),
		Ip:           user.IP,
		RegisteredIp: user.RegisteredIP,
		LastServer:   user.LastServer,
		Auth:         user.Authorized,
		LocaleId:     0,
	}, nil
}

func (r *routerService) UpdateLastServer(ctx context.Context, res *auth.UpdateLastServerRequest) (*emptypb.Empty, error) {
	err := r.service.UpdateLastServer(int64(res.GetUserId()), res.GetLastServer())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
