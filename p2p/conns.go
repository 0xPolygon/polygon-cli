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

// BlockCache stores the actual block data to avoid duplicate fetches and database queries.
type BlockCache struct {
	Header *types.Header
	Body   *eth.BlockBody
	TD     *big.Int
}

// ConnsOptions contains configuration options for creating a new Conns manager.
type ConnsOptions struct {
	BlocksCache                CacheOptions
	TxsCache                   CacheOptions
	KnownTxsCache              CacheOptions
	KnownBlocksCache           CacheOptions
	Head                       eth.NewBlockPacket
	ShouldBroadcastTx          bool
	ShouldBroadcastTxHashes    bool
	ShouldBroadcastBlocks      bool
	ShouldBroadcastBlockHashes bool
}

// Conns manages a collection of active peer connections for transaction broadcasting.
// It also maintains a global cache of blocks written to the database.
type Conns struct {
	conns map[string]*conn
	mu    sync.RWMutex

	// blocks tracks blocks written to the database across all peers
	// to avoid duplicate writes and requests.
	blocks *Cache[common.Hash, BlockCache]

	// txs caches transactions for serving to peers and duplicate detection
	txs *Cache[common.Hash, *types.Transaction]

	// knownTxsOpts and knownBlocksOpts store cache options for per-peer caches
	knownTxsOpts    CacheOptions
	knownBlocksOpts CacheOptions

	// oldest stores the first block the sensor has seen so when fetching
	// parent blocks, it does not request blocks older than this.
	oldest *Locked[*types.Header]

	// head keeps track of the current head block of the chain.
	head *Locked[eth.NewBlockPacket]

	// Broadcast flags control what gets cached and rebroadcasted
	shouldBroadcastTx          bool
	shouldBroadcastTxHashes    bool
	shouldBroadcastBlocks      bool
	shouldBroadcastBlockHashes bool
}

// NewConns creates a new connection manager with a blocks cache.
func NewConns(opts ConnsOptions) *Conns {
	head := &Locked[eth.NewBlockPacket]{}
	head.Set(opts.Head)

	oldest := &Locked[*types.Header]{}
	oldest.Set(opts.Head.Block.Header())

	return &Conns{
		conns:                      make(map[string]*conn),
		blocks:                     NewCache[common.Hash, BlockCache](opts.BlocksCache),
		txs:                        NewCache[common.Hash, *types.Transaction](opts.TxsCache),
		knownTxsOpts:               opts.KnownTxsCache,
		knownBlocksOpts:            opts.KnownBlocksCache,
		oldest:                     oldest,
		head:                       head,
		shouldBroadcastTx:          opts.ShouldBroadcastTx,
		shouldBroadcastTxHashes:    opts.ShouldBroadcastTxHashes,
		shouldBroadcastBlocks:      opts.ShouldBroadcastBlocks,
		shouldBroadcastBlockHashes: opts.ShouldBroadcastBlockHashes,
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

// BroadcastTxs broadcasts multiple transactions to all connected peers,
// filtering out transactions that each peer already knows about.
// Returns the number of peers the transactions were successfully sent to.
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
		cn.countMsgSent(packet.Name(), float64(len(unknownTxs)))
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

		cn.countMsgSent(packet.Name(), float64(len(unknownHashes)))
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

		cn.countMsgSent(packet.Name(), 1)
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

		cn.countMsgSent(packet.Name(), float64(len(unknownHashes)))
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

// PeerConnectedAt returns the connection time for a peer by their ID.
// Returns zero time if the peer is not found.
func (c *Conns) PeerConnectedAt(peerID string) time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cn, ok := c.conns[peerID]; ok {
		return cn.connectedAt
	}

	return time.Time{}
}

// AddTx adds a transaction to the shared cache for duplicate detection and serving.
func (c *Conns) AddTx(hash common.Hash, tx *types.Transaction) {
	c.txs.Add(hash, tx)
}

// GetTx retrieves a transaction from the shared cache.
func (c *Conns) GetTx(hash common.Hash) (*types.Transaction, bool) {
	return c.txs.Get(hash)
}

// Blocks returns the global blocks cache.
func (c *Conns) Blocks() *Cache[common.Hash, BlockCache] {
	return c.blocks
}

// OldestBlock returns the oldest block the sensor will fetch parents for.
// This is set once at initialization to the head block and acts as a floor
// to prevent the sensor from crawling backwards indefinitely.
func (c *Conns) OldestBlock() *types.Header {
	return c.oldest.Get()
}

// HeadBlock returns the current head block packet.
func (c *Conns) HeadBlock() eth.NewBlockPacket {
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

// KnownTxsOpts returns the cache options for per-peer known tx caches.
func (c *Conns) KnownTxsOpts() CacheOptions {
	return c.knownTxsOpts
}

// KnownBlocksOpts returns the cache options for per-peer known block caches.
func (c *Conns) KnownBlocksOpts() CacheOptions {
	return c.knownBlocksOpts
}

// ShouldBroadcastTx returns whether full transaction broadcasting is enabled.
func (c *Conns) ShouldBroadcastTx() bool {
	return c.shouldBroadcastTx
}

// ShouldBroadcastTxHashes returns whether transaction hash broadcasting is enabled.
func (c *Conns) ShouldBroadcastTxHashes() bool {
	return c.shouldBroadcastTxHashes
}

// ShouldBroadcastBlocks returns whether full block broadcasting is enabled.
func (c *Conns) ShouldBroadcastBlocks() bool {
	return c.shouldBroadcastBlocks
}

// ShouldBroadcastBlockHashes returns whether block hash broadcasting is enabled.
func (c *Conns) ShouldBroadcastBlockHashes() bool {
	return c.shouldBroadcastBlockHashes
}

// GetPeerMessages returns a snapshot of message counts for a specific peer.
// Returns nil if the peer is not found.
func (c *Conns) GetPeerMessages(peerID string) *PeerMessages {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cn, ok := c.conns[peerID]; ok {
		msgs := cn.messages.Load()
		return &msgs
	}

	return nil
}

// GetPeerName returns the fullname (client identifier) for a specific peer.
// Returns empty string if the peer is not found.
func (c *Conns) GetPeerName(peerID string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cn, ok := c.conns[peerID]; ok {
		return cn.peer.Fullname()
	}

	return ""
}
