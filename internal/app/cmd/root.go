package cmd

import (
	"os"

	"github.com/nkoporec/pmd/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

var RootCmd = &cobra.Command{
		Use:   "pmd",
		Short: "Poor's man debugger.",
		// @todo: Add description.
		// Long: `Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
}

func Execute(c *config.Config) {
	cfg = c
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
