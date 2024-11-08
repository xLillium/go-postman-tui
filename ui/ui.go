package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xlillium/go-postman-tui/handlers"
)

// UI struct holds all UI components
type UI struct {
	App            *tview.Application
	Flex           *tview.Flex
	LeftBox        *tview.Box
	MethodDropdown *tview.DropDown
	UrlInput       *tview.InputField
	BodyInput      *tview.InputField
	ResponseBox    *tview.TextView
	ConsoleBox     *tview.TextView
	RightBox       *tview.Box
}

// NewUI creates and configures the UI components
func NewUI(app *tview.Application) *UI {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	// Left Box
	leftBox := createBox("Left (1/2 x width of Top)")

	// HTTP Method DropDown
	methodDropDown := tview.NewDropDown().
		SetOptions([]string{"GET", "POST"}, nil).
		SetCurrentOption(0).
		SetFieldWidth(4).
		SetFieldBackgroundColor(tcell.ColorDefault)
	methodDropDown.
		SetBorder(true).
		SetTitle("Method").
		SetTitleAlign(tview.AlignLeft)

		// Url Input Field (URL input)
	urlInput := tview.NewInputField().
		SetPlaceholder("http://api.github.com").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
		SetFieldBackgroundColor(tcell.ColorDefault)
	urlInput.
		SetBorder(true).
		SetTitle("URL").
		SetTitleAlign(tview.AlignLeft)

	// Url Input Row (Method Dropdown and URL Input)
	urlInputRow := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(methodDropDown, 10, 0, false).
		AddItem(urlInput, 0, 1, true)

	// Request Body Input Field (for POST)
	bodyInput := tview.NewInputField().
		SetLabel("").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDefault)

	bodyInput.
		SetBorder(true).
		SetTitle("Body").
		SetTitleAlign(tview.AlignLeft)

	// Request Box (Input Row and Body Input)
	requestBox := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(urlInputRow, 3, 1, true)

	// Middle TextView
	responseBox := createTextView("Response", true)

	// Bottom TextView
	consoleBox := createTextView("Console", true)

	// Right Box
	rightBox := createBox("Right (20 cols)")

	// Main Flex Layout
	innerFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(requestBox, 0, 1, true).
		AddItem(responseBox, 0, 3, false).
		AddItem(consoleBox, 5, 1, false)

	flex := tview.NewFlex().
		AddItem(leftBox, 0, 1, false).
		AddItem(innerFlex, 0, 2, true).
		AddItem(rightBox, 20, 1, false)

	// Initialize the UI struct
	ui := &UI{
		App:            app,
		Flex:           flex,
		LeftBox:        leftBox,
		MethodDropdown: methodDropDown,
		UrlInput:       urlInput,
		BodyInput:      bodyInput,
		ResponseBox:    responseBox,
		ConsoleBox:     consoleBox,
		RightBox:       rightBox,
	}

	// Function to perform the request
	performRequest := func() {
		url := urlInput.GetText()
		methodIndex, _ := methodDropDown.GetCurrentOption()
		methodText := []string{"GET", "POST"}[methodIndex]
		body := bodyInput.GetText()

		go func() {
			var formattedResponse, status string
			var err error

			formattedResponse, status, err = handlers.PerformRequest(methodText, url, body)

			ui.App.QueueUpdateDraw(func() {
				if err != nil {
					ui.ConsoleBox.SetText(fmt.Sprintf("[red]%v", err))
					ui.ResponseBox.SetText("")
				} else {
					ui.ResponseBox.SetText(formattedResponse)
					ui.ConsoleBox.SetText(fmt.Sprintf("[green]Request successful (%s)", status))
					ui.App.SetFocus(ui.ResponseBox)
				}
			})
		}()
	}

	// Update the visibility of bodyInput based on the selected method
	updateInputFlex := func() {
		method, _ := methodDropDown.GetCurrentOption()
		if method == 1 { // POST method
			if requestBox.GetItemCount() < 2 {
				requestBox.AddItem(bodyInput, 3, 1, false)
			}
		} else {
			if requestBox.GetItemCount() == 2 {
				requestBox.RemoveItem(bodyInput)
			}
		}
		// app.Draw()
	}

	// Set a handler for when the selected option changes
	methodDropDown.SetSelectedFunc(func(text string, index int) {
		go func() {
			app.QueueUpdateDraw(func() {
				updateInputFlex()
			})
		}()
	})

	// Initial call to set the correct state
	updateInputFlex()

	// InputCapture for methodDropDown
	methodDropDown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyRight:
			app.SetFocus(urlInput)
			return nil
			// case tcell.KeyEnter:
			// 	return nil
		}
		return event
	})

	// InputCapture for topInput
	urlInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyLeft:
			app.SetFocus(methodDropDown)
			return nil
		case tcell.KeyEnter:
			method, _ := methodDropDown.GetCurrentOption()
			if method == 1 { // POST method
				app.SetFocus(bodyInput)
			} else {
				performRequest()
			}
			return nil
		}
		return event
	})

	// InputCapture for bodyInput
	bodyInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(responseBox)
			return nil
		case tcell.KeyEnter:
			performRequest()
			return nil
		}
		return event
	})

	// InputCapture for methodDropDown
	responseBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyRight:
			app.SetFocus(urlInput)
			return nil
			// case tcell.KeyEnter:
			// 	return nil
		}
		return event
	})

	// Return the UI struct
	return ui
}

// Helper functions
func createBox(title string) *tview.Box {
	return tview.NewBox().
		SetBorder(true).
		SetTitle(title)
}

func createTextView(title string, wrap bool) *tview.TextView {
	textView := tview.NewTextView().
		SetWrap(wrap).
		SetDynamicColors(true).
		SetScrollable(true)

	textView.
		SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft)

	return textView
}

func createInputField(label, title string, placeholder string) *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel(label).
		SetPlaceholder(placeholder).
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorDefault)).
		SetFieldBackgroundColor(tcell.ColorDefault)

	inputField.
		SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft)

	return inputField
}
