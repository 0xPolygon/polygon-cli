package database

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// NoopDatabase is a database implementation that does nothing.
// It's useful when you want to run the sensor without persisting any data.
type NoopDatabase struct{}

// NewNoopDatabase creates a new NoopDatabase instance.
func NewNoopDatabase() Database {
	return &NoopDatabase{}
}

// WriteBlock does nothing.
func (n *NoopDatabase) WriteBlock(ctx context.Context, peer *enode.Node, block *types.Block, td *big.Int, tfs time.Time) {
}

// WriteBlockHeaders does nothing.
func (n *NoopDatabase) WriteBlockHeaders(ctx context.Context, headers []*types.Header, tfs time.Time) {
}

// WriteBlockHashes does nothing.
func (n *NoopDatabase) WriteBlockHashes(ctx context.Context, peer *enode.Node, hashes []common.Hash, tfs time.Time) {
}

// WriteBlockBody does nothing.
func (n *NoopDatabase) WriteBlockBody(ctx context.Context, body *eth.BlockBody, hash common.Hash, tfs time.Time) {
}

// WriteTransactions does nothing.
func (n *NoopDatabase) WriteTransactions(ctx context.Context, peer *enode.Node, txs []*types.Transaction, tfs time.Time) {
}

// WritePeers does nothing.
func (n *NoopDatabase) WritePeers(ctx context.Context, peers []*p2p.Peer, tls time.Time) {
}

// HasBlock always returns true to avoid re-fetching.
func (n *NoopDatabase) HasBlock(ctx context.Context, hash common.Hash) bool {
	return true
}

// MaxConcurrentWrites returns 1 as no actual writes occur.
func (n *NoopDatabase) MaxConcurrentWrites() int {
	return 1
}

// ShouldWriteBlocks returns false.
func (n *NoopDatabase) ShouldWriteBlocks() bool {
	return false
}

// ShouldWriteBlockEvents returns false.
func (n *NoopDatabase) ShouldWriteBlockEvents() bool {
	return false
}

// ShouldWriteTransactions returns false.
func (n *NoopDatabase) ShouldWriteTransactions() bool {
	return false
}

// ShouldWriteTransactionEvents returns false.
func (n *NoopDatabase) ShouldWriteTransactionEvents() bool {
	return false
}

// ShouldWritePeers returns false.
func (n *NoopDatabase) ShouldWritePeers() bool {
	return false
}

// NodeList returns an empty list.
func (n *NoopDatabase) NodeList(ctx context.Context, limit int) ([]string, error) {
	return []string{}, nil
}