package views

import (
	"fmt"

	"github.com/0xPolygon/polygon-cli/cmd/monitor_tview/router"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewTxView(app *tview.Application, r *router.Router, id string) tview.Primitive {
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	text.SetText(fmt.Sprintf(
		"[::b]Transaction Info[::-]\n\nHash: %s\nMethod: transfer\nTo: 0xaaa...bbb\nAmount: 1.25 ETH\nGas Price: 35 Gwei",
		id,
	))
	text.SetBorder(true)
	text.SetTitle(" Tx View ")

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
