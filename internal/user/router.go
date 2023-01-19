package user

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"hyneo/protos/auth"
)

type userRouter struct {
	service UserService
	auth.UnimplementedUserServiceServer
}

func NewUserRouter(service UserService) auth.UserServiceServer {
	return &userRouter{
		service: service,
	}
}

func (r userRouter) UpdateUser(_ context.Context, res *auth.UpdateUserRequest) (*auth.User, error) {
	u, err := r.service.UpdateUser(convertGRPUserToUser(res.GetUser()))
	if err != nil {
		return nil, err
	}
	return convertUserToGRPCUser(u), nil
}

func (r userRouter) GetUser(ctx context.Context, res *auth.GetUserRequest) (*auth.GetUserResponse, error) {
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

func (r userRouter) AddIgnore(ctx context.Context, res *auth.AddIgnoreRequest) (*emptypb.Empty, error) {
	err := r.service.AddIgnore(int64(res.GetUserId()), res.GetIgnoreId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (r userRouter) RemoveIgnore(ctx context.Context, res *auth.RemoveIgnoreRequest) (*emptypb.Empty, error) {
	err := r.service.RemoveIgnore(int64(res.GetUserId()), res.GetIgnoreId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (r userRouter) GetIgnoreList(ctx context.Context, res *auth.GetIgnoreListRequest) (*auth.GetIgnoreListResponse, error) {
	users, err := r.service.IgnoreList(int64(res.GetUserId()))
	if err != nil {
		return nil, err
	}
	return &auth.GetIgnoreListResponse{
		IgnoreList: convertIgnoreListToGRPC(*users),
	}, nil
}

//convert ignoreList to auth.IgnoreUser[]
func convertIgnoreListToGRPC(ignoreList []IgnoreUser) []*auth.IgnoreUser {
	var ignoreListGRPC []*auth.IgnoreUser
	for _, ignore := range ignoreList {
		ignoreListGRPC = append(ignoreListGRPC,
			&auth.IgnoreUser{
				Id:       ignore.IgnoreID,
				Username: ignore.IgnoreUser.Username,
			})
	}
	return ignoreListGRPC
}

func convertUserToGRPCUser(user *User) *auth.User {
	return &auth.User{
		Id:           int64(user.ID),
		Username:     user.Username,
		LastLogin:    timestamppb.New(user.LastJoin),
		Ip:           user.IP,
		RegisteredIp: user.RegisteredIP,
		LastServer:   user.LastServer,
		Auth:         user.Authorized,
		LocaleId:     0,
	}
}

func convertGRPUserToUser(authUser *auth.User) *User {
	return &User{
		ID:         int64(authUser.Id),
		Username:   authUser.Username,
		LastJoin:   authUser.LastLogin.AsTime(),
		Authorized: authUser.Auth,
		IP:         authUser.Ip,
		LastServer: authUser.LastServer,
	}
}
