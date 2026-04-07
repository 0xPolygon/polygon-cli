package p2p

import (
	"bytes"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"

	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

// BlockCache stores the actual block data to avoid duplicate fetches and database queries.
type BlockCache struct {
	Header *types.Header
	Body   *eth.BlockBody
	TD     *big.Int
}

// ConnsOptions contains configuration options for creating a new Conns manager.
type ConnsOptions struct {
	BlocksCache                ds.LRUOptions
	TxsCache                   ds.LRUOptions
	KnownTxsBloom              ds.BloomSetOptions
	KnownBlocksMax             int
	Head                       NewBlockPacket
	ShouldBroadcastTx          bool
	ShouldBroadcastTxHashes    bool
	ShouldBroadcastBlocks      bool
	ShouldBroadcastBlockHashes bool
	BroadcastWorkers           int
	TxBatchTimeout             time.Duration
	TxBroadcastQueueSize       int
}

// Conns manages a collection of active peer connections for transaction broadcasting.
// It also maintains a global cache of blocks written to the database.
type Conns struct {
	conns map[string]*conn
	mu    sync.RWMutex

	// blocks tracks blocks written to the database across all peers
	// to avoid duplicate writes and requests.
	blocks *ds.LRU[common.Hash, BlockCache]

	// txs caches transactions for serving to peers and duplicate detection
	txs *ds.LRU[common.Hash, *types.Transaction]

	// knownTxsOpts stores bloom filter options for per-peer known tx tracking
	knownTxsOpts ds.BloomSetOptions
	// knownBlocksMax stores the maximum size for per-peer known block caches
	knownBlocksMax int

	// oldest stores the first block the sensor has seen so when fetching
	// parent blocks, it does not request blocks older than this.
	oldest *ds.Locked[*types.Header]

	// head keeps track of the current head block of the chain.
	head *ds.Locked[NewBlockPacket]

	// Broadcast flags control what gets cached and rebroadcasted
	shouldBroadcastTx          bool
	shouldBroadcastTxHashes    bool
	shouldBroadcastBlocks      bool
	shouldBroadcastBlockHashes bool

	// txBroadcastCh is a buffered channel for async transaction broadcast
	txBroadcastCh chan types.Transactions

	// txBatchTimeout is the timeout for batching transactions before broadcast
	txBatchTimeout time.Duration
}

// NewConns creates a new connection manager with a blocks cache.
func NewConns(opts ConnsOptions) *Conns {
	head := &ds.Locked[NewBlockPacket]{}
	head.Set(opts.Head)

	oldest := &ds.Locked[*types.Header]{}
	oldest.Set(opts.Head.Block.Header())

	txBatchTimeout := opts.TxBatchTimeout
	if txBatchTimeout <= 0 {
		txBatchTimeout = 500 * time.Millisecond
	}

	txBroadcastQueueSize := opts.TxBroadcastQueueSize
	if txBroadcastQueueSize <= 0 {
		txBroadcastQueueSize = 100000
	}

	c := &Conns{
		conns:                      make(map[string]*conn),
		blocks:                     ds.NewLRU[common.Hash, BlockCache](opts.BlocksCache),
		txs:                        ds.NewLRU[common.Hash, *types.Transaction](opts.TxsCache),
		knownTxsOpts:               opts.KnownTxsBloom,
		knownBlocksMax:             opts.KnownBlocksMax,
		oldest:                     oldest,
		head:                       head,
		shouldBroadcastTx:          opts.ShouldBroadcastTx,
		shouldBroadcastTxHashes:    opts.ShouldBroadcastTxHashes,
		shouldBroadcastBlocks:      opts.ShouldBroadcastBlocks,
		shouldBroadcastBlockHashes: opts.ShouldBroadcastBlockHashes,
		txBroadcastCh:              make(chan types.Transactions, txBroadcastQueueSize),
		txBatchTimeout:             txBatchTimeout,
	}

	workers := opts.BroadcastWorkers
	if workers <= 0 {
		workers = 4
	}
	for i := 0; i < workers; i++ {
		go c.txBroadcastLoop()
	}

	return c
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

// snapshotPeers returns a copy of current peer connections.
// The caller can safely iterate without holding the lock.
func (c *Conns) snapshotPeers() []*conn {
	c.mu.RLock()
	defer c.mu.RUnlock()

	peers := make([]*conn, 0, len(c.conns))
	for _, cn := range c.conns {
		peers = append(peers, cn)
	}
	return peers
}

// BroadcastTx broadcasts a single transaction to all connected peers.
// Returns the number of peers the transaction was successfully sent to.
func (c *Conns) BroadcastTx(tx *types.Transaction) int {
	return c.BroadcastTxs(types.Transactions{tx})
}

// BroadcastTxs rebroadcasts transactions received from peers to all connected peers,
// filtering out transactions that each peer already knows about.
// Returns the number of peers the transactions were successfully sent to.
// If broadcast flags are disabled, this is a no-op.
// Note: For RPC-submitted transactions, use EnqueueTxBroadcast instead.
func (c *Conns) BroadcastTxs(txs types.Transactions) int {
	if !c.shouldBroadcastTx || len(txs) == 0 {
		return 0
	}

	// Pre-compute transaction hashes once to avoid redundant Keccak256 computations
	hashes := make([]common.Hash, len(txs))
	txByHash := make(map[common.Hash]*types.Transaction, len(txs))
	for i, tx := range txs {
		h := tx.Hash()
		hashes[i] = h
		txByHash[h] = tx
	}

	peers := c.snapshotPeers()
	if len(peers) == 0 {
		return 0
	}

	// Broadcast concurrently to all peers
	var count atomic.Int32
	var wg sync.WaitGroup

	for _, peer := range peers {
		wg.Go(func() {
			// Filter transactions this peer doesn't know about using batch bloom operation
			unknownHashes := peer.filterUnknownTxHashes(hashes)
			if len(unknownHashes) == 0 {
				return
			}

			// Build transaction list from pre-computed map
			unknownTxs := make(types.Transactions, 0, len(unknownHashes))
			for _, h := range unknownHashes {
				if tx, ok := txByHash[h]; ok {
					unknownTxs = append(unknownTxs, tx)
				}
			}
			if len(unknownTxs) == 0 {
				return
			}

			rawList, err := rlp.EncodeToRawList([]*types.Transaction(unknownTxs))
			if err != nil {
				peer.logger.Debug().Err(err).Msg("Failed to encode transactions")
				return
			}
			packet := &eth.TransactionsPacket{RawList: rawList}
			peer.countMsgSent(packet.Name(), float64(len(unknownTxs)))
			if err := ethp2p.Send(peer.rw, eth.TransactionsMsg, packet); err != nil {
				peer.logger.Debug().Err(err).Msg("Failed to send transactions")
				return
			}

			peer.addKnownTxHashes(unknownHashes)
			count.Add(1)
		})
	}

	wg.Wait()

	finalCount := int(count.Load())
	if finalCount > 0 {
		log.Debug().
			Int("peers", finalCount).
			Int("txs", len(txs)).
			Msg("Broadcasted transactions")
	}

	return finalCount
}

// broadcastTxs sends RPC-submitted transactions to all peers via TransactionsMsg.
// Used by txBroadcastLoop to process transactions from EnqueueTxBroadcast.
func (c *Conns) broadcastTxs(txs types.Transactions, hashes []common.Hash, peers []*conn) {
	rawList, err := rlp.EncodeToRawList([]*types.Transaction(txs))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to encode transactions")
		return
	}
	packet := &eth.TransactionsPacket{RawList: rawList}

	// Pre-encode the entire message once to avoid re-encoding for each peer.
	encodedMsg, err := rlp.EncodeToBytes(packet)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to encode message")
		return
	}
	msgSize := uint32(len(encodedMsg))

	// Cache loop-invariant values
	msgName := packet.Name()
	txCount := float64(len(txs))
	rebroadcasting := c.shouldBroadcastTx || c.shouldBroadcastTxHashes

	for _, peer := range peers {
		go func(peer *conn) {
			peer.countMsgSent(msgName, txCount)

			// Use WriteMsg directly with pre-encoded bytes instead of Send()
			msg := ethp2p.Msg{
				Code:    eth.TransactionsMsg,
				Size:    msgSize,
				Payload: bytes.NewReader(encodedMsg),
			}
			if err := peer.rw.WriteMsg(msg); err != nil {
				peer.logger.Debug().Err(err).Msg("Failed to send transactions")
				return
			}

			// Only track known hashes if rebroadcasting is enabled
			if rebroadcasting {
				peer.addKnownTxHashes(hashes)
			}
		}(peer)
	}
}

// txBroadcastLoop is the worker that drains the broadcast channel and sends
// transactions to all peers, batching up to max packet size per iteration.
func (c *Conns) txBroadcastLoop() {
	for {
		batch, hashes := c.pullTxBatch()
		if batch == nil {
			return
		}

		peers := c.snapshotPeers()
		if len(peers) == 0 {
			continue
		}

		c.broadcastTxs(batch, hashes, peers)

		log.Info().
			Int("txs", len(batch)).
			Int("peers", len(peers)).
			Msg("Broadcasted transaction batch")
	}
}

// pullTxBatch pulls transactions from the channel up to maxTxPacketSize.
// Returns nil, nil when the channel is closed.
func (c *Conns) pullTxBatch() (types.Transactions, []common.Hash) {
	// Block for first item
	txs, ok := <-c.txBroadcastCh
	if !ok {
		return nil, nil
	}

	// Pre-allocate with initial capacity
	batch := make(types.Transactions, 0, len(txs))
	hashes := make([]common.Hash, 0, len(txs))
	var batchSize uint64

	// Add transactions from first receive
	for _, tx := range txs {
		batch = append(batch, tx)
		hashes = append(hashes, tx.Hash())
		batchSize += tx.Size()
	}

	// Drain more until max size or timeout
	timer := time.NewTimer(c.txBatchTimeout)
	defer timer.Stop()

	for batchSize < maxTxPacketSize {
		select {
		case txs, ok := <-c.txBroadcastCh:
			if !ok {
				return batch, hashes
			}
			for _, tx := range txs {
				batch = append(batch, tx)
				hashes = append(hashes, tx.Hash())
				batchSize += tx.Size()
			}
		case <-timer.C:
			return batch, hashes
		}
	}

	return batch, hashes
}

