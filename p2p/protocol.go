package p2p

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0xPolygon/polygon-cli/p2p/database"
	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

const (
	// maxTxPacketSize is the target size for transaction announcement packets.
	// Matches Bor's limit of 100KB.
	maxTxPacketSize = 100 * 1024

	// maxQueuedTxAnns is the maximum number of transaction announcements to
	// queue before dropping oldest. Matches Bor.
	maxQueuedTxAnns = 4096

	// maxQueuedBlockAnns is the maximum number of block announcements to queue
	// before dropping. Matches Bor.
	maxQueuedBlockAnns = 4
)

// protocolLengths maps protocol versions to their message counts.
var protocolLengths = map[uint]uint64{
	66: 17,
	67: 17,
	68: 17,
	69: 18,
}

// conn represents an individual connection with a peer.
type conn struct {
	sensorID string
	node     *enode.Node
	logger   zerolog.Logger
	rw       ethp2p.MsgReadWriter
	db       database.Database
	peer     *ethp2p.Peer

	// requests is used to store the request ID and the block hash. This is used
	// when fetching block bodies because the eth protocol block bodies do not
	// contain information about the block hash.
	requests   *ds.LRU[uint64, common.Hash]
	requestNum uint64

	// parents tracks hashes of blocks requested as parents to mark them
	// with IsParent=true when writing to the database.
	parents *ds.LRU[common.Hash, struct{}]

	// conns provides access to the global connection manager, which includes
	// the blocks cache shared across all peers.
	conns *Conns

	// connectedAt stores when this peer connection was established.
	connectedAt time.Time

	// peerURL is cached to avoid repeated URLv4() calls.
	peerURL string

	// Broadcast flags control what gets rebroadcasted to other peers
	shouldBroadcastTx          bool
	shouldBroadcastTxHashes    bool
	shouldBroadcastBlocks      bool
	shouldBroadcastBlockHashes bool

	// Known caches track what this peer has seen to avoid redundant sends.
	// knownTxs uses a bloom filter for memory efficiency (~40KB vs ~4MB per peer).
	// knownBlocks uses a simple bounded set for lower memory overhead than the generic LRU.
	knownTxs    *ds.BloomSet
	knownBlocks *ds.BoundedSet[common.Hash]

	// messages tracks per-peer message counts for API visibility.
	messages *PeerMessages

	// Broadcast queues for per-peer rate limiting. These decouple message
	// reception from broadcasting to prevent flooding peers with immediate
	// broadcasts.
	txAnnounce    chan []common.Hash
	blockAnnounce chan NewBlockHashesPacket
	closeCh       chan struct{}

	// version stores the negotiated eth protocol version (e.g., 68 or 69).
	version uint
}

// EthProtocolOptions is the options used when creating a new eth protocol.
type EthProtocolOptions struct {
	Context     context.Context
	Database    database.Database
	GenesisHash common.Hash
	RPC         string
	SensorID    string
	NetworkID   uint64
	Conns       *Conns
	ForkID      forkid.ID

	// Cache configurations
	RequestsCache ds.LRUOptions
	ParentsCache  ds.LRUOptions

	// Broadcast flags control what gets rebroadcasted to other peers
	ShouldBroadcastTx          bool
	ShouldBroadcastTxHashes    bool
	ShouldBroadcastBlocks      bool
	ShouldBroadcastBlockHashes bool
}

