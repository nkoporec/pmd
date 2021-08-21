package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/nkoporec/dump/internal/http"
	"github.com/spf13/cobra"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Starts a debugging server.",
		Run: func(cmd *cobra.Command, args []string) {
			go listen()
			gui()
		},
}

func init() {
  RootCmd.AddCommand(listenCmd)
}

func listen() {
	// Start http.
	handler := gin.New()
	http.NewRouter(handler)
	handler.Run()
}


func gui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
