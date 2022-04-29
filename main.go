package main

import (
	"errors"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/nkoporec/pmd/config"
	"github.com/nkoporec/pmd/internal/app"
	"github.com/dgraph-io/ristretto"
)

var cfg config.Config

func main() {
	if _, err := os.Stat(cfg.ConfigPath()); errors.Is(err, os.ErrNotExist) {
		_, err := cfg.CreateConfig()
		if err != nil {
			panic(err)
		}
	}

	err := cleanenv.ReadConfig(cfg.ConfigPath(), &cfg.Yaml)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// @TODO: Logger

	// Init cache.
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	// Run
	app.Run(&cfg, cache)
}
