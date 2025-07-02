package renderer

import (
	"context"
	"encoding/json"
	"os"

	"github.com/0xPolygon/polygon-cli/indexer"
	"github.com/rs/zerolog/log"
)

// JSONRenderer outputs blockchain data as line-delimited JSON
type JSONRenderer struct {
	BaseRenderer
}

// NewJSONRenderer creates a new JSON renderer
func NewJSONRenderer(indexer *indexer.Indexer) *JSONRenderer {
	return &JSONRenderer{
		BaseRenderer: NewBaseRenderer(indexer),
	}
}

// Start begins rendering JSON output
func (j *JSONRenderer) Start(ctx context.Context) error {
	log.Info().Msg("Starting JSON renderer")

	// Consume blocks from the indexer's channel
	blockChan := j.indexer.BlockChannel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case block, ok := <-blockChan:
			if !ok {
				log.Info().Msg("Block channel closed, stopping JSON renderer")
				return nil
			}
			if err := j.outputBlock(block); err != nil {
				log.Error().Err(err).Msg("Error outputting block")
			}
		}
	}
}

// Stop gracefully stops the JSON renderer
func (j *JSONRenderer) Stop() error {
	return nil
}

// outputBlock outputs a block as JSON
func (j *JSONRenderer) outputBlock(block interface{}) error {
	// Output as line-delimited JSON
	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(block)
}
