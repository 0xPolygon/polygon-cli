package renderer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/indexer"
	"github.com/0xPolygon/polygon-cli/indexer/metrics"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// Comparison functions for different data types
func compareNumbers(a, b interface{}) int {
	aNum := a.(*big.Int)
	bNum := b.(*big.Int)
	return aNum.Cmp(bNum)
}

func compareUint64(a, b interface{}) int {
	aNum := a.(uint64)
	bNum := b.(uint64)
	if aNum < bNum {
		return -1
	} else if aNum > bNum {
		return 1
	}
	return 0
}

func compareStrings(a, b interface{}) int {
	aStr := a.(string)
	bStr := b.(string)
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// createColumnDefinitions creates all sortable column definitions
func createColumnDefinitions() []ColumnDef {
	return []ColumnDef{
		{
			Name: "BLOCK #", Key: "number", Align: tview.AlignRight, Expansion: 1,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Number() },
			CompareFunc: compareNumbers,
		},
		{
			Name: "TIME", Key: "time", Align: tview.AlignLeft, Expansion: 3,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Time() },
			CompareFunc: compareUint64,
		},
		{
			Name: "INTERVAL", Key: "interval", Align: tview.AlignRight, Expansion: 1,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Time() }, // Will be calculated separately
			CompareFunc: compareUint64,
		},
		{
			Name: "HASH", Key: "hash", Align: tview.AlignLeft, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Hash().Hex() },
			CompareFunc: compareStrings,
		},
		{
			Name: "TXS", Key: "txs", Align: tview.AlignRight, Expansion: 1,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return uint64(len(block.Transactions())) },
			CompareFunc: compareUint64,
		},
		{
			Name: "SIZE", Key: "size", Align: tview.AlignRight, Expansion: 1,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Size() },
			CompareFunc: compareUint64,
		},
		{
			Name: "GAS USED", Key: "gasused", Align: tview.AlignRight, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.GasUsed() },
			CompareFunc: compareUint64,
		},
		{
			Name: "GAS %", Key: "gaspct", Align: tview.AlignRight, Expansion: 1,
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
			Name: "STATE ROOT", Key: "stateroot", Align: tview.AlignLeft, Expansion: 2,
			SortFunc:    func(block rpctypes.PolyBlock) interface{} { return block.Root().Hex() },
			CompareFunc: compareStrings,
		},
	}
}

// ColumnDef defines a sortable column with its properties
type ColumnDef struct {
	Name        string                               // Display name
	Key         string                               // Internal identifier
	Align       int                                  // Text alignment
	Expansion   int                                  // Column width allocation
	SortFunc    func(rpctypes.PolyBlock) interface{} // Custom sort extraction
	CompareFunc func(interface{}, interface{}) int   // Custom comparison
}

// ViewState tracks the current view preferences and selection state
type ViewState struct {
	followMode      bool   // Auto-follow newest block vs manual navigation
	sortColumn      string // Current sort column key
	sortColumnIndex int    // Index of current sort column (0-based)
	sortAscending   bool   // Sort direction (true=asc, false=desc)
	selectedBlock   string // Hash of currently selected block (empty = none)
	manualSelect    bool   // User made manual selection (disables auto-follow)
}

// TviewRenderer provides a terminal UI using the tview library
type TviewRenderer struct {
	BaseRenderer
	app    *tview.Application
	pages  *tview.Pages
	blocks []rpctypes.PolyBlock
	// Map to store blocks by hash for parent lookup
	blocksByHash map[string]rpctypes.PolyBlock
	// Column definitions for sorting
	columns []ColumnDef
	// View state management
	viewState ViewState
	// Mutex for thread-safe access to blocks and viewState
	blocksMu    sync.RWMutex
	viewStateMu sync.RWMutex

	// Pages
	homePage         *tview.Flex     // Changed to Flex to hold multiple sections
	homeTopSection   *tview.Flex     // Flex container for 2-column top section
	homeStatusPane   *tview.TextView // Left pane: Status information (1/3 width)
	homeMetricsPane  *tview.Table    // Right pane: Metrics table (2/3 width)
	homeTable        *tview.Table
	blockDetailPage  *tview.Flex     // Changed to Flex for side-by-side layout
	blockDetailLeft  *tview.Table    // Left pane: Transaction table
	blockDetailRight *tview.TextView // Right pane: Raw JSON
	txDetailPage     *tview.Flex     // Transaction detail with human-readable left, stacked JSON right
	txDetailLeft     *tview.TextView // Left pane: Human-readable transaction properties
	txDetailRight    *tview.Flex     // Right pane: Container for stacked JSON views
	txDetailTxJSON   *tview.TextView // Top right: Transaction JSON
	txDetailRcptJSON *tview.TextView // Bottom right: Receipt JSON
	infoPage         *tview.TextView
	helpPage         *tview.TextView

	// Block info for metrics display
	latestBlockNum    *big.Int
	safeBlockNum      *big.Int
	finalizedBlockNum *big.Int
	blockInfoMu       sync.RWMutex

	// Network info for metrics display
	gasPrice      string
	txPoolPending string
	txPoolQueued  string
	peerCount     string
	networkInfoMu sync.RWMutex

	// Current block being viewed in detail (for transaction selection)
	currentBlock   rpctypes.PolyBlock
	currentBlockMu sync.RWMutex

	// Throttling for UI updates
	lastDrawTime    time.Time
	drawMu          sync.Mutex
	minDrawInterval time.Duration

	// Modals
	quitModal *tview.Modal
}

