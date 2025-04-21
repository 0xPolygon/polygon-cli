package views

import (
	"fmt"

	"github.com/0xPolygon/polygon-cli/cmd/monitor_tview/router"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewNetworkView(app *tview.Application, router *router.Router, rpcURL string) tview.Primitive {
	root := tview.NewFlex().SetDirection(tview.FlexRow)

	// ─ Header
	rpcView := tview.NewTextView().SetText("RPC: https://...").SetWrap(false)
	chainIDView := tview.NewTextView().SetText("Chain ID: 137").SetWrap(false)
	latestBlockView := tview.NewTextView().SetText("Latest: 50192821").SetWrap(false)
	safeBlockView := tview.NewTextView().SetText("Safe: 50192819").SetWrap(false)
	finalizedView := tview.NewTextView().SetText("Finalized: 50192818").SetWrap(false)
	gasPriceView := tview.NewTextView().SetText("Gas: 28 Gwei").SetWrap(false)

	row1 := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(rpcView, 0, 1, false).
		AddItem(chainIDView, 0, 1, false).
		AddItem(latestBlockView, 0, 1, false)

	row2 := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(safeBlockView, 0, 1, false).
		AddItem(finalizedView, 0, 1, false).
		AddItem(gasPriceView, 0, 1, false)

	header := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(row1, 1, 1, false).
		AddItem(row2, 1, 1, false)
	header.SetBorder(true).SetTitle(" Network Info ")

	// header := tview.NewTextView().
	// 	SetDynamicColors(true).
	// 	SetTextAlign(tview.AlignLeft)
	// header.SetText(fmt.Sprintf("RPC URL: %s\nChain ID: 137\nLatest Block: 50192821", rpcURL))
	// header.SetBorder(true).SetTitle(" Network Info ")

	// ─ Block Table
	blocks := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false)
	blocks.SetBorder(true).SetTitle(" Latest Blocks ")

	blockHeaders := []string{"Block#", "Time", "Tx Count", "Gas Used", "Hash", "Miner"}
	blockData := [][]string{
		{"50192821", "13:51:12", "143", "14.3M", "0xabc...123", "0xdead...beef"},
		{"50192820", "13:51:05", "159", "15.1M", "0xdef...456", "0xbeef...cafe"},
	}

	for col, title := range blockHeaders {
		cell := tview.NewTableCell(fmt.Sprintf("[::b]%s", title)).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetExpansion(1)
		blocks.SetCell(0, col, cell)
	}

	for rowIdx, row := range blockData {
		for colIdx, val := range row {
			cell := tview.NewTableCell(val).
				SetAlign(tview.AlignCenter).
				SetExpansion(1).
				SetMaxWidth(1)
			blocks.SetCell(rowIdx+1, colIdx, cell)
		}
	}

	blocks.SetSelectedFunc(func(row, col int) {
		if row == 0 {
			return // skip header
		}
		blockID := blocks.GetCell(row, 0).Text // Block number
		router.Navigate("block", blockID)
	})

	// ─ Transaction Table
	txs := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(false, false)
	txs.SetBorder(true).SetTitle(" Latest Transactions ")

	txHeaders := []string{"Tx Hash", "Method", "To", "Amount", "Gas Price"}
	txData := [][]string{
		{"0xabc...001", "transfer", "0xaaa...bbb", "1.25 ETH", "35 Gwei"},
		{"0xabc...002", "0xa9059cbb", "0xccc...ddd", "500 USDC", "28 Gwei"},
	}

	for col, title := range txHeaders {
		cell := tview.NewTableCell(fmt.Sprintf("[::b]%s", title)).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetExpansion(1)
		txs.SetCell(0, col, cell)
	}

	for rowIdx, row := range txData {
		for colIdx, val := range row {
			cell := tview.NewTableCell(val).
				SetAlign(tview.AlignCenter).
				SetExpansion(1).
				SetMaxWidth(1)
			txs.SetCell(rowIdx+1, colIdx, cell)
		}
	}

	txs.SetSelectedFunc(func(row, col int) {
		if row == 0 {
			return
		}
		txHash := txs.GetCell(row, 0).Text
		router.Navigate("tx", txHash)
	})

	// ─ Tab to switch focus
	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			if blocks.HasFocus() {
				blocks.SetSelectable(false, false)
				txs.SetSelectable(true, false)
				app.SetFocus(txs)
			} else {
				txs.SetSelectable(false, false)
				blocks.SetSelectable(true, false)
				app.SetFocus(blocks)
			}
			return nil
		}
		return event
	})

	// ─ Layout
	root.
		AddItem(header, 3, 0, false).
		AddItem(blocks, 0, 1, true).
		AddItem(txs, 0, 1, false)

	app.SetFocus(blocks)
	return root
}
