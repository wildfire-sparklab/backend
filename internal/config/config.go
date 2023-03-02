package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MapKey string `yaml:"mapKey" env:"MAP_KEY"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
		}
	})
	return instance
}