// NewEthProtocol creates the new eth protocol. This will handle writing the
// status exchange, message handling, and writing blocks/txs to the database.
func NewEthProtocol(version uint, opts EthProtocolOptions) ethp2p.Protocol {
	return ethp2p.Protocol{
		Name:    "eth",
		Version: version,
		Length:  protocolLengths[version],
		Run: func(p *ethp2p.Peer, rw ethp2p.MsgReadWriter) error {
			peerURL := p.Node().URLv4()
			c := &conn{
				sensorID:                   opts.SensorID,
				node:                       p.Node(),
				logger:                     log.With().Str("peer", peerURL).Logger(),
				rw:                         rw,
				db:                         opts.Database,
				requests:                   ds.NewLRU[uint64, common.Hash](opts.RequestsCache),
				requestNum:                 0,
				parents:                    ds.NewLRU[common.Hash, struct{}](opts.ParentsCache),
				peer:                       p,
				conns:                      opts.Conns,
				connectedAt:                time.Now(),
				peerURL:                    peerURL,
				shouldBroadcastTx:          opts.ShouldBroadcastTx,
				shouldBroadcastTxHashes:    opts.ShouldBroadcastTxHashes,
				shouldBroadcastBlocks:      opts.ShouldBroadcastBlocks,
				shouldBroadcastBlockHashes: opts.ShouldBroadcastBlockHashes,
				knownTxs:                   ds.NewBloomSet(opts.Conns.KnownTxsOpts()),
				knownBlocks:                ds.NewBoundedSet[common.Hash](opts.Conns.KnownBlocksMax()),
				messages:                   NewPeerMessages(),
				txAnnounce:                 make(chan []common.Hash),
				blockAnnounce:              make(chan NewBlockHashesPacket, maxQueuedBlockAnns),
				closeCh:                    make(chan struct{}),
				version:                    version,
			}

			// Ensure cleanup happens on any exit path (including statusExchange failure)
			defer func() {
				close(c.closeCh)
				opts.Conns.Remove(c)
			}()

			// Start broadcast loops for per-peer queued broadcasting
			if opts.ShouldBroadcastTxHashes {
				go c.txAnnouncementLoop()
			}
			if opts.ShouldBroadcastBlockHashes {
				go c.blockAnnouncementLoop()
			}

			if err := c.statusExchange(version, opts); err != nil {
				return err
			}

			// Update logger with peer name now that status exchange is complete
			c.logger = log.With().Str("peer", peerURL).Str("peer_name", c.peer.Fullname()).Logger()

			// Send the connection object to the conns manager for RPC broadcasting
			opts.Conns.Add(c)

			ctx := opts.Context

			// Handle all the of the messages here.
			for {
				msg, err := rw.ReadMsg()
				if err != nil {
					return err
				}

				switch msg.Code {
				case eth.NewBlockHashesMsg:
					err = c.handleNewBlockHashes(ctx, msg)
				case eth.TransactionsMsg:
					err = c.handleTransactions(ctx, msg)
				case eth.GetBlockHeadersMsg:
					err = c.handleGetBlockHeaders(msg)
				case eth.BlockHeadersMsg:
					err = c.handleBlockHeaders(ctx, msg)
				case eth.GetBlockBodiesMsg:
					err = c.handleGetBlockBodies(msg)
				case eth.BlockBodiesMsg:
					err = c.handleBlockBodies(ctx, msg)
				case eth.NewBlockMsg:
					err = c.handleNewBlock(ctx, msg)
				case eth.NewPooledTransactionHashesMsg:
					err = c.handleNewPooledTransactionHashes(version, msg)
				case eth.GetPooledTransactionsMsg:
					err = c.handleGetPooledTransactions(msg)
				case eth.PooledTransactionsMsg:
					err = c.handlePooledTransactions(ctx, msg)
				case eth.GetReceiptsMsg:
					err = c.handleGetReceipts(msg)
				case eth.BlockRangeUpdateMsg:
					err = c.handleBlockRangeUpdate(msg)
				default:
					c.logger.Trace().Interface("msg", msg).Send()
				}

				// All the handler functions are built in a way where returning an error
				// should drop the connection. If the connection shouldn't be dropped,
				// then return nil and log the error instead.
				if err != nil {
					c.logger.Error().Err(err).Send()
					return err
				}

				if err = msg.Discard(); err != nil {
					return err
				}
			}
		},
	}
}

// statusExchange performs the eth protocol handshake, using the appropriate
// status packet format based on the negotiated protocol version.
func (c *conn) statusExchange(version uint, opts EthProtocolOptions) error {
	head := c.conns.HeadBlock()

	if version >= eth.ETH69 {
		status := BorStatusPacket69{
			ProtocolVersion: uint32(version),
			NetworkID:       opts.NetworkID,
			TD:              head.TD,
			Genesis:         opts.GenesisHash,
			ForkID:          opts.ForkID,
			EarliestBlock:   head.Block.NumberU64(),
			LatestBlock:     head.Block.NumberU64(),
			LatestBlockHash: head.Block.Hash(),
		}

		return c.statusExchange69(&status)
	}

	status := StatusPacket68{
		ProtocolVersion: uint32(version),
		NetworkID:       opts.NetworkID,
		Genesis:         opts.GenesisHash,
		ForkID:          opts.ForkID,
		Head:            head.Block.Hash(),
		TD:              head.TD,
	}

	return c.statusExchange68(&status)
}

// statusExchange68 will exchange status message for ETH68 and below.
func (c *conn) statusExchange68(packet *StatusPacket68) error {
	errc := make(chan error, 2)

	go func() {
		c.countMsgSent((&StatusPacket68{}).Name(), 1)
		errc <- ethp2p.Send(c.rw, eth.StatusMsg, packet)
	}()

	go func() {
		errc <- c.readStatus68(packet)
	}()

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	for range 2 {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return ethp2p.DiscReadTimeout
		}
	}

	return nil
}

