package main

import (
	"log"

	"github.com/rivo/tview"
	"github.com/xlillium/go-postman-tui/ui"
)

func main() {
	app := tview.NewApplication()
	ui := ui.InitializeUI(app)

	// Start the application
	if err := app.SetRoot(ui.Layout, true).SetFocus(ui.URLInputField).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
