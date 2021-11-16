package ui

import (
	"encoding/json"
	"log"
	"net/http"
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

	termui.Render()

	termWidth, termHeight, err := terminal.GetSize(0)
	if err != nil {
		log.Fatalf("failed to get terminal size: %v", err)
	}

	term := &Term{
		Width: termWidth,
		Height: termHeight,
	}

	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				term.Width = payload.Width
				term.Height = payload.Height
				getUpdates(term)
			}
		case <-ticker:
			getUpdates(term)
		}
	}
}

func getUpdates(term *Term) {
	var displayData *DisplayData

	request, err := http.Get("http://" + cfg.Server.Host + ":" + cfg.Server.Port + "/api/get")
	if err != nil {
		panic(err)
	}

 	err = json.NewDecoder(request.Body).Decode(&displayData)
	if err != nil {
		panic(err)
	}

	draw(displayData, term.Width, term.Height)
}


func draw(data *DisplayData, width int, height int) {
	//  Table.
	table := widgets.NewTable()
	table.Rows = [][]string{
		{"Type", "File", "Timestamp", "Payload"},
	}

	lastPayload := "";
	for _, elem := range data.Data {
		row := []string{
			elem.Type,
			elem.File,
			elem.Timestamp,
			elem.Payload,
		}
		table.Rows = append(table.Rows, row)

		// Set payload.
		lastPayload = elem.Payload
	}

	table.TextStyle = termui.NewStyle(termui.ColorWhite)
	table.SetRect(0, 0, width, (height/4))
	termui.Render(table)

	// Payload
	paragraph := widgets.NewParagraph()
	paragraph.Title = "Last payload"
	paragraph.Text = lastPayload
	paragraph.SetRect(0, (height/4), width, height)
	termui.Render(paragraph)
}
