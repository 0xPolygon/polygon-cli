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

// nodb is a database implementation that does nothing.
// It's useful when you want to run the sensor without persisting any data.
type nodb struct{}

// NoDatabase creates a new nodb instance.
func NoDatabase() Database {
	return &nodb{}
}

// WriteBlock does nothing.
func (n *nodb) WriteBlock(ctx context.Context, _ GoroutineTracker, peer *enode.Node, block *types.Block, td *big.Int, tfs time.Time) {
}

// WriteBlockHeaders does nothing.
func (n *nodb) WriteBlockHeaders(ctx context.Context, _ GoroutineTracker, headers []*types.Header, tfs time.Time) {
}

// WriteBlockHashes does nothing.
func (n *nodb) WriteBlockHashes(ctx context.Context, _ GoroutineTracker, peer *enode.Node, hashes []common.Hash, tfs time.Time) {
}

// WriteBlockBody does nothing.
func (n *nodb) WriteBlockBody(ctx context.Context, _ GoroutineTracker, body *eth.BlockBody, hash common.Hash, tfs time.Time) {
}

// WriteTransactions does nothing.
func (n *nodb) WriteTransactions(ctx context.Context, _ GoroutineTracker, peer *enode.Node, txs []*types.Transaction, tfs time.Time) {
}

// WritePeers does nothing.
func (n *nodb) WritePeers(ctx context.Context, peers []*p2p.Peer, tls time.Time) {
}

// HasBlock always returns true to avoid re-fetching.
func (n *nodb) HasBlock(ctx context.Context, hash common.Hash) bool {
	return true
}

// MaxConcurrentWrites returns 0 as no actual writes occur.
func (n *nodb) MaxConcurrentWrites() int {
	return 0
}

// ShouldWriteBlocks returns false.
func (n *nodb) ShouldWriteBlocks() bool {
	return false
}

// ShouldWriteBlockEvents returns false.
func (n *nodb) ShouldWriteBlockEvents() bool {
	return false
}

// ShouldWriteTransactions returns false.
func (n *nodb) ShouldWriteTransactions() bool {
	return false
}

// ShouldWriteTransactionEvents returns false.
func (n *nodb) ShouldWriteTransactionEvents() bool {
	return false
}

// ShouldWritePeers returns false.
func (n *nodb) ShouldWritePeers() bool {
	return false
}

// NodeList returns an empty list.
func (n *nodb) NodeList(ctx context.Context, limit int) ([]string, error) {
	return []string{}, nil
}
