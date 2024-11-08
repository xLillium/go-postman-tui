package ui

import (
    "github.com/rivo/tview"
)

func (ui *UI) setupLayout() {
    // Combine Method Dropdown and URL Input Field horizontally
    requestHeaderFlex := tview.NewFlex().
        SetDirection(tview.FlexColumn).
        AddItem(ui.MethodDropdown, 10, 0, false).
        AddItem(ui.URLInputField, 0, 1, true)

    // Vertical layout for request components
    ui.RequestFlex = tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(requestHeaderFlex, 3, 1, true)

    // Main layout
    ui.RootLayout = tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(ui.RequestFlex, 0, 1, true).
        AddItem(ui.ResponseTextView, 0, 3, false).
        AddItem(ui.ConsoleTextView, 5, 1, false)

    ui.InitialFocus = ui.URLInputField
}