// NewTviewRenderer creates a new TUI renderer using tview
func NewTviewRenderer(indexer *indexer.Indexer) *TviewRenderer {
	app := tview.NewApplication()

	columns := createColumnDefinitions()

	renderer := &TviewRenderer{
		BaseRenderer: NewBaseRenderer(indexer),
		app:          app,
		blocks:       make([]rpctypes.PolyBlock, 0),
		blocksByHash: make(map[string]rpctypes.PolyBlock),
		columns:      columns,
		viewState: ViewState{
			followMode:      true,     // Start in follow mode
			sortColumn:      "number", // Default sort by block number
			sortColumnIndex: 0,        // Block number is first column
			sortAscending:   false,    // Descending (newest first)
			selectedBlock:   "",       // No selection initially
			manualSelect:    false,    // Auto-follow enabled
		},
		minDrawInterval: 50 * time.Millisecond, // Limit updates to 20 FPS
	}

	// Create all pages
	renderer.createPages()

	// Set up keyboard shortcuts
	renderer.setupKeyboardShortcuts()

	// Set the pages as the root of the application
	app.SetRoot(renderer.pages, true)

	return renderer
}

// createPages initializes all the different pages/views
func (t *TviewRenderer) createPages() {
	// Create pages container
	t.pages = tview.NewPages()

	// Create Home page (main block table)
	t.createHomePage()

	// Create Block Detail page
	t.createBlockDetailPage()

	// Create Transaction Detail page
	t.createTransactionDetailPage()

	// Create Info page
	t.createInfoPage()

	// Create Help page
	t.createHelpPage()

	// Create Quit confirmation modal
	t.createQuitModal()

	// Add all pages to the container
	t.pages.AddPage("home", t.homePage, true, true)
	t.pages.AddPage("block-detail", t.blockDetailPage, true, false)
	t.pages.AddPage("tx-detail", t.txDetailPage, true, false)
	t.pages.AddPage("info", t.infoPage, true, false)
	t.pages.AddPage("help", t.helpPage, true, false)
	t.pages.AddPage("quit", t.quitModal, true, false)
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

	// Initialize with 8 rows of placeholders
	for row := 0; row < 8; row++ {
		t.homeMetricsPane.SetCell(row, 0,
			tview.NewTableCell(fmt.Sprintf("Row %d placeholder ", row+1)). // Add trailing space
											SetAlign(tview.AlignLeft).
											SetExpansion(1))
	}

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
				// Mark as manual selection if this wasn't triggered by auto-follow
				if t.viewState.followMode && row != 1 {
					t.viewState.manualSelect = true
					t.viewState.selectedBlock = blockHash
					log.Debug().Str("hash", blockHash).Msg("Manual selection detected, disabling auto-follow")
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

// createTransactionDetailPage creates the transaction detail view with human-readable left pane and stacked JSON right panes
func (t *TviewRenderer) createTransactionDetailPage() {
	// Create left pane for human-readable transaction properties
	t.txDetailLeft = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	t.txDetailLeft.SetBorder(true).SetTitle(" Transaction Details ")
	t.txDetailLeft.SetText("Transaction details will be displayed here")

	// Create top right pane for transaction JSON
	t.txDetailTxJSON = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	t.txDetailTxJSON.SetBorder(true).SetTitle(" Transaction JSON ")
	t.txDetailTxJSON.SetText("Select a transaction to view its JSON representation")

	// Create bottom right pane for receipt JSON
	t.txDetailRcptJSON = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	t.txDetailRcptJSON.SetBorder(true).SetTitle(" Receipt JSON ")
	t.txDetailRcptJSON.SetText("Select a transaction to view its receipt JSON")

	// Create right flex container to stack the JSON views vertically
	t.txDetailRight = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.txDetailTxJSON, 0, 1, false).   // Top: Transaction JSON (50% height)
		AddItem(t.txDetailRcptJSON, 0, 1, false) // Bottom: Receipt JSON (50% height)

	// Create main flex container to hold left pane and right stack side by side
	t.txDetailPage = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.txDetailLeft, 0, 1, true).  // Left pane: 50% width, focusable
		AddItem(t.txDetailRight, 0, 1, true) // Right stack: 50% width, focusable
}

// createInfoPage creates the application info page
func (t *TviewRenderer) createInfoPage() {
	t.infoPage = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	t.infoPage.SetTitle(" Information ")
	t.infoPage.SetBorder(true)
	t.infoPage.SetText("Application Information\n\nMonitorv2 - Blockchain Monitor\nVersion: 1.0.0\n\nPress 'Esc' to go back to home")
}

// createHelpPage creates the help/shortcuts page
func (t *TviewRenderer) createHelpPage() {
	t.helpPage = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	t.helpPage.SetTitle(" Help ")
	t.helpPage.SetBorder(true)

	helpText := `Keyboard Shortcuts:

q - Quit (with confirmation)
h - Show this help page
i - Show information page
Esc - Go back to home page
Enter - View block details (on home page)

Navigation:
↑↓ - Scroll through blocks
Home/End - Jump to top/bottom

Sorting (on home page):
< - Move sort column left
> - Move sort column right
R - Reverse sort direction

Press 'Esc' to go back to home`

	t.helpPage.SetText(helpText)
}

// createQuitModal creates the quit confirmation dialog
func (t *TviewRenderer) createQuitModal() {
	t.quitModal = tview.NewModal().
		SetText("Are you sure you want to quit?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				t.app.Stop()
			} else {
				t.pages.HidePage("quit")
			}
		})
}

