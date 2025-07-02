package renderer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

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

	// Reset scroll position to the top
	t.txDetailLeft.ScrollToBeginning()
	t.txDetailTxJSON.ScrollToBeginning()
	t.txDetailRcptJSON.ScrollToBeginning()

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
	details = append(details, fmt.Sprintf("Gas Price: %s", formatBaseFee(tx.GasPrice())))
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
			details = append(details, fmt.Sprintf("Max Fee Per Gas: %s", formatBaseFee(maxFeeBig)))
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
			details = append(details, fmt.Sprintf("Max Priority Fee Per Gas: %s", formatBaseFee(maxPriorityBig)))
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
	details = append(details, fmt.Sprintf("Gas Price: %s", formatBaseFee(tx.GasPrice())))
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
			details = append(details, fmt.Sprintf("Max Fee Per Gas: %s", formatBaseFee(maxFeeBig)))
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
			details = append(details, fmt.Sprintf("Max Priority Fee Per Gas: %s", formatBaseFee(maxPriorityBig)))
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
