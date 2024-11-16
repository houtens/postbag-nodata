package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Env    string `default:"dev"`
	Host   string `default:"localhost"`
	Port   string `default:"3000"`
	Secret string `default:"somethingsupersecret"`
}

type DBConfig struct {
	DBDriver string `default:"postgres"`
	DBName   string `default:"postbag"`
	DBUser   string `default:"postgres"`
	DBPass   string `default:""`
	DBHost   string `default:"localhost"`
	DBPort   string `default:"5432"`
}

type Config struct {
	AppConfig
	DBConfig
}

func Load() Config {
	var c Config
	err := envconfig.Process("postbag", &c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
