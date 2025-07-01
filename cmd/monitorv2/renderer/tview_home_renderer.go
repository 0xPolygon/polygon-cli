package renderer

// createColumnDefinitions creates all sortable column definitions
func createColumnDefinitions() []ColumnDef {
	return []ColumnDef{
		{
			Name: "BLOCK #", Key: "number", Align: tview.AlignLeft, Expansion: 1,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Number() },
			CompareFunc: compareNumbers,
		},
		{
			Name: "TIME", Key: "time", Align: tview.AlignLeft, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Time() },
			CompareFunc: compareUint64,
		},
		{
			Name: "INTERVAL", Key: "interval", Align: tview.AlignCenter, Expansion: 1,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Time() }, // Will be calculated separately
			CompareFunc: compareUint64,
		},
		{
			Name: "HASH", Key: "hash", Align: tview.AlignCenter, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Hash().Hex() },
			CompareFunc: compareStrings,
		},
		{
			Name: "TXS", Key: "txs", Align: tview.AlignCenter, Expansion: 1,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return uint64(len(block.Transactions())) },
			CompareFunc: compareUint64,
		},
		{
			Name: "SIZE", Key: "size", Align: tview.AlignCenter, Expansion: 1,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Size() },
			CompareFunc: compareUint64,
		},
		{
			Name: "BASE FEE", Key: "basefee", Align: tview.AlignCenter, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.BaseFee() },
			CompareFunc: compareNumbers,
		},
		{
			Name: "GAS USED", Key: "gasused", Align: tview.AlignCenter, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.GasUsed() },
			CompareFunc: compareUint64,
		},
		{
			Name: "GAS %", Key: "gaspct", Align: tview.AlignCenter, Expansion: 1,
			SortFunc: func(block rpctypes.PolyBlock) interface{} {
				if block.GasLimit() == 0 {
					return uint64(0)
				}
				return uint64(float64(block.GasUsed()) / float64(block.GasLimit()) * 10000) // *10000 for precision
			},
			CompareFunc: compareUint64,
		},
		{
			Name: "GAS LIMIT", Key: "gaslimit", Align: tview.AlignRight, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.GasLimit() },
			CompareFunc: compareUint64,
		},
		{
			Name: "STATE ROOT", Key: "stateroot", Align: tview.AlignRight, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Root().Hex() },
			CompareFunc: compareStrings,
		},
	}
}

// createHomePage creates the main page with 2-column top section and block listing table
func (t *TviewRenderer) createHomePage() {
	// Create the two top panes
	t.homeStatusPane = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	t.homeStatusPane.SetBorder(true).SetTitle(" Status ")
	t.homeStatusPane.SetText("Initializing...")

	// Create metrics table
	t.homeMetricsPane = tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false). // Non-selectable
		SetFixed(0, 0).              // No fixed rows/columns
		SetSeparator('|')            // Vertical pipe separator
	t.homeMetricsPane.SetBorder(true).SetTitle(" Metrics ")

	// Configure table to use full width
	// The expansion will be handled by individual cells

	// Initialize with 1 rows of placeholders
	t.homeMetricsPane.SetCell(0, 0,
		tview.NewTableCell("Initializing... "). // Add trailing space
							SetAlign(tview.AlignLeft).
							SetExpansion(1))

	// Create horizontal flex container for the 2 top sections
	t.homeTopSection = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.homeStatusPane, 0, 1, false). // Left: Status (1/3 width)
		AddItem(t.homeMetricsPane, 0, 2, false) // Right: Metrics (2/3 width)

	// Create the block table
	t.homeTable = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0).   // Fix the header row so it doesn't scroll
		SetSeparator(' ') // Use space as separator instead of borders

	// Add border and title to the table
	t.homeTable.SetBorder(true).SetTitle(" Blocks ")

	// Set up table headers with sort indicators
	t.updateTableHeaders()

	// Set up selection handler for Enter key
	t.homeTable.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(t.blocks) { // Skip header row
			// Navigate to block detail page
			t.showBlockDetail(t.blocks[row-1])
		}
	})

	// Set up selection change handler to track manual selection
	t.homeTable.SetSelectionChangedFunc(func(row, column int) {
		if row > 0 {
			// Get block data safely
			t.blocksMu.RLock()
			if row-1 < len(t.blocks) {
				selectedBlock := t.blocks[row-1]
				blockHash := selectedBlock.Hash().Hex()
				t.blocksMu.RUnlock()

				// Update view state safely
				t.viewStateMu.Lock()
				if t.viewState.followMode {
					if row == 1 {
						// User selected the first row - re-enable auto-follow
						if t.viewState.manualSelect {
							t.viewState.manualSelect = false
							t.viewState.selectedBlock = ""
							log.Debug().Msg("First row selected, re-enabling auto-follow")
						}
					} else {
						// User selected a different row - disable auto-follow
						t.viewState.manualSelect = true
						t.viewState.selectedBlock = blockHash
						log.Debug().Str("hash", blockHash).Int("row", row).Msg("Manual selection detected, disabling auto-follow")
					}
				}
				t.viewStateMu.Unlock()
			} else {
				t.blocksMu.RUnlock()
			}
		}
	})

	// Create flex container to hold both sections
	t.homePage = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.homeTopSection, 10, 0, false). // Top section: 10 lines high
		AddItem(t.homeTable, 0, 1, true)         // Table: takes remaining space
}

