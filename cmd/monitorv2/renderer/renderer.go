package renderer

import (
	"context"

	"github.com/0xPolygon/polygon-cli/indexer"
)

// Renderer defines the interface for different output renderers (TUI, JSON, etc.)
type Renderer interface {
	// Start begins rendering output
	Start(ctx context.Context) error

	// Stop gracefully stops the renderer
	Stop() error
}

// BaseRenderer contains common fields that all renderers will need
type BaseRenderer struct {
	indexer *indexer.Indexer
}

// NewBaseRenderer creates a new base renderer with the given dependencies
func NewBaseRenderer(indexer *indexer.Indexer) BaseRenderer {
	return BaseRenderer{
		indexer: indexer,
	}
}