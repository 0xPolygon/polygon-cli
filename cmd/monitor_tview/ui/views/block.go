package views

import (
	"fmt"

	"github.com/0xPolygon/polygon-cli/cmd/monitor_tview/router"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewBlockView(app *tview.Application, r *router.Router, id string) tview.Primitive {
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	text.SetText(fmt.Sprintf(
		"[::b]Block Info[::-]\n\nHash/Number: %s\nTimestamp: 13:51:12\nTx Count: 143\nGas Used: 14.3M\nMiner: 0xdead...beef",
		id,
	))
	text.SetBorder(true)
	text.SetTitle(" Block View ")

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			r.Navigate("network", "")
			return nil
		}
		return event
	})

	return text
}
