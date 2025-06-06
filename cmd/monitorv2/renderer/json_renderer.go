package renderer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/0xPolygon/polygon-cli/blockstore"
	"github.com/0xPolygon/polygon-cli/indexer"
)

// JSONRenderer outputs blockchain data as line-delimited JSON
type JSONRenderer struct {
	BaseRenderer
}

// NewJSONRenderer creates a new JSON renderer
func NewJSONRenderer(store blockstore.BlockStore, indexer *indexer.Indexer) *JSONRenderer {
	return &JSONRenderer{
		BaseRenderer: NewBaseRenderer(store, indexer),
	}
}

// Start begins rendering JSON output
func (j *JSONRenderer) Start(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "Starting JSON renderer...")

	// For now, just fetch and output the latest block every few seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := j.outputLatestBlock(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching latest block: %v\n", err)
			}
		}
	}
}

// Stop gracefully stops the JSON renderer
func (j *JSONRenderer) Stop() error {
	return nil
}

// outputLatestBlock fetches and outputs the latest block as JSON
func (j *JSONRenderer) outputLatestBlock(ctx context.Context) error {
	block, err := j.store.GetLatestBlock(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	// Create a JSON event
	event := map[string]interface{}{
		"type":      "block",
		"timestamp": time.Now().Unix(),
		"data":      block,
	}

	// Output as line-delimited JSON
	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(event)
}