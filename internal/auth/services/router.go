package services

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (r *serviceRouter) NotifyServer(ctx context.Context, res *service.NotifyServerRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		_ = s.NotifyServer(res.UserId, res.Server)
	}
	return &emptypb.Empty{}, nil
}

func (r *serviceRouter) Join(ctx context.Context, res *service.JoinRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		_ = s.Join(res.UserId, res.Ip)
	}
	return &emptypb.Empty{}, nil
}

func (r *serviceRouter) CheckCode(ctx context.Context, res *service.CheckCodeRequest) (*emptypb.Empty, error) {
	for _, s := range r.services {
		err := s.CheckCode(res.Username, res.Code)
		if err != nil {
			return &emptypb.Empty{}, err
		}
	}
	return &emptypb.Empty{}, nil
}
