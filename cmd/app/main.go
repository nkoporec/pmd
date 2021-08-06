package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/nkoporec/dump/config"
	"github.com/nkoporec/dump/internal/app"
)

func main() {
	// Configuration
	var cfg config.Config

	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// @TODO: Logger

	// Run
	app.Run(&cfg)
}