// fetchNetworkInfo fetches and caches network information
func (t *TviewRenderer) fetchNetworkInfo(ctx context.Context) {
	t.networkInfoMu.Lock()
	defer t.networkInfoMu.Unlock()

	// Get gas price
	if gasPrice, err := t.indexer.GetGasPrice(ctx); err == nil {
		t.gasPrice = "price " + weiToGwei(gasPrice) + " gwei"
	} else {
		t.gasPrice = "price N/A"
	}

	// Get txpool status
	if txPoolStatus, err := t.indexer.GetTxPoolStatus(ctx); err == nil {
		// Parse pending count
		if pending, ok := txPoolStatus["pending"]; ok {
			if pendingBig, err := hexToDecimal(pending); err == nil {
				t.txPoolPending = "pending " + formatNumber(pendingBig.Uint64())
			} else {
				t.txPoolPending = "pending N/A"
			}
		} else {
			t.txPoolPending = "pending N/A"
		}

		// Parse queued count
		if queued, ok := txPoolStatus["queued"]; ok {
			if queuedBig, err := hexToDecimal(queued); err == nil {
				t.txPoolQueued = "queued " + formatNumber(queuedBig.Uint64())
			} else {
				t.txPoolQueued = "queued N/A"
			}
		} else {
			t.txPoolQueued = "queued N/A"
		}
	} else {
		t.txPoolPending = "pending N/A"
		t.txPoolQueued = "queued N/A"
	}

	// Get peer count
	if peerCount, err := t.indexer.GetNetPeerCount(ctx); err == nil {
		t.peerCount = "peers " + peerCount.String()
	} else {
		t.peerCount = "peers unknown"
	}

	// Measure connection latency
	if latency, err := t.indexer.MeasureConnectionLatency(ctx); err == nil {
		t.connectionStatus = formatConnectionStatus(latency)
	} else {
		t.connectionStatus = "[ERROR] Connection failed"
		log.Debug().Err(err).Msg("Failed to measure connection latency")
	}

	log.Debug().
		Str("gasPrice", t.gasPrice).
		Str("txPoolPending", t.txPoolPending).
		Str("txPoolQueued", t.txPoolQueued).
		Str("peerCount", t.peerCount).
		Str("connectionStatus", t.connectionStatus).
		Msg("Updated network info cache")
}

