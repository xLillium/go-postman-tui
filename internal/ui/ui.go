package ui

import (
    "github.com/rivo/tview"
)

type UI struct {
    App          *tview.Application
    RootLayout   *tview.Flex
    InitialFocus tview.Primitive

    // UI Components
    MethodDropdown   *tview.DropDown
    URLInputField    *tview.InputField
    BodyInputField   *tview.InputField
    ResponseTextView *tview.TextView
    ConsoleTextView  *tview.TextView
    RequestFlex      *tview.Flex
}

func Initialize(app *tview.Application) *UI {
    ui := &UI{
        App: app,
    }

    ui.setupComponents()
    ui.setupLayout()
    ui.setupEventHandlers()

    return ui
}

