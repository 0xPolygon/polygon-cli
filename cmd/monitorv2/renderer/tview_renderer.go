package renderer

import (
	"container/list"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-cli/chainstore"
	"github.com/0xPolygon/polygon-cli/indexer"
	polymetrics "github.com/0xPolygon/polygon-cli/metrics"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// Zero address constant for efficient comparisons
var zeroAddress = common.Address{}

// Maximum number of blocks to store and display
const maxBlocks = 1000

// Maximum number of signer cache entries (LRU eviction when exceeded)
const maxSignerCacheSize = 1000

// signerCacheEntry represents an entry in the LRU cache
type signerCacheEntry struct {
	key   string
	value string
}

// signerLRUCache implements a simple LRU cache for signer addresses
type signerLRUCache struct {
	capacity int
	cache    map[string]*list.Element
	lruList  *list.List
	mu       sync.RWMutex
}

// newSignerLRUCache creates a new LRU cache with the given capacity
func newSignerLRUCache(capacity int) *signerLRUCache {
	return &signerLRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		lruList:  list.New(),
	}
}

// get retrieves a value from the cache and marks it as recently used
func (c *signerLRUCache) get(key string) (string, bool) {
	c.mu.RLock()
	elem, exists := c.cache[key]
	c.mu.RUnlock()

	if !exists {
		return "", false
	}

	c.mu.Lock()
	c.lruList.MoveToFront(elem)
	c.mu.Unlock()

	return elem.Value.(*signerCacheEntry).value, true
}

// put adds or updates a value in the cache
func (c *signerLRUCache) put(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if key already exists
	if elem, exists := c.cache[key]; exists {
		// Update existing entry and move to front
		c.lruList.MoveToFront(elem)
		elem.Value.(*signerCacheEntry).value = value
		return
	}

	// Add new entry
	entry := &signerCacheEntry{key: key, value: value}
	elem := c.lruList.PushFront(entry)
	c.cache[key] = elem

	// Evict least recently used if over capacity
	if c.lruList.Len() > c.capacity {
		c.evictOldest()
	}
}

// evictOldest removes the least recently used entry (must be called with lock held)
func (c *signerLRUCache) evictOldest() {
	elem := c.lruList.Back()
	if elem != nil {
		c.lruList.Remove(elem)
		entry := elem.Value.(*signerCacheEntry)
		delete(c.cache, entry.key)
	}
}

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
	// Signer LRU cache to avoid expensive Ecrecover operations
	signerCache *signerLRUCache
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
	isModalActive    bool
	activeModalName  string
	previousPageName string // Track page before modal was opened
	modalStateMu     sync.RWMutex
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
		signerCache:  newSignerLRUCache(maxSignerCacheSize),
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
			} else if currentPage == "home" {
				// On home page, reset table selection to top
				if t.homeTable != nil {
					t.homeTable.Select(1, 0) // Row 1 (first data row, since 0 is header)
					t.app.SetFocus(t.homeTable)
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

// getCachedSigner gets the signer for a block, using LRU cache to avoid expensive Ecrecover calls
func (t *TviewRenderer) getCachedSigner(block rpctypes.PolyBlock) string {
	blockHash := block.Hash().Hex()

	// Check cache first
	if cached, exists := t.signerCache.get(blockHash); exists {
		return cached
	}

	// If miner is non-zero, use the miner
	zeroAddr := common.Address{}
	if block.Miner() != zeroAddr {
		result := truncateHash(block.Miner().Hex(), 6, 4)
		t.signerCache.put(blockHash, result)
		return result
	}

	// If miner is zero, try to extract signer from extra data (EXPENSIVE - Ecrecover)
	if signer, err := polymetrics.Ecrecover(&block); err == nil {
		signerAddr := common.HexToAddress("0x" + hex.EncodeToString(signer))
		result := truncateHash(signerAddr.Hex(), 6, 4)
		t.signerCache.put(blockHash, result)
		return result
	}

	// If can't extract signer, cache N/A result
	result := "N/A"
	t.signerCache.put(blockHash, result)
	return result
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

// formatBaseFee formats base fee in wei to human-readable format using appropriate SI units
func formatBaseFee(baseFee *big.Int) string {
	if baseFee == nil || baseFee.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	// Define units with their divisors, names, thresholds, and decimal precision
	type unit struct {
		divisor   *big.Int
		name      string
		threshold *big.Int
		decimals  int
	}

	units := []unit{
		{new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), "ether", new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), 3}, // 10^18
		{new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil), "milli", new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil), 3}, // 10^15 milliether
		{new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil), "micro", new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil), 3}, // 10^12 microether
		{new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil), "gwei", new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil), 3},    // 10^9 gwei
		{new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil), "mwei", new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil), 3},    // 10^6 megawei
		{new(big.Int).Exp(big.NewInt(10), big.NewInt(3), nil), "kwei", new(big.Int).Exp(big.NewInt(10), big.NewInt(3), nil), 3},    // 10^3 kilowei
		{big.NewInt(1), "wei", big.NewInt(1), 0}, // wei (no decimals)
	}

	// Find the appropriate unit (largest unit where value >= threshold)
	for _, u := range units {
		if baseFee.Cmp(u.threshold) >= 0 {
			if u.name == "wei" {
				// For wei, add thousand separators for values >= 1,000
				weiStr := baseFee.String()
				if len(weiStr) > 3 {
					return fmt.Sprintf("%s wei", addThousandSeparators(weiStr))
				}
				return fmt.Sprintf("%s wei", weiStr)
			}

			// Convert to the selected unit using big.Float for precision
			value := new(big.Float).SetInt(baseFee)
			divisor := new(big.Float).SetInt(u.divisor)
			result := new(big.Float).Quo(value, divisor)

			// Format with appropriate precision
			formatStr := fmt.Sprintf("%%.%df %%s", u.decimals)
			resultFloat, _ := result.Float64()

			// Remove trailing zeros from decimal representation
			formatted := fmt.Sprintf(formatStr, resultFloat, u.name)
			return removeTrailingZeros(formatted)
		}
	}

	// Fallback (should never reach here)
	return baseFee.String() + " wei"
}

