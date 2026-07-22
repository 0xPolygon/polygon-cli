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

// Database represents a database solution to write block and transaction data
// to. To use another database solution, just implement these methods and
// update the sensor to use the new connection.
type Database interface {
	// WriteBlock will write the both the block and block event to the database
	// if ShouldWriteBlocks and ShouldWriteBlockEvents return true, respectively.
	WriteBlock(context.Context, *enode.Node, *types.Block, *big.Int, time.Time)

	// WriteBlockHeaders will write the block headers if ShouldWriteBlocks
	// returns true.
	WriteBlockHeaders(context.Context, []*types.Header, time.Time, bool)

	// WriteBlockEvents appends an inbound block event (peer, hash, time) for
	// each hash — one per peer we received the announcement from. The caller
	// decides which hashes to pass (every announcement for the full per-peer
	// stream, or just the first-seen ones); the backend only appends.
	WriteBlockEvents(context.Context, *enode.Node, []common.Hash, time.Time)

	// WriteBlockHashFirstSeen records the earliest sighting of a block hash on
	// the block entity itself (Datastore's TimeFirstSeenHash). Backends that
	// derive first-seen from the event stream (e.g. ClickHouse) treat it as a
	// no-op.
	WriteBlockHashFirstSeen(context.Context, *enode.Node, common.Hash, time.Time)

	// WriteBlockBody writes the transactions carried in the block body (the block
	// row itself comes from WriteBlock/WriteBlockHeaders). Backends with a
	// separate transactions table (e.g. ClickHouse) gate this on
	// ShouldWriteTransactions; the Datastore backend gates on ShouldWriteBlocks
	// because it links the transactions and uncles onto the block entity.
	WriteBlockBody(context.Context, *eth.BlockBody, common.Hash, time.Time)

	// WriteTransactions writes the transaction bodies if ShouldWriteTransactions
	// returns true. Transaction events are recorded separately via
	// WriteTransactionEvents at announcement time.
	WriteTransactions(context.Context, *enode.Node, []*types.Transaction, time.Time)

	// WriteTransactionEvents appends an inbound transaction event (peer, hash,
	// time) for each hash — the tx mirror of WriteBlockEvents. The caller
	// decides which hashes to pass (every announcement for the full per-peer
	// stream, or just the first-seen ones); the backend only appends.
	WriteTransactionEvents(context.Context, *enode.Node, []common.Hash, time.Time)

	// WritePeers will write the connected peers to the database.
	WritePeers(context.Context, []*p2p.Peer, time.Time)

	// HasBlock will return whether the block is in the database. If the database
	// client has not been initialized this will always return true.
	HasBlock(context.Context, common.Hash) bool

	MaxConcurrentWrites() int
	ShouldWriteBlocks() bool
	ShouldWriteBlockEvents() bool
	ShouldWriteFirstBlockEvent() bool
	ShouldWriteTransactions() bool
	ShouldWriteTransactionEvents() bool
	ShouldWriteFirstTransactionEvent() bool
	ShouldWritePeers() bool

	// NodeList will return a list of enode URLs.
	NodeList(ctx context.Context, limit int) ([]string, error)

	// Close flushes any buffered writes and releases the underlying database
	// connection. It blocks until in-flight writes have drained (or their
	// shutdown flush times out). Implementations with no resources to release
	// return nil.
	Close() error
}
