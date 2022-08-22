package services

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"hyneo/internal/auth"
	"hyneo/pkg/mysql"
	"hyneo/protos/service"
	"strconv"
)

type serviceRouter struct {
	Client   mysql.Client
	services []Service
	service.UnimplementedServiceServer
}

func NewServiceRouter(client mysql.Client, services []Service) service.ServiceServer {
	return &serviceRouter{
		Client:   client,
		services: services,
	}
}

func (r *serviceRouter) NotifyServer(ctx context.Context, res *service.NotifyServerRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		ser := s.GetService()
		var User auth.LinkUser
		err := ser.Client.DB.Joins("User", ser.Client.DB.Where("id = ?", res.GetUserId())).First(&User).Error
		if err != nil {
			return nil, UserNotFound
		}
		userIdInt, _ := strconv.ParseInt(res.GetUserId(), 10, 64)
		s.SendMessage("Вы подключились к серверу "+res.GetServer(), userIdInt)
	}
	return &emptypb.Empty{}, nil
}

func (r *serviceRouter) Join(ctx context.Context, res *service.JoinRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		ser := s.GetService()
		var User auth.LinkUser
		err := ser.Client.DB.Joins("User", ser.Client.DB.Where("id = ?", res.GetUserId())).First(&User).Error
		if err != nil {
			return nil, UserNotFound
		}
		userIdInt, _ := strconv.ParseInt(res.GetUserId(), 10, 64)
		s.SendMessage("Вы подключились к серверу с "+res.GetIp(), userIdInt)
	}
	return &emptypb.Empty{}, nil
}

func (r *serviceRouter) CheckCode(ctx context.Context, res *service.CheckCodeRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		ser := s.GetService()
		user, err := s.GetMCUser(res.GetUsername())
		if err != nil {
			return nil, UserNotFound
		}
		VkID := ser.Code.GetCode(res.GetUsername())
		if VkID == nil {
			return nil, InvalidCode
		}
		if VkID.Service != ser.ServiceID {
			return &emptypb.Empty{}, nil
		}
		if !ser.Code.CompareCode(res.GetUsername(), res.GetCode()) {
			return nil, InvalidCode
		}
		vkUser := &auth.LinkUser{
			ServiceUserID: VkID.UserID,
			User:          *user,
			ServiceId:     ser.ServiceID,
		}
		err = ser.Client.DB.Save(vkUser).Error
		if err != nil {
			return nil, err
		}
		ser.Code.RemoveCode(res.GetUsername())
		s.SendKeyboard("Вы успешно привязали аккаунт "+user.Username, VkID.UserID)
	}
	return &emptypb.Empty{}, nil
}
