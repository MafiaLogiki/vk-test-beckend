package config

import (
	"marketplace-service/internal/logger"
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
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		DBName   string `env:"DB_NAME"`
	}
}

var instance *Config
var once sync.Once

func GetConfig(l logger.Logger) *Config {
	once.Do(func () {
		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			l.Info(help)
			l.Fatal(err)
		}
	})

	return instance
}
