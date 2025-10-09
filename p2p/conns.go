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
	"github.com/rs/zerolog/log"
)

// ConnsOptions contains configuration options for the connection manager.
type ConnsOptions struct {
	MaxCachedTxs               int
	MaxCachedBlocks            int
	CacheTTL                   time.Duration
	ShouldBroadcastTx          bool
	ShouldBroadcastTxHashes    bool
	ShouldBroadcastBlocks      bool
	ShouldBroadcastBlockHashes bool
}

// Conns manages a collection of active peer connections for transaction broadcasting.
type Conns struct {
	conns map[string]*conn
	mu    sync.RWMutex

	// Shared caches for serving broadcast data to peers
	txs    *Cache[common.Hash, *types.Transaction]
	blocks *Cache[common.Hash, *types.Block]

	// Broadcast flags control what gets cached and rebroadcasted
	shouldBroadcastTx          bool
	shouldBroadcastTxHashes    bool
	shouldBroadcastBlocks      bool
	shouldBroadcastBlockHashes bool
}

// NewConns creates a new connection manager with the specified options.
func NewConns(opts ConnsOptions) *Conns {
	// Create caches with configured TTL for data freshness
	txCache := NewCache[common.Hash, *types.Transaction](opts.MaxCachedTxs, opts.CacheTTL)
	blockCache := NewCache[common.Hash, *types.Block](opts.MaxCachedBlocks, opts.CacheTTL)

	return &Conns{
		conns:                      make(map[string]*conn),
		txs:                        txCache,
		blocks:                     blockCache,
		shouldBroadcastTx:          opts.ShouldBroadcastTx,
		shouldBroadcastTxHashes:    opts.ShouldBroadcastTxHashes,
		shouldBroadcastBlocks:      opts.ShouldBroadcastBlocks,
		shouldBroadcastBlockHashes: opts.ShouldBroadcastBlockHashes,
	}
}

// AddConn adds a connection to the manager.
func (c *Conns) AddConn(cn *conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[cn.node.ID().String()] = cn
	cn.logger.Debug().Msg("Added connection")
}

// RemoveConn removes a connection from the manager when a peer disconnects.
func (c *Conns) RemoveConn(cn *conn) {
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
// If broadcast flags are disabled, this is a no-op.
func (c *Conns) BroadcastTxs(txs types.Transactions) int {
	if !c.shouldBroadcastTx {
		return 0
	}

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
// sent to. If broadcast flags are disabled, this is a no-op.
func (c *Conns) BroadcastTxHashes(hashes []common.Hash) int {
	if !c.shouldBroadcastTxHashes {
		return 0
	}

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
// If broadcast flags are disabled, this is a no-op.
func (c *Conns) BroadcastBlock(block *types.Block, td *big.Int) int {
	if !c.shouldBroadcastBlocks {
		return 0
	}

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
// of peers the hashes were successfully sent to. If broadcast flags are disabled, this is a no-op.
func (c *Conns) BroadcastBlockHashes(hashes []common.Hash, numbers []uint64) int {
	if !c.shouldBroadcastBlockHashes {
		return 0
	}

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

// AddTx adds a transaction to the shared cache for duplicate detection and serving.
func (c *Conns) AddTx(hash common.Hash, tx *types.Transaction) {
	c.txs.Add(hash, tx)
}

// AddBlock adds a block to the shared cache for duplicate detection and serving.
func (c *Conns) AddBlock(hash common.Hash, block *types.Block) {
	c.blocks.Add(hash, block)
}

// AddBlockHeader adds a block header to the cache. If a block already exists with a real header, does nothing.
// If a block exists with an empty header (body received first), replaces it with the real header.
// Otherwise creates a new block with just the header.
func (c *Conns) AddBlockHeader(header *types.Header) {
	hash := header.Hash()

	// Check if block already exists in cache
	block, ok := c.blocks.Get(hash)
	if !ok {
		// No block exists, create new one with header only
		c.AddBlock(hash, types.NewBlockWithHeader(header))
		return
	}

	// Check if existing block has a real header already
	if block.Number() != nil && block.Number().Uint64() > 0 {
		// Block already has a real header, don't overwrite
		return
	}

	// Block has empty header (body came first), replace with real header + keep body
	b := types.NewBlockWithHeader(header).WithBody(types.Body{
		Transactions: block.Transactions(),
		Uncles:       block.Uncles(),
		Withdrawals:  block.Withdrawals(),
	})
	c.AddBlock(hash, b)
}

// AddBlockBody adds a body to an existing block in the cache. If no block exists for this hash,
// creates a block with an empty header and the body. If a block exists with only a header, updates it with the body.
func (c *Conns) AddBlockBody(hash common.Hash, body *eth.BlockBody) {
	// Get existing block from cache
	block, ok := c.blocks.Get(hash)
	if !ok {
		// No header yet, create block with empty header and body
		blockWithBody := types.NewBlockWithHeader(&types.Header{}).WithBody(types.Body(*body))
		c.AddBlock(hash, blockWithBody)
		return
	}

	// Check if block already has a body
	if len(block.Transactions()) > 0 || len(block.Uncles()) > 0 || len(block.Withdrawals()) > 0 {
		// Block already has a body, no need to update
		return
	}

	// Reconstruct full block with existing header and body
	c.AddBlock(hash, block.WithBody(types.Body(*body)))
}

// GetTx retrieves a transaction from the shared cache.
func (c *Conns) GetTx(hash common.Hash) (*types.Transaction, bool) {
	return c.txs.Get(hash)
}

// GetBlock retrieves a block from the shared cache.
func (c *Conns) GetBlock(hash common.Hash) (*types.Block, bool) {
	return c.blocks.Get(hash)
}

// HasBlockHeader checks if we have at least a header for a block in the cache.
// Returns true if we have a block with a real header (number > 0).
func (c *Conns) HasBlockHeader(hash common.Hash) bool {
	block, ok := c.blocks.Get(hash)
	if !ok {
		return false
	}

	// Check if block has a real header (not empty)
	return block.Number() != nil && block.Number().Uint64() > 0
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
