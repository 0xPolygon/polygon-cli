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
	WriteBlock(context.Context, *enode.Node, *types.Block, *big.Int)
	WriteBlockHeaders(context.Context, []*types.Header)
	WriteBlockHashes(context.Context, *enode.Node, []common.Hash)
	WriteBlockBody(context.Context, *eth.BlockBody, common.Hash)
	WriteTransactions(context.Context, *enode.Node, []*types.Transaction)
	HasParentBlock(context.Context, common.Hash) bool

	MaxConcurrentWrites() int
	ShouldWriteBlocks() bool
	ShouldWriteBlockEvents() bool
	ShouldWriteTransactions() bool
	ShouldWriteTransactionEvents() bool
}
