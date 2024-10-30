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
	WriteBlock(context.Context, *enode.Node, *types.Block, *big.Int, *time.Time)

	// WriteBlockHeaders will write the block headers if ShouldWriteBlocks
	// returns true.
	WriteBlockHeaders(context.Context, []*types.Header, *time.Time)

	// WriteBlockHashes will write the block hashes if ShouldWriteBlockEvents
	// returns true.
	WriteBlockHashes(context.Context, *enode.Node, []common.Hash, *time.Time)

	// WriteBlockBodies will write the block bodies if ShouldWriteBlocks returns
	// true.
	WriteBlockBody(context.Context, *eth.BlockBody, common.Hash, *time.Time)

	// WriteTransactions will write the both the transaction and transaction
	// event to the database if ShouldWriteTransactions and
	// ShouldWriteTransactionEvents return true, respectively.
	WriteTransactions(context.Context, *enode.Node, []*types.Transaction, *time.Time)

	// WritePeers will write the connected peers to the database.
	WritePeers(context.Context, []*p2p.Peer, *time.Time)

	// HasBlock will return whether the block is in the database. If the database
	// client has not been initialized this will always return true.
	HasBlock(context.Context, common.Hash) bool

	MaxConcurrentWrites() int
	ShouldWriteBlocks() bool
	ShouldWriteBlockEvents() bool
	ShouldWriteTransactions() bool
	ShouldWriteTransactionEvents() bool
	ShouldWritePeers() bool

	// NodeList will return a list of enode URLs.
	NodeList(ctx context.Context, limit int) ([]string, error)
}
