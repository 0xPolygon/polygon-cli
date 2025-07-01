package renderer

// createBlockDetailPage creates the block detail view with side-by-side panes
func (t *TviewRenderer) createBlockDetailPage() {
	// Create left pane as transaction table
	t.blockDetailLeft = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0).   // Fix the header row
		SetSeparator(' ') // Use space as separator
	t.blockDetailLeft.SetBorder(true).SetTitle(" Transactions ")

	// Set up transaction table headers
	headers := []string{"INDEX", "FROM", "TO", "GAS LIMIT", "INPUT"}
	aligns := []int{tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignLeft}
	expansions := []int{1, 3, 3, 2, 2}

	for col, header := range headers {
		t.blockDetailLeft.SetCell(0, col, tview.NewTableCell(header).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(aligns[col]).
			SetExpansion(expansions[col]).
			SetAttributes(tcell.AttrBold))
	}

	// Create right pane for raw JSON
	t.blockDetailRight = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	t.blockDetailRight.SetBorder(true).SetTitle(" Raw JSON ")
	t.blockDetailRight.SetText("Select a block to view its JSON representation")

	// Create flex container to hold both panes side by side
	t.blockDetailPage = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.blockDetailLeft, 0, 1, true). // Left pane: 50% width, focusable
		AddItem(t.blockDetailRight, 0, 1, true) // Right pane: 50% width, focusable
}

// showBlockDetail navigates to block detail page and populates it
func (t *TviewRenderer) showBlockDetail(block rpctypes.PolyBlock) {
	log.Debug().
		Str("blockHash", block.Hash().Hex()).
		Str("blockNumber", block.Number().String()).
		Msg("showBlockDetail called")

	// Store the current block for transaction selection
	t.currentBlockMu.Lock()
	t.currentBlock = block
	t.currentBlockMu.Unlock()

	// Clear existing table rows (except header)
	rowCount := t.blockDetailLeft.GetRowCount()
	for row := 1; row < rowCount; row++ {
		for col := 0; col < 5; col++ {
			t.blockDetailLeft.SetCell(row, col, nil)
		}
	}

	// Populate transaction table
	transactions := block.Transactions()
	for i, tx := range transactions {
		row := i + 1 // +1 to account for header row

		// Column 0: Transaction index
		t.blockDetailLeft.SetCell(row, 0, tview.NewTableCell(strconv.Itoa(i)).SetAlign(tview.AlignRight))

		// Column 1: From address (truncated)
		fromAddr := truncateHash(tx.From().Hex(), 6, 4)
		t.blockDetailLeft.SetCell(row, 1, tview.NewTableCell(fromAddr).SetAlign(tview.AlignLeft))

		// Column 2: To address (truncated), handle contract creation (empty address)
		var toAddr string
		if tx.To().Hex() != "0x0000000000000000000000000000000000000000" {
			toAddr = truncateHash(tx.To().Hex(), 6, 4)
		} else {
			toAddr = "CONTRACT"
		}
		t.blockDetailLeft.SetCell(row, 2, tview.NewTableCell(toAddr).SetAlign(tview.AlignLeft))

		// Column 3: Gas limit
		gasLimit := formatNumber(tx.Gas())
		t.blockDetailLeft.SetCell(row, 3, tview.NewTableCell(gasLimit).SetAlign(tview.AlignRight))

		// Column 4: First 4 bytes of input data
		inputData := "N/A"
		if len(tx.Data()) >= 4 {
			inputData = fmt.Sprintf("0x%x", tx.Data()[:4])
		} else if len(tx.Data()) > 0 {
			inputData = fmt.Sprintf("0x%x", tx.Data())
		}
		t.blockDetailLeft.SetCell(row, 4, tview.NewTableCell(inputData).SetAlign(tview.AlignLeft))
	}

	// Update table title with transaction count
	title := fmt.Sprintf(" Transactions (%d) ", len(transactions))
	t.blockDetailLeft.SetTitle(title)

	// Right pane shows pretty-printed JSON of the block
	blockJSON, err := rpctypes.PolyBlockToPrettyJSON(block)
	if err != nil {
		t.blockDetailRight.SetText(fmt.Sprintf("Error marshaling block JSON: %v", err))
	} else {
		// Pretty print the JSON
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, blockJSON, "", "  "); err != nil {
			t.blockDetailRight.SetText(fmt.Sprintf("Error formatting JSON: %v", err))
		} else {
			t.blockDetailRight.SetText(prettyJSON.String())
		}
	}

	t.pages.SwitchToPage("block-detail")

	// Reset table selection to the first row (row 1, since row 0 is header)
	if len(transactions) > 0 {
		t.blockDetailLeft.Select(1, 0)
	}

	// Set focus to the left pane (transaction table) by default
	t.app.SetFocus(t.blockDetailLeft)
}
