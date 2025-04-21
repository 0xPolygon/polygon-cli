package ui

import "github.com/rivo/tview"

func Start() {
	// Initialize the UI application
	app := NewApp()

	// Create the main layout
	layout := NewLayout(app)

	// Set the root view of the application
	app.SetRoot(layout, true)

	// Start the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// NewApp creates a new Tview application
func NewApp() *tview.Application {
	return tview.NewApplication()
}

// NewLayout creates the main layout for the application
func NewLayout(app *tview.Application) tview.Primitive {
	// Create a text view for the header
	header := tview.NewTextView().
		SetText("Welcome to the Tview Application!").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle(" Header ")

	// Create a text view for the body
	body := tview.NewTextView().
		SetText("This is the main content area.").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle(" Body ")

	// Create a vertical layout with header and body
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(body, 0, 1, true)

	return layout
}
