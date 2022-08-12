package main

import (
	"context"
	"google.golang.org/grpc"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/mc"
	service2 "hyneo/internal/auth/mc/service"
	"hyneo/internal/auth/password"
	"hyneo/internal/auth/services"
	"hyneo/internal/auth/services/command"
	"hyneo/internal/config"
	"hyneo/pkg/logging"
	"hyneo/pkg/mysql"
	"hyneo/protos/auth"
	serviceRouter "hyneo/protos/service"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {
	logging.Init()
	_ = logging.GetLogger()
	cfg := config.GetConfig()
	client := mysql.NewClient(context.Background(), 5, cfg.MySQL)
	codeService := code.CodeService{}
	servicess := RunServices(cfg, codeService, client)
	command.RegisterCommands()
	runGRPCServer(servicess, *client, *cfg)
	rand.Seed(time.Now().UnixNano())

}

func runGRPCServer(servicess []services.Service, client mysql.Client, cfg config.Config) {
	addr := "0.0.0.0:" + cfg.GRPCPort
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	passwordService := password.NewPasswordService()
	serviceAuth := service2.NewMCService(&client, passwordService)
	authService := mc.NewAuthRouter(&client, serviceAuth)
	auth.RegisterAuthServer(s, authService)

	service := services.NewServiceRouter(client, servicess)
	serviceRouter.RegisterServiceServer(s, service)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