// setupKeyboardShortcuts configures global keyboard shortcuts
func (t *TviewRenderer) setupKeyboardShortcuts() {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Check current page
		currentPage, _ := t.pages.GetFrontPage()

		// Handle quit modal
		if currentPage == "quit" {
			switch event.Rune() {
			case 'y', 'Y':
				t.app.Stop()
				return nil
			case 'n', 'N':
				t.pages.HidePage("quit")
				return nil
			case 'q', 'Q':
				return nil // Ignore additional q presses
			}
			return event // Let modal handle other keys
		}

		// Global shortcuts (work from any page)
		switch event.Rune() {
		case 'q', 'Q':
			t.pages.ShowPage("quit")
			return nil
		case 'h', 'H':
			t.pages.SwitchToPage("help")
			return nil
		case 'i', 'I':
			t.pages.SwitchToPage("info")
			return nil
		}

		// Handle Escape key to go back to home
		if event.Key() == tcell.KeyEscape {
			t.pages.SwitchToPage("home")
			// Set focus to the table when returning to home
			if t.homeTable != nil {
				t.app.SetFocus(t.homeTable)
			}
			return nil
		}

		// Page-specific shortcuts
		switch currentPage {
		case "home":
			switch event.Key() {
			case tcell.KeyEnter:
				// Handle Enter key on home page (block selection)
				if t.homeTable != nil {
					row, _ := t.homeTable.GetSelection()
					if row > 0 && row-1 < len(t.blocks) {
						t.showBlockDetail(t.blocks[row-1])
					}
				}
				return nil
			}

			// Handle sorting shortcuts with immediate feedback
			switch event.Rune() {
			case '<':
				// Move sort column left and redraw immediately
				t.changeSortColumn(-1)
				t.resortBlocks()
				t.updateTable()
				t.updateTableHeaders()
				return nil
			case '>':
				// Move sort column right and redraw immediately
				t.changeSortColumn(1)
				t.resortBlocks()
				t.updateTable()
				t.updateTableHeaders()
				return nil
			case 'r', 'R':
				// Reverse sort direction and redraw immediately
				t.toggleSortDirection()
				t.resortBlocks()
				t.updateTable()
				t.updateTableHeaders()
				return nil
			}
		case "block-detail":
			switch event.Key() {
			case tcell.KeyTab:
				// Switch focus between left and right panes
				focused := t.app.GetFocus()
				if focused == t.blockDetailLeft {
					t.app.SetFocus(t.blockDetailRight)
				} else {
					t.app.SetFocus(t.blockDetailLeft)
				}
				return nil
			case tcell.KeyEnter:
				// Handle Enter on transaction table to show transaction detail
				focused := t.app.GetFocus()
				if focused == t.blockDetailLeft {
					if row, _ := t.blockDetailLeft.GetSelection(); row > 0 {
						txIndex := row - 1 // -1 to account for header row
						// Get the current block and its transactions
						t.currentBlockMu.RLock()
						currentBlock := t.currentBlock
						t.currentBlockMu.RUnlock()

						if currentBlock != nil {
							transactions := currentBlock.Transactions()
							if txIndex < len(transactions) {
								// Navigate to transaction detail page with actual transaction
								t.showTransactionDetail(transactions[txIndex], txIndex)
							}
						}
					}
				}
				return nil
			}
		case "tx-detail":
			switch event.Key() {
			case tcell.KeyTab:
				// Cycle focus between left pane, transaction JSON, and receipt JSON
				focused := t.app.GetFocus()
				if focused == t.txDetailLeft {
					t.app.SetFocus(t.txDetailTxJSON)
				} else if focused == t.txDetailTxJSON {
					t.app.SetFocus(t.txDetailRcptJSON)
				} else {
					t.app.SetFocus(t.txDetailLeft)
				}
				return nil
			}
		}

		return event
	})
}

// showBlockDetail navigates to block detail page and populates it
func (t *TviewRenderer) showBlockDetail(block rpctypes.PolyBlock) {
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
		toAddr := "N/A"
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
	blockJSON, err := block.MarshalJSON()
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

	// Set focus to the left pane (transaction table) by default
	t.app.SetFocus(t.blockDetailLeft)
}

// showTransactionDetail navigates to transaction detail page and populates it asynchronously
func (t *TviewRenderer) showTransactionDetail(tx rpctypes.PolyTransaction, txIndex int) {
	// Update pane titles to reflect the transaction content
	t.txDetailLeft.SetTitle(fmt.Sprintf(" Transaction Details (Index: %d) ", txIndex))
	t.txDetailTxJSON.SetTitle(fmt.Sprintf(" Transaction JSON (Hash: %s) ", truncateHash(tx.Hash().Hex(), 8, 8)))
	t.txDetailRcptJSON.SetTitle(" Receipt JSON ")

	// Set loading states for all panes
	t.txDetailLeft.SetText("Loading transaction details...")
	t.txDetailTxJSON.SetText("Loading transaction JSON...")
	t.txDetailRcptJSON.SetText("Loading receipt JSON...")

	// Switch to transaction detail page immediately
	t.pages.SwitchToPage("tx-detail")
	
	// Set focus to the left pane by default
	t.app.SetFocus(t.txDetailLeft)

	// Start async operations
	go t.loadTransactionJSONAsync(tx)
	go t.loadTransactionDetailsAsync(tx, txIndex)
	go t.loadReceiptJSONAsync(tx)
}