// updateMetricsPane updates the metrics display with new metric data
func (t *TviewRenderer) updateMetricsPane(update metrics.MetricUpdate) {
	if t.homeMetricsPane == nil {
		return
	}

	// Get cached network info
	t.networkInfoMu.RLock()
	gasPriceStr := t.gasPrice
	txPoolPendingStr := t.txPoolPending
	txPoolQueuedStr := t.txPoolQueued
	peerCountStr := t.peerCount
	t.networkInfoMu.RUnlock()

	// For now, just update with placeholder content
	// Later we'll format metrics data into the cells using atop-style format
	for row := 0; row < 8; row++ {
		var cells [5]string
		switch row {
		case 0: // BLK
			// Get current block info with mutex protection
			t.blockInfoMu.RLock()
			latestStr := formatBlockNumber(t.latestBlockNum)
			safeStr := formatBlockNumber(t.safeBlockNum)
			finalStr := formatBlockNumber(t.finalizedBlockNum)
			t.blockInfoMu.RUnlock()

			cells = [5]string{"BLOK", "late " + latestStr, "safe " + safeStr, "final " + finalStr, "[placeholder]"}
		case 1: // THR
			// Get throughput metrics
			if throughputValue, ok := t.indexer.GetMetric("throughput"); ok {
				stats := throughputValue.(metrics.ThroughputStats)
				tps10Str := "tps10 " + formatThroughput(stats.TPS10, "")
				tps30Str := "tps30 " + formatThroughput(stats.TPS30, "")
				gps10Str := "gps10 " + formatThroughput(stats.GPS10, "")
				gps30Str := "gps30 " + formatThroughput(stats.GPS30, "")
				cells = [5]string{"THRU", tps10Str, tps30Str, gps10Str, gps30Str}
			} else {
				cells = [5]string{"THRU", "tps10 N/A", "tps30 N/A", "gps10 N/A", "gps30 N/A"}
			}
		case 2: // GAS
			// Get base fee metrics
			if baseFeeValue, ok := t.indexer.GetMetric("basefee"); ok {
				stats := baseFeeValue.(metrics.BaseFeeStats)
				base10Str := "base10 " + formatBaseFee(stats.BaseFee10)
				base30Str := "base30 " + formatBaseFee(stats.BaseFee30)
				cells = [5]string{"GAS ", base10Str, base30Str, gasPriceStr, "[placeholder]"}
			} else {
				cells = [5]string{"GAS ", "base10 N/A", "base30 N/A", gasPriceStr, "[placeholder]"}
			}
		case 3: // POOL
			// Use cached txpool status
			cells = [5]string{"POOL", txPoolPendingStr, txPoolQueuedStr, "[placeholder]", "[placeholder]"}
		case 4: // SIG (1)
			// Calculate transaction counters
			eoaCount, deployCount := t.calculateTransactionCounters()
			eoaStr := "EOA " + formatNumber(eoaCount)
			deployStr := "deploy " + formatNumber(deployCount)
			cells = [5]string{"SIG1", eoaStr, deployStr, "[placeholder]", "[placeholder]"}
		case 5: // SIG (2)
			// Calculate ERC20 and NFT transaction counters
			erc20Count, nftCount := t.calculateERC20NFTCounters()
			erc20Str := "ERC20 " + formatNumber(erc20Count)
			nftStr := "NFT " + formatNumber(nftCount)
			cells = [5]string{"SIG2", erc20Str, nftStr, "Other [placeholder]", "[placeholder]"}
		case 6: // ACC
			// Calculate unique address counters
			fromCount, toCount := t.calculateUniqueAddressCounters()
			fromStr := "from " + formatNumber(fromCount)
			toStr := "to " + formatNumber(toCount)
			cells = [5]string{"ACCO", fromStr, toStr, "[placeholder]", "[placeholder]"}
		case 7: // PER
			// Use cached peer count
			cells = [5]string{"PEER", peerCountStr, "[placeholder]", "[placeholder]", "[placeholder]"}
		}

		// Set all 5 cells for this row
		for col := 0; col < 5; col++ {
			// Add padding around cell content for better spacing with separator
			cellContent := cells[col]
			if col > 0 { // Add leading space for non-first columns
				cellContent = " " + cellContent
			}
			if col < 4 { // Add trailing space for non-last columns
				cellContent = cellContent + " "
			}

			t.homeMetricsPane.SetCell(row, col,
				tview.NewTableCell(cellContent).
					SetAlign(tview.AlignLeft).
					SetExpansion(1)) // Make each cell expand to use available space
		}
	}
}

// updateTableHeaders updates the table headers with sort indicators
func (t *TviewRenderer) updateTableHeaders() {
	if t.homeTable == nil {
		return
	}

	t.viewStateMu.RLock()
	sortColIndex := t.viewState.sortColumnIndex
	sortAsc := t.viewState.sortAscending
	t.viewStateMu.RUnlock()

	// Update headers with sort indicators
	for col, column := range t.columns {
		headerText := column.Name

		// Add sort indicator if this is the active sort column
		if col == sortColIndex {
			if sortAsc {
				headerText += " ↑"
			} else {
				headerText += " ↓"
			}
		}

		t.homeTable.SetCell(0, col, tview.NewTableCell(headerText).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(column.Align).
			SetExpansion(column.Expansion).
			SetAttributes(tcell.AttrBold))
	}
}

// updateTable refreshes the home page table with current blocks
func (t *TviewRenderer) updateTable() {
	if t.homeTable == nil {
		return
	}

	t.blocksMu.RLock()
	blocks := make([]rpctypes.PolyBlock, len(t.blocks))
	copy(blocks, t.blocks) // Copy for thread safety
	t.blocksMu.RUnlock()

	// Clear existing rows (except header)
	rowCount := t.homeTable.GetRowCount()
	for row := 1; row < rowCount; row++ {
		for col := 0; col < 10; col++ { // Updated to 10 columns
			t.homeTable.SetCell(row, col, nil)
		}
	}

	// Add blocks to table (newest first)
	for i, block := range blocks {
		if i >= maxBlocks { // Limit blocks for performance
			break
		}

		row := i + 1 // +1 to account for header row

		// Column 0: Block number
		blockNum := block.Number().String()
		t.homeTable.SetCell(row, 0, tview.NewTableCell(blockNum).SetAlign(t.columns[0].Align))

		// Column 1: Time (absolute and relative)
		timeStr := formatBlockTime(block.Time())
		t.homeTable.SetCell(row, 1, tview.NewTableCell(timeStr).SetAlign(t.columns[1].Align))

		// Column 2: Block interval
		intervalStr := t.calculateBlockInterval(block, i, blocks)
		t.homeTable.SetCell(row, 2, tview.NewTableCell(intervalStr).SetAlign(t.columns[2].Align))

		// Column 3: Block hash (truncated for display)
		hashStr := truncateHash(block.Hash().Hex(), 10, 10)
		t.homeTable.SetCell(row, 3, tview.NewTableCell(hashStr).SetAlign(t.columns[3].Align))

		// Column 4: Number of transactions
		txCount := len(block.Transactions())
		t.homeTable.SetCell(row, 4, tview.NewTableCell(strconv.Itoa(txCount)).SetAlign(t.columns[4].Align))

		// Column 5: Block size
		sizeStr := formatBytes(block.Size())
		t.homeTable.SetCell(row, 5, tview.NewTableCell(sizeStr).SetAlign(t.columns[5].Align))

		// Column 6: Base fee
		baseFeeStr := formatBaseFee(block.BaseFee())
		t.homeTable.SetCell(row, 6, tview.NewTableCell(baseFeeStr).SetAlign(t.columns[6].Align))

		// Column 7: Gas used
		gasUsedStr := formatNumber(block.GasUsed())
		t.homeTable.SetCell(row, 7, tview.NewTableCell(gasUsedStr).SetAlign(t.columns[7].Align))

		// Column 8: Gas percentage
		gasPercentStr := formatGasPercentage(block.GasUsed(), block.GasLimit())
		t.homeTable.SetCell(row, 8, tview.NewTableCell(gasPercentStr).SetAlign(t.columns[8].Align))

		// Column 9: Gas limit
		gasLimitStr := formatNumber(block.GasLimit())
		t.homeTable.SetCell(row, 9, tview.NewTableCell(gasLimitStr).SetAlign(t.columns[9].Align))

		// Column 10: State root (truncated)
		stateRootStr := truncateHash(block.Root().Hex(), 8, 8)
		t.homeTable.SetCell(row, 10, tview.NewTableCell(stateRootStr).SetAlign(t.columns[10].Align))
	}

	// Update table title with current block count
	title := fmt.Sprintf(" Blocks (%d) ", len(blocks))
	t.homeTable.SetTitle(title)

	// Update headers with current sort indicators
	t.updateTableHeaders()
}

