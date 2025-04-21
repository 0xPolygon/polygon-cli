package views

import (
	"github.com/rivo/tview"
)

func NewHelpModal(app *tview.Application, dismiss func()) tview.Primitive {
	helpText := `[::b]?[::-] Show this help
[::b]↑↓[::-] Navigate block/tx lists
[::b]Tab[::-] Switch focus
[::b]Enter[::-] View details
[::b]ESC[::-] Back
[::b]q[::-] Quit
`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(_ int, _ string) {
			if dismiss != nil {
				dismiss()
			}
		})

	modal.SetTitle(" Help ")
	modal.SetBorder(true)

	return modal
}
