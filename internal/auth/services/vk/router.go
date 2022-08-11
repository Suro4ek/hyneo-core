package vk

import (
	"context"
	"hyneo/eu.suro/hyneo/protos/vk"
	"hyneo/pkg/mysql"
)

type vkRouter struct {
	Client  mysql.Client
	service VKService
	vk.UnimplementedVKServer
}

func NewVKRouter() vk.VKServer {
	return &vkRouter{}
}

func (r *vkRouter) NotifyServer(ctx context.Context, res *vk.NotifyServerRequest) (*vk.Empty, error) {
	err := r.service.NotifyServer(res.UserId, res.Server)
	if err != nil {
		return nil, err
	}
	return &vk.Empty{}, nil
}

func (r *vkRouter) Join(ctx context.Context, res *vk.JoinRequest) (*vk.Empty, error) {
	err := r.service.Join(res.UserId, res.Ip)
	if err != nil {
		return nil, err
	}
	return &vk.Empty{}, nil
}

func (r *vkRouter) CheckCode(ctx context.Context, res *vk.CheckCodeRequest) (*vk.Empty, error) {
	err := r.service.CheckCode(res.Username, res.Code)
	if err != nil {
		return nil, err
	}
	return &vk.Empty{}, nil
}
