package p2p

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
)

// Conns manages a collection of peer connections for transaction broadcasting.
// It keeps a historical record of all peers that have ever connected.
type Conns struct {
	conns map[string]*conn
	mu    sync.RWMutex
}

// NewConns creates a new connection manager.
func NewConns() *Conns {
	return &Conns{
		conns: make(map[string]*conn),
	}
}

// Add adds a connection to the manager.
// Connections are kept in the map even after they disconnect for historical tracking.
func (c *Conns) Add(cn *conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[cn.node.ID().String()] = cn
	log.Debug().Msg("Added connection for broadcasting")
}

// Remove is a no-op. Connections remain in the map for historical tracking.
// Active connections are identified by checking if the peer is still connected during operations.
func (c *Conns) Remove(cn *conn) {
	// Intentionally empty - keep all connections for historical record
}

// BroadcastTx broadcasts a single transaction to all connected peers.
// Returns the number of peers the transaction was successfully sent to.
func (c *Conns) BroadcastTx(tx *types.Transaction) int {
	return c.BroadcastTxs(types.Transactions{tx})
}

// BroadcastTxs broadcasts multiple transactions to all currently active peers.
// Returns the number of peers the transactions were successfully sent to.
// Silently skips disconnected peers from the historical connection list.
func (c *Conns) BroadcastTxs(txs types.Transactions) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(txs) == 0 {
		return 0
	}

	count := 0
	for _, cn := range c.conns {
		// Try to send, failures indicate disconnected peers
		if err := ethp2p.Send(cn.rw, 0x02, txs); err != nil {
			// Silently skip - this peer is in our history but disconnected
			continue
		}
		count++
	}

	return count
}

// Nodes returns all peer nodes that have ever connected (historical record).
// This includes both currently active and previously disconnected peers.
func (c *Conns) Nodes() []*enode.Node {
	c.mu.RLock()
	defer c.mu.RUnlock()

	nodes := make([]*enode.Node, 0, len(c.conns))
	for _, cn := range c.conns {
		nodes = append(nodes, cn.node)
	}

	return nodes
}
