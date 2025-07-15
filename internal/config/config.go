package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen struct {
		BindIp string `env:"BIND_IP"`
		Port   int    `env:"PORT"`
	}

	Postgres struct {
		Host     string `env:"DB_HOST"`
		Port     string `env:"DB_PORT"`
		HostName string `env:"DB_HOST_NAME"`
		Password string `env:"DB_PASSWORD"`
		DBName   string `env:"DB_NAME"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func () {
		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			print(help)
		}
	})

	return instance
}
