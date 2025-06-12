package renderer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/chainstore"
	"github.com/0xPolygon/polygon-cli/indexer"
	"github.com/0xPolygon/polygon-cli/indexer/metrics"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// Zero address constant for efficient comparisons
var zeroAddress = common.Address{}

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
	gasPrice         string
	txPoolPending    string
	txPoolQueued     string
	peerCount        string
	connectionStatus string
	networkInfoMu    sync.RWMutex

	// Current block being viewed in detail (for transaction selection)
	currentBlock   rpctypes.PolyBlock
	currentBlockMu sync.RWMutex

	// Throttling for UI updates
	lastDrawTime    time.Time
	drawMu          sync.Mutex
	minDrawInterval time.Duration

	// Transaction counters for metrics - removed unused fields

	// Modals
	quitModal  *tview.Modal
	searchForm *tview.Form

	// Modal state management
	isModalActive   bool
	activeModalName string
	modalStateMu    sync.RWMutex
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

	// Create Search modal
	t.createSearchModal()

	// Add all pages to the container
	t.pages.AddPage("home", t.homePage, true, true)
	t.pages.AddPage("block-detail", t.blockDetailPage, true, false)
	t.pages.AddPage("tx-detail", t.txDetailPage, true, false)
	t.pages.AddPage("info", t.infoPage, true, false)
	t.pages.AddPage("help", t.helpPage, true, false)
	t.pages.AddPage("quit", t.quitModal, true, false)
	t.pages.AddPage("search", t.searchForm, true, false)
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
		AddItem(t.txDetailTxJSON, 0, 1, false).  // Top: Transaction JSON (50% height)
		AddItem(t.txDetailRcptJSON, 0, 1, false) // Bottom: Receipt JSON (50% height)

	// Create main flex container to hold left pane and right stack side by side
	t.txDetailPage = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.txDetailLeft, 0, 1, true). // Left pane: 50% width, focusable
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
/ or s - Open search modal
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
				t.hideModal()
			}
		})
}

