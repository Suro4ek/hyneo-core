package config

import (
	"hyneo/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MySQL    MySQL    `yaml:"mysql"`
	GRPCPort string   `yaml:"grpc_port" env:"GRPC_PORT"`
	VK       VKConfig `yaml:"vk"`
}

type MySQL struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	DB       string `yaml:"db"`
}

type VKConfig struct {
	GroupID int64  `yaml:"group_id"`
	Token   string `yaml:"token"`
}

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
