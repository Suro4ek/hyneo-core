package services

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"hyneo/internal/user"
	"hyneo/pkg/mysql"
	"hyneo/protos/service"
)

type serviceRouter struct {
	Client   *mysql.Client
	services []Service
	service.UnimplementedServiceServer
}

func NewServiceRouter(client *mysql.Client, services []Service) service.ServiceServer {
	return &serviceRouter{
		Client:   client,
		services: services,
	}
}

func (r *serviceRouter) NotifyServer(_ context.Context, res *service.NotifyServerRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		u, err := s.GetUserID(int64(res.UserId))
		if err != nil {
			return nil, UserNotFound
		}
		if !u.Notificated {
			return &emptypb.Empty{}, nil
		}
		if u.ServiceId != s.GetService().ServiceID {
			continue
		}
		s.SendMessage("Вы подключились к серверу "+res.GetServer(), u.ServiceUserID)
	}
	return &emptypb.Empty{}, nil
}

func (r *serviceRouter) Join(_ context.Context, res *service.JoinRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		u, err := s.GetUserID(res.UserId)
		if err != nil {
			return nil, UserNotFound
		}
		if !u.Notificated {
			return &emptypb.Empty{}, nil
		}
		if u.ServiceId != s.GetService().ServiceID {
			continue
		}
		s.SendMessage("Вы подключились к серверу с ip: "+res.GetIp(), u.ServiceUserID)
	}
	return &emptypb.Empty{}, nil
}

func (r *serviceRouter) CheckCode(_ context.Context, res *service.CheckCodeRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		ser := s.GetService()
		u, err := s.GetMCUser(res.GetUsername())
		if err != nil {
			return nil, UserNotFound
		}
		VkID := ser.Code.GetCode(res.GetUsername())
		if VkID == nil {
			return nil, InvalidCode
		}
		if VkID.Service != ser.ServiceID {
			continue
		}
		if !ser.Code.CompareCode(res.GetUsername(), res.GetCode()) {
			return nil, InvalidCode
		}
		vkUser := &user.LinkUser{
			ServiceUserID: VkID.UserID,
			User:          *u,
			ServiceId:     ser.ServiceID,
			Notificated:   true,
			Banned:        false,
			DoubleAuth:    false,
		}
		err = ser.Client.DB.Save(vkUser).Error
		if err != nil {
			return nil, err
		}
		ser.Code.RemoveCode(res.GetUsername())
		s.SendKeyboard("Вы успешно привязали аккаунт "+u.Username, VkID.UserID)
	}
	return &emptypb.Empty{}, nil
}
