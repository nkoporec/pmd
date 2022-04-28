package main

import (
	"errors"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/nkoporec/pmd/config"
	"github.com/nkoporec/pmd/internal/app"
)

func main() {
	// Configuration
	var cfg config.Config
	var cfgYaml config.ConfigYaml

	if _, err := os.Stat(cfg.ConfigPath()); errors.Is(err, os.ErrNotExist) {
		_, err := cfg.CreateConfig()
		if err != nil {
			panic(err)
		}
	}

	err := cleanenv.ReadConfig(cfg.ConfigPath(), &cfgYaml)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// @TODO: Logger

	// Run
	app.Run(&cfg)
}
