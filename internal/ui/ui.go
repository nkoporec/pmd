package ui

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func Display() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termtermui: %v", err)
	}
	defer termui.Close()

	p := widgets.NewParagraph()
	p.Text = "Hello World!"
	p.SetRect(0, 0, 25, 5)

	termui.Render(p)

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
			drawFunction()
		}
	}
}

func drawFunction() {
	resp, err := http.Get("http://127.0.0.1:8080/api/get")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	p := widgets.NewParagraph()
	p.Text = string(body)
	p.SetRect(0, 0, 100, 5)
	termui.Render(p)
}
