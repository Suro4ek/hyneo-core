package mc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"hyneo/internal/user"
	"hyneo/protos/auth"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type routerService struct {
	service Service
	auth.UnimplementedAuthServer
}

func NewAuthRouter(service Service) auth.AuthServer {
	return &routerService{
		service: service,
	}
}

func (r *routerService) Login(_ context.Context, res *auth.LoginRequest) (*auth.User, error) {
	u, err := r.service.Login(res.User.Username, res.Password)
	if err != nil {
		return nil, err
	}
	return convertUserToGRPCUser(u), nil
}

func (r *routerService) Register(_ context.Context, res *auth.RegisterRequest) (*auth.User, error) {
	authUser := convertGRPUserToUser(res.GetUser())
	authUser.LastJoin = time.Now()
	authUser.Session = time.Now().Add(24 * time.Hour)
	authUser.Authorized = true
	u, err := r.service.Register(authUser)
	if err != nil {
		return nil, err
	}
	return convertUserToGRPCUser(u), nil
}

func (r *routerService) ChangePassword(_ context.Context, res *auth.ChangePasswordRequest) (*emptypb.Empty, error) {
	err := r.service.ChangePassword(res.Username, res.OldPassword, res.NewPassword)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (r *routerService) ChangePasswordConsole(_ context.Context, res *auth.ChangePasswordConsoleRequest) (*emptypb.Empty, error) {
	err := r.service.ChangePasswordConsole(res.Username, res.NewPassword)
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
	u, err := r.service.GetUser(res.Username)
	if err != nil {
		return nil, err
	}
	linked := false
	users, err := r.service.GetLinkedUsers(int64(u.ID))
	if err == nil && len(users) > 0 {
		linked = true
	}
	authUser := convertUserToGRPCUser(u)
	authUser.Linked = linked
	return &auth.GetUserResponse{
		User: authUser,
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
	//TODO check is null??? GetUser()
	u, err := r.service.UpdateUser(convertGRPUserToUser(res.GetUser()))
	if err != nil {
		return nil, err
	}
	return convertUserToGRPCUser(u), nil
}

func (r *routerService) UpdateLastServer(ctx context.Context, res *auth.UpdateLastServerRequest) (*emptypb.Empty, error) {
	err := r.service.UpdateLastServer(int64(res.GetUserId()), res.GetLastServer())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func convertUserToGRPCUser(user *user.User) *auth.User {
	return &auth.User{
		Id:           user.ID,
		Username:     user.Username,
		LastLogin:    timestamppb.New(user.LastJoin),
		Ip:           user.IP,
		RegisteredIp: user.RegisteredIP,
		LastServer:   user.LastServer,
		Auth:         user.Authorized,
		LocaleId:     0,
	}
}

func convertGRPUserToUser(authUser *auth.User) *user.User {
	return &user.User{
		ID:         authUser.Id,
		Username:   authUser.Username,
		LastJoin:   authUser.LastLogin.AsTime(),
		Authorized: authUser.Auth,
		IP:         authUser.Ip,
		LastServer: authUser.LastServer,
	}
}