// EnqueueTxBroadcast adds transactions to the broadcast channel for async sending.
func (c *Conns) EnqueueTxBroadcast(txs types.Transactions) {
	c.txBroadcastCh <- txs
}

// Close stops the broadcast worker.
func (c *Conns) Close() {
	close(c.txBroadcastCh)
}

// BroadcastTxHashes enqueues transaction hashes to per-peer broadcast queues.
// Each peer has a dedicated goroutine that drains the queue and batches sends.
// Returns the number of peers the hashes were enqueued to.
// If broadcast flags are disabled, this is a no-op.
func (c *Conns) BroadcastTxHashes(hashes []common.Hash) int {
	if !c.shouldBroadcastTxHashes || len(hashes) == 0 {
		return 0
	}

	// Copy peers to avoid holding lock during sends
	c.mu.RLock()
	peers := make([]*conn, 0, len(c.conns))
	for _, cn := range c.conns {
		if cn.txAnnounce != nil {
			peers = append(peers, cn)
		}
	}
	c.mu.RUnlock()

	count := 0
	for _, cn := range peers {
		// Block until announcement loop is ready or peer closes (matches Bor)
		select {
		case cn.txAnnounce <- hashes:
			count++
		case <-cn.closeCh:
			// Peer closing, skip
		}
	}

	return count
}

// BroadcastBlock broadcasts a full block to peers that don't already know
// about it and returns the number of peers the block was successfully sent to.
// If broadcast flags are disabled, this is a no-op.
func (c *Conns) BroadcastBlock(block *types.Block, td *big.Int) int {
	if !c.shouldBroadcastBlocks || block == nil {
		return 0
	}

	hash := block.Hash()

	peers := c.snapshotPeers()
	if len(peers) == 0 {
		return 0
	}

	// Broadcast concurrently to all peers
	var count atomic.Int32
	var wg sync.WaitGroup

	for _, peer := range peers {
		wg.Go(func() {
			if peer.hasKnownBlock(hash) {
				return
			}

			packet := NewBlockPacket{Block: block, TD: td}
			peer.countMsgSent(packet.Name(), 1)
			if err := ethp2p.Send(peer.rw, eth.NewBlockMsg, &packet); err != nil {
				peer.logger.Debug().Err(err).Uint64("number", block.Number().Uint64()).Msg("Failed to send block")
				return
			}

			peer.addKnownBlock(hash)
			count.Add(1)
		})
	}

	wg.Wait()

	finalCount := int(count.Load())
	if finalCount > 0 {
		log.Debug().
			Int("peers", finalCount).
			Uint64("number", block.NumberU64()).
			Msg("Broadcasted block")
	}

	return finalCount
}

