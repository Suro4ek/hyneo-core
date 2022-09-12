package services

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"hyneo/internal/auth"
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
		user, err := s.GetUserID(int64(res.UserId))
		if err != nil {
			return nil, UserNotFound
		}
		if user.ServiceId != s.GetService().ServiceID {
			continue
		}
		s.SendMessage("Вы подключились к серверу "+res.GetServer(), user.ServiceUserID)
	}
	return &emptypb.Empty{}, nil
}

func (r *serviceRouter) Join(_ context.Context, res *service.JoinRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		user, err := s.GetUserID(int64(res.UserId))
		if err != nil {
			return nil, UserNotFound
		}
		if user.ServiceId != s.GetService().ServiceID {
			continue
		}
		s.SendMessage("Вы подключились к серверу с ip: "+res.GetIp(), user.ServiceUserID)
	}
	return &emptypb.Empty{}, nil
}

func (r *serviceRouter) CheckCode(_ context.Context, res *service.CheckCodeRequest) (*emptypb.Empty, error) {
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
