package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MapKey  string `yaml:"mapKey" env:"MAP_KEY"`
	WindKey string `yaml:"windKey" env:"WIND_KEY"`
	MySQL   MySQL  `yaml:"mysql"`
	AMQP    string `yaml:"amqp" env:"AMQP"`
	S3      S3     `yaml:"s3"`
}

type MySQL struct {
	Host     string `yaml:"host" env:"MYSQL_HOST"`
	Port     string `yaml:"port" env:"MYSQL_PORT"`
	User     string `yaml:"user" env:"MYSQL_USER"`
	Password string `yaml:"pass" env:"MYSQL_PASS"`
	DB       string `yaml:"db" env:"MYSQL_DB"`
}

type S3 struct {
	Bucket    string `yaml:"bucket" env:"BUCKET"`
	AccessKey string `yaml:"access_key" env:"ACCESS_KEY"`
	SecretKey string `yaml:"secret_key" env:"SECRET_KEY"`
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
