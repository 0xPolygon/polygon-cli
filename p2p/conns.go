package p2p

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/rs/zerolog/log"
)

const (
	// maxCachedTxs is the maximum number of transactions to cache for serving to peers
	maxCachedTxs = 10000

	// maxCachedBlocks is the maximum number of blocks to cache for serving to peers
	maxCachedBlocks = 1000
)

// Conns manages a collection of active peer connections for transaction broadcasting.
type Conns struct {
	conns map[string]*conn
	mu    sync.RWMutex

	// Shared LRU caches for serving broadcast data to peers
	txs    *lru.Cache[common.Hash, *types.Transaction]
	blocks *lru.Cache[common.Hash, *types.Block]
}

// NewConns creates a new connection manager.
func NewConns() *Conns {
	txCache, err := lru.New[common.Hash, *types.Transaction](maxCachedTxs)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create transaction cache")
	}

	blockCache, err := lru.New[common.Hash, *types.Block](maxCachedBlocks)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create block cache")
	}

	return &Conns{
		conns:  make(map[string]*conn),
		txs:    txCache,
		blocks: blockCache,
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

// BroadcastTx broadcasts a single transaction to all connected peers and
// returns the number of peers the transaction was successfully sent to.
func (c *Conns) BroadcastTx(tx *types.Transaction) int {
	return c.BroadcastTxs(types.Transactions{tx})
}

// BroadcastTxs broadcasts multiple transactions to all connected peers,
// filtering out transactions that each peer already knows about, and returns
// the number of peers the transactions were successfully sent to.
func (c *Conns) BroadcastTxs(txs types.Transactions) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(txs) == 0 {
		return 0
	}

	count := 0
	for _, cn := range c.conns {
		// Filter transactions this peer doesn't know about
		unknownTxs := make(types.Transactions, 0, len(txs))
		for _, tx := range txs {
			if !cn.hasKnownTx(tx.Hash()) {
				unknownTxs = append(unknownTxs, tx)
			}
		}

		if len(unknownTxs) == 0 {
			continue
		}

		// Send as TransactionsPacket
		packet := eth.TransactionsPacket(unknownTxs)
		cn.AddCountSent(packet.Name(), 1)
		if err := ethp2p.Send(cn.rw, eth.TransactionsMsg, packet); err != nil {
			cn.logger.Debug().
				Err(err).
				Msg("Failed to send transactions")
			continue
		}

		// Mark transactions as known for this peer
		for _, tx := range unknownTxs {
			cn.addKnownTx(tx.Hash())
		}

		count++
	}

	if count > 0 {
		log.Debug().
			Int("peers", count).
			Int("txs", len(txs)).
			Msg("Broadcasted transactions")
	}

	return count
}

// BroadcastTxHashes broadcasts transaction hashes to peers that don't already
// know about them and returns the number of peers the hashes were successfully
// sent to.
func (c *Conns) BroadcastTxHashes(hashes []common.Hash) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(hashes) == 0 {
		return 0
	}

	count := 0
	for _, cn := range c.conns {
		// Filter hashes this peer doesn't know about
		unknownHashes := make([]common.Hash, 0, len(hashes))
		for _, hash := range hashes {
			if !cn.hasKnownTx(hash) {
				unknownHashes = append(unknownHashes, hash)
			}
		}

		if len(unknownHashes) == 0 {
			continue
		}

		// Send NewPooledTransactionHashesPacket
		packet := eth.NewPooledTransactionHashesPacket{
			Types:  make([]byte, len(unknownHashes)),
			Sizes:  make([]uint32, len(unknownHashes)),
			Hashes: unknownHashes,
		}

		cn.AddCountSent(packet.Name(), 1)
		if err := ethp2p.Send(cn.rw, eth.NewPooledTransactionHashesMsg, packet); err != nil {
			cn.logger.Debug().
				Err(err).
				Msg("Failed to send transaction hashes")
			continue
		}

		// Mark hashes as known for this peer
		for _, hash := range unknownHashes {
			cn.addKnownTx(hash)
		}

		count++
	}

	if count > 0 {
		log.Debug().
			Int("peers", count).
			Int("hashes", len(hashes)).
			Msg("Broadcasted transaction hashes")
	}

	return count
}

// BroadcastBlock broadcasts a full block to peers that don't already know
// about it and returns the number of peers the block was successfully sent to.
func (c *Conns) BroadcastBlock(block *types.Block, td *big.Int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if block == nil {
		return 0
	}

	hash := block.Hash()
	count := 0

	for _, cn := range c.conns {
		// Skip if peer already knows about this block
		if cn.hasKnownBlock(hash) {
			continue
		}

		// Send NewBlockPacket
		packet := eth.NewBlockPacket{
			Block: block,
			TD:    td,
		}

		cn.AddCountSent(packet.Name(), 1)
		if err := ethp2p.Send(cn.rw, eth.NewBlockMsg, &packet); err != nil {
			cn.logger.Debug().
				Err(err).
				Uint64("number", block.Number().Uint64()).
				Msg("Failed to send block")
			continue
		}

		// Mark block as known for this peer
		cn.addKnownBlock(hash)
		count++
	}

	if count > 0 {
		log.Debug().
			Int("peers", count).
			Uint64("number", block.NumberU64()).
			Msg("Broadcasted block")
	}

	return count
}

// BroadcastBlockHashes broadcasts block hashes with their corresponding block
// numbers to peers that don't already know about them and returns the number
// of peers the hashes were successfully sent to.
func (c *Conns) BroadcastBlockHashes(hashes []common.Hash, numbers []uint64) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(hashes) == 0 || len(hashes) != len(numbers) {
		return 0
	}

	count := 0

	for _, cn := range c.conns {
		// Filter hashes this peer doesn't know about
		unknownHashes := make([]common.Hash, 0, len(hashes))
		unknownNumbers := make([]uint64, 0, len(numbers))

		for i, hash := range hashes {
			if !cn.hasKnownBlock(hash) {
				unknownHashes = append(unknownHashes, hash)
				unknownNumbers = append(unknownNumbers, numbers[i])
			}
		}

		if len(unknownHashes) == 0 {
			continue
		}

		// Send NewBlockHashesPacket
		packet := make(eth.NewBlockHashesPacket, len(unknownHashes))
		for i := range unknownHashes {
			packet[i].Hash = unknownHashes[i]
			packet[i].Number = unknownNumbers[i]
		}

		cn.AddCountSent(packet.Name(), 1)
		if err := ethp2p.Send(cn.rw, eth.NewBlockHashesMsg, packet); err != nil {
			cn.logger.Debug().
				Err(err).
				Msg("Failed to send block hashes")
			continue
		}

		// Mark hashes as known for this peer
		for _, hash := range unknownHashes {
			cn.addKnownBlock(hash)
		}

		count++
	}

	if count > 0 {
		log.Debug().
			Int("peers", count).
			Int("hashes", len(hashes)).
			Msg("Broadcasted block hashes")
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

// GetPeerConnectedAt returns the time when a peer connected, or zero time if not found.
func (c *Conns) GetPeerConnectedAt(url string) time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cn, ok := c.conns[url]; ok {
		return cn.connectedAt
	}

	return time.Time{}
}