// statusExchange69 will exchange status message for ETH69.
func (c *conn) statusExchange69(packet *BorStatusPacket69) error {
	errc := make(chan error, 2)

	go func() {
		c.countMsgSent(packet.Name(), 1)
		errc <- ethp2p.Send(c.rw, eth.StatusMsg, packet)
	}()

	go func() {
		errc <- c.readStatus69(packet)
	}()

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	for range 2 {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return ethp2p.DiscReadTimeout
		}
	}

	return nil
}

// countMsgReceived increments the global Prometheus counter and per-peer message tracking for received messages.
func (c *conn) countMsgReceived(messageName string, count float64) {
	// Increment global Prometheus counter (low cardinality)
	SensorMsgCounter.WithLabelValues(messageName, string(MsgReceived)).Add(count)
	SensorMsgCounter.WithLabelValues(messageName+PacketSuffix, string(MsgReceived)).Add(1)

	// Increment per-peer message tracking (for API visibility)
	c.messages.IncrementReceived(messageName, int64(count))
}

// countMsgSent increments the global Prometheus counter and per-peer message tracking for sent messages.
func (c *conn) countMsgSent(messageName string, count float64) {
	// Increment global Prometheus counter (low cardinality)
	SensorMsgCounter.WithLabelValues(messageName, string(MsgSent)).Add(count)
	SensorMsgCounter.WithLabelValues(messageName+PacketSuffix, string(MsgSent)).Add(1)

	// Increment per-peer message tracking (for API visibility)
	c.messages.IncrementSent(messageName, int64(count))
}

func (c *conn) readStatus68(packet *StatusPacket68) error {
	msg, err := c.rw.ReadMsg()
	if err != nil {
		return err
	}

	if msg.Code != eth.StatusMsg {
		return errors.New("expected status message code")
	}

	var status StatusPacket68
	if err := msg.Decode(&status); err != nil {
		return err
	}

	if status.NetworkID != packet.NetworkID {
		return fmt.Errorf("network ID mismatch: %d (!= %d)", status.NetworkID, packet.NetworkID)
	}
	if status.Genesis != packet.Genesis {
		return fmt.Errorf("genesis mismatch: %v (!= %v)", status.Genesis, packet.Genesis)
	}
	if status.ForkID.Hash != packet.ForkID.Hash {
		return fmt.Errorf("fork ID mismatch: %v (!= %v)", status.ForkID, packet.ForkID)
	}

	c.logger.Info().
		Interface("status", status).
		Str("fork_id", hex.EncodeToString(status.ForkID.Hash[:])).
		Msg("New peer")

	return nil
}

func (c *conn) readStatus69(packet *BorStatusPacket69) error {
	msg, err := c.rw.ReadMsg()
	if err != nil {
		return err
	}

	if msg.Code != eth.StatusMsg {
		return errors.New("expected status message code")
	}

	var status BorStatusPacket69
	if err := msg.Decode(&status); err != nil {
		return err
	}

	if status.NetworkID != packet.NetworkID {
		return fmt.Errorf("network ID mismatch: %d (!= %d)", status.NetworkID, packet.NetworkID)
	}
	if status.Genesis != packet.Genesis {
		return fmt.Errorf("genesis mismatch: %v (!= %v)", status.Genesis, packet.Genesis)
	}
	if status.ForkID.Hash != packet.ForkID.Hash {
		return fmt.Errorf("fork ID mismatch: %v (!= %v)", status.ForkID, packet.ForkID)
	}

	c.logger.Info().
		Interface("status", status).
		Str("fork_id", hex.EncodeToString(status.ForkID.Hash[:])).
		Uint64("earliest_block", status.EarliestBlock).
		Uint64("latest_block", status.LatestBlock).
		Msg("New peer")

	return nil
}

// handleBlockRangeUpdate handles BlockRangeUpdateMsg (ETH69).
// This message announces the peer's available block range.
func (c *conn) handleBlockRangeUpdate(msg ethp2p.Msg) error {
	var packet eth.BlockRangeUpdatePacket
	if err := msg.Decode(&packet); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to decode BlockRangeUpdate")
		return nil
	}

	c.countMsgReceived(packet.Name(), 1)
	c.logger.Debug().
		Uint64("earliest", packet.EarliestBlock).
		Uint64("latest", packet.LatestBlock).
		Hex("hash", packet.LatestBlockHash[:]).
		Msg("Received BlockRangeUpdate")

	return nil
}

