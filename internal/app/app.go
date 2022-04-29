package app

import (
	"github.com/dgraph-io/ristretto"
	"github.com/nkoporec/pmd/config"
	"github.com/nkoporec/pmd/internal/app/cmd"
)

func Run(cfg *config.Config, cache *ristretto.Cache) {
	cmd.Execute(cfg, cache);
}
