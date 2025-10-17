package p2p

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Conns manages a collection of active peer connections for transaction broadcasting.
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

// GetPeerConnectedAt returns the connection time for a peer by their ID.
// Returns zero time if the peer is not found.
func (c *Conns) GetPeerConnectedAt(peerID string) time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cn, ok := c.conns[peerID]; ok {
		return cn.connectedAt
	}

	return time.Time{}
}

// GetPeerActiveGoroutines returns the number of active goroutines for a peer by their ID.
// Returns 0 if the peer is not found.
func (c *Conns) GetPeerActiveGoroutines(peerID string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cn, ok := c.conns[peerID]; ok {
		return cn.ActiveGoroutines()
	}

	return 0
}