// getBlockData will send GetBlockHeaders and/or GetBlockBodies requests to the
// peer based on what parts of the block we already have. It will return an error
// if sending either of the requests failed. The isParent parameter indicates if
// this block is being fetched as a parent block.
func (c *conn) getBlockData(hash common.Hash, cache BlockCache, isParent bool) error {
	// Only request header if we don't have it
	if cache.Header == nil {
		headersRequest := &GetBlockHeaders{
			GetBlockHeadersRequest: &eth.GetBlockHeadersRequest{
				// Providing both the hash and number will result in a `both origin
				// hash and number` error.
				Origin: eth.HashOrNumber{Hash: hash},
				Amount: 1,
			},
		}

		if isParent {
			c.parents.Add(hash, struct{}{})
		}

		c.countMsgSent(headersRequest.Name(), 1)
		if err := ethp2p.Send(c.rw, eth.GetBlockHeadersMsg, headersRequest); err != nil {
			return err
		}
	}

	// Only request body if we don't have it
	if cache.Body == nil {
		c.requestNum++
		c.requests.Add(c.requestNum, hash)

		bodiesRequest := &GetBlockBodies{
			RequestId:             c.requestNum,
			GetBlockBodiesRequest: []common.Hash{hash},
		}

		c.countMsgSent(bodiesRequest.Name(), 1)
		if err := ethp2p.Send(c.rw, eth.GetBlockBodiesMsg, bodiesRequest); err != nil {
			return err
		}
	}

	return nil
}

// getParentBlock will send a request to the peer if the parent of the header
// does not exist in the database. It only fetches parents back to the oldest
// block (initialized to the head block at sensor startup).
func (c *conn) getParentBlock(ctx context.Context, header *types.Header) error {
	if !c.db.ShouldWriteBlocks() {
		return nil
	}

	oldestBlock := c.conns.OldestBlock()
	if oldestBlock == nil {
		return nil
	}

	// Check cache first before querying the database
	cache, ok := c.conns.Blocks().Peek(header.ParentHash)
	if ok && cache.Header != nil && cache.Body != nil {
		return nil
	}

	// Don't fetch parents older than our starting point (oldest block)
	if c.db.HasBlock(ctx, header.ParentHash) || header.Number.Cmp(oldestBlock.Number) != 1 {
		return nil
	}

	c.logger.Info().
		Str("hash", header.ParentHash.Hex()).
		Str("number", new(big.Int).Sub(header.Number, big.NewInt(1)).String()).
		Msg("Fetching missing parent block")

	return c.getBlockData(header.ParentHash, cache, true)
}

func (c *conn) handleNewBlockHashes(ctx context.Context, msg ethp2p.Msg) error {
	var packet NewBlockHashesPacket
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	tfs := time.Now()

	c.countMsgReceived(packet.Name(), float64(len(packet)))

	// Collect unique hashes and numbers for database write and broadcasting.
	uniqueHashes := make([]common.Hash, 0, len(packet))
	uniqueNumbers := make([]uint64, 0, len(packet))

	for _, entry := range packet {
		hash := entry.Hash

		// Mark as known from this peer
		c.addKnownBlock(hash)

		// Atomically check and add to cache to prevent duplicate writes from
		// concurrent peers receiving the same block hash.
		ok := c.conns.Blocks().Update(hash, func(cache BlockCache) BlockCache {
			if cache != (BlockCache{}) {
				return cache
			}
			return BlockCache{}
		})

		if ok {
			continue
		}

		// Write hash first seen time immediately for new blocks
		c.db.WriteBlockHashFirstSeen(ctx, c.node, hash, tfs)

		// Request only the parts we don't have
		if err := c.getBlockData(hash, BlockCache{}, false); err != nil {
			return err
		}

		uniqueHashes = append(uniqueHashes, hash)
		uniqueNumbers = append(uniqueNumbers, entry.Number)
	}

	// Write only unique hashes to the database.
	if len(uniqueHashes) == 0 {
		return nil
	}

	c.db.WriteBlockHashes(ctx, c.node, uniqueHashes, tfs)

	// Broadcast block hashes to other peers asynchronously
	go c.conns.BroadcastBlockHashes(uniqueHashes, uniqueNumbers)

	return nil
}

// addKnownTx adds a transaction hash to the known tx cache.
func (c *conn) addKnownTx(hash common.Hash) {
	if !c.shouldBroadcastTx && !c.shouldBroadcastTxHashes {
		return
	}

	c.knownTxs.Add(hash)
}

// addKnownBlock adds a block hash to the known block cache.
func (c *conn) addKnownBlock(hash common.Hash) {
	if !c.shouldBroadcastBlocks && !c.shouldBroadcastBlockHashes {
		return
	}

	c.knownBlocks.Add(hash)
}

// hasKnownTx checks if a transaction hash is in the known tx cache.
func (c *conn) hasKnownTx(hash common.Hash) bool {
	if !c.shouldBroadcastTx && !c.shouldBroadcastTxHashes {
		return false
	}

	return c.knownTxs.Contains(hash)
}

// hasKnownBlock checks if a block hash is in the known block cache.
func (c *conn) hasKnownBlock(hash common.Hash) bool {
	if !c.shouldBroadcastBlocks && !c.shouldBroadcastBlockHashes {
		return false
	}

	return c.knownBlocks.Contains(hash)
}

