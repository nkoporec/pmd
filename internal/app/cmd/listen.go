package cmd

import (
	"fmt"
	nhttp "net/http"
	"os"

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

	// Check if port is free.
	// For some reason gin is allowing
	// multiple connections to the same port
	// so we use a simple GET request to check if the
	// host is free.
	resp, err := nhttp.Get("http://" + cfg.Yaml.Host + ":" + cfg.Yaml.Port)
	if err != nil {
		// Start http server.
		handler := gin.New()
		http.NewRouter(handler, messages, cch)
		handler.Run(cfg.Yaml.Host + ":" + cfg.Yaml.Port)
	}
	defer resp.Body.Close()

	fmt.Println("Host is already in use.")
	os.Exit(1)
}

func displayUi(messages chan interface{}) {
	ui.Display(messages, cch, cfg)
}
