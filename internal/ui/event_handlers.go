package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/xlillium/go-postman-tui/internal/handlers"
	"github.com/xlillium/go-postman-tui/internal/utils"
)

func (ui *UI) setupEventHandlers() {
	httpHandler := handlers.NewHandler(utils.DefaultHTTPClient)

	// Function to perform the HTTP request
	performRequest := func() {
		url := ui.URLInputField.GetText()
		methodIndex, _ := ui.MethodDropdown.GetCurrentOption()
		method := utils.SupportedMethods[methodIndex]
		body := ui.BodyInputField.GetText()

		go func() {
			formattedResponse, status, err := httpHandler.PerformRequest(method, url, body)

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
				ui.RequestFlex.AddItem(ui.BodyInputField, 0, 1, false)
			}
		} else {
			if ui.RequestFlex.GetItemCount() >= 2 {
				ui.RequestFlex.RemoveItem(ui.BodyInputField)
			}
		}
	}

	// Method Dropdown event handlers
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

	// URL Input Field event handlers
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

	// Body Input Field event handlers
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

	// Response Text View event handlers
	ui.ResponseTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			ui.App.SetFocus(ui.MethodDropdown)
			return nil
		}
		return event
	})

	// Initial call to set the correct state
	updateRequestFlex()
}
