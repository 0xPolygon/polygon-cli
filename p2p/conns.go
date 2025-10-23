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

	// blocks tracks blocks written to the database across all peers
	// to avoid duplicate writes and requests.
	blocks *Cache[common.Hash, BlockCache]

	// oldest stores the first block the sensor has seen so when fetching
	// parent blocks, it does not request blocks older than this.
	oldest *Locked[*types.Header]

	// head keeps track of the current head block of the chain.
	head *Locked[eth.NewBlockPacket]
}

// ConnsOptions contains configuration options for creating a new Conns manager.
type ConnsOptions struct {
	BlocksCache CacheOptions
	Head        eth.NewBlockPacket
}

// NewConns creates a new connection manager with a blocks cache.
func NewConns(opts ConnsOptions) *Conns {
	head := &Locked[eth.NewBlockPacket]{}
	head.Set(opts.Head)

	oldest := &Locked[*types.Header]{}
	oldest.Set(opts.Head.Block.Header())

	return &Conns{
		conns:  make(map[string]*conn),
		blocks: NewCache[common.Hash, BlockCache](opts.BlocksCache),
		oldest: oldest,
		head:   head,
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

// Blocks returns the global blocks cache.
func (c *Conns) Blocks() *Cache[common.Hash, BlockCache] {
	return c.blocks
}

// GetOldestBlock returns the oldest block the sensor will fetch parents for.
// This is set once at initialization to the head block and acts as a floor
// to prevent the sensor from crawling backwards indefinitely.
func (c *Conns) GetOldestBlock() *types.Header {
	return c.oldest.Get()
}

// GetHeadBlock returns the current head block packet.
func (c *Conns) GetHeadBlock() eth.NewBlockPacket {
	return c.head.Get()
}

// UpdateHeadBlock updates the head block if the provided block is newer.
// Returns true if the head block was updated, false otherwise.
func (c *Conns) UpdateHeadBlock(packet eth.NewBlockPacket) bool {
	return c.head.Update(func(current eth.NewBlockPacket) (eth.NewBlockPacket, bool) {
		if current.Block == nil || (packet.Block.NumberU64() > current.Block.NumberU64() && packet.TD.Cmp(current.TD) == 1) {
			return packet, true
		}
		return current, false
	})
}