// txAnnouncementLoop schedules transaction hash announcements to the peer.
// Matches Bor's announceTransactions pattern: async sends with internal queue.
func (c *conn) txAnnouncementLoop() {
	var (
		queue  []common.Hash         // Queue of hashes to announce
		done   chan struct{}         // Non-nil if background announcer is running
		fail   = make(chan error, 1) // Channel used to receive network error
		failed bool                  // Flag whether a send failed
	)

	for {
		// If there's no in-flight announce running, check if a new one is needed
		if done == nil && len(queue) > 0 {
			var pending []common.Hash
			pending, queue = c.prepareTxAnnouncements(queue)

			// If there's anything available to transfer, fire up an async writer
			if len(pending) > 0 {
				done = make(chan struct{})
				go func() {
					if err := c.sendTxAnnouncements(pending); err != nil {
						fail <- err
						return
					}
					close(done)
				}()
			}
		}

		// Transfer goroutine may or may not have been started, listen for events
		select {
		case hashes := <-c.txAnnounce:
			if !failed {
				queue = c.enqueueTxHashes(queue, hashes)
			}

		case <-done:
			done = nil

		case <-fail:
			failed = true

		case <-c.closeCh:
			return
		}
	}
}

// prepareTxAnnouncements extracts a batch of unknown tx hashes from the queue
// up to maxTxPacketSize bytes. Returns the pending hashes and remaining queue.
func (c *conn) prepareTxAnnouncements(queue []common.Hash) (pending, remaining []common.Hash) {
	// Calculate max hashes we can send based on packet size limit
	maxHashes := min(maxTxPacketSize/common.HashLength, len(queue))

	// Filter out known hashes in a single lock operation
	pending = c.knownTxs.FilterNotContained(queue[:maxHashes])
	remaining = queue[:copy(queue, queue[maxHashes:])]
	return pending, remaining
}

// enqueueTxHashes adds hashes to the queue, dropping oldest if over capacity.
func (c *conn) enqueueTxHashes(queue, hashes []common.Hash) []common.Hash {
	queue = append(queue, hashes...)
	if len(queue) > maxQueuedTxAnns {
		queue = queue[:copy(queue, queue[len(queue)-maxQueuedTxAnns:])]
	}
	return queue
}

// sendTxAnnouncements sends a batch of transaction hashes to the peer.
// It looks up each transaction from the cache to populate Types and Sizes
// as required by the ETH68 protocol.
func (c *conn) sendTxAnnouncements(hashes []common.Hash) error {
	// Batch lookup all transactions in a single lock operation.
	// Skip hashes where the transaction is no longer in cache.
	pending, txs := c.conns.PeekTxsWithHashes(hashes)
	if len(pending) == 0 {
		return nil
	}

	// Build Types and Sizes from the found transactions.
	pendingTypes := make([]byte, len(txs))
	pendingSizes := make([]uint32, len(txs))
	for i, tx := range txs {
		pendingTypes[i] = tx.Type()
		pendingSizes[i] = uint32(tx.Size())
	}

	packet := eth.NewPooledTransactionHashesPacket{
		Types:  pendingTypes,
		Sizes:  pendingSizes,
		Hashes: pending,
	}
	c.countMsgSent(packet.Name(), float64(len(pending)))
	if err := ethp2p.Send(c.rw, eth.NewPooledTransactionHashesMsg, packet); err != nil {
		c.logger.Debug().Err(err).Msg("Failed to send tx announcements")
		return err
	}

	// Mark all hashes as known in a single lock operation
	c.knownTxs.AddMany(pending)
	return nil
}

// blockAnnouncementLoop drains the blockAnnounce queue and sends block
// announcements. Matches Bor's broadcastBlocks pattern.
func (c *conn) blockAnnouncementLoop() {
	for {
		select {
		case packet := <-c.blockAnnounce:
			if c.sendBlockAnnouncements(packet) != nil {
				return
			}
		case <-c.closeCh:
			return
		}
	}
}

