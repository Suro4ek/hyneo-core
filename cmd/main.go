package main

import (
	"context"
	"google.golang.org/grpc"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/mc"
	service2 "hyneo/internal/auth/mc/service"
	"hyneo/internal/auth/password"
	"hyneo/internal/config"
	"hyneo/internal/logs"
	storage2 "hyneo/internal/logs/storage"
	services2 "hyneo/internal/social/services"
	"hyneo/internal/social/services/command"
	"hyneo/internal/user"
	"hyneo/internal/user/storage"
	"hyneo/pkg/logging"
	"hyneo/pkg/mysql"
	"hyneo/pkg/redis"
	"hyneo/protos/auth"
	logs2 "hyneo/protos/logs"
	serviceRouter "hyneo/protos/service"
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
	logger.Info("Register services")
	passwordService := password.NewPasswordService()

	storageUser := storage.CreateStorageUser(client, redisClient)
	storageLogs := storage2.NewLogsStorage(client)

	userService := user.CreateUserService(storageUser, &logger)
	logsService := logs.NewLogsService(storageLogs)

	runServices := RunServices(cfg, codeService, redisClient, &logger, passwordService, storageUser, client)
	logger.Info("Register commands")
	command.RegisterCommands()
	logger.Info("Run GRPC Server to " + cfg.GRPCPort)
	runGRPCServer(runServices, client, *cfg, &logger, passwordService, storageUser, *userService, *logsService)
}

func migrate(client *mysql.Client) {
	err := client.DB.AutoMigrate(&user.User{})
	if err != nil {
		return
	}
	err = client.DB.AutoMigrate(&user.LinkUser{})
	if err != nil {
		return
	}
	err = client.DB.AutoMigrate(&user.IgnoreUser{})
	if err != nil {
		return
	}
	err = client.DB.AutoMigrate(&logs.Logs{})
	if err != nil {
		return
	}
}

func runGRPCServer(
	servicess []services2.Service,
	client *mysql.Client,
	cfg config.Config,
	log *logging.Logger,
	passwordService password.Service,
	storageUser user.Service,
	userService user.UserService,
	logsService logs.Service,
) {
	addr := "0.0.0.0:" + cfg.GRPCPort
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	serviceAuth := service2.NewMCService(passwordService, log, storageUser)
	authService := mc.NewAuthRouter(serviceAuth)
	auth.RegisterAuthServer(s, authService)
	userRouter := user.NewUserRouter(userService)
	auth.RegisterUserServiceServer(s, userRouter)

	logsRouter := logs.NewLogsRouter(logsService)
	logs2.RegisterLogsServer(s, logsRouter)

	service := services2.NewServiceRouter(client, servicess)

	serviceRouter.RegisterServiceServer(s, service)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
