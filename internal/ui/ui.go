package ui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgraph-io/ristretto"
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/ilyakaznacheev/cleanenv"
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

func Display(messages chan interface{}, cch *ristretto.Cache) {
	// Init config.
	var cfg config.Config
	err := cleanenv.ReadConfig(cfg.ConfigPath(), &cfg.Yaml)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

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

	l, p := elements(term.Width, term.Height)
	termui.Render(l, p)

	breakpoint_pos := 0

	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<C-r>":
				l.Rows = []string{}
				p.Text = ""

				cch.Set("breakpoints", []*http.RequestData{}, 1)
				termui.Render(l, p)
			case "j", "<Down>":
				if breakpoint_pos < len(displayedData)-1 {
					l.ScrollDown()
					breakpoint_pos++

					p.Text = formatPayload(displayedData[breakpoint_pos].Payload)
					termui.Render(l, p)
				}
			case "k", "<Up>":
				if breakpoint_pos <= 0 {
					l.ScrollUp()
				} else {
					l.ScrollUp()
					breakpoint_pos--
					p.Text = formatPayload(displayedData[breakpoint_pos].Payload)
				}
				termui.Render(l, p)
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				term.Width = payload.Width
				term.Height = payload.Height
				termui.Render(l, p)
			}
		case <-ticker:
			select {
			case msg := <-messages:
				l.Rows = []string{}
				p.Text = ""
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

					// Add to list.
					l.Rows = append(l.Rows, row)
				}

				p.Text = formatPayload(displayedData[breakpoint_pos].Payload)

				termui.Render(l, p)
			default:
			}

		}
	}
}

func elements(width int, height int) (*widgets.List, *widgets.Paragraph) {
	l := widgets.NewList()
	l.Title = "Breakpoints"
	l.Rows = []string{}

	l.TextStyle = termui.NewStyle(termui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, width, (height / 4))

	// Payload
	paragraph := widgets.NewParagraph()
	paragraph.Title = "Payload"
	paragraph.Text = ""
	paragraph.SetRect(0, (height / 4), width, height)

	return l, paragraph
}

func formatPayload(payload string) string {
	return string(pretty.Pretty([]byte(payload)))
}