// sendBlockAnnouncements sends a batch of block hashes to the peer,
// filtering out blocks the peer already knows about.
func (c *conn) sendBlockAnnouncements(packet NewBlockHashesPacket) error {
	// Filter to only unknown blocks
	var filtered NewBlockHashesPacket
	for _, entry := range packet {
		if !c.hasKnownBlock(entry.Hash) {
			filtered = append(filtered, entry)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	c.countMsgSent(filtered.Name(), float64(len(filtered)))
	if err := ethp2p.Send(c.rw, eth.NewBlockHashesMsg, filtered); err != nil {
		c.logger.Debug().Err(err).Msg("Failed to send block announcements")
		return err
	}
	for _, entry := range filtered {
		c.addKnownBlock(entry.Hash)
	}
	return nil
}

// decodeTx attempts to decode a transaction from an RLP-encoded raw value.
func (c *conn) decodeTx(raw []byte) *types.Transaction {
	if len(raw) == 0 {
		return nil
	}

	// Try decoding as RLP-wrapped bytes first (legacy format)
	var bytes []byte
	if rlp.DecodeBytes(raw, &bytes) == nil {
		tx := new(types.Transaction)
		err := tx.UnmarshalBinary(bytes)
		if err == nil {
			return tx
		}

		c.logger.Warn().
			Err(err).
			Uint8("type", bytes[0]).
			Int("size", len(bytes)).
			Str("hash", crypto.Keccak256Hash(bytes).Hex()).
			Msg("Failed to decode transaction")

		return nil
	}

	// Try decoding as raw binary (typed transaction format)
	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(raw)
	if err == nil {
		return tx
	}

	c.logger.Warn().
		Err(err).
		Uint8("prefix", raw[0]).
		Int("size", len(raw)).
		Str("hash", crypto.Keccak256Hash(raw).Hex()).
		Msg("Failed to decode transaction")

	return nil
}

// decodeTxs decodes a list of transactions, returning only successfully decoded ones.
func (c *conn) decodeTxs(rawTxs []rlp.RawValue) []*types.Transaction {
	var txs []*types.Transaction

	for _, raw := range rawTxs {
		if tx := c.decodeTx(raw); tx != nil {
			txs = append(txs, tx)
		}
	}

	return txs
}

// encodeBlockBody converts a block to an eth.BlockBody with RLP-encoded fields.
func encodeBlockBody(block *types.Block) (*eth.BlockBody, error) {
	txList, err := rlp.EncodeToRawList([]*types.Transaction(block.Transactions()))
	if err != nil {
		return nil, fmt.Errorf("failed to encode transactions: %w", err)
	}
	uncleList, err := rlp.EncodeToRawList(block.Uncles())
	if err != nil {
		return nil, fmt.Errorf("failed to encode uncles: %w", err)
	}
	var withdrawalList *rlp.RawList[*types.Withdrawal]
	if withdrawals := block.Withdrawals(); withdrawals != nil {
		wl, err := rlp.EncodeToRawList([]*types.Withdrawal(withdrawals))
		if err != nil {
			return nil, fmt.Errorf("failed to encode withdrawals: %w", err)
		}
		withdrawalList = &wl
	}
	return &eth.BlockBody{
		Transactions: txList,
		Uncles:       uncleList,
		Withdrawals:  withdrawalList,
	}, nil
}

func (c *conn) handleTransactions(ctx context.Context, msg ethp2p.Msg) error {
	payload, err := io.ReadAll(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to read transactions payload: %w", err)
	}

	var rawTxs []rlp.RawValue
	if err := rlp.DecodeBytes(payload, &rawTxs); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to decode transactions")
		return nil
	}

	txs := c.decodeTxs(rawTxs)
	tfs := time.Now()

	c.countMsgReceived((&eth.TransactionsPacket{}).Name(), float64(len(txs)))

	// Mark transactions as known from this peer
	for _, tx := range txs {
		c.addKnownTx(tx.Hash())
	}

	if len(txs) > 0 {
		c.db.WriteTransactions(ctx, c.node, txs, tfs)
	}

	// Cache transactions for duplicate detection and serving to peers (single lock)
	hashes := c.conns.AddTxs(txs)

	// Broadcast transactions or hashes to other peers asynchronously
	go c.conns.BroadcastTxs(types.Transactions(txs))
	go c.conns.BroadcastTxHashes(hashes)

	return nil
}

func (c *conn) handleGetBlockHeaders(msg ethp2p.Msg) error {
	var request eth.GetBlockHeadersPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), 1)

	// Try to serve from cache if we have the block
	var headers []*types.Header
	if cache, ok := c.conns.Blocks().Peek(request.Origin.Hash); ok && cache.Header != nil {
		headers = []*types.Header{cache.Header}
	}

	headerList, err := rlp.EncodeToRawList(headers)
	if err != nil {
		return fmt.Errorf("failed to encode headers: %w", err)
	}
	response := &eth.BlockHeadersPacket{
		RequestId: request.RequestId,
		List:      headerList,
	}
	c.countMsgSent((*eth.BlockHeadersRequest)(nil).Name(), float64(len(headers)))
	return ethp2p.Send(c.rw, eth.BlockHeadersMsg, response)
}

