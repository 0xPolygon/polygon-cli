package database

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Database represents a database solution to write block and transaction data
// to. To use another database solution, just implement these methods and
// update the sensor to use the new connection.
type Database interface {
	WriteBlock(*enode.Node, *types.Block)
	WriteBlockHeaders([]*types.Header)
	WriteBlockHashes(*enode.Node, []common.Hash)
	WriteBlockBody(*eth.BlockBody, common.Hash)
	WriteTransactions(*enode.Node, []*types.Transaction)

	MaxConcurrentWrites() int
	ShouldWriteBlocks() bool
	ShouldWriteTransactions() bool
}
