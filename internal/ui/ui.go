package ui

import (
	"github.com/rivo/tview"
	"github.com/xlillium/go-postman-tui/internal/storage"
)

type UI struct {
	App          *tview.Application
	Pages        *tview.Pages
	RootLayout   *tview.Flex
	InitialFocus tview.Primitive
	Storage      *storage.Storage

	// Saved Requests Panel
	RequestList *tview.List
	// Dialog Components
	requestNameInputfield   *tview.InputField
	saveRequestButton       *tview.Button
	cancelSaveRequestButton *tview.Button
	// Root Components
	MethodDropdown   *tview.DropDown
	URLInputField    *tview.InputField
	BodyInputField   *tview.InputField
	ResponseTextView *tview.TextView
	ConsoleTextView  *tview.TextView
	RequestFlex      *tview.Flex
}

func Initialize(app *tview.Application) *UI {
	ui := &UI{
		App:     app,
		Storage: storage.NewStorage("requests.json"),
	}

	ui.setupComponents()
	ui.setupLayout()
	ui.setupEventHandlers()
	ui.setupGlobalKeybindings()
	requests := ui.Storage.GetRequests()
	ui.RequestList.Clear()
	for _, req := range requests {
		ui.RequestList.AddItem(req.Name, "", 0, nil)
	}

	return ui
}