func (c *conn) handleBlockHeaders(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.BlockHeadersPacket
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	tfs := time.Now()

	headers, err := packet.List.Items()
	if err != nil {
		return fmt.Errorf("failed to decode block headers: %w", err)
	}
	if len(headers) == 0 {
		return nil
	}

	c.countMsgReceived((*eth.BlockHeadersRequest)(nil).Name(), float64(len(headers)))

	for _, header := range headers {
		if err := c.getParentBlock(ctx, header); err != nil {
			return err
		}
	}

	// Check if any of these headers were requested as parent blocks
	_, isParent := c.parents.Remove(headers[0].Hash())

	c.db.WriteBlockHeaders(ctx, headers, tfs, isParent)

	// Update cache to store headers
	for _, header := range headers {
		hash := header.Hash()
		c.conns.Blocks().Update(hash, func(cache BlockCache) BlockCache {
			cache.Header = header
			return cache
		})
	}

	return nil
}

func (c *conn) handleGetBlockBodies(msg ethp2p.Msg) error {
	var request eth.GetBlockBodiesPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), float64(len(request.GetBlockBodiesRequest)))

	// Try to serve from cache
	var bodies []*eth.BlockBody
	for _, hash := range request.GetBlockBodiesRequest {
		if cache, ok := c.conns.Blocks().Peek(hash); ok && cache.Body != nil {
			bodies = append(bodies, cache.Body)
		}
	}

	// Convert to non-pointer slice for encoding
	bodiesForRLP := make([]eth.BlockBody, len(bodies))
	for i, b := range bodies {
		bodiesForRLP[i] = *b
	}
	bodyList, err := rlp.EncodeToRawList(bodiesForRLP)
	if err != nil {
		return fmt.Errorf("failed to encode block bodies: %w", err)
	}
	response := &eth.BlockBodiesPacket{
		RequestId: request.RequestId,
		List:      bodyList,
	}
	c.countMsgSent((*eth.BlockBodiesResponse)(nil).Name(), float64(len(bodies)))
	return ethp2p.Send(c.rw, eth.BlockBodiesMsg, response)
}

func (c *conn) handleBlockBodies(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.BlockBodiesRLPPacket
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	tfs := time.Now()

	if len(packet.BlockBodiesRLPResponse) == 0 {
		return nil
	}

	c.countMsgReceived((*eth.BlockBodiesResponse)(nil).Name(), float64(len(packet.BlockBodiesRLPResponse)))

	hash, ok := c.requests.Get(packet.RequestId)
	if !ok {
		c.logger.Warn().Msg("No block hash found for block body")
		return nil
	}
	c.requests.Remove(packet.RequestId)

	// Check if we already have the body in the cache
	if cache, ok := c.conns.Blocks().Peek(hash); ok && cache.Body != nil {
		return nil
	}

	var decoded rawBlockBody
	if err := rlp.DecodeBytes(packet.BlockBodiesRLPResponse[0], &decoded); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to decode block body")
		return nil
	}

	txs := c.decodeTxs(decoded.Transactions)
	txList, err := rlp.EncodeToRawList(txs)
	if err != nil {
		c.logger.Warn().Err(err).Msg("Failed to encode transactions")
		return nil
	}
	uncleList, err := rlp.EncodeToRawList(decoded.Uncles)
	if err != nil {
		c.logger.Warn().Err(err).Msg("Failed to encode uncles")
		return nil
	}
	var withdrawalList *rlp.RawList[*types.Withdrawal]
	if decoded.Withdrawals != nil {
		wl, err := rlp.EncodeToRawList(decoded.Withdrawals)
		if err != nil {
			c.logger.Warn().Err(err).Msg("Failed to encode withdrawals")
			return nil
		}
		withdrawalList = &wl
	}
	body := &eth.BlockBody{
		Transactions: txList,
		Uncles:       uncleList,
		Withdrawals:  withdrawalList,
	}

	c.db.WriteBlockBody(ctx, body, hash, tfs)

	// Update cache to store body
	c.conns.Blocks().Update(hash, func(cache BlockCache) BlockCache {
		cache.Body = body
		return cache
	})

	return nil
}

