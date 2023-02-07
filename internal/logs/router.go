package logs

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"hyneo/protos/logs"
)

type logsRouter struct {
	service Service
	logs.UnimplementedLogsServer
}

func NewLogsRouter(service Service) logs.LogsServer {
	return &logsRouter{
		service: service,
	}
}

func (l logsRouter) Join(ctx context.Context, request *logs.Request) (*emptypb.Empty, error) {
	err := l.service.Join(request.GetPlayerId(), request.GetServerName(), request.GetMessage())
	return &emptypb.Empty{}, err
}

func (l logsRouter) Quit(ctx context.Context, request *logs.Request) (*emptypb.Empty, error) {
	err := l.service.Quit(request.GetPlayerId(), request.GetServerName(), request.GetMessage())
	return &emptypb.Empty{}, err
}

func (l logsRouter) Message(ctx context.Context, request *logs.Request) (*emptypb.Empty, error) {
	err := l.service.Message(request.GetPlayerId(), request.GetServerName(), request.GetMessage())
	return &emptypb.Empty{}, err
}

func (l logsRouter) Command(ctx context.Context, request *logs.Request) (*emptypb.Empty, error) {
	err := l.service.Command(request.GetPlayerId(), request.GetServerName(), request.GetMessage())
	return &emptypb.Empty{}, err
}