// createBasicTransactionDetails creates basic transaction details without signature lookup
func (t *TviewRenderer) createBasicTransactionDetails(tx rpctypes.PolyTransaction, txIndex int) string {
	var details []string

	// Basic transaction information
	details = append(details, fmt.Sprintf("Transaction Index: %d", txIndex))
	details = append(details, fmt.Sprintf("Hash: %s", tx.Hash().Hex()))
	details = append(details, fmt.Sprintf("Block Number: %s", tx.BlockNumber().String()))
	details = append(details, fmt.Sprintf("Chain ID: %d", tx.ChainID()))
	details = append(details, "")

	// Transaction details
	details = append(details, fmt.Sprintf("From: %s", tx.From().Hex()))
	if tx.To().Hex() != "0x0000000000000000000000000000000000000000" {
		details = append(details, fmt.Sprintf("To: %s", tx.To().Hex()))
	} else {
		details = append(details, "To: [Contract Creation]")
	}
	details = append(details, fmt.Sprintf("Value: %s ETH", weiToEther(tx.Value())))
	details = append(details, fmt.Sprintf("Gas: %s", formatNumber(tx.Gas())))
	details = append(details, fmt.Sprintf("Gas Price: %s gwei", weiToGwei(tx.GasPrice())))
	details = append(details, fmt.Sprintf("Nonce: %d", tx.Nonce()))
	details = append(details, "")

	// Transaction type and data (without signature lookup)
	details = append(details, fmt.Sprintf("Type: %d", tx.Type()))
	details = append(details, fmt.Sprintf("Data Size: %d bytes", len(tx.Data())))
	if len(tx.Data()) > 0 {
		if len(tx.Data()) >= 4 {
			// Display method signature without lookup (loading)
			methodSig := fmt.Sprintf("0x%x", tx.Data()[:4])
			details = append(details, fmt.Sprintf("Method Signature: %s (loading...)", methodSig))
		}
		if len(tx.Data()) <= 32 {
			details = append(details, fmt.Sprintf("Full Data: 0x%x", tx.Data()))
		} else {
			details = append(details, fmt.Sprintf("Data Preview: 0x%x...", tx.Data()[:32]))
		}
	} else {
		details = append(details, "Data: [Empty]")
	}
	details = append(details, "")

	// EIP-1559 fields (if applicable)
	if tx.Type() >= 2 {
		if tx.MaxFeePerGas() > 0 {
			maxFeeBig := big.NewInt(int64(tx.MaxFeePerGas()))
			details = append(details, fmt.Sprintf("Max Fee Per Gas: %s gwei", weiToGwei(maxFeeBig)))
		}
		if tx.MaxPriorityFeePerGas() > 0 {
			maxPriorityBig := big.NewInt(int64(tx.MaxPriorityFeePerGas()))
			details = append(details, fmt.Sprintf("Max Priority Fee Per Gas: %s gwei", weiToGwei(maxPriorityBig)))
		}
		details = append(details, "")
	}

	// Signature details
	details = append(details, "Signature:")
	details = append(details, fmt.Sprintf("  V: %s", tx.V().String()))
	details = append(details, fmt.Sprintf("  R: %s", tx.R().String()))
	details = append(details, fmt.Sprintf("  S: %s", tx.S().String()))

	// Combine all details into a single string
	detailText := ""
	for _, detail := range details {
		detailText += detail + "\n"
	}

	return detailText
}

// createHumanReadableTransactionDetailsSync creates a human-readable view of transaction details with signature lookup
func (t *TviewRenderer) createHumanReadableTransactionDetailsSync(tx rpctypes.PolyTransaction, txIndex int) string {
	var details []string

	// Basic transaction information
	details = append(details, fmt.Sprintf("Transaction Index: %d", txIndex))
	details = append(details, fmt.Sprintf("Hash: %s", tx.Hash().Hex()))
	details = append(details, fmt.Sprintf("Block Number: %s", tx.BlockNumber().String()))
	details = append(details, fmt.Sprintf("Chain ID: %d", tx.ChainID()))
	details = append(details, "")

	// Transaction details
	details = append(details, fmt.Sprintf("From: %s", tx.From().Hex()))
	if tx.To().Hex() != "0x0000000000000000000000000000000000000000" {
		details = append(details, fmt.Sprintf("To: %s", tx.To().Hex()))
	} else {
		details = append(details, "To: [Contract Creation]")
	}
	details = append(details, fmt.Sprintf("Value: %s ETH", weiToEther(tx.Value())))
	details = append(details, fmt.Sprintf("Gas: %s", formatNumber(tx.Gas())))
	details = append(details, fmt.Sprintf("Gas Price: %s gwei", weiToGwei(tx.GasPrice())))
	details = append(details, fmt.Sprintf("Nonce: %d", tx.Nonce()))
	details = append(details, "")

	// Transaction type and data
	details = append(details, fmt.Sprintf("Type: %d", tx.Type()))
	details = append(details, fmt.Sprintf("Data Size: %d bytes", len(tx.Data())))
	if len(tx.Data()) > 0 {
		if len(tx.Data()) >= 4 {
			// Display method signature with human-readable lookup
			methodSig := fmt.Sprintf("0x%x", tx.Data()[:4])
			sigDetails := t.getMethodSignatureDetails(methodSig)
			details = append(details, fmt.Sprintf("Method Signature: %s", sigDetails))
		}
		if len(tx.Data()) <= 32 {
			details = append(details, fmt.Sprintf("Full Data: 0x%x", tx.Data()))
		} else {
			details = append(details, fmt.Sprintf("Data Preview: 0x%x...", tx.Data()[:32]))
		}
	} else {
		details = append(details, "Data: [Empty]")
	}
	details = append(details, "")

	// EIP-1559 fields (if applicable)
	if tx.Type() >= 2 {
		if tx.MaxFeePerGas() > 0 {
			maxFeeBig := big.NewInt(int64(tx.MaxFeePerGas()))
			details = append(details, fmt.Sprintf("Max Fee Per Gas: %s gwei", weiToGwei(maxFeeBig)))
		}
		if tx.MaxPriorityFeePerGas() > 0 {
			maxPriorityBig := big.NewInt(int64(tx.MaxPriorityFeePerGas()))
			details = append(details, fmt.Sprintf("Max Priority Fee Per Gas: %s gwei", weiToGwei(maxPriorityBig)))
		}
		details = append(details, "")
	}

	// Signature details
	details = append(details, "Signature:")
	details = append(details, fmt.Sprintf("  V: %s", tx.V().String()))
	details = append(details, fmt.Sprintf("  R: %s", tx.R().String()))
	details = append(details, fmt.Sprintf("  S: %s", tx.S().String()))

	// Combine all details into a single string
	detailText := ""
	for _, detail := range details {
		detailText += detail + "\n"
	}

	return detailText
}

