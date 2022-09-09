package config

import (
	"hyneo/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MySQL    MySQL  `yaml:"mysql"`
	GRPCPort string `yaml:"grpc_port" env:"GRPC_PORT"`
	Redis    Redis  `yaml:"redis"`
	Social   Social `yaml:"social"`
}

type Social struct {
	VK       VKConfig         `yaml:"vk"`
	Telegram TelegramConfig   `yaml:"telegram"`
	Keyboard []KeyboardConfig `yaml:"keyboard"`
}

type MySQL struct {
	Host     string `yaml:"host" env:"MYSQL_HOST"`
	Port     string `yaml:"port" env:"MYSQL_PORT"`
	User     string `yaml:"user" env:"MYSQL_USER"`
	Password string `yaml:"pass" env:"MYSQL_PASSWORD"`
	DB       string `yaml:"db" env:"MYSQL_DB"`
}

type Redis struct {
	Host string `yaml:"host" env:"REDIS_HOST"`
	Port string `yaml:"port" env:"REDIS_PORT"`
	Pass string `yaml:"pass" env:"REDIS_PASS"`
}

type VKConfig struct {
	GroupID int64  `yaml:"group_id" env:"VK_GROUP_ID"`
	Token   string `yaml:"token" env:"VK_TOKEN"`
}

type TelegramConfig struct {
	Token string `yaml:"token" env:"TELEGRAM_TOKEN"`
}

type (
	KeyboardConfig struct {
		KeyboardButtons []KeyboardButtons `yaml:"buttons"`
	}
	KeyboardButtons struct {
		Name    string `yaml:"name"`
		Payload string `yaml:"payload"`
		Color   string `yaml:"color"`
	}
)

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
