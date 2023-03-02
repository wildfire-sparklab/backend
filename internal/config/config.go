package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MapKey   string   `yaml:"mapKey" env:"MAP_KEY"`
	Postgres Postgres `yaml:"postgres"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     string `yaml:"port" env:"POSTGRES_PORT"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	Password string `yaml:"pass" env:"POSTGRES_PASS"`
	DB       string `yaml:"db" env:"POSTGRES_DB"`
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
