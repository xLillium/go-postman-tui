package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) setupComponents() {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	// HTTP Method Dropdown
	ui.MethodDropdown = tview.NewDropDown().
		SetOptions([]string{"GET", "POST"}, nil).
		SetCurrentOption(0).
		SetFieldWidth(6).
		SetFieldBackgroundColor(tcell.ColorDefault)
	ui.MethodDropdown.
		SetBorder(true).
		SetTitle("Method").
		SetTitleAlign(tview.AlignLeft)

	// URL Input Field
	ui.URLInputField = tview.NewInputField().
		SetPlaceholder("http://example.com").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorGrey)).
		SetFieldBackgroundColor(tcell.ColorDefault)
	ui.URLInputField.
		SetBorder(true).
		SetTitle("URL").
		SetTitleAlign(tview.AlignLeft)

	// Request Body Input Field (for POST)
	ui.BodyInputField = tview.NewInputField().
		SetPlaceholder(`{"key": "value"}`).
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorGrey)).
		SetFieldBackgroundColor(tcell.ColorDefault)
	ui.BodyInputField.
		SetBorder(true).
		SetTitle("Body").
		SetTitleAlign(tview.AlignLeft)

	// Response TextView
	ui.ResponseTextView = createTextView("Response", true)

	// Console TextView
	ui.ConsoleTextView = createTextView("Console", true)
}

func createTextView(title string, wrap bool) *tview.TextView {
	newTview := tview.NewTextView().
		SetWrap(wrap).
		SetDynamicColors(true).
		SetScrollable(true)
	newTview.
		SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft)
	return newTview
}
