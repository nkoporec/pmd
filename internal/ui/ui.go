package ui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgraph-io/ristretto"
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/nkoporec/pmd/config"
	"github.com/nkoporec/pmd/internal/http"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/tidwall/pretty" 
)

const (
	timeFormat = "2006-01-02 20:00:00"
)

type Term struct {
	Width  int
	Height int
}

var displayedData []*http.RequestData

func Display(messages chan interface{}, cch *ristretto.Cache, cfg *config.Config) {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termtermui: %v", err)
	}
	defer termui.Close()

	termWidth, termHeight, err := terminal.GetSize(0)
	if err != nil {
		log.Fatalf("failed to get terminal size: %v", err)
	}

	term := &Term{
		Width:  termWidth,
		Height: termHeight,
	}

	breakpointsWidget, payloadWidget := elements(term.Width, term.Height)
	termui.Render(breakpointsWidget, payloadWidget)

	selectedLine := 0

	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<C-r>":
				breakpointsWidget.Rows = []string{}
				payloadWidget.Text = ""

				cch.Set("breakpoints", []*http.RequestData{}, 1)
				termui.Render(breakpointsWidget, payloadWidget)
			case "j", "<Down>":
				if selectedLine < len(displayedData)-1 {
					breakpointsWidget.ScrollDown()
					selectedLine++

					payloadWidget.Text = formatPayload(displayedData[selectedLine].Payload)
					termui.Render(breakpointsWidget, payloadWidget)
				}
			case "k", "<Up>":
				if selectedLine <= 0 {
					breakpointsWidget.ScrollUp()
				} else {
					breakpointsWidget.ScrollUp()
					selectedLine--
					payloadWidget.Text = formatPayload(displayedData[selectedLine].Payload)
				}
				termui.Render(breakpointsWidget, payloadWidget)
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				term.Width = payload.Width
				term.Height = payload.Height
				termui.Render(breakpointsWidget, payloadWidget)
			}
		case <-ticker:
			select {
			case msg := <-messages:
				breakpointsWidget.Rows = []string{}
				payloadWidget.Text = ""
				data := msg.([]*http.RequestData)
				displayedData = data

				for _, elem := range data {
					i, err := strconv.ParseInt(elem.Timestamp, 10, 64)
					if err != nil {
						panic(err)
					}

					row := fmt.Sprintf(
						"[%s] [%s] %s:[%s](fg:white,bg:red)",
						elem.Type,
						time.Unix(i, 0),
						elem.Line,
						elem.File,
					)

					breakpointsWidget.Rows = append(breakpointsWidget.Rows, row)
				}

				payloadWidget.Text = formatPayload(displayedData[selectedLine].Payload)

				termui.Render(breakpointsWidget, payloadWidget)
			default:
			}

		}
	}
}

func elements(width int, height int) (*widgets.List, *widgets.Paragraph) {
	breakpointsWidget := widgets.NewList()
	breakpointsWidget.Title = "Breakpoints"
	breakpointsWidget.Rows = []string{}

	breakpointsWidget.TextStyle = termui.NewStyle(termui.ColorYellow)
	breakpointsWidget.WrapText = false
	breakpointsWidget.SetRect(0, 0, width, (height / 4))

	// Payload
	payloadWidget := widgets.NewParagraph()
	payloadWidget.Title = "Payload"
	payloadWidget.Text = ""
	payloadWidget.SetRect(0, (height / 4), width, height)

	return breakpointsWidget, payloadWidget
}

func formatPayload(payload string) string {
	return string(pretty.Pretty([]byte(payload)))
}
