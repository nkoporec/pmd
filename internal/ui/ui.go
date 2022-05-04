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
	"github.com/tidwall/pretty"
	"golang.org/x/crypto/ssh/terminal"
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

	breakpointsWidget, callstackWidget, payloadWidget := elements(term.Width, term.Height)
	termui.Render(breakpointsWidget, callstackWidget, payloadWidget)

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
				callstackWidget.Rows = []string{}
				payloadWidget.Text = ""
				selectedLine = 0

				cch.Set("breakpoints", []*http.RequestData{}, 1)
				termui.Render(breakpointsWidget, callstackWidget, payloadWidget)
			case "j", "<Down>":
				if selectedLine < len(displayedData)-1 {
					breakpointsWidget.ScrollDown()
					selectedLine++

					// Show the selected breakpoint callstack.
					callstackWidget.Rows = []string{}
					for _, call := range displayedData[selectedLine].Callstack {
						callstackRow := fmt.Sprintf(
							"%s:%s:%s",
							call.Line,
							call.File,
							call.Function,
						)
						callstackWidget.Rows = append(callstackWidget.Rows, callstackRow)
					}
					payloadWidget.Text = formatPayload(displayedData[selectedLine].Payload)
					termui.Render(breakpointsWidget, callstackWidget, payloadWidget)
				}
			case "k", "<Up>":
				if selectedLine <= 0 {
					breakpointsWidget.ScrollUp()
				} else {
					breakpointsWidget.ScrollUp()
					selectedLine--

					// Show the selected breakpoint callstack.
					callstackWidget.Rows = []string{}
					for _, call := range displayedData[selectedLine].Callstack {
						callstackRow := fmt.Sprintf(
							"%s:%s:%s",
							call.Line,
							call.File,
							call.Function,
						)
						callstackWidget.Rows = append(callstackWidget.Rows, callstackRow)
					}

					payloadWidget.Text = formatPayload(displayedData[selectedLine].Payload)
				}
				termui.Render(breakpointsWidget, callstackWidget, payloadWidget)
			case "J":
				callstackWidget.ScrollDown()
				termui.Render(breakpointsWidget, callstackWidget, payloadWidget)
			case "K":
				callstackWidget.ScrollUp()
				termui.Render(breakpointsWidget, callstackWidget, payloadWidget)
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				term.Width = payload.Width
				term.Height = payload.Height
				termui.Render(breakpointsWidget, callstackWidget, payloadWidget)
			}
		case <-ticker:
			select {
			case msg := <-messages:
				breakpointsWidget.Rows = []string{}
				callstackWidget.Rows = []string{}
				payloadWidget.Text = ""

				data := msg.([]*http.RequestData)
				// so we can use it later.
				displayedData = data

				for _, elem := range data {
					i, err := strconv.ParseInt(elem.Timestamp, 10, 64)
					if err != nil {
						panic(err)
					}

					breakpoint := fmt.Sprintf(
						"[%s] [%s] %s:[%s](fg:white,bg:red)",
						elem.Type,
						time.Unix(i, 0),
						elem.Line,
						elem.File,
					)
					breakpointsWidget.Rows = append(breakpointsWidget.Rows, breakpoint)

					for _, call := range elem.Callstack {
						callstackRow := fmt.Sprintf(
							"%s:%s:%s",
							call.Line,
							call.File,
							call.Function,
						)
						callstackWidget.Rows = append(callstackWidget.Rows, callstackRow)
					}

				}

				payloadWidget.Text = formatPayload(displayedData[selectedLine].Payload)

				// Show the selected breakpoint callstack.
				callstackWidget.Rows = []string{}
				for _, call := range displayedData[selectedLine].Callstack {
					callstackRow := fmt.Sprintf(
						"%s:%s:%s",
						call.Line,
						call.File,
						call.Function,
					)
					callstackWidget.Rows = append(callstackWidget.Rows, callstackRow)
				}

				termui.Render(breakpointsWidget, callstackWidget, payloadWidget)
			default:
			}

		}
	}
}

func elements(width int, height int) (*widgets.List, *widgets.List, *widgets.Paragraph) {
	breakpointsWidget := widgets.NewList()
	breakpointsWidget.Title = "Breakpoints"
	breakpointsWidget.Rows = []string{}
	breakpointsWidget.TextStyle = termui.NewStyle(termui.ColorYellow)
	breakpointsWidget.WrapText = false
	breakpointsWidget.SetRect(0, 0, (width / 2), (height / 4))

	callstackWidget := widgets.NewList()
	callstackWidget.Title = "Call stack"
	callstackWidget.Rows = []string{}
	callstackWidget.TextStyle = termui.NewStyle(termui.ColorYellow)
	callstackWidget.WrapText = false
	callstackWidget.SetRect(width, 0, (width / 2), (height / 4))

	payloadWidget := widgets.NewParagraph()
	payloadWidget.Title = "Payload"
	payloadWidget.Text = ""
	payloadWidget.SetRect(0, (height / 4), width, height)

	return breakpointsWidget, callstackWidget, payloadWidget
}

func formatPayload(payload string) string {
	return string(pretty.Pretty([]byte(payload)))
}
