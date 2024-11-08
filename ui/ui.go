package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xlillium/go-postman-tui/handlers"
)

// UI holds all the UI components
type UI struct {
	App              *tview.Application
	Layout           *tview.Flex
	MethodDropdown   *tview.DropDown
	URLInputField    *tview.InputField
	RequestFlex      *tview.Flex
	BodyInputField   *tview.InputField
	ResponseTextView *tview.TextView
	ConsoleTextView  *tview.TextView
}

// InitializeUI creates and configures the UI components
func InitializeUI(app *tview.Application) *UI {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	// HTTP Method DropDown
	methodDropdown := tview.NewDropDown().
		SetOptions([]string{"GET", "POST"}, nil).
		SetCurrentOption(0).
		SetFieldWidth(6).
		SetFieldBackgroundColor(tcell.ColorDefault)
	methodDropdown.
		SetBorder(true).
		SetTitle("Method").
		SetTitleAlign(tview.AlignLeft)

	// Url Input Field
	urlInputField := tview.NewInputField().
		SetPlaceholder("http://api.github.com").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
		SetFieldBackgroundColor(tcell.ColorDefault)
	urlInputField.
		SetBorder(true).
		SetTitle("URL").
		SetTitleAlign(tview.AlignLeft)

	// Request Body Input Field (for POST)
	bodyInputField := tview.NewInputField().
		SetPlaceholder(`{"key": "value"}`).
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDefault)
	bodyInputField.
		SetBorder(true).
		SetTitle("Body").
		SetTitleAlign(tview.AlignLeft)

	// Response TextView
	responseTextView := createTextView("Response", true)

	// Console TextView
	consoleTextView := createTextView("Console", true)

	// Combine Method Dropdown and URL Input Field horizontally
	requestHeaderFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(methodDropdown, 10, 0, false).
		AddItem(urlInputField, 0, 1, true)

	// Vertical layout for request components
	requestFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(requestHeaderFlex, 3, 1, true)

	// Main layout
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(requestFlex, 0, 1, true).
		AddItem(responseTextView, 0, 3, false).
		AddItem(consoleTextView, 5, 1, false)

	// Initialize the UI struct
	ui := &UI{
		App:              app,
		Layout:           mainFlex,
		MethodDropdown:   methodDropdown,
		URLInputField:    urlInputField,
		RequestFlex:      requestFlex,
		BodyInputField:   bodyInputField,
		ResponseTextView: responseTextView,
		ConsoleTextView:  consoleTextView,
	}

	// Setup event handlers
	ui.setupEventHandlers()

	return ui
}

// setupEventHandlers sets up the event handlers for the UI components.
func (ui *UI) setupEventHandlers() {
	// Function to perform the HTTP request
	performRequest := func() {
		url := ui.URLInputField.GetText()
		methodIndex, _ := ui.MethodDropdown.GetCurrentOption()
		method := []string{"GET", "POST"}[methodIndex]
		body := ui.BodyInputField.GetText()

		go func() {
			formattedResponse, status, err := handlers.PerformRequest(method, url, body)

			ui.App.QueueUpdateDraw(func() {
				if err != nil {
					ui.ConsoleTextView.SetText(fmt.Sprintf("[red]%v", err))
					ui.ResponseTextView.SetText("")
				} else {
					ui.ResponseTextView.SetText(formattedResponse)
					ui.ConsoleTextView.SetText(fmt.Sprintf("[green]Request successful (%s)", status))
					ui.App.SetFocus(ui.ResponseTextView)
				}
			})
		}()
	}

	// Update the visibility of BodyInputField based on the selected method
	updateRequestFlex := func() {
		methodIndex, _ := ui.MethodDropdown.GetCurrentOption()
		if methodIndex == 1 { // POST method
			if ui.RequestFlex.GetItemCount() < 2 {
				ui.RequestFlex.AddItem(ui.BodyInputField, 3, 1, false)
			}
		} else {
			if ui.Layout.GetItemCount() >= 2 {
				ui.RequestFlex.RemoveItem(ui.BodyInputField)
			}
		}
	}

	// Method Dropdown event handler
	ui.MethodDropdown.SetSelectedFunc(func(text string, index int) {
		go func() {
			ui.App.QueueUpdateDraw(func() {
				updateRequestFlex()
			})
		}()
	})

	ui.MethodDropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			ui.App.SetFocus(ui.URLInputField)
			return nil
		}
		return event
	})

	// URL Input Field event handler
	ui.URLInputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			performRequest()
		}
	})

	ui.URLInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			methodIndex, _ := ui.MethodDropdown.GetCurrentOption()
			if methodIndex == 1 { // POST method
				ui.App.SetFocus(ui.BodyInputField)
			} else {
				if ui.ResponseTextView.GetText(true) != "" {
					ui.App.SetFocus(ui.ResponseTextView)
				} else {
					ui.App.SetFocus(ui.MethodDropdown)
				}
			}
			return nil
		}
		return event
	})

	// Body Input Field event handler
	ui.BodyInputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			performRequest()
		}
	})

	ui.BodyInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
				if ui.ResponseTextView.GetText(true) != "" {
					ui.App.SetFocus(ui.ResponseTextView)
				} else {
					ui.App.SetFocus(ui.MethodDropdown)
				}
			return nil
		}
		return event
	})

	ui.ResponseTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			ui.App.SetFocus(ui.MethodDropdown)
			return nil
		}
		return event
	})
}

// createTextView creates a configured TextView.
func createTextView(title string, wrap bool) *tview.TextView {
	newTextView := tview.NewTextView().
		SetWrap(wrap).
		SetDynamicColors(true).
		SetScrollable(true)
	newTextView.
		SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft)
	return newTextView
}
