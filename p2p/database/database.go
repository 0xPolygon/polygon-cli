package database

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Database represents a database solution to write block and transaction data
// to. To use another database solution, just implement these methods and
// update the sensor to use the new connection.
type Database interface {
	// WriteBlock will write the both the block and block event to the database
	// if ShouldWriteBlocks and ShouldWriteBlockEvents return true, respectively.
	WriteBlock(context.Context, *enode.Node, *types.Block, *big.Int)

	// WriteBlockHeaders will write the block headers if ShouldWriteBlocks
	// returns true.
	WriteBlockHeaders(context.Context, []*types.Header)

	// WriteBlockHashes will write the block hashes if ShouldWriteBlockEvents
	// returns true.
	WriteBlockHashes(context.Context, *enode.Node, []common.Hash)

	// WriteBlockBodies will write the block bodies if ShouldWriteBlocks returns
	// true.
	WriteBlockBody(context.Context, *eth.BlockBody, common.Hash)

	// WriteTransactions will write the both the transaction and transaction
	// event to the database if ShouldWriteTransactions and
	// ShouldWriteTransactionEvents return true, respectively.
	WriteTransactions(context.Context, *enode.Node, []*types.Transaction)

	// HasBlock will return whether the block is in the database. If the database
	// client has not been initialized this will always return true.
	HasBlock(context.Context, common.Hash) bool

	MaxConcurrentWrites() int
	ShouldWriteBlocks() bool
	ShouldWriteBlockEvents() bool
	ShouldWriteTransactions() bool
	ShouldWriteTransactionEvents() bool

	// NodeList will return a list of enode URLs.
	NodeList(ctx context.Context, limit int) ([]string, error)
}
