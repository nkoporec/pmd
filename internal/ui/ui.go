package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/nkoporec/dump/config"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	timeFormat = "2006-01-02 20:00:00"
)

var cfg *config.Config

type DisplayData struct {
		Data []struct {
			Payload     string  `json:"payload"`
			File  string  `json:"file"`
			Type string  `json:"type"`
			Timestamp  string `json:"timestamp"`
		} `json:"data"`
}

type Term struct {
	Width int
	Height  int
}

func Display() {
	// Init config.
	cfg = config.InitConfig()

	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termtermui: %v", err)
	}
	defer termui.Close()

	termWidth, termHeight, err := terminal.GetSize(0)
	if err != nil {
		log.Fatalf("failed to get terminal size: %v", err)
	}

	term := &Term{
		Width: termWidth,
		Height: termHeight,
	}

	l, p := elements(term.Width, term.Height)
	termui.Render(l,p)

	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "j", "<Down>":
				l.ScrollDown()
			case "k", "<Up>":
				l.ScrollUp()
			case "<C-d>":
				l.ScrollHalfPageDown()
			case "<C-u>":
				l.ScrollHalfPageUp()
			case "<C-f>":
				l.ScrollPageDown()
			case "<C-b>":
				l.ScrollPageUp()
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				term.Width = payload.Width
				term.Height = payload.Height
				l,p = getUpdates(l,p)
				termui.Render(l,p)
			}
		case <-ticker:
			l,p = getUpdates(l,p)
			termui.Render(l,p)
		}
	}
}

func getUpdates(list *widgets.List,  paragraph *widgets.Paragraph) (l *widgets.List,  p *widgets.Paragraph) {
	var displayData *DisplayData

	request, err := http.Get("http://" + cfg.Server.Host + ":" + cfg.Server.Port + "/api/get")
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()

 	err = json.NewDecoder(request.Body).Decode(&displayData)
	if err != nil {
		panic(err)
	}

	// Clear list and paragraph
	list.Rows = []string{}
	paragraph.Text = ""

	payload := "";
	for _, elem := range displayData.Data {
    	i, err := strconv.ParseInt(elem.Timestamp, 10, 64)
		if err != nil {
			panic(err)
		}

		row := fmt.Sprintf("[%s] [%s](fg:white,bg:red)", time.Unix(i, 0).Format(timeFormat), elem.File);

		// Add to list.
		list.Rows = append(list.Rows, row)

		// Set payload.
		payload = elem.Payload
	}

	paragraph.Text = payload
	
	return list, paragraph
}

func elements(width int, height int) (*widgets.List,  *widgets.Paragraph) {
	l := widgets.NewList()
	l.Title = "Breakpoints"
	l.Rows = []string{}

	l.TextStyle = termui.NewStyle(termui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, width, (height/4))

	// Payload
	paragraph := widgets.NewParagraph()
	paragraph.Title = "Payload"
	paragraph.Text = ""
	paragraph.SetRect(0, (height/4), width, height)

	return l, paragraph
}
