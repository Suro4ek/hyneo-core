package main

import (
	"context"
	"google.golang.org/grpc"
	auth2 "hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/mc"
	service2 "hyneo/internal/auth/mc/service"
	"hyneo/internal/auth/password"
	"hyneo/internal/auth/services"
	"hyneo/internal/auth/services/command"
	"hyneo/internal/config"
	"hyneo/pkg/logging"
	"hyneo/pkg/mysql"
	"hyneo/pkg/redis"
	"hyneo/protos/auth"
	serviceRouter "hyneo/protos/service"
	"log"
	"net"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	client := mysql.NewClient(context.Background(), 5, cfg.MySQL)
	redisClient, err := redis.NewClient(context.Background(), cfg.Redis)
	if err != nil {
		logger.Fatal(err)
	}
	migrate(client)
	codeService := &code.Service{
		Client: redisClient,
	}
	runServices := RunServices(cfg, codeService, client)
	command.RegisterCommands()
	runGRPCServer(runServices, *client, *cfg)
}

func migrate(client *mysql.Client) {
	err := client.DB.AutoMigrate(&auth2.LinkUser{})
	if err != nil {
		return
	}
	err = client.DB.AutoMigrate(&auth2.User{})
	if err != nil {
		return
	}
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