// addThousandSeparators adds commas to a numeric string for readability
func addThousandSeparators(numStr string) string {
	if len(numStr) <= 3 {
		return numStr
	}

	var result []rune
	for i, char := range numStr {
		if i > 0 && (len(numStr)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, char)
	}
	return string(result)
}

// removeTrailingZeros removes trailing zeros from decimal numbers in formatted strings
func removeTrailingZeros(formatted string) string {
	// Split on space to separate number from unit
	parts := strings.Split(formatted, " ")
	if len(parts) != 2 {
		return formatted
	}

	numberPart := parts[0]
	unitPart := parts[1]

	// If it contains a decimal point, remove trailing zeros
	if strings.Contains(numberPart, ".") {
		numberPart = strings.TrimRight(numberPart, "0")
		numberPart = strings.TrimRight(numberPart, ".")
	}

	return numberPart + " " + unitPart
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
	// Get current page before showing modal
	currentPage, _ := t.pages.GetFrontPage()

	t.modalStateMu.Lock()
	t.isModalActive = true
	t.activeModalName = name
	t.previousPageName = currentPage // Track previous page
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
	previousPage := t.previousPageName
	t.isModalActive = false
	t.activeModalName = ""
	t.previousPageName = ""
	t.modalStateMu.Unlock()

	// Hide the current modal page
	if modalName != "" {
		t.pages.HidePage(modalName)
	}

	// Return focus to the appropriate UI element based on the previous page
	switch previousPage {
	case "home":
		if t.homeTable != nil {
			t.app.SetFocus(t.homeTable)
		}
	case "block-detail":
		if t.blockDetailLeft != nil {
			t.app.SetFocus(t.blockDetailLeft)
		}
	case "tx-detail":
		if t.txDetailLeft != nil {
			t.app.SetFocus(t.txDetailLeft)
		}
	case "info":
		if t.infoPage != nil {
			t.app.SetFocus(t.infoPage)
		}
	case "help":
		if t.helpPage != nil {
			t.app.SetFocus(t.helpPage)
		}
	default:
		// Fallback to home page
		if t.homeTable != nil {
			t.app.SetFocus(t.homeTable)
		}
	}
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
