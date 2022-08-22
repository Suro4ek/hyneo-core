package main

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"hyneo/internal/auth/code"
	"hyneo/internal/auth/services"
	"hyneo/internal/auth/services/telegram"
	vk2 "hyneo/internal/auth/services/vk"
	"hyneo/internal/config"
	"hyneo/pkg/mysql"
	"log"
)

func RunServices(cfg *config.Config, service *code.Service, client *mysql.Client) []services.Service {
	servicess := make([]services.Service, 0)
	servicess = append(servicess, runVKLongServer(client, cfg, service))
	servicess = append(servicess, runTGServer(client, cfg, service))
	return servicess
}

func runVKLongServer(Client *mysql.Client, cfg *config.Config, code *code.Service) services.Service {
	token := cfg.VK.Token // use os.Getenv("TOKEN")
	vk := api.NewVK(token)

	service := vk2.NewVkService(Client, vk, code, 0)
	// get information about the group
	group, err := vk.GroupsGetByID(api.Params{
		"group_id": cfg.VK.GroupID,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Initializing Long Poll
	go func() {
		lp, err := longpoll.NewLongPoll(vk, group[0].ID)
		if err != nil {
			log.Fatal(err)
		}
		handler := vk2.NewVKHandler(lp, &service)
		handler.Message()
	}()
	return service
}

func runTGServer(Client *mysql.Client, cfg *config.Config, code *code.Service) services.Service {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	service := telegram.NewTelegramService(Client, bot, code, 1)
	go func() {
		handler := telegram.NewTelegramHandler(bot, &service)
		handler.Message()
	}()
	return service
}
