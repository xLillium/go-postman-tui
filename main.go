package main

import (
	"fmt"
	"log"

	"github.com/rivo/tview"
	"github.com/xlillium/go-postman-tui/ui"
)

func main() {
	fmt.Println("Starting application...")

	app := tview.NewApplication()
	ui := ui.NewUI(app)

	// Start the application
	if err := app.SetRoot(ui.Flex, true).SetFocus(ui.UrlInput).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
