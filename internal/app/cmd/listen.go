package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/nkoporec/pmd/internal/http"
	"github.com/nkoporec/pmd/internal/ui"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Starts a debugging server.",
		Run: func(cmd *cobra.Command, args []string) {
			messages := make(chan interface{})

			go startServer(messages)
			displayUi(messages)
		},
}

func init() {
  RootCmd.AddCommand(listenCmd)
}

func startServer(messages chan interface{}) {
	gin.SetMode(gin.ReleaseMode)

	// Start http server.
	handler := gin.New()
	http.NewRouter(handler, messages, cch)
	handler.Run(cfg.Yaml.Host + ":" + cfg.Yaml.Port)
}

func displayUi(messages chan interface{}) {
	ui.Display(messages, cch)
}