// weiToEther converts wei to ether with reasonable precision
func weiToEther(wei *big.Int) string {
	if wei == nil {
		return "0"
	}

	// Convert wei to ether (divide by 10^18)
	ether := new(big.Float).SetInt(wei)
	ether = ether.Quo(ether, big.NewFloat(1e18))

	// Format with 6 decimal places
	return fmt.Sprintf("%.6f", ether)
}

// getMethodSignatureDetails fetches and formats method signature information
func (t *TviewRenderer) getMethodSignatureDetails(hexSignature string) string {
	// First check if we have access to the indexer and it has a store
	if t.indexer == nil {
		return hexSignature
	}

	// Try to get signature from 4byte.directory
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	signatures, err := t.indexer.GetSignature(ctx, hexSignature)
	if err != nil {
		// Log error but don't fail - fallback to hex signature
		log.Debug().Err(err).Str("signature", hexSignature).Msg("Failed to lookup method signature")
		return hexSignature
	}

	if len(signatures) == 0 {
		return fmt.Sprintf("%s (unknown)", hexSignature)
	}

	// Use the first signature (most common)
	firstSig := signatures[0]
	if len(signatures) == 1 {
		return fmt.Sprintf("%s (%s)", hexSignature, firstSig.TextSignature)
	} else {
		return fmt.Sprintf("%s (%s +%d more)", hexSignature, firstSig.TextSignature, len(signatures)-1)
	}
}

// loadTransactionJSONAsync loads and formats transaction JSON asynchronously
func (t *TviewRenderer) loadTransactionJSONAsync(tx rpctypes.PolyTransaction) {
	// Marshal transaction JSON
	txJSON, err := tx.MarshalJSON()
	if err != nil {
		t.app.QueueUpdateDraw(func() {
			t.txDetailTxJSON.SetText(fmt.Sprintf("Error marshaling transaction JSON: %v", err))
		})
		return
	}

	// Pretty print the JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, txJSON, "", "  "); err != nil {
		t.app.QueueUpdateDraw(func() {
			t.txDetailTxJSON.SetText(fmt.Sprintf("Error formatting JSON: %v", err))
		})
		return
	}

	// Update UI on the main thread
	t.app.QueueUpdateDraw(func() {
		t.txDetailTxJSON.SetText(prettyJSON.String())
	})
}

// loadTransactionDetailsAsync loads human-readable transaction details asynchronously
func (t *TviewRenderer) loadTransactionDetailsAsync(tx rpctypes.PolyTransaction, txIndex int) {
	// Create basic transaction details without signature lookup first
	basicDetails := t.createBasicTransactionDetails(tx, txIndex)
	
	// Update UI with basic details immediately
	t.app.QueueUpdateDraw(func() {
		t.txDetailLeft.SetText(basicDetails)
	})

	// Now fetch signatures and update with enhanced details
	enhancedDetails := t.createHumanReadableTransactionDetailsSync(tx, txIndex)
	
	// Update UI with enhanced details including signatures
	t.app.QueueUpdateDraw(func() {
		t.txDetailLeft.SetText(enhancedDetails)
	})
}

// loadReceiptJSONAsync loads and formats receipt JSON asynchronously
func (t *TviewRenderer) loadReceiptJSONAsync(tx rpctypes.PolyTransaction) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch receipt
	receipt, err := t.indexer.GetReceipt(ctx, tx.Hash())
	if err != nil {
		t.app.QueueUpdateDraw(func() {
			t.txDetailRcptJSON.SetText(fmt.Sprintf("Error fetching receipt: %v\n\n(Receipt may not be available for pending transactions)", err))
		})
		return
	}

	// Marshal receipt to JSON
	receiptJSON, err := json.Marshal(receipt)
	if err != nil {
		t.app.QueueUpdateDraw(func() {
			t.txDetailRcptJSON.SetText(fmt.Sprintf("Error marshaling receipt JSON: %v", err))
		})
		return
	}

	// Pretty print the receipt JSON
	var prettyReceiptJSON bytes.Buffer
	if err := json.Indent(&prettyReceiptJSON, receiptJSON, "", "  "); err != nil {
		t.app.QueueUpdateDraw(func() {
			t.txDetailRcptJSON.SetText(fmt.Sprintf("Error formatting receipt JSON: %v", err))
		})
		return
	}

	// Update UI on the main thread
	t.app.QueueUpdateDraw(func() {
		t.txDetailRcptJSON.SetText(prettyReceiptJSON.String())
	})
}

