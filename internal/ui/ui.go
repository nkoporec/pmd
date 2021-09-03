package ui

import (
	"encoding/json"
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
			switch e.ID { // event string/identifier
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

	// Don't show anything, if we don't have any data.
	// if len(displayData.Data.Payload) <= 0 {
	// 	return
	// }

	//  Table.
	table := widgets.NewTable()
	table.Rows = [][]string{
		{"Type", "File", "Timestamp", "Payload"},
	}

	table.TextStyle = termui.NewStyle(termui.ColorWhite)
	table.SetRect(0, 0, 239, 20)
	termui.Render(table)

	// Payload
	paragraph := widgets.NewParagraph()
	paragraph.Title = "Last payload"
	paragraph.Text = "payload"
	paragraph.SetRect(0, 50, 239, 20)
	termui.Render(paragraph)
}
