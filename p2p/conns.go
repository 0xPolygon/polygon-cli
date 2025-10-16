package p2p

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Conns manages a collection of active peer connections for transaction broadcasting.
// It also maintains a global cache of blocks written to the database.
type Conns struct {
	conns map[string]*conn
	mu    sync.RWMutex

	// BlocksCache tracks blocks written to the database across all peers
	// to avoid duplicate writes and requests.
	BlocksCache *Cache[common.Hash, BlockWriteState]
}

// NewConns creates a new connection manager with a blocks cache.
func NewConns(maxBlocks int, blocksCacheTTL time.Duration) *Conns {
	return &Conns{
		conns:       make(map[string]*conn),
		BlocksCache: NewCache[common.Hash, BlockWriteState](maxBlocks, blocksCacheTTL),
	}
}

// Add adds a connection to the manager.
func (c *Conns) Add(cn *conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[cn.node.ID().String()] = cn
	cn.logger.Debug().Msg("Added connection")
}

// Remove removes a connection from the manager when a peer disconnects.
func (c *Conns) Remove(cn *conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.conns, cn.node.ID().String())
	cn.logger.Debug().Msg("Removed connection")
}

// BroadcastTx broadcasts a single transaction to all connected peers.
// Returns the number of peers the transaction was successfully sent to.
func (c *Conns) BroadcastTx(tx *types.Transaction) int {
	return c.BroadcastTxs(types.Transactions{tx})
}

// BroadcastTxs broadcasts multiple transactions to all connected peers.
// Returns the number of peers the transactions were successfully sent to.
func (c *Conns) BroadcastTxs(txs types.Transactions) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(txs) == 0 {
		return 0
	}

	count := 0
	for _, cn := range c.conns {
		if err := ethp2p.Send(cn.rw, eth.TransactionsMsg, txs); err != nil {
			continue
		}
		count++
	}

	return count
}

// Nodes returns all currently connected peer nodes.
func (c *Conns) Nodes() []*enode.Node {
	c.mu.RLock()
	defer c.mu.RUnlock()

	nodes := make([]*enode.Node, 0, len(c.conns))
	for _, cn := range c.conns {
		nodes = append(nodes, cn.node)
	}

	return nodes
}