// createSearchModal creates the search dialog using Form
func (t *TviewRenderer) createSearchModal() {
	t.searchForm = tview.NewForm().
		AddInputField("Search", "", 50, nil, nil).
		AddButton("Search", func() {
			// Get the search query from the input field
			query := t.searchForm.GetFormItem(0).(*tview.InputField).GetText()
			if strings.TrimSpace(query) != "" {
				t.performSearch(query)
			}
			t.hideModal()
		}).
		AddButton("Cancel", func() {
			t.hideModal()
		})

	// Set border and title
	t.searchForm.SetBorder(true).SetTitle(" Search ")

	// Handle Escape key to close modal
	t.searchForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			t.hideModal()
			return nil
		}
		return event
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

		// Handle search modal
		if currentPage == "search" {
			// Let the modal handle all input for now
			return event
		}

		// Global shortcuts (work from any page)
		switch event.Rune() {
		case 'q', 'Q':
			t.showModal("quit")
			return nil
		case 'h', 'H':
			t.pages.SwitchToPage("help")
			return nil
		case 'i', 'I':
			t.pages.SwitchToPage("info")
			return nil
		case '/', 's', 'S':
			t.showModal("search")
			return nil
		}

		// Handle Escape key for breadcrumb-style navigation
		if event.Key() == tcell.KeyEscape {
			currentPage, _ = t.pages.GetFrontPage()
			if currentPage == "tx-detail" {
				// From transaction detail, go back to block detail
				t.pages.SwitchToPage("block-detail")
				if t.blockDetailLeft != nil {
					t.app.SetFocus(t.blockDetailLeft)
				}
			} else {
				// From all other pages, go back to home
				t.pages.SwitchToPage("home")
				if t.homeTable != nil {
					t.app.SetFocus(t.homeTable)
				}
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

	// Set focus to the left pane (transaction table) by default
	t.app.SetFocus(t.blockDetailLeft)
}

// showTransactionDetail navigates to transaction detail page and populates it asynchronously
func (t *TviewRenderer) showTransactionDetail(tx rpctypes.PolyTransaction, txIndex int) {
	log.Debug().
		Str("txHash", tx.Hash().Hex()).
		Int("txIndex", txIndex).
		Msg("showTransactionDetail called")

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
			maxFeeGas := tx.MaxFeePerGas()
			var maxFeeBig *big.Int
			if maxFeeGas > math.MaxInt64 {
				log.Error().Uint64("max_fee_per_gas", maxFeeGas).Msg("MaxFeePerGas exceeds int64 range, using MaxInt64")
				maxFeeBig = big.NewInt(math.MaxInt64)
			} else {
				maxFeeBig = big.NewInt(int64(maxFeeGas))
			}
			details = append(details, fmt.Sprintf("Max Fee Per Gas: %s gwei", weiToGwei(maxFeeBig)))
		}
		if tx.MaxPriorityFeePerGas() > 0 {
			maxPriorityGas := tx.MaxPriorityFeePerGas()
			var maxPriorityBig *big.Int
			if maxPriorityGas > math.MaxInt64 {
				log.Error().Uint64("max_priority_fee_per_gas", maxPriorityGas).Msg("MaxPriorityFeePerGas exceeds int64 range, using MaxInt64")
				maxPriorityBig = big.NewInt(math.MaxInt64)
			} else {
				maxPriorityBig = big.NewInt(int64(maxPriorityGas))
			}
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
			maxFeeGas := tx.MaxFeePerGas()
			var maxFeeBig *big.Int
			if maxFeeGas > math.MaxInt64 {
				log.Error().Uint64("max_fee_per_gas", maxFeeGas).Msg("MaxFeePerGas exceeds int64 range, using MaxInt64")
				maxFeeBig = big.NewInt(math.MaxInt64)
			} else {
				maxFeeBig = big.NewInt(int64(maxFeeGas))
			}
			details = append(details, fmt.Sprintf("Max Fee Per Gas: %s gwei", weiToGwei(maxFeeBig)))
		}
		if tx.MaxPriorityFeePerGas() > 0 {
			maxPriorityGas := tx.MaxPriorityFeePerGas()
			var maxPriorityBig *big.Int
			if maxPriorityGas > math.MaxInt64 {
				log.Error().Uint64("max_priority_fee_per_gas", maxPriorityGas).Msg("MaxPriorityFeePerGas exceeds int64 range, using MaxInt64")
				maxPriorityBig = big.NewInt(math.MaxInt64)
			} else {
				maxPriorityBig = big.NewInt(int64(maxPriorityGas))
			}
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

// extractEventSignatures extracts unique event signature hashes from receipt logs
func (t *TviewRenderer) extractEventSignatures(receipt rpctypes.PolyReceipt) []string {
	if receipt == nil {
		return nil
	}

	logs := receipt.Logs()
	if len(logs) == 0 {
		return nil
	}

	// Use a map to collect unique event signatures
	uniqueSigs := make(map[string]bool)

	for _, logEntry := range logs {
		// Check if the log has topics and the first topic exists (event signature)
		if len(logEntry.Topics) > 0 {
			// Get the event signature hash from the first topic
			eventSigHash := logEntry.Topics[0].ToHash().Hex()
			uniqueSigs[eventSigHash] = true
		}
	}

	// Convert map keys to slice
	var signatures []string
	for sig := range uniqueSigs {
		signatures = append(signatures, sig)
	}

	return signatures
}

// findBestSignature returns the signature with the minimum ID (earliest submission)
func findBestSignature(signatures []chainstore.Signature) chainstore.Signature {
	if len(signatures) == 0 {
		return chainstore.Signature{}
	}

	if len(signatures) == 1 {
		return signatures[0]
	}

	// Find signature with minimum ID (earliest submission, more likely to be correct)
	bestSig := signatures[0]
	for _, sig := range signatures[1:] {
		if sig.ID < bestSig.ID {
			bestSig = sig
		}
	}

	return bestSig
}

// getEventSignatureDetails fetches and formats event signature information
func (t *TviewRenderer) getEventSignatureDetails(eventSignatures []string) map[string]string {
	if t.indexer == nil || len(eventSignatures) == 0 {
		return nil
	}

	eventDetails := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Look up each unique event signature
	for _, eventSig := range eventSignatures {
		signatures, err := t.indexer.GetSignature(ctx, eventSig)
		if err != nil {
			log.Debug().Err(err).Str("signature", eventSig).Msg("Failed to lookup event signature")
			eventDetails[eventSig] = fmt.Sprintf("%s (unknown)", eventSig[:10]+"...")
			continue
		}

		if len(signatures) == 0 {
			eventDetails[eventSig] = fmt.Sprintf("%s (unknown)", eventSig[:10]+"...")
		} else {
			// Use the signature with minimum ID (earliest submission, most likely correct)
			bestSig := findBestSignature(signatures)
			if len(signatures) == 1 {
				eventDetails[eventSig] = bestSig.TextSignature
			} else {
				eventDetails[eventSig] = fmt.Sprintf("%s (+%d more)", bestSig.TextSignature, len(signatures)-1)
			}
		}
	}

	return eventDetails
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

	// Use the signature with minimum ID (earliest submission, most likely correct)
	bestSig := findBestSignature(signatures)
	if len(signatures) == 1 {
		return fmt.Sprintf("%s (%s)", hexSignature, bestSig.TextSignature)
	} else {
		return fmt.Sprintf("%s (%s +%d more)", hexSignature, bestSig.TextSignature, len(signatures)-1)
	}
}

// loadTransactionJSONAsync loads and formats transaction JSON asynchronously
func (t *TviewRenderer) loadTransactionJSONAsync(tx rpctypes.PolyTransaction) {
	// Marshal transaction JSON
	txJSON, err := rpctypes.PolyTransactionToPrettyJSON(tx)
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

	// Start coordinated async loading for method signatures and event logs
	go t.loadEnhancedTransactionDetailsAsync(tx, txIndex)
}

// loadEnhancedTransactionDetailsAsync coordinates method signature and event log loading
func (t *TviewRenderer) loadEnhancedTransactionDetailsAsync(tx rpctypes.PolyTransaction, txIndex int) {
	// Channels to receive results
	methodSigChan := make(chan string, 1)
	eventLogsChan := make(chan string, 1)

	// Start method signature lookup
	go func() {
		enhancedDetails := t.createHumanReadableTransactionDetailsSync(tx, txIndex)
		methodSigChan <- enhancedDetails
	}()

	// Start event logs lookup
	go func() {
		// Create context with timeout for receipt fetching
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Fetch receipt
		receipt, err := t.indexer.GetReceipt(ctx, tx.Hash())
		if err != nil {
			eventLogsChan <- "" // No event logs available
			return
		}

		// Extract and look up event signatures from logs
		eventSignatures := t.extractEventSignatures(receipt)
		if len(eventSignatures) == 0 {
			eventLogsChan <- "" // No events to process
			return
		}

		// Look up event signature details
		eventDetails := t.getEventSignatureDetails(eventSignatures)

		// Build event logs section for display
		eventLogsText := t.buildEventLogsText(receipt, eventDetails)
		eventLogsChan <- eventLogsText
	}()

	// Wait for both responses and combine them
	var methodDetails, eventLogs string
	for i := 0; i < 2; i++ {
		select {
		case methodDetails = <-methodSigChan:
			// Method signature details received
		case eventLogs = <-eventLogsChan:
			// Event logs received
		}
	}

	// Combine the results
	finalDetails := methodDetails
	if eventLogs != "" {
		finalDetails += "\n" + eventLogs
	}

	// Update UI with complete details
	t.app.QueueUpdateDraw(func() {
		t.txDetailLeft.SetText(finalDetails)
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
	receiptJSON, err := rpctypes.PolyReceiptToPrettyJSON(receipt)
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

// buildEventLogsText creates a formatted display of event logs with resolved signatures
func (t *TviewRenderer) buildEventLogsText(receipt rpctypes.PolyReceipt, eventDetails map[string]string) string {
	logs := receipt.Logs()
	if len(logs) == 0 {
		return ""
	}

	var logLines []string
	logLines = append(logLines, "Event Logs:")

	for i, logEntry := range logs {
		// Get the contract address that emitted the event
		contractAddr := logEntry.Address.ToAddress().Hex()
		contractAddrShort := truncateHash(contractAddr, 6, 4)

		if len(logEntry.Topics) > 0 {
			// Get the event signature hash and look up its human-readable name
			eventSigHash := logEntry.Topics[0].ToHash().Hex()
			eventName := "Unknown"

			if eventDetails != nil {
				if name, exists := eventDetails[eventSigHash]; exists {
					eventName = name
				}
			}

			// Format: "  [index] EventName from 0x1234...5678"
			logLine := fmt.Sprintf("  [%d] %s from %s", i, eventName, contractAddrShort)

			// Add topic count if there are indexed parameters
			if len(logEntry.Topics) > 1 {
				logLine += fmt.Sprintf(" (%d indexed args)", len(logEntry.Topics)-1)
			}

			logLines = append(logLines, logLine)
		} else {
			// Anonymous event (no topics)
			logLine := fmt.Sprintf("  [%d] Anonymous Event from %s", i, contractAddrShort)
			logLines = append(logLines, logLine)
		}
	}

	// Add summary line
	if len(logs) > 0 {
		logLines = append(logLines, fmt.Sprintf("  Total: %d event(s)", len(logs)))
	}

	return strings.Join(logLines, "\n")
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

	// Save current focus if a modal is active
	var currentFocus tview.Primitive
	if t.isModalCurrentlyActive() {
		currentFocus = t.app.GetFocus()
	}

	t.lastDrawTime = now
	t.app.Draw()

	// Restore focus to modal if it was stolen during draw
	if t.isModalCurrentlyActive() && currentFocus != nil {
		// Small delay to ensure the draw is complete before restoring focus
		go func() {
			time.Sleep(1 * time.Millisecond)
			t.app.QueueUpdateDraw(func() {
				t.app.SetFocus(currentFocus)
			})
		}()
	}
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
		if i >= 100 { // Limit to 100 blocks for performance
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

		// Column 6: Gas used
		gasUsedStr := formatNumber(block.GasUsed())
		t.homeTable.SetCell(row, 6, tview.NewTableCell(gasUsedStr).SetAlign(t.columns[6].Align))

		// Column 7: Gas percentage
		gasPercentStr := formatGasPercentage(block.GasUsed(), block.GasLimit())
		t.homeTable.SetCell(row, 7, tview.NewTableCell(gasPercentStr).SetAlign(t.columns[7].Align))

		// Column 8: Gas limit
		gasLimitStr := formatNumber(block.GasLimit())
		t.homeTable.SetCell(row, 8, tview.NewTableCell(gasLimitStr).SetAlign(t.columns[8].Align))

		// Column 9: State root (truncated)
		stateRootStr := truncateHash(block.Root().Hex(), 8, 8)
		t.homeTable.SetCell(row, 9, tview.NewTableCell(stateRootStr).SetAlign(t.columns[9].Align))
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

// ChainInfo represents a blockchain network from chainlist
type ChainInfo struct {
	Name    string `json:"name"`
	ChainID int64  `json:"chainId"`
	Chain   string `json:"chain"`
}

// chainlistCache holds the cached chain information
var (
	chainlistCache map[int64]ChainInfo
	chainlistMu    sync.RWMutex
	chainlistFetch sync.Once
)

// getNetworkName returns the human-readable network name for a given chain ID
func getNetworkName(chainID *big.Int) string {
	if chainID == nil {
		return "Unknown"
	}

	chainIDInt64 := chainID.Int64()

	// Try to get from chainlist cache
	chainlistMu.RLock()
	if chainlistCache != nil {
		if chain, exists := chainlistCache[chainIDInt64]; exists {
			chainlistMu.RUnlock()
			return chain.Name
		}
	}
	chainlistMu.RUnlock()

	// Initialize chainlist cache on first use
	chainlistFetch.Do(initChainlist)

	// Try again after initialization
	chainlistMu.RLock()
	if chainlistCache != nil {
		if chain, exists := chainlistCache[chainIDInt64]; exists {
			chainlistMu.RUnlock()
			return chain.Name
		}
	}
	chainlistMu.RUnlock()

	// Fallback to static mapping for common chains
	staticNames := map[int64]string{
		1:        "Ethereum Mainnet",
		137:      "Polygon PoS",
		56:       "BNB Smart Chain",
		10:       "Optimism",
		42161:    "Arbitrum One",
		43114:    "Avalanche C-Chain",
		250:      "Fantom Opera",
		8453:     "Base",
		100:      "Gnosis Chain",
		324:      "zkSync Era",
		1101:     "Polygon zkEVM",
		80001:    "Polygon Mumbai",
		11155111: "Sepolia Testnet",
		5:        "Goerli Testnet",
	}

	if name, exists := staticNames[chainIDInt64]; exists {
		return name
	}

	// Return chain ID if name not found
	return fmt.Sprintf("Chain %s", chainID.String())
}

// initChainlist downloads and caches chain information from chainlist.org
func initChainlist() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://chainid.network/chains.json", nil)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create chainlist request")
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to fetch chainlist")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Debug().Int("status", resp.StatusCode).Msg("Chainlist request failed")
		return
	}

	var chains []ChainInfo
	if err := json.NewDecoder(resp.Body).Decode(&chains); err != nil {
		log.Debug().Err(err).Msg("Failed to decode chainlist")
		return
	}

	// Build cache map
	cache := make(map[int64]ChainInfo, len(chains))
	for _, chain := range chains {
		cache[chain.ChainID] = chain
	}

	chainlistMu.Lock()
	chainlistCache = cache
	chainlistMu.Unlock()

	log.Debug().Int("chains", len(chains)).Msg("Loaded chainlist cache")
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

// hexToBigInt converts hex string to big.Int (helper for sync status)
func hexToBigInt(hex string) (*big.Int, error) {
	if len(hex) >= 2 && hex[:2] == "0x" {
		hex = hex[2:]
	}
	result := big.NewInt(0)
	result, ok := result.SetString(hex, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string")
	}
	return result, nil
}

// formatConnectionStatus converts latency duration to human-readable connection status
func formatConnectionStatus(latency time.Duration) string {
	ms := latency.Milliseconds()

	switch {
	case ms < 50:
		return fmt.Sprintf("[OK] Excellent (%dms)", ms)
	case ms < 150:
		return fmt.Sprintf("[OK] Good (%dms)", ms)
	case ms < 300:
		return fmt.Sprintf("[OK] Fair (%dms)", ms)
	case ms < 500:
		return fmt.Sprintf("[SLOW] Slow (%dms)", ms)
	case ms < 1000:
		return fmt.Sprintf("[SLOW] Poor (%dms)", ms)
	default:
		return fmt.Sprintf("[POOR] Very Poor (%dms)", ms)
	}
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

	// Check if timestamp can be safely converted to int64
	if timestamp > math.MaxInt64 {
		log.Error().Uint64("timestamp", timestamp).Msg("Timestamp exceeds int64 range")
		return "invalid"
	}

	// Safe conversion after bounds check
	timestampInt64 := int64(timestamp)

	// Handle case where timestamp is in the future
	if timestampInt64 > now {
		return "future"
	}

	diff := now - timestampInt64

	if diff < 60 {
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
	// Check if timestamp can be safely converted to int64
	if timestamp > math.MaxInt64 {
		log.Error().Uint64("timestamp", timestamp).Msg("Timestamp exceeds int64 range")
		// Return a special format for invalid timestamps
		return "invalid timestamp - invalid"
	}

	// Safe conversion after bounds check
	t := time.Unix(int64(timestamp), 0).UTC()
	absolute := t.Format("2006-01-02T15:04:05Z")
	relative := formatRelativeTime(timestamp)
	return fmt.Sprintf("%s - %s", absolute, relative)
}

// calculateBlockInterval calculates the time interval between a block and its parent
func (t *TviewRenderer) calculateBlockInterval(block rpctypes.PolyBlock, index int, blocks []rpctypes.PolyBlock) string {
	// First try to look up parent block by hash (most accurate)
	parentHash := block.ParentHash().Hex()
	if parentBlock, exists := t.blocksByHash[parentHash]; exists {
		// Calculate interval in seconds
		blockTime := block.Time()
		parentTime := parentBlock.Time()

		// Check if times can be safely converted to int64
		if blockTime > math.MaxInt64 || parentTime > math.MaxInt64 {
			log.Error().
				Uint64("block_time", blockTime).
				Uint64("parent_time", parentTime).
				Msg("Time values exceed int64 range")
			return "N/A"
		}

		// Safe conversion after bounds check
		interval := int64(blockTime) - int64(parentTime)
		return fmt.Sprintf("%ds", interval)
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

			// Check if times can be safely converted to int64
			if blockTime > math.MaxInt64 || prevTime > math.MaxInt64 {
				log.Error().
					Uint64("block_time", blockTime).
					Uint64("prev_time", prevTime).
					Msg("Time values exceed int64 range")
				return "N/A"
			}

			// Safe conversion after bounds check
			interval := int64(blockTime) - int64(prevTime)

			// For non-consecutive blocks, show average interval
			if currentBlockNum-prevBlockNum > 1 {
				blockDiff := currentBlockNum - prevBlockNum
				// Check if block difference can be safely converted to int64
				if blockDiff > math.MaxInt64 {
					log.Error().Uint64("block_diff", blockDiff).Msg("Block difference exceeds int64 range")
					return "N/A"
				}
				avgInterval := interval / int64(blockDiff)
				return fmt.Sprintf("~%ds", avgInterval)
			}
			return fmt.Sprintf("%ds", interval)
		}
	}

	// Can't calculate interval (first block in the list or large gap)
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
			// Only set focus if no modal is currently active
			if !t.isModalCurrentlyActive() {
				t.app.SetFocus(t.homeTable)
			}
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
		if v > math.MaxInt64 || v < math.MinInt64 {
			log.Error().Float64("value", v).Msg("Float64 value exceeds int64 range, clamping")
			if v > math.MaxInt64 {
				v = math.MaxInt64
			} else {
				v = math.MinInt64
			}
		}
		return big.NewInt(int64(v)), nil
	case uint64:
		if v > math.MaxInt64 {
			log.Error().Uint64("value", v).Msg("Uint64 value exceeds int64 range, using MaxInt64")
			v = math.MaxInt64
		}
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	case int:
		return big.NewInt(int64(v)), nil
	default:
		return nil, fmt.Errorf("unsupported number type: %T", value)
	}
}

// isValidBlock checks if a block response contains valid data
func (t *TviewRenderer) isValidBlock(block rpctypes.PolyBlock) bool {
	if block == nil {
		log.Debug().Msg("isValidBlock: block is nil")
		return false
	}

	// Check if the block has a valid hash (not zero)
	blockHash := block.Hash()
	if blockHash == (common.Hash{}) {
		log.Debug().Msg("isValidBlock: block hash is zero")
		return false
	}

	// Check if the block has a valid number (not zero)
	blockNum := block.Number()
	if blockNum == nil {
		log.Debug().Msg("isValidBlock: block number is nil")
		return false
	}

	if blockNum.Cmp(big.NewInt(0)) < 0 {
		log.Debug().Str("blockNum", blockNum.String()).Msg("isValidBlock: block number is negative")
		return false
	}

	log.Debug().
		Str("hash", blockHash.Hex()).
		Str("number", blockNum.String()).
		Msg("isValidBlock: block is valid")
	return true
}

// isValidTransaction checks if a transaction response contains valid data
func (t *TviewRenderer) isValidTransaction(tx rpctypes.PolyTransaction) bool {
	if tx == nil {
		log.Debug().Msg("isValidTransaction: transaction is nil")
		return false
	}

	// Check if the transaction has a valid hash (not zero)
	txHash := tx.Hash()
	if txHash == (common.Hash{}) {
		log.Debug().Msg("isValidTransaction: transaction hash is zero")
		return false
	}

	// Check if the transaction has a from address (all transactions must have a sender)
	fromAddr := tx.From()
	if fromAddr == (common.Address{}) {
		log.Debug().Msg("isValidTransaction: from address is zero")
		return false
	}

	// Check if block number is set (transaction has been mined)
	blockNum := tx.BlockNumber()
	if blockNum == nil || blockNum.Cmp(big.NewInt(0)) <= 0 {
		log.Debug().
			Bool("blockNumNil", blockNum == nil).
			Msg("isValidTransaction: invalid block number (pending or invalid)")
		return false
	}

	log.Debug().
		Str("hash", txHash.Hex()).
		Str("from", fromAddr.Hex()).
		Str("blockNumber", blockNum.String()).
		Msg("isValidTransaction: transaction is valid")
	return true
}

// Search functionality methods

// performSearch parses the search query and determines the search type
func (t *TviewRenderer) performSearch(query string) {
	query = strings.TrimSpace(query)
	log.Debug().Str("query", query).Msg("performSearch called")

	if query == "" {
		log.Debug().Msg("Empty search query, returning")
		return
	}

	// Try to parse as a number (block number)
	if blockNum, err := strconv.ParseUint(query, 10, 64); err == nil {
		log.Debug().Uint64("blockNum", blockNum).Msg("Detected block number search")
		go t.searchBlockByNumber(blockNum)
		return
	}

	// Check if it looks like a hash (0x prefix and 66 characters for full hash)
	if strings.HasPrefix(query, "0x") && len(query) == 66 {
		log.Debug().Str("hash", query).Msg("Detected hash search")
		go t.searchByHash(query)
		return
	}

	// Show error for invalid format
	log.Debug().Str("query", query).Msg("Invalid search format")
	t.showSearchError("Invalid input. Enter a block number or hash (0x...)")
}

// searchBlockByNumber searches for a block by its number
func (t *TviewRenderer) searchBlockByNumber(blockNum uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert to big.Int for the indexer call
	var constrainedBlockNum int64
	if blockNum > math.MaxInt64 {
		log.Error().Uint64("block_num", blockNum).Msg("Block number exceeds int64 range, using MaxInt64")
		constrainedBlockNum = math.MaxInt64
	} else {
		constrainedBlockNum = int64(blockNum)
	}
	blockNumber := big.NewInt(constrainedBlockNum)

	block, err := t.indexer.GetBlock(ctx, blockNumber)
	if err != nil {
		log.Debug().Err(err).Uint64("blockNum", blockNum).Msg("Block not found")
		t.showSearchError(fmt.Sprintf("Block #%d not found", blockNum))
		return
	}

	// Navigate to block detail page
	t.app.QueueUpdateDraw(func() {
		t.showBlockDetail(block)
	})
}

// searchByHash searches for a block or transaction by hash
func (t *TviewRenderer) searchByHash(hash string) {
	log.Debug().Str("hash", hash).Msg("searchByHash started")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string to common.Hash
	hashBytes := common.HexToHash(hash)
	log.Debug().Str("hashBytes", hashBytes.Hex()).Msg("Converted to common.Hash")

	// Try to find as a transaction first (more common than block hash searches)
	log.Debug().Msg("Attempting to find as transaction")
	if tx, err := t.indexer.GetTransaction(ctx, hashBytes); err == nil {
		log.Debug().
			Str("txHash", hashBytes.Hex()).
			Interface("tx", tx).
			Msg("Found as transaction")

		// Validate the transaction has actual data
		if t.isValidTransaction(tx) {
			log.Debug().Msg("Transaction is valid, showing detail")
			// Found as transaction - need to get the containing block
			go t.showTransactionWithBlock(tx, hashBytes)
			return
		} else {
			log.Debug().Msg("Transaction response is empty/invalid, trying as block")
		}
	} else {
		log.Debug().Err(err).Str("hash", hash).Msg("Not found as transaction")
	}

	// Try to find as a block
	log.Debug().Msg("Attempting to find as block")
	block, err := t.indexer.GetBlock(ctx, hashBytes)
	if err != nil {
		log.Debug().Err(err).Str("hash", hash).Msg("Error getting block")
	} else {
		log.Debug().
			Str("blockHash", hash).
			Bool("blockNotNil", block != nil).
			Msg("GetBlock returned")

		if block != nil {
			blockHash := block.Hash()
			blockNum := block.Number()
			log.Debug().
				Str("returnedHash", blockHash.Hex()).
				Bool("hashIsZero", blockHash == (common.Hash{})).
				Interface("blockNumber", blockNum).
				Bool("numberIsNil", blockNum == nil).
				Msg("Block details")
		}

		isValid := t.isValidBlock(block)
		log.Debug().Bool("isValidBlock", isValid).Msg("Block validation result")

		if isValid {
			// Found as block and it's valid (not empty) - navigate to block detail
			log.Debug().Msg("Showing block detail")
			t.app.QueueUpdateDraw(func() {
				t.showBlockDetail(block)
			})
			return
		}
	}

	// Not found as either block or transaction
	log.Debug().Str("hash", hash).Msg("Hash not found as block or transaction")
	t.showSearchError(fmt.Sprintf("Hash %s not found", truncateHash(hash, 8, 8)))
}

// showTransactionWithBlock gets a transaction's block and shows the transaction detail
func (t *TviewRenderer) showTransactionWithBlock(tx rpctypes.PolyTransaction, txHash common.Hash) {
	log.Debug().Str("txHash", txHash.Hex()).Msg("showTransactionWithBlock started")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the receipt to find the block number
	log.Debug().Msg("Getting receipt for transaction")
	receipt, err := t.indexer.GetReceipt(ctx, txHash)
	if err != nil {
		log.Debug().Err(err).Str("txHash", txHash.Hex()).Msg("Failed to get receipt")
		t.showSearchError("Transaction found but could not load details")
		return
	}

	// Get the block number from the receipt
	blockNumber := receipt.BlockNumber()
	log.Debug().Str("blockNumber", blockNumber.String()).Msg("Got block number from receipt")

	// Get the full block
	log.Debug().Msg("Getting block for transaction")
	block, err := t.indexer.GetBlock(ctx, blockNumber)
	if err != nil {
		log.Debug().Err(err).Str("blockNum", blockNumber.String()).Msg("Failed to get block for transaction")
		t.showSearchError("Transaction found but could not load block context")
		return
	}

	// Find the transaction index within the block
	transactions := block.Transactions()
	log.Debug().Int("txCount", len(transactions)).Msg("Searching for transaction in block")

	txIndex := -1
	for i, blockTx := range transactions {
		if blockTx.Hash().Hex() == txHash.Hex() {
			txIndex = i
			log.Debug().Int("txIndex", i).Msg("Found transaction at index")
			break
		}
	}

	if txIndex == -1 {
		log.Debug().Str("txHash", txHash.Hex()).Msg("Transaction not found in its block")
		t.showSearchError("Transaction found but could not locate in block")
		return
	}

	// Update UI on main thread
	log.Debug().Msg("Updating UI with transaction detail")
	t.app.QueueUpdateDraw(func() {
		// Set the block as current so back navigation works
		t.currentBlockMu.Lock()
		t.currentBlock = block
		t.currentBlockMu.Unlock()

		// Show transaction detail
		t.showTransactionDetail(tx, txIndex)
		log.Debug().Msg("showTransactionDetail called")
	})
}

// showSearchError displays an error message to the user
func (t *TviewRenderer) showSearchError(message string) {
	// For now, just log the error. In the future, we could show a toast or status message
	log.Info().Str("searchError", message).Msg("Search error")
	log.Debug().Str("errorMessage", message).Msg("showSearchError called")

	// TODO: Could implement a status bar or toast notification here
}

// Modal management methods

// showModal displays a modal and tracks its state
func (t *TviewRenderer) showModal(name string) {
	t.modalStateMu.Lock()
	t.isModalActive = true
	t.activeModalName = name
	t.modalStateMu.Unlock()

	// Show the modal page
	t.pages.ShowPage(name)

	// Set focus to the modal based on its name
	switch name {
	case "quit":
		t.quitModal.SetFocus(0) // Always start on "Yes"
		t.app.SetFocus(t.quitModal)
	case "search":
		// Clear any existing text and focus on the input field
		inputField := t.searchForm.GetFormItem(0).(*tview.InputField)
		inputField.SetText("")
		t.app.SetFocus(inputField)
	}
}

// hideModal hides the currently active modal and clears state
func (t *TviewRenderer) hideModal() {
	t.modalStateMu.Lock()
	modalName := t.activeModalName
	t.isModalActive = false
	t.activeModalName = ""
	t.modalStateMu.Unlock()

	// Hide the current modal page
	if modalName != "" {
		t.pages.HidePage(modalName)
	}

	// Return focus to the home page
	t.app.SetFocus(t.homeTable)
}

// isModalCurrentlyActive returns true if any modal is currently active
func (t *TviewRenderer) isModalCurrentlyActive() bool {
	t.modalStateMu.RLock()
	defer t.modalStateMu.RUnlock()
	return t.isModalActive
}

// Stop gracefully stops the TUI renderer
func (t *TviewRenderer) Stop() error {
	log.Info().Msg("Stopping Tview renderer")
	t.app.Stop()
	return nil
}