func (c *conn) handleNewBlock(ctx context.Context, msg ethp2p.Msg) error {
	payload, err := io.ReadAll(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to read new block payload: %w", err)
	}

	var raw rawNewBlockPacket
	if err = rlp.DecodeBytes(payload, &raw); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to decode new block")
		return nil
	}

	block := types.NewBlockWithHeader(raw.Block.Header).WithBody(types.Body{
		Transactions: c.decodeTxs(raw.Block.Txs),
		Uncles:       raw.Block.Uncles,
		Withdrawals:  raw.Block.Withdrawals,
	})
	packet := &NewBlockPacket{Block: block, TD: raw.TD}

	tfs := time.Now()
	hash := packet.Block.Hash()

	c.countMsgReceived(packet.Name(), 1)

	// Mark block as known from this peer
	c.addKnownBlock(hash)

	// Set the head block if newer.
	if c.conns.UpdateHeadBlock(*packet) {
		c.logger.Info().
			Str("hash", hash.Hex()).
			Uint64("number", packet.Block.Number().Uint64()).
			Str("td", packet.TD.String()).
			Msg("Updated head block")
	}

	if err = c.getParentBlock(ctx, packet.Block.Header()); err != nil {
		return err
	}

	// Create BlockBody with encoded RawLists
	blockBody, err := encodeBlockBody(packet.Block)
	if err != nil {
		return fmt.Errorf("failed to encode block body: %w", err)
	}

	// Atomically check and add to cache to prevent duplicate writes from
	// concurrent peers receiving the same block.
	var exists bool
	ok := c.conns.Blocks().Update(hash, func(cache BlockCache) BlockCache {
		if cache.TD != nil {
			exists = true
			return cache
		}

		return BlockCache{
			Header: packet.Block.Header(),
			Body:   blockBody,
			TD:     packet.TD,
		}
	})

	if exists {
		return nil
	}

	// Write first-seen event for blocks arriving directly (not announced via hash first)
	if !ok {
		c.db.WriteBlockHashFirstSeen(ctx, c.node, hash, tfs)
	}

	c.db.WriteBlock(ctx, c.node, packet.Block, packet.TD, tfs)

	// Broadcast block or block hash to other peers asynchronously
	go c.conns.BroadcastBlock(packet.Block, packet.TD)
	go c.conns.BroadcastBlockHashes(
		[]common.Hash{hash},
		[]uint64{packet.Block.Number().Uint64()},
	)

	return nil
}

func (c *conn) handleGetPooledTransactions(msg ethp2p.Msg) error {
	var request eth.GetPooledTransactionsPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), float64(len(request.GetPooledTransactionsRequest)))

	// Try to serve from cache using batch lookup (single read lock operation)
	txs := c.conns.PeekTxs(request.GetPooledTransactionsRequest)

	txList, err := rlp.EncodeToRawList(txs)
	if err != nil {
		return fmt.Errorf("failed to encode pooled transactions: %w", err)
	}
	response := &eth.PooledTransactionsPacket{
		RequestId: request.RequestId,
		List:      txList,
	}
	c.countMsgSent(response.Name(), float64(len(txs)))
	return ethp2p.Send(c.rw, eth.PooledTransactionsMsg, response)
}

func (c *conn) handleNewPooledTransactionHashes(version uint, msg ethp2p.Msg) error {
	var hashes []common.Hash
	var name string

	switch version {
	case 67, 68, 69:
		var txs eth.NewPooledTransactionHashesPacket
		if err := msg.Decode(&txs); err != nil {
			return err
		}
		hashes = txs.Hashes
		name = txs.Name()
	default:
		return errors.New("protocol version not found")
	}

	c.countMsgReceived(name, float64(len(hashes)))

	if !c.db.ShouldWriteTransactions() || !c.db.ShouldWriteTransactionEvents() {
		return nil
	}

	request := &eth.GetPooledTransactionsPacket{GetPooledTransactionsRequest: hashes}
	c.countMsgSent(request.Name(), float64(len(hashes)))
	return ethp2p.Send(c.rw, eth.GetPooledTransactionsMsg, request)
}

func (c *conn) handlePooledTransactions(ctx context.Context, msg ethp2p.Msg) error {
	payload, err := io.ReadAll(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to read pooled transactions payload: %w", err)
	}

	var raw rawPooledTransactionsPacket
	if err := rlp.DecodeBytes(payload, &raw); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to decode pooled transactions")
		return nil
	}

	txs := c.decodeTxs(raw.Txs)

	tfs := time.Now()

	c.countMsgReceived((*eth.PooledTransactionsPacket)(nil).Name(), float64(len(txs)))

	// Mark transactions as known from this peer
	for _, tx := range txs {
		c.addKnownTx(tx.Hash())
	}

	if len(txs) > 0 {
		c.db.WriteTransactions(ctx, c.node, txs, tfs)
	}

	// Cache transactions for duplicate detection and serving to peers (single lock)
	hashes := c.conns.AddTxs(txs)

	// Broadcast transactions or hashes to other peers asynchronously
	go c.conns.BroadcastTxs(types.Transactions(txs))
	go c.conns.BroadcastTxHashes(hashes)

	return nil
}

func (c *conn) handleGetReceipts(msg ethp2p.Msg) error {
	var request eth.GetReceiptsPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), float64(len(request.GetReceiptsRequest)))

	response := &ReceiptsRLPPacket{RequestId: request.RequestId}
	c.countMsgSent((&eth.ReceiptsRLPResponse{}).Name(), 0)
	return ethp2p.Send(c.rw, eth.ReceiptsMsg, response)
}
