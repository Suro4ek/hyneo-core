package services

import (
	"context"
	"hyneo/pkg/mysql"
	"hyneo/protos/service"
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

func (r *serviceRouter) NotifyServer(ctx context.Context, res *service.NotifyServerRequest) (*service.Empty, error) {
	for _, s := range r.services {
		_ = s.NotifyServer(res.UserId, res.Server)
	}
	return &service.Empty{}, nil
}

func (r *serviceRouter) Join(ctx context.Context, res *service.JoinRequest) (*service.Empty, error) {
	for _, s := range r.services {
		_ = s.Join(res.UserId, res.Ip)
	}
	return &service.Empty{}, nil
}

func (r *serviceRouter) CheckCode(ctx context.Context, res *service.CheckCodeRequest) (*service.Empty, error) {
	for _, s := range r.services {
		err := s.CheckCode(res.Username, res.Code)
		if err != nil {
			return &service.Empty{}, err
		}
	}
	return &service.Empty{}, nil
}