// Start begins the TUI rendering
func (t *TviewRenderer) Start(ctx context.Context) error {
	log.Info().Msg("Starting Tview renderer")

	// Start consuming blocks in a separate goroutine
	go t.consumeBlocks(ctx)

	// Start consuming metrics in a separate goroutine
	go t.consumeMetrics(ctx)

	// Start periodic status updates
	go t.updateChainInfo(ctx)

	// Start periodic block info updates
	go t.updateBlockInfo(ctx)

	// Start periodic network info updates
	go t.updateNetworkInfo(ctx)

	// Table selection is handled automatically by view state logic

	// Start the TUI application
	// This will block until the application is stopped
	if err := t.app.Run(); err != nil {
		log.Error().Err(err).Msg("Error running tview application")
		return err
	}

	return nil
}

// throttledDraw performs a Draw() operation with throttling to prevent overwhelming the UI
func (t *TviewRenderer) throttledDraw() {
	t.drawMu.Lock()
	defer t.drawMu.Unlock()

	now := time.Now()
	elapsed := now.Sub(t.lastDrawTime)

	if elapsed < t.minDrawInterval {
		// Too soon since last draw, skip this one
		return
	}

	t.lastDrawTime = now
	t.app.Draw()
}

// consumeBlocks consumes blocks from the indexer and updates the table
func (t *TviewRenderer) consumeBlocks(ctx context.Context) {
	blockChan := t.indexer.BlockChannel()

	for {
		select {
		case <-ctx.Done():
			return
		case block, ok := <-blockChan:
			if !ok {
				log.Info().Msg("Block channel closed, stopping Tview renderer")
				return
			}

			// Insert block in sorted order (always maintains descending order by block number)
			t.insertBlockSorted(block)

			// Update the table and apply view state
			// Direct Draw() call is safe from any goroutine according to tview docs
			t.updateTable()
			t.applyViewState()
			t.throttledDraw()
		}
	}
}

// consumeMetrics consumes metrics updates from the indexer and updates the metrics pane
func (t *TviewRenderer) consumeMetrics(ctx context.Context) {
	metricsChan := t.indexer.MetricsChannel()

	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-metricsChan:
			if !ok {
				log.Info().Msg("Metrics channel closed")
				return
			}

			// Update the metrics pane
			// Direct Draw() call is safe from any goroutine according to tview docs
			t.updateMetricsPane(update)
			t.throttledDraw()
		}
	}
}