// updateChainInfo periodically updates the status section with chain information
func (t *TviewRenderer) updateChainInfo(ctx context.Context) {
	// Update immediately on start
	t.refreshChainInfo(ctx)

	// Set up ticker for periodic updates
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.refreshChainInfo(ctx)
		}
	}
}

// refreshChainInfo fetches and displays current chain information
func (t *TviewRenderer) refreshChainInfo(ctx context.Context) {
	if t.homeStatusPane == nil {
		return
	}

	// Get chain store from indexer (need to access the store directly)
	// For now, we'll call the indexer methods which delegate to the store

	var statusLines []string

	// Add current time as first status line
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	statusLines = append(statusLines, fmt.Sprintf("Time: %s", currentTime))

	// Add RPC URL as second status line
	rpcURL := t.indexer.GetRPCURL()
	statusLines = append(statusLines, fmt.Sprintf("RPC: %s", rpcURL))

	// Try to get chain ID and network name
	if chainID, err := t.indexer.GetChainID(ctx); err == nil {
		networkName := getNetworkName(chainID)
		statusLines = append(statusLines, fmt.Sprintf("Network: %s", networkName))
		statusLines = append(statusLines, fmt.Sprintf("Chain ID: %s", chainID.String()))
	} else {
		statusLines = append(statusLines, "Network: Unknown")
		statusLines = append(statusLines, "Chain ID: N/A")
	}

	// Try to get client version
	if clientVersion, err := t.indexer.GetClientVersion(ctx); err == nil {
		statusLines = append(statusLines, fmt.Sprintf("Client: %s", clientVersion))
	} else {
		statusLines = append(statusLines, "Client: N/A")
	}

	// Try to get sync status
	if syncStatus, err := t.indexer.GetSyncStatus(ctx); err == nil {
		syncStr := formatSyncStatus(syncStatus)
		statusLines = append(statusLines, fmt.Sprintf("Sync: %s", syncStr))
	} else {
		statusLines = append(statusLines, "Sync: N/A")
	}

	// Get cached connection status
	t.networkInfoMu.RLock()
	connectionStatusStr := t.connectionStatus
	t.networkInfoMu.RUnlock()

	if connectionStatusStr != "" {
		statusLines = append(statusLines, fmt.Sprintf("Connection: %s", connectionStatusStr))
	} else {
		statusLines = append(statusLines, "Connection: Measuring...")
	}

	// Format status text
	statusText := ""
	for i, line := range statusLines {
		if i > 0 {
			statusText += "\n"
		}
		statusText += line
	}

	// Update the status section
	// Direct Draw() call is safe from any goroutine according to tview docs
	t.homeStatusPane.SetText(statusText)
	t.throttledDraw()

	log.Debug().Msg("Updated chain info in status section")
}

// formatSyncStatus converts sync status response to human-readable string
func formatSyncStatus(syncStatus interface{}) string {
	switch v := syncStatus.(type) {
	case bool:
		if v {
			return "[SYNC] Syncing"
		}
		return "[OK] Synced"
	case map[string]interface{}:
		// Parse sync progress from object response
		current, ok1 := v["currentBlock"].(string)
		highest, ok2 := v["highestBlock"].(string)

		if ok1 && ok2 {
			// Parse hex strings to get numeric values
			if currentBig, err1 := hexToBigInt(current); err1 == nil {
				if highestBig, err2 := hexToBigInt(highest); err2 == nil {
					if highestBig.Cmp(big.NewInt(0)) > 0 {
						// Calculate percentage
						progress := new(big.Float).SetInt(currentBig)
						total := new(big.Float).SetInt(highestBig)
						percentage := new(big.Float).Quo(progress, total)
						percentage.Mul(percentage, big.NewFloat(100))

						pct, _ := percentage.Float64()
						return fmt.Sprintf("[SYNC] Syncing (%.1f%%)", pct)
					}
				}
			}
		}
		return "[SYNC] Syncing"
	default:
		return "[?] Unknown"
	}
}

