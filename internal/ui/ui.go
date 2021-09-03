package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type DisplayData struct {
		Data []struct {
			Payload     string  `json:"payload"`
			File  string  `json:"file"`
			Type string  `json:"type"`
			Timestamp  int `json:"timestamp"`
		} `json:"data"`
}

const (
	timeFormat = "2006-01-02 20:00:00"
)

func Display() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termtermui: %v", err)
	}
	defer termui.Close()

	termui.Render()

	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			}
		// use Go's built-in tickers for updating and drawing data
		case <-ticker:
			getUpdates()
		}
	}
}

func getUpdates() {
	request, err := http.Get("http://127.0.0.1:8080/api/get")
	if err != nil {
		panic(err)
	}

	var displayData *DisplayData
 	err = json.NewDecoder(request.Body).Decode(&displayData)
	if err != nil {
		panic(err)
	}

	//  Table.
	table := widgets.NewTable()
	table.Rows = [][]string{
		{"Type", "File", "Timestamp", "Payload"},
	}

	lastPayload := "";
	for _, elem := range displayData.Data {
		timestamp := time.Unix(int64(elem.Timestamp), 0)

		row := []string{
			elem.Type,
			elem.File,
			fmt.Sprint(timestamp),
			elem.Payload,
		}
		table.Rows = append(table.Rows, row)

		// Set payload.
		lastPayload = elem.Payload
	}

	table.TextStyle = termui.NewStyle(termui.ColorWhite)
	table.SetRect(0, 0, 239, 20)
	termui.Render(table)

	// Payload
	paragraph := widgets.NewParagraph()
	paragraph.Title = "Last payload"
	paragraph.Text = lastPayload
	paragraph.SetRect(0, 50, 239, 20)
	termui.Render(paragraph)
}
