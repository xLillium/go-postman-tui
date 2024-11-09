package app

import (
	"github.com/rivo/tview"
	"github.com/xlillium/go-postman-tui/internal/ui"
)

func Run() error {
	app := tview.NewApplication()
	ui := ui.Initialize(app)

	// Start the application
	return app.SetRoot(ui.Pages, true).SetFocus(ui.InitialFocus).EnableMouse(true).Run()
}
