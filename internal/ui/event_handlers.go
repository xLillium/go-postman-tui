package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/xlillium/go-postman-tui/internal/domain"
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

	// Saved Requests event handler
	ui.RequestList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			requests := ui.Storage.GetRequests()
			index := ui.RequestList.GetCurrentItem()
			req := requests[index]
			ui.MethodDropdown.SetCurrentOption(req.Method)
			ui.URLInputField.SetText(req.URL)
			ui.BodyInputField.SetText(req.Body)
			ui.ResponseTextView.SetText(req.Response)
			ui.ConsoleTextView.SetText(fmt.Sprintf("[green]Loaded request '%s'", req.Name))
			return nil
		}
		return event
	})

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

	//
	//  Save Request Dialog event handlers
	//

	// Request Name Input Field event handlers
	ui.requestNameInputfield.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			ui.App.SetFocus(ui.saveRequestButton)
			return nil
		} else if event.Key() == tcell.KeyEnter {
			methodIndex, _ := ui.MethodDropdown.GetCurrentOption()
			req := domain.Request{
				Name:     ui.requestNameInputfield.GetText(),
				Method:   methodIndex,
				URL:      ui.URLInputField.GetText(),
				Body:     ui.BodyInputField.GetText(),
				Response: ui.ResponseTextView.GetText(false),
			}
			ui.Storage.AddRequest(req)
			ui.Storage.Save()
			ui.RequestList.AddItem(req.Name, "", 0, nil)
			ui.Pages.ShowPage("main")
			ui.Pages.HidePage("saveRequestDialog")
			ui.ConsoleTextView.SetText("Saved Request : " + req.Name)
			ui.requestNameInputfield.SetText("")
		}

		return event
	})

	// Save Request Button event handlers
	ui.saveRequestButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			ui.App.SetFocus(ui.cancelSaveRequestButton)
			return nil
		}
		return event
	})

	// Cancel Request Button event handlers
	ui.cancelSaveRequestButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			ui.App.SetFocus(ui.requestNameInputfield)
			return nil
		}
		return event
	})

	// Initial call to set the correct state
	updateRequestFlex()
}

func (ui *UI) setupGlobalKeybindings() {
	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModCtrl {
			ui.ConsoleTextView.SetText(event.Name())
			switch event.Name() {
			case "Ctrl+S":
				ui.Pages.ShowPage("saveRequestDialog")
				ui.Pages.HidePage("main")
				return nil
			case "Ctrl+A":
				ui.ConsoleTextView.SetText("request list")
				ui.App.SetFocus(ui.RequestList)
			}
		}
		return event
	})
}
