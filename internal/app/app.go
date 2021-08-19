package app

import (
	"github.com/nkoporec/dump/config"
	"github.com/nkoporec/dump/internal/app/cmd"
)

func Run(cfg *config.Config) {
	cmd.Execute();
}