// updateNetworkInfo periodically updates network information to reduce RPC calls
func (t *TviewRenderer) updateNetworkInfo(ctx context.Context) {
	// Update immediately on start
	t.fetchNetworkInfo(ctx)

	// Set up ticker for periodic updates (every 5 seconds)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.fetchNetworkInfo(ctx)
		}
	}
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

	log.Debug().
		Str("gasPrice", t.gasPrice).
		Str("txPoolPending", t.txPoolPending).
		Str("txPoolQueued", t.txPoolQueued).
		Str("peerCount", t.peerCount).
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
				base10Str := "base10 " + weiToGwei(stats.BaseFee10) + " gwei"
				base30Str := "base30 " + weiToGwei(stats.BaseFee30) + " gwei"
				cells = [5]string{"GAS ", base10Str, base30Str, gasPriceStr, "[placeholder]"}
			} else {
				cells = [5]string{"GAS ", "base10 N/A", "base30 N/A", gasPriceStr, "[placeholder]"}
			}
		case 3: // POOL
			// Use cached txpool status
			cells = [5]string{"POOL", txPoolPendingStr, txPoolQueuedStr, "[placeholder]", "[placeholder]"}
		case 4: // SIG (1)
			cells = [5]string{"SIG1", "EOA [placeholder]", "ERC20 [placeholder]", "NFT [placeholder]", "[placeholder]"}
		case 5: // SIG (2)
			cells = [5]string{"SIG2", "Contract Deploy [placeholder]", "Uniswap [placeholder]", "Other [placeholder]", "[placeholder]"}
		case 6: // ACC
			cells = [5]string{"ACCO", "Unique From [placeholder]", "Unique To [placeholder]", "[placeholder]", "[placeholder]"}
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
		if i >= 100 { // Limit to 100 blocks for performance
			break
		}

		row := i + 1 // +1 to account for header row

		// Column 0: Block number
		blockNum := block.Number().String()
		t.homeTable.SetCell(row, 0, tview.NewTableCell(blockNum).SetAlign(tview.AlignRight))

		// Column 1: Time (absolute and relative)
		timeStr := formatBlockTime(block.Time())
		t.homeTable.SetCell(row, 1, tview.NewTableCell(timeStr).SetAlign(tview.AlignLeft))

		// Column 2: Block interval
		intervalStr := t.calculateBlockInterval(block)
		t.homeTable.SetCell(row, 2, tview.NewTableCell(intervalStr).SetAlign(tview.AlignRight))

		// Column 3: Block hash (truncated for display)
		hashStr := truncateHash(block.Hash().Hex(), 10, 10)
		t.homeTable.SetCell(row, 3, tview.NewTableCell(hashStr).SetAlign(tview.AlignLeft))

		// Column 4: Number of transactions
		txCount := len(block.Transactions())
		t.homeTable.SetCell(row, 4, tview.NewTableCell(strconv.Itoa(txCount)).SetAlign(tview.AlignRight))

		// Column 5: Block size
		sizeStr := formatBytes(block.Size())
		t.homeTable.SetCell(row, 5, tview.NewTableCell(sizeStr).SetAlign(tview.AlignRight))

		// Column 6: Gas used
		gasUsedStr := formatNumber(block.GasUsed())
		t.homeTable.SetCell(row, 6, tview.NewTableCell(gasUsedStr).SetAlign(tview.AlignRight))

		// Column 7: Gas percentage
		gasPercentStr := formatGasPercentage(block.GasUsed(), block.GasLimit())
		t.homeTable.SetCell(row, 7, tview.NewTableCell(gasPercentStr).SetAlign(tview.AlignRight))

		// Column 8: Gas limit
		gasLimitStr := formatNumber(block.GasLimit())
		t.homeTable.SetCell(row, 8, tview.NewTableCell(gasLimitStr).SetAlign(tview.AlignRight))

		// Column 9: State root (truncated)
		stateRootStr := truncateHash(block.Root().Hex(), 8, 8)
		t.homeTable.SetCell(row, 9, tview.NewTableCell(stateRootStr).SetAlign(tview.AlignLeft))
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

	// Try to get chain ID
	if chainID, err := t.indexer.GetChainID(ctx); err == nil {
		statusLines = append(statusLines, fmt.Sprintf("Chain ID: %s", chainID.String()))
	} else {
		statusLines = append(statusLines, "Chain ID: N/A")
	}

	// Try to get gas price
	if gasPrice, err := t.indexer.GetGasPrice(ctx); err == nil {
		gasPriceGwei := weiToGwei(gasPrice)
		statusLines = append(statusLines, fmt.Sprintf("Gas Price: %s gwei", gasPriceGwei))
	} else {
		statusLines = append(statusLines, "Gas Price: N/A")
	}

	// Try to get pending transaction count
	if pendingTxs, err := t.indexer.GetPendingTransactionCount(ctx); err == nil {
		statusLines = append(statusLines, fmt.Sprintf("Pending TXs: %s", pendingTxs.String()))
	} else {
		statusLines = append(statusLines, "Pending TXs: N/A")
	}

	// Try to get safe block
	if safeBlock, err := t.indexer.GetSafeBlock(ctx); err == nil {
		statusLines = append(statusLines, fmt.Sprintf("Safe Block: #%s", safeBlock.String()))
	} else {
		statusLines = append(statusLines, "Safe Block: N/A")
	}

	// Get latest block info (use store method)
	if latestBlock, err := t.indexer.GetBlock(ctx, "latest"); err == nil {
		statusLines = append(statusLines, fmt.Sprintf("Latest Block: #%s", latestBlock.Number().String()))

		// Add base fee if available
		if baseFee := latestBlock.BaseFee(); baseFee != nil {
			baseFeeGwei := weiToGwei(baseFee)
			statusLines = append(statusLines, fmt.Sprintf("Base Fee: %s gwei", baseFeeGwei))
		}
	} else {
		statusLines = append(statusLines, "Latest Block: N/A")
	}

	// Format status text
	statusText := ""
	for i, line := range statusLines {
		if i > 0 {
			if i%2 == 0 {
				statusText += "\n"
			} else {
				statusText += " | "
			}
		}
		statusText += line
	}

	// Update the status section
	// Direct Draw() call is safe from any goroutine according to tview docs
	t.homeStatusPane.SetText(statusText)
	t.throttledDraw()

	log.Debug().Msg("Updated chain info in status section")
}

// weiToGwei converts wei to gwei with reasonable precision
func weiToGwei(wei *big.Int) string {
	if wei == nil {
		return "0"
	}

	// Convert wei to gwei (divide by 10^9)
	gwei := new(big.Float).SetInt(wei)
	gwei = gwei.Quo(gwei, big.NewFloat(1e9))

	// Format with 2 decimal places
	return fmt.Sprintf("%.2f", gwei)
}

// formatRelativeTime converts Unix timestamp to human-readable relative time
func formatRelativeTime(timestamp uint64) string {
	now := time.Now().Unix()
	diff := now - int64(timestamp)

	if diff < 0 {
		return "future"
	} else if diff < 60 {
		return fmt.Sprintf("%ds ago", diff)
	} else if diff < 3600 {
		return fmt.Sprintf("%dm ago", diff/60)
	} else if diff < 86400 {
		return fmt.Sprintf("%dh ago", diff/3600)
	} else {
		return fmt.Sprintf("%dd ago", diff/86400)
	}
}

// formatBlockTime formats block timestamp as "2006-01-02T15:04:05Z - 6m ago"
func formatBlockTime(timestamp uint64) string {
	t := time.Unix(int64(timestamp), 0).UTC()
	absolute := t.Format("2006-01-02T15:04:05Z")
	relative := formatRelativeTime(timestamp)
	return fmt.Sprintf("%s - %s", absolute, relative)
}

// calculateBlockInterval calculates the time interval between a block and its parent
func (t *TviewRenderer) calculateBlockInterval(block rpctypes.PolyBlock) string {
	// Look up parent block by hash
	parentHash := block.ParentHash().Hex()
	if parentBlock, exists := t.blocksByHash[parentHash]; exists {
		// Calculate interval in seconds
		interval := int64(block.Time()) - int64(parentBlock.Time())
		return fmt.Sprintf("%ds", interval)
	}
	// Parent not found (might be initial blocks or missing data)
	return "N/A"
}

// formatBytes converts bytes to human-readable size
func formatBytes(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
	}
}

