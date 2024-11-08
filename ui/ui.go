package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI struct holds all UI components
type UI struct {
	App       *tview.Application
	Flex      *tview.Flex
	LeftBox   *tview.Box
	TopInput  *tview.InputField
	MiddleBox *tview.TextView
	BottomBox *tview.TextView
	RightBox  *tview.Box
}

// NewUI creates and configures the UI components
func NewUI(app *tview.Application) *UI {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	// Left Box
	leftBox := createBox("Left (1/2 x width of Top)")

	// Top Input Field
	topInput := createInputField("URL: ", "Enter URL and press Enter", 50)

	// Middle TextView
	middleBox := createTextView("Response", true)
	middleBox.SetWrap(true)

	// Bottom TextView
	bottomBox := createTextView("Console", true)

	// Right Box
	rightBox := createBox("Right (20 cols)")

	// Main Layout
	flex := tview.NewFlex()

	// Nested Flex
	innerFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Add items to the inner flex
	innerFlex.AddItem(topInput, 3, 1, true) // Set focus to the input field
	innerFlex.AddItem(middleBox, 0, 3, false)
	innerFlex.AddItem(bottomBox, 5, 1, false)

	// Add items to the main flex
	flex.AddItem(leftBox, 0, 1, false)
	flex.AddItem(innerFlex, 0, 2, true) // Set focus to enable input
	flex.AddItem(rightBox, 20, 1, false)

	return &UI{
		App:       app,
		Flex:      flex,
		LeftBox:   leftBox,
		TopInput:  topInput,
		MiddleBox: middleBox,
		BottomBox: bottomBox,
		RightBox:  rightBox,
	}
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
		SetDynamicColors(true)

	textView.
		SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft)

	return textView
}

func createInputField(label, title string, fieldWidth int) *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel(label).
		// SetLabelColor(tcell.ColorDefault).
		SetFieldWidth(fieldWidth).
		SetFieldBackgroundColor(tcell.ColorDefault)

	inputField.
		SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft)

	return inputField
}
