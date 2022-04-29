package cmd

import (
	"os"

	"github.com/dgraph-io/ristretto"
	"github.com/nkoporec/pmd/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config
var cch *ristretto.Cache

var RootCmd = &cobra.Command{
		Use:   "pmd",
		Short: "Poor's man debugger.",
}

func Execute(c *config.Config, cache *ristretto.Cache) {
	// For later use.
	cfg = c
	cch = cache

	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