// BroadcastBlockHashes enqueues block hashes to per-peer broadcast queues.
// Each peer has a dedicated goroutine that drains the queue and sends.
// Returns the number of peers the hashes were enqueued to.
// If broadcast flags are disabled, this is a no-op.
func (c *Conns) BroadcastBlockHashes(hashes []common.Hash, numbers []uint64) int {
	if !c.shouldBroadcastBlockHashes || len(hashes) == 0 || len(hashes) != len(numbers) {
		return 0
	}

	// Build packet once, share across all peers
	packet := make(NewBlockHashesPacket, len(hashes))
	for i := range hashes {
		packet[i].Hash = hashes[i]
		packet[i].Number = numbers[i]
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	count := 0
	for _, cn := range c.conns {
		if cn.blockAnnounce == nil {
			continue
		}
		// Non-blocking send, drop if queue full (matches Bor)
		select {
		case cn.blockAnnounce <- packet:
			count++
		default:
		}
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

// AddTxs adds multiple transactions to the shared cache in a single lock operation.
// Returns the computed hashes for reuse by the caller.
func (c *Conns) AddTxs(txs []*types.Transaction) []common.Hash {
	if len(txs) == 0 {
		return nil
	}
	hashes := make([]common.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = tx.Hash()
	}
	c.txs.AddBatch(hashes, txs)
	return hashes
}

// GetTx retrieves a transaction from the shared cache and updates LRU ordering.
func (c *Conns) GetTx(hash common.Hash) (*types.Transaction, bool) {
	return c.txs.Get(hash)
}

// PeekTxs retrieves multiple transactions from the shared cache without updating LRU ordering.
// Uses a single read lock for better concurrency when LRU ordering is not needed.
func (c *Conns) PeekTxs(hashes []common.Hash) []*types.Transaction {
	return c.txs.PeekMany(hashes)
}

// PeekTxsWithHashes retrieves multiple transactions with their hashes from the cache.
// Returns parallel slices of found hashes and transactions. Uses a single read lock.
func (c *Conns) PeekTxsWithHashes(hashes []common.Hash) ([]common.Hash, []*types.Transaction) {
	return c.txs.PeekManyWithKeys(hashes)
}

// Blocks returns the global blocks cache.
func (c *Conns) Blocks() *ds.LRU[common.Hash, BlockCache] {
	return c.blocks
}

// OldestBlock returns the oldest block the sensor will fetch parents for.
// This is set once at initialization to the head block and acts as a floor
// to prevent the sensor from crawling backwards indefinitely.
func (c *Conns) OldestBlock() *types.Header {
	return c.oldest.Get()
}

// HeadBlock returns the current head block packet.
func (c *Conns) HeadBlock() NewBlockPacket {
	return c.head.Get()
}

// UpdateHeadBlock updates the head block if the provided block is newer.
// Returns true if the head block was updated, false otherwise.
func (c *Conns) UpdateHeadBlock(packet NewBlockPacket) bool {
	return c.head.Update(func(current NewBlockPacket) (NewBlockPacket, bool) {
		if current.Block == nil || (packet.Block.NumberU64() > current.Block.NumberU64() && packet.TD.Cmp(current.TD) == 1) {
			return packet, true
		}
		return current, false
	})
}

// KnownTxsOpts returns the bloom filter options for per-peer known tx tracking.
func (c *Conns) KnownTxsOpts() ds.BloomSetOptions {
	return c.knownTxsOpts
}

// KnownBlocksMax returns the maximum size for per-peer known block caches.
func (c *Conns) KnownBlocksMax() int {
	return c.knownBlocksMax
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

// GetBlockByNumber iterates through the cache to find a block by its number.
// Returns the hash, block cache, and true if found; empty values and false otherwise.
func (c *Conns) GetBlockByNumber(number uint64) (common.Hash, BlockCache, bool) {
	for _, hash := range c.blocks.Keys() {
		if cache, ok := c.blocks.Peek(hash); ok && cache.Header != nil {
			if cache.Header.Number.Uint64() == number {
				return hash, cache, true
			}
		}
	}
	return common.Hash{}, BlockCache{}, false
}

// GetPeerVersion returns the negotiated eth protocol version for a specific peer.
// Returns 0 if the peer is not found.
func (c *Conns) GetPeerVersion(peerID string) uint {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cn, ok := c.conns[peerID]; ok {
		return cn.version
	}

	return 0
}

// GetPeerLatestBlock returns the latest block hash and number for a peer.
// Returns zero hash and 0 if the peer is not found or no block has been received.
func (c *Conns) GetPeerLatestBlock(peerID string) (common.Hash, uint64) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cn, ok := c.conns[peerID]; ok {
		info := cn.latestBlock.Get()
		return info.Hash, info.Number
	}

	return common.Hash{}, 0
}
