package app

import (
	"github.com/nkoporec/pmd/config"
	"github.com/nkoporec/pmd/internal/app/cmd"
)

func Run(cfg *config.Config) {
	cmd.Execute(cfg);
}
