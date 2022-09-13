package main

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/go-redis/redis/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/password"
	"hyneo/internal/config"
	"hyneo/internal/social/services"
	telegram2 "hyneo/internal/social/services/telegram"
	vk3 "hyneo/internal/social/services/vk"
	"hyneo/pkg/logging"
	"hyneo/pkg/mysql"
)

func RunServices(cfg *config.Config, service *code.Service, client *mysql.Client, redis *redis.Client, log *logging.Logger, passwordService password.Service) []services.Service {
	servicess := make([]services.Service, 0)
	servicess = append(servicess, runVKLongServer(client, cfg, service, redis, log, passwordService))
	servicess = append(servicess, runTGServer(client, cfg, service, redis, log, passwordService))
	return servicess
}

func runVKLongServer(Client *mysql.Client,
	cfg *config.Config,
	code *code.Service,
	redis *redis.Client,
	log *logging.Logger,
	passwordService password.Service) services.Service {
	token := cfg.Social.VK.Token // use os.Getenv("TOKEN")
	vk := api.NewVK(token)

	service := vk3.NewVkService(Client, vk, code, redis, 0, log, passwordService)
	// get information about the group
	group, err := vk.GroupsGetByID(api.Params{
		"group_id": cfg.Social.VK.GroupID,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Run LongPoll VK")
	// Initializing Long Poll
	go func() {
		lp, err := longpoll.NewLongPoll(vk, group[0].ID)
		if err != nil {
			log.Fatal(err)
		}
		handler := vk3.NewVKHandler(lp, &service)
		handler.Message()
	}()
	return service
}

func runTGServer(Client *mysql.Client,
	cfg *config.Config,
	code *code.Service,
	redis *redis.Client,
	log *logging.Logger,
	passwordService password.Service) services.Service {
	bot, err := tgbotapi.NewBotAPI(cfg.Social.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	service := telegram2.NewTelegramService(Client, bot, code, redis, 1, log, passwordService)
	log.Info("Run Listen message Telegram")
	go func() {
		handler := telegram2.NewTelegramHandler(bot, &service)
		handler.Message()
	}()
	return service
}