// calculateBlockInterval calculates the time interval between a block and its parent
func (t *TviewRenderer) calculateBlockInterval(block rpctypes.PolyBlock, index int, blocks []rpctypes.PolyBlock) string {
	// First try to look up parent block by hash (most accurate)
	parentHash := block.ParentHash().Hex()
	if parentBlock, exists := t.blocksByHash[parentHash]; exists {
		// Calculate interval in seconds
		blockTime := block.Time()
		parentTime := parentBlock.Time()

		// Check if blockTime is after parentTime (normal case)
		if blockTime >= parentTime {
			interval := blockTime - parentTime
			return fmt.Sprintf("%ds", interval)
		}

		// Handle edge case where parent time is after block time
		log.Warn().
			Uint64("block_time", blockTime).
			Uint64("parent_time", parentTime).
			Msg("Parent block time is after child block time")
		return "N/A"
	}

	// Parent not found by hash, use the next block in the slice
	// Since blocks are sorted newest first, the next block (index+1) is the previous block in time
	if index+1 < len(blocks) {
		prevBlock := blocks[index+1]
		// Make sure it's actually the previous block number or close to it
		currentBlockNum := block.Number().Uint64()
		prevBlockNum := prevBlock.Number().Uint64()

		// If the blocks are consecutive or reasonably close (within 100 blocks)
		// we can calculate a meaningful interval
		if currentBlockNum > prevBlockNum && currentBlockNum-prevBlockNum <= 100 {
			blockTime := block.Time()
			prevTime := prevBlock.Time()

			// Check if blockTime is after prevTime (normal case)
			if blockTime < prevTime {
				log.Warn().
					Uint64("block_time", blockTime).
					Uint64("prev_time", prevTime).
					Msg("Previous block time is after current block time")
				return "N/A"
			}

			// Calculate interval as uint64
			interval := blockTime - prevTime

			// For non-consecutive blocks, show average interval
			if currentBlockNum-prevBlockNum > 1 {
				blockDiff := currentBlockNum - prevBlockNum
				avgInterval := interval / blockDiff
				return fmt.Sprintf("~%ds", avgInterval)
			}
			return fmt.Sprintf("%ds", interval)
		}
	}

	// Can't calculate interval (first block in the list or large gap)
	return "N/A"
}

// calculateTransactionCounters calculates EOA and contract deployment counters from all blocks
func (t *TviewRenderer) calculateTransactionCounters() (uint64, uint64) {
	t.blocksMu.RLock()
	defer t.blocksMu.RUnlock()

	var eoaCount, deployCount uint64

	for _, block := range t.blocks {
		transactions := block.Transactions()
		for _, tx := range transactions {
			toAddr := tx.To()
			dataStr := tx.DataStr()
			hasInputData := len(dataStr) > 2 // More than just "0x"

			if !hasInputData && toAddr != zeroAddress {
				// EOA transaction: no input data and not sent to zero address
				eoaCount++
			} else if hasInputData && toAddr == zeroAddress {
				// Contract deployment: has input data and sent to zero address
				deployCount++
			}
		}
	}

	return eoaCount, deployCount
}

// calculateERC20NFTCounters calculates ERC20 and NFT transaction counters from all blocks
func (t *TviewRenderer) calculateERC20NFTCounters() (uint64, uint64) {
	t.blocksMu.RLock()
	defer t.blocksMu.RUnlock()

	var erc20Count, nftCount uint64

	for _, block := range t.blocks {
		transactions := block.Transactions()
		for _, tx := range transactions {
			dataStr := tx.DataStr()
			// Check if transaction has enough input data for a method selector
			// "0x" + 8 hex chars = 4 bytes
			if len(dataStr) >= 10 {
				// Direct string prefix matching with 0x prefix
				if strings.HasPrefix(dataStr, "0xa9059cbb") {
					// ERC20 transfer(address,uint256)
					erc20Count++
				} else if strings.HasPrefix(dataStr, "0xb88d4fde") || // safeTransferFrom(address,address,uint256,bytes)
					strings.HasPrefix(dataStr, "0x42842e0e") || // safeTransferFrom(address,address,uint256)
					strings.HasPrefix(dataStr, "0xa22cb465") { // setApprovalForAll(address,bool)
					// NFT methods
					nftCount++
				}
			}
		}
	}

	return erc20Count, nftCount
}

