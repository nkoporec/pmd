package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	App  `yaml:"app"`
	Server  `yaml:"server"`
}

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Server struct {
	Host    string `yaml:"host" env-default: "127.0.0.1"`
	Port string `yaml:"port" env-default: "8080"`
}

func InitConfig() *Config {
	config := Config{}
	err := cleanenv.ReadConfig("config/config.yml", &config)

	if err != nil {
		panic(err)
	}

	return &config
}
