package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/nkoporec/dump/internal/http"
	"github.com/nkoporec/dump/internal/ui"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Starts a debugging server.",
		Run: func(cmd *cobra.Command, args []string) {
			go startServer()
			displayUi()
		},
}

func init() {
  RootCmd.AddCommand(listenCmd)
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)

	// Start http server.
	handler := gin.New()
	http.NewRouter(handler)
	handler.Run()
}

func displayUi() {
	ui.Display()
}
