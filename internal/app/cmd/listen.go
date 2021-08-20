package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/nkoporec/dump/internal/http"
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
	// Start http.
	handler := gin.Default()
	http.NewRouter(handler)
	handler.Run()
}
