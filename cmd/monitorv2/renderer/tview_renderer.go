package renderer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/0xPolygon/polygon-cli/indexer"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// TviewRenderer provides a terminal UI using the tview library
type TviewRenderer struct {
	BaseRenderer
	app   *tview.Application
	pages *tview.Pages
	blocks []rpctypes.PolyBlock
	
	// Pages
	homePage        *tview.Table
	blockDetailPage *tview.TextView
	txDetailPage    *tview.TextView
	infoPage        *tview.TextView
	helpPage        *tview.TextView
	
	// Modals
	quitModal *tview.Modal
}

// NewTviewRenderer creates a new TUI renderer using tview
func NewTviewRenderer(indexer *indexer.Indexer) *TviewRenderer {
	app := tview.NewApplication()
	
	renderer := &TviewRenderer{
		BaseRenderer: NewBaseRenderer(indexer),
		app:          app,
		blocks:       make([]rpctypes.PolyBlock, 0),
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

// createHomePage creates the main block listing table
func (t *TviewRenderer) createHomePage() {
	t.homePage = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0).  // Fix the header row so it doesn't scroll
		SetSeparator(' ') // Use space as separator instead of borders
	
	// Set up table headers
	t.homePage.SetCell(0, 0, tview.NewTableCell("BLOCK NUMBER").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignRight).
		SetExpansion(1).
		SetAttributes(tcell.AttrBold))
	t.homePage.SetCell(0, 1, tview.NewTableCell("BLOCK HASH").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignLeft).
		SetExpansion(2).
		SetAttributes(tcell.AttrBold))
	t.homePage.SetCell(0, 2, tview.NewTableCell("TXS").
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetAlign(tview.AlignRight).
		SetExpansion(1).
		SetAttributes(tcell.AttrBold))
	
	// Set up selection handler for Enter key
	t.homePage.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(t.blocks) { // Skip header row
			// Navigate to block detail page
			t.showBlockDetail(t.blocks[row-1])
		}
	})
}

// createBlockDetailPage creates the block detail view
func (t *TviewRenderer) createBlockDetailPage() {
	t.blockDetailPage = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	
	t.blockDetailPage.SetTitle(" Block Detail ")
	t.blockDetailPage.SetBorder(true)
	t.blockDetailPage.SetText("Block detail view - placeholder\n\nPress 'Esc' to go back to home")
}

// createTransactionDetailPage creates the transaction detail view
func (t *TviewRenderer) createTransactionDetailPage() {
	t.txDetailPage = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	
	t.txDetailPage.SetTitle(" Transaction Detail ")
	t.txDetailPage.SetBorder(true)
	t.txDetailPage.SetText("Transaction detail view - placeholder\n\nPress 'Esc' to go back")
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
			return nil
		}
		
		// Page-specific shortcuts
		switch currentPage {
		case "home":
			switch event.Key() {
			case tcell.KeyEnter:
				// Handle Enter key on home page (block selection)
				if t.homePage != nil {
					row, _ := t.homePage.GetSelection()
					if row > 0 && row-1 < len(t.blocks) {
						t.showBlockDetail(t.blocks[row-1])
					}
				}
				return nil
			}
		}
		
		return event
	})
}

// showBlockDetail navigates to block detail page and populates it
func (t *TviewRenderer) showBlockDetail(block rpctypes.PolyBlock) {
	// Update block detail content
	detailText := fmt.Sprintf(`Block Details:

Block Number: %s
Block Hash: %s
Parent Hash: %s
Transactions: %d

Press 'Esc' to go back to home`, 
		block.Number().String(),
		block.Hash().Hex(),
		block.ParentHash().Hex(),
		len(block.Transactions()))
	
	t.blockDetailPage.SetText(detailText)
	t.pages.SwitchToPage("block-detail")
}

// Start begins the TUI rendering
func (t *TviewRenderer) Start(ctx context.Context) error {
	log.Info().Msg("Starting Tview renderer")
	
	// Start consuming blocks in a separate goroutine
	go t.consumeBlocks(ctx)
	
	// Start the TUI application
	// This will block until the application is stopped
	if err := t.app.Run(); err != nil {
		log.Error().Err(err).Msg("Error running tview application")
		return err
	}
	
	return nil
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
			
			// Add block to the beginning of the slice (descending order)
			t.blocks = append([]rpctypes.PolyBlock{block}, t.blocks...)
			
			// Update the table in the main thread
			t.app.QueueUpdateDraw(func() {
				t.updateTable()
			})
		}
	}
}

// updateTable refreshes the home page table with current blocks
func (t *TviewRenderer) updateTable() {
	if t.homePage == nil {
		return
	}
	
	// Clear existing rows (except header)
	rowCount := t.homePage.GetRowCount()
	for row := 1; row < rowCount; row++ {
		for col := 0; col < 3; col++ {
			t.homePage.SetCell(row, col, nil)
		}
	}
	
	// Add blocks to table (newest first)
	for i, block := range t.blocks {
		if i >= 100 { // Limit to 100 blocks for performance
			break
		}
		
		row := i + 1 // +1 to account for header row
		
		// Block number
		blockNum := block.Number().String()
		t.homePage.SetCell(row, 0, tview.NewTableCell(blockNum).SetAlign(tview.AlignRight))
		
		// Block hash (truncated for display)
		hashStr := block.Hash().Hex()
		if len(hashStr) > 20 {
			hashStr = hashStr[:10] + "..." + hashStr[len(hashStr)-10:]
		}
		t.homePage.SetCell(row, 1, tview.NewTableCell(hashStr).SetAlign(tview.AlignLeft))
		
		// Number of transactions
		txCount := len(block.Transactions())
		t.homePage.SetCell(row, 2, tview.NewTableCell(strconv.Itoa(txCount)).SetAlign(tview.AlignRight))
	}
	
	// Set table title with current block count
	title := fmt.Sprintf(" Blockchain Monitor (%d blocks) ", len(t.blocks))
	t.homePage.SetTitle(title)
}

// Stop gracefully stops the TUI renderer
func (t *TviewRenderer) Stop() error {
	log.Info().Msg("Stopping Tview renderer")
	t.app.Stop()
	return nil
}