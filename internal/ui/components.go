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

	//
	//  Save Request Dialog Components
	//

	// Request Name Input Field
	ui.requestNameInputfield = tview.NewInputField().
		SetLabel("Name: ").
		SetPlaceholder("My Request Name").
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorGrey)).
		SetFieldBackgroundColor(tcell.ColorDefault)

	// Save Request Button
	ui.saveRequestButton = tview.NewButton("Save").
		SetStyle(tcell.StyleDefault.Background(tcell.ColorLightGrey).Foreground(tcell.ColorGrey)).
		SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorLightGrey)).
		SetSelectedFunc(func() {
			ui.Pages.ShowPage("main")
			ui.Pages.HidePage("saveRequestDialog")
			ui.ConsoleTextView.SetText("Saved Request : " + ui.requestNameInputfield.GetText())
			ui.requestNameInputfield.SetText("")
		})

	// Cancel Save Request Button
	ui.cancelSaveRequestButton = tview.NewButton("Cancel").
		SetStyle(tcell.StyleDefault.Background(tcell.ColorLightGray).Foreground(tcell.ColorGrey)).
		SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorLightGrey)).
		SetSelectedFunc(func() {
			ui.Pages.ShowPage("main")
			ui.Pages.HidePage("saveRequestDialog")
			ui.ConsoleTextView.SetText("Canceled Save Request : " + ui.requestNameInputfield.GetText())
			ui.requestNameInputfield.SetText("")
		})

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