// calculateUniqueAddressCounters calculates unique from and to address counters from all blocks
func (t *TviewRenderer) calculateUniqueAddressCounters() (uint64, uint64) {
	t.blocksMu.RLock()
	defer t.blocksMu.RUnlock()

	uniqueFrom := make(map[common.Address]bool)
	uniqueTo := make(map[common.Address]bool)

	for _, block := range t.blocks {
		transactions := block.Transactions()
		for _, tx := range transactions {
			// Track unique from addresses
			fromAddr := tx.From()
			uniqueFrom[fromAddr] = true

			// Track unique to addresses (if not contract creation)
			toAddr := tx.To()
			if toAddr != zeroAddress {
				uniqueTo[toAddr] = true
			}
		}
	}

	return uint64(len(uniqueFrom)), uint64(len(uniqueTo))
}

// fetchBlockInfo retrieves the latest, safe, and finalized block numbers
func (t *TviewRenderer) fetchBlockInfo(ctx context.Context) {
	if t.indexer == nil {
		return
	}

	// Get latest block
	latestBlock, err := t.indexer.GetBlock(ctx, "latest")
	var latestNum *big.Int
	if err == nil {
		latestNum = latestBlock.Number()
	} else {
		log.Debug().Err(err).Msg("Failed to get latest block")
	}

	// Get safe block
	safeNum, err := t.indexer.GetSafeBlock(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get safe block")
		safeNum = nil
	}

	// Get finalized block
	finalizedNum, err := t.indexer.GetFinalizedBlock(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get finalized block")
		finalizedNum = nil
	}

	// Update stored values with mutex protection
	t.blockInfoMu.Lock()
	t.latestBlockNum = latestNum
	t.safeBlockNum = safeNum
	t.finalizedBlockNum = finalizedNum
	t.blockInfoMu.Unlock()

	log.Debug().
		Str("latest", formatBlockNumber(latestNum)).
		Str("safe", formatBlockNumber(safeNum)).
		Str("finalized", formatBlockNumber(finalizedNum)).
		Msg("Updated block info")
}

// updateBlockInfo periodically updates block information for the metrics display
func (t *TviewRenderer) updateBlockInfo(ctx context.Context) {
	// Update immediately on start
	t.fetchBlockInfo(ctx)

	// Set up ticker for periodic updates
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.fetchBlockInfo(ctx)
			// Trigger a metrics pane update to reflect new block info
			// Direct Draw() call is safe from any goroutine according to tview docs
			if t.homeMetricsPane != nil {
				// Force a metrics update by creating a dummy metrics update
				// This will cause updateMetricsPane to be called with fresh block data
				t.updateMetricsPane(metrics.MetricUpdate{
					Name:  "blockInfo",
					Value: "update",
					Time:  time.Now(),
				})
			}
			t.throttledDraw()
		}
	}
}