// formatNumber adds thousand separators to large numbers
func formatNumber(num uint64) string {
	str := fmt.Sprintf("%d", num)
	if len(str) <= 3 {
		return str
	}

	// Add commas every 3 digits from right
	var result []rune
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, char)
	}
	return string(result)
}

// formatGasPercentage calculates and formats gas usage percentage
func formatGasPercentage(gasUsed, gasLimit uint64) string {
	if gasLimit == 0 {
		return "0.0%"
	}
	percentage := float64(gasUsed) / float64(gasLimit) * 100
	return fmt.Sprintf("%.1f%%", percentage)
}

// truncateHash shortens a hash for display
func truncateHash(hash string, prefixLen, suffixLen int) string {
	if len(hash) <= prefixLen+suffixLen+3 {
		return hash
	}
	return hash[:prefixLen] + "..." + hash[len(hash)-suffixLen:]
}

// formatDuration formats a duration for human-readable display
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d.Nanoseconds())/1000000)
	} else if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
}

// formatBlockNumber formats a block number for display with thousand separators
func formatBlockNumber(blockNum *big.Int) string {
	if blockNum == nil {
		return "N/A"
	}

	// Format the number with thousand separators
	str := blockNum.String()
	if len(str) <= 3 {
		return "#" + str
	}

	// Add commas every 3 digits from right
	var result []rune
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, char)
	}
	return "#" + string(result)
}

// formatThroughput formats throughput values for display with appropriate units
func formatThroughput(value float64, unit string) string {
	if value == 0 {
		if unit == "" {
			return "0"
		}
		return "0 " + unit
	}

	var formatted string
	if value >= 1000000000 {
		formatted = fmt.Sprintf("%.1fG", value/1000000000)
	} else if value >= 1000000 {
		formatted = fmt.Sprintf("%.1fM", value/1000000)
	} else if value >= 1000 {
		formatted = fmt.Sprintf("%.1fK", value/1000)
	} else {
		formatted = fmt.Sprintf("%.1f", value)
	}

	if unit == "" {
		return formatted
	}
	return formatted + " " + unit
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
	if len(t.blocks) > 1000 {
		// Remove oldest blocks (at the end of the array)
		for i := 1000; i < len(t.blocks); i++ {
			delete(t.blocksByHash, t.blocks[i].Hash().Hex())
		}
		t.blocks = t.blocks[:1000]
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

// getCurrentSortColumn returns the current sort column definition
func (t *TviewRenderer) getCurrentSortColumn() ColumnDef {
	t.viewStateMu.RLock()
	defer t.viewStateMu.RUnlock()

	if t.viewState.sortColumnIndex < 0 || t.viewState.sortColumnIndex >= len(t.columns) {
		return t.columns[0] // Default to first column
	}
	return t.columns[t.viewState.sortColumnIndex]
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

// applyViewState applies the current view state to the table selection
func (t *TviewRenderer) applyViewState() {
	if t.homeTable == nil {
		return
	}

	// Get current view state safely
	t.viewStateMu.RLock()
	followMode := t.viewState.followMode
	manualSelect := t.viewState.manualSelect
	selectedBlock := t.viewState.selectedBlock
	t.viewStateMu.RUnlock()

	// Get blocks data safely
	t.blocksMu.RLock()
	hasBlocks := len(t.blocks) > 0
	var newestBlockHash string
	if hasBlocks {
		newestBlockHash = t.blocks[0].Hash().Hex()
	}

	// Find selected block index without nested locks
	selectedIndex := -1
	if selectedBlock != "" {
		for i, block := range t.blocks {
			if block.Hash().Hex() == selectedBlock {
				selectedIndex = i
				break
			}
		}
	}
	t.blocksMu.RUnlock()

	// Apply view state logic
	if followMode && !manualSelect {
		// Auto-follow mode: always select newest block (index 0, table row 1)
		if hasBlocks {
			t.homeTable.Select(1, 0)
			t.app.SetFocus(t.homeTable)
			log.Debug().Msg("Auto-follow: selected newest block")
		}
	} else if selectedBlock != "" {
		// Manual selection: find the selected block and maintain selection
		if selectedIndex >= 0 {
			t.homeTable.Select(selectedIndex+1, 0) // +1 for header row
			log.Debug().Int("index", selectedIndex).Str("hash", selectedBlock).Msg("Maintained manual selection")
		} else {
			// Selected block no longer exists, fall back to newest
			if hasBlocks {
				t.homeTable.Select(1, 0)
				// Update view state safely
				t.viewStateMu.Lock()
				t.viewState.selectedBlock = newestBlockHash
				t.viewStateMu.Unlock()
				log.Debug().Msg("Selected block not found, fallback to newest")
			}
		}
	}
}

// hexToDecimal converts various hex number formats to big.Int
func hexToDecimal(value interface{}) (*big.Int, error) {
	switch v := value.(type) {
	case string:
		// Handle hex string format
		if len(v) >= 2 && v[:2] == "0x" {
			v = v[2:]
		}
		result := big.NewInt(0)
		result, ok := result.SetString(v, 16)
		if !ok {
			return nil, fmt.Errorf("invalid hex string: %v", value)
		}
		return result, nil
	case float64:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	case int:
		return big.NewInt(int64(v)), nil
	default:
		return nil, fmt.Errorf("unsupported number type: %T", value)
	}
}

// Stop gracefully stops the TUI renderer
func (t *TviewRenderer) Stop() error {
	log.Info().Msg("Stopping Tview renderer")
	t.app.Stop()
	return nil
}
