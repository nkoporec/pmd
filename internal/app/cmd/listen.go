package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Starts a debugging server.",
		Run: func(cmd *cobra.Command, args []string) {
			listen()
		},
}

func init() {
  RootCmd.AddCommand(listenCmd)
}

func listen() {
	log.Fatal("listen cmd")
}