// insertBlockSorted inserts a block in the correct position to maintain current sort order
func (t *TviewRenderer) insertBlockSorted(block rpctypes.PolyBlock) {
	blockNum := block.Number()
	blockHash := block.Hash().Hex()

	// Get current sort settings first, outside of locks
	t.viewStateMu.RLock()
	sortColIndex := t.viewState.sortColumnIndex
	sortAsc := t.viewState.sortAscending
	t.viewStateMu.RUnlock()

	// Get the sort column definition
	if sortColIndex < 0 || sortColIndex >= len(t.columns) {
		sortColIndex = 0 // Default to first column
	}
	column := t.columns[sortColIndex]

	// Now acquire blocks lock and insert
	t.blocksMu.Lock()
	defer t.blocksMu.Unlock()

	// Check if block already exists
	if _, exists := t.blocksByHash[blockHash]; exists {
		log.Debug().Str("hash", blockHash).Msg("Block already exists, skipping")
		return
	}

	// Find insertion point using binary search with current sort order
	left, right := 0, len(t.blocks)
	for left < right {
		mid := (left + right) / 2

		// Extract values for comparison
		midVal := column.SortFunc(t.blocks[mid])
		newVal := column.SortFunc(block)

		// Compare using the column's comparison function
		cmp := column.CompareFunc(midVal, newVal)

		// Apply sort direction logic
		if sortAsc {
			// Ascending: if mid < new, search right half
			if cmp < 0 {
				left = mid + 1
			} else {
				right = mid
			}
		} else {
			// Descending: if mid > new, search right half
			if cmp > 0 {
				left = mid + 1
			} else {
				right = mid
			}
		}
	}

	// Insert at the found position
	t.blocks = append(t.blocks, nil)         // Expand slice
	copy(t.blocks[left+1:], t.blocks[left:]) // Shift elements right
	t.blocks[left] = block                   // Insert new block

	// Update hash map
	t.blocksByHash[blockHash] = block

	// Limit blocks to prevent memory issues
	if len(t.blocks) > maxBlocks {
		// Remove oldest blocks (at the end of the array)
		for i := maxBlocks; i < len(t.blocks); i++ {
			delete(t.blocksByHash, t.blocks[i].Hash().Hex())
		}
		t.blocks = t.blocks[:maxBlocks]
	}

	log.Debug().
		Str("hash", blockHash).
		Str("number", blockNum.String()).
		Int("position", left).
		Int("totalBlocks", len(t.blocks)).
		Msg("Block inserted in sorted order")
}

// resortBlocks sorts the existing blocks array using the current sort settings
func (t *TviewRenderer) resortBlocks() {
	// Get sort settings first, outside of any locks
	t.viewStateMu.RLock()
	sortColIndex := t.viewState.sortColumnIndex
	sortAsc := t.viewState.sortAscending
	t.viewStateMu.RUnlock()

	if sortColIndex < 0 || sortColIndex >= len(t.columns) {
		log.Error().Int("index", sortColIndex).Msg("Invalid sort column index")
		return
	}

	column := t.columns[sortColIndex]

	// Now acquire blocks lock and sort
	t.blocksMu.Lock()
	defer t.blocksMu.Unlock()

	sort.Slice(t.blocks, func(i, j int) bool {
		// Extract sort values
		valI := column.SortFunc(t.blocks[i])
		valJ := column.SortFunc(t.blocks[j])

		// Compare using the column's comparison function
		cmp := column.CompareFunc(valI, valJ)

		// Apply sort direction
		if sortAsc {
			return cmp < 0
		} else {
			return cmp > 0
		}
	})

	log.Debug().
		Str("column", column.Key).
		Bool("ascending", sortAsc).
		Int("blocks", len(t.blocks)).
		Msg("Resorted blocks")
}

// changeSortColumn changes the sort column by delta (-1 for left, +1 for right)
func (t *TviewRenderer) changeSortColumn(delta int) {
	t.viewStateMu.Lock()
	defer t.viewStateMu.Unlock()

	newIndex := t.viewState.sortColumnIndex + delta
	if newIndex < 0 {
		newIndex = len(t.columns) - 1 // Wrap to last column
	} else if newIndex >= len(t.columns) {
		newIndex = 0 // Wrap to first column
	}

	t.viewState.sortColumnIndex = newIndex
	t.viewState.sortColumn = t.columns[newIndex].Key

	log.Debug().
		Int("newIndex", newIndex).
		Str("newColumn", t.viewState.sortColumn).
		Msg("Changed sort column")
}

// toggleSortDirection reverses the current sort direction
func (t *TviewRenderer) toggleSortDirection() {
	t.viewStateMu.Lock()
	defer t.viewStateMu.Unlock()

	t.viewState.sortAscending = !t.viewState.sortAscending

	log.Debug().
		Bool("ascending", t.viewState.sortAscending).
		Msg("Toggled sort direction")
}
