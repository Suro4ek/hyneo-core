package main

import (
	"context"
	"google.golang.org/grpc"
	auth2 "hyneo/internal/auth"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/mc"
	service2 "hyneo/internal/auth/mc/service"
	"hyneo/internal/auth/password"
	"hyneo/internal/config"
	"hyneo/internal/social/keyboard"
	services2 "hyneo/internal/social/services"
	"hyneo/internal/social/services/command"
	"hyneo/pkg/logging"
	"hyneo/pkg/mysql"
	"hyneo/pkg/redis"
	"hyneo/protos/auth"
	serviceRouter "hyneo/protos/service"
	"net"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	client := mysql.NewClient(context.Background(), 5, cfg.MySQL)
	keyboard.Init(cfg.Social.Keyboard)
	redisClient, err := redis.NewClient(context.Background(), cfg.Redis)
	if err != nil {
		logger.Fatal(err)
	}
	migrate(client)
	codeService := &code.Service{
		Client: redisClient,
	}
	logger.Info("Register services")
	passwordService := password.NewPasswordService()
	runServices := RunServices(cfg, codeService, client, redisClient, &logger, passwordService)
	logger.Info("Register commands")
	command.RegisterCommands()
	logger.Info("Run GRPC Server")
	runGRPCServer(runServices, client, *cfg, &logger, passwordService)
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

func runGRPCServer(servicess []services2.Service, client *mysql.Client, cfg config.Config, log *logging.Logger, passwordService password.Service) {
	addr := "0.0.0.0:" + cfg.GRPCPort
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	serviceAuth := service2.NewMCService(client, passwordService, log)
	authService := mc.NewAuthRouter(client, serviceAuth)
	auth.RegisterAuthServer(s, authService)

	service := services2.NewServiceRouter(client, servicess)
	serviceRouter.RegisterServiceServer(s, service)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
