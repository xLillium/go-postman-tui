package ui

import (
	"github.com/rivo/tview"
)

func (ui *UI) setupLayout() {
	// Left panel for saved requests
	ui.RequestList = tview.NewList()
	ui.RequestList.
		SetBorder(true).
		SetTitle("Saved Requests").
		SetTitleAlign(tview.AlignLeft)
	// Combine Method Dropdown and URL Input Field horizontally
	requestHeaderFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ui.MethodDropdown, 10, 0, false).
		AddItem(ui.URLInputField, 0, 1, true)

	// Vertical layout for request components
	ui.RequestFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(requestHeaderFlex, 3, 1, true)

	// Main layout (request form, response, console)
	mainContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.RequestFlex, 0, 1, true).
		AddItem(ui.ResponseTextView, 0, 3, false).
		AddItem(ui.ConsoleTextView, 5, 1, false)

	// Root layout with left panel
	ui.RootLayout = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ui.RequestList, 30, 0, false). // Left panel width
		AddItem(mainContent, 0, 1, true)

	buttonsFlex := tview.NewFlex().
		AddItem(nil, 2, 1, false).
		AddItem(ui.saveRequestButton, 0, 1, false).
		AddItem(nil, 2, 1, false).
		AddItem(ui.cancelSaveRequestButton, 0, 1, false).
		AddItem(nil, 2, 1, false)

	modalContent := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Top spacing to center the input
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
							AddItem(nil, 2, 1, false).
							AddItem(ui.requestNameInputfield, 0, 1, true). // Make the input have same margin as button line
							AddItem(nil, 2, 1, false), 0, 1, true).
		AddItem(buttonsFlex, 1, 1, false) // Button line
	modalContent.SetBorder(true).
		SetTitle("Save Request").
		SetTitleAlign(tview.AlignCenter)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modalContent, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	// Set up pages to include the modal
	ui.Pages = tview.NewPages().
		AddPage("main", ui.RootLayout, true, true).
		AddPage("saveRequestDialog", modal, true, false)

	ui.InitialFocus = ui.URLInputField
}
