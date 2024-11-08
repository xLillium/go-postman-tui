package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xlillium/go-postman-tui/handlers"
	"github.com/xlillium/go-postman-tui/ui"
)

func main() {
	app := tview.NewApplication()
	ui := ui.NewUI(app)

	// Event handler for the input field
	ui.TopInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			url := ui.TopInput.GetText()
			go func() {
				handlers.PerformGetRequest(url, ui)
			}()
		}
	})

	// Start the application
	if err := app.SetRoot(ui.Flex, true).SetFocus(ui.TopInput).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
