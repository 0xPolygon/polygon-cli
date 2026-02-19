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
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0xPolygon/polygon-cli/p2p/database"
)

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
	requests   *Cache[uint64, common.Hash]
	requestNum uint64

	// parents tracks hashes of blocks requested as parents to mark them
	// with IsParent=true when writing to the database.
	parents *Cache[common.Hash, struct{}]

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
	knownTxs    *Cache[common.Hash, struct{}]
	knownBlocks *Cache[common.Hash, struct{}]

	// messages tracks per-peer message counts for API visibility.
	messages *PeerMessages
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
	RequestsCache CacheOptions
	ParentsCache  CacheOptions

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
		Length:  17,
		Run: func(p *ethp2p.Peer, rw ethp2p.MsgReadWriter) error {
			peerURL := p.Node().URLv4()
			c := &conn{
				sensorID:                   opts.SensorID,
				node:                       p.Node(),
				logger:                     log.With().Str("peer", peerURL).Logger(),
				rw:                         rw,
				db:                         opts.Database,
				requests:                   NewCache[uint64, common.Hash](opts.RequestsCache),
				requestNum:                 0,
				parents:                    NewCache[common.Hash, struct{}](opts.ParentsCache),
				peer:                       p,
				conns:                      opts.Conns,
				connectedAt:                time.Now(),
				peerURL:                    peerURL,
				shouldBroadcastTx:          opts.ShouldBroadcastTx,
				shouldBroadcastTxHashes:    opts.ShouldBroadcastTxHashes,
				shouldBroadcastBlocks:      opts.ShouldBroadcastBlocks,
				shouldBroadcastBlockHashes: opts.ShouldBroadcastBlockHashes,
				knownTxs:                   NewCache[common.Hash, struct{}](opts.Conns.KnownTxsOpts()),
				knownBlocks:                NewCache[common.Hash, struct{}](opts.Conns.KnownBlocksOpts()),
				messages:                   NewPeerMessages(),
			}

			head := c.conns.HeadBlock()
			status := eth.StatusPacket68{
				ProtocolVersion: uint32(version),
				NetworkID:       opts.NetworkID,
				Genesis:         opts.GenesisHash,
				ForkID:          opts.ForkID,
				Head:            head.Block.Hash(),
				TD:              head.TD,
			}
			err := c.statusExchange(&status)
			if err != nil {
				return err
			}

			// Send the connection object to the conns manager for RPC broadcasting
			opts.Conns.Add(c)
			defer opts.Conns.Remove(c)

			ctx := opts.Context

			// Disconnect peer when context is cancelled to unblock ReadMsg.
			go func() {
				<-ctx.Done()
				p.Disconnect(ethp2p.DiscQuitting)
			}()

			// Handle all the of the messages here.
			for {
				// Check for context cancellation before processing next message.
				select {
				case <-ctx.Done():
					return nil
				default:
				}

				msg, err := rw.ReadMsg()
				if err != nil {
					// Return nil on context cancellation to avoid error logging.
					if ctx.Err() != nil {
						return nil
					}
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

// statusExchange will exchange status message between the nodes. It will return
// an error if the nodes are incompatible.
func (c *conn) statusExchange(packet *eth.StatusPacket68) error {
	errc := make(chan error, 2)

	go func() {
		c.countMsgSent((&eth.StatusPacket68{}).Name(), 1)
		errc <- ethp2p.Send(c.rw, eth.StatusMsg, &packet)
	}()

	go func() {
		errc <- c.readStatus(packet)
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

func (c *conn) readStatus(packet *eth.StatusPacket68) error {
	msg, err := c.rw.ReadMsg()
	if err != nil {
		return err
	}

	defer func() {
		if msgErr := msg.Discard(); msgErr != nil {
			c.logger.Error().Err(msgErr).Msg("Failed to discard message")
		}
	}()

	if msg.Code != eth.StatusMsg {
		return errors.New("expected status message code")
	}

	var status eth.StatusPacket68
	err = msg.Decode(&status)
	if err != nil {
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
	var packet eth.NewBlockHashesPacket
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

		// Check what parts of the block we already have
		cache, ok := c.conns.Blocks().Get(hash)
		if ok {
			continue
		}

		// Write hash first seen time immediately for new blocks
		c.db.WriteBlockHashFirstSeen(ctx, hash, tfs)

		// Request only the parts we don't have
		if err := c.getBlockData(hash, cache, false); err != nil {
			return err
		}

		c.conns.Blocks().Add(hash, BlockCache{})
		uniqueHashes = append(uniqueHashes, hash)
		uniqueNumbers = append(uniqueNumbers, entry.Number)
	}

	// Write only unique hashes to the database.
	if len(uniqueHashes) == 0 {
		return nil
	}

	c.db.WriteBlockHashes(ctx, c.node, uniqueHashes, tfs)

	// Broadcast block hashes to other peers
	c.conns.BroadcastBlockHashes(uniqueHashes, uniqueNumbers)

	return nil
}

// addKnownTx adds a transaction hash to the known tx cache.
func (c *conn) addKnownTx(hash common.Hash) {
	if !c.shouldBroadcastTx && !c.shouldBroadcastTxHashes {
		return
	}

	c.knownTxs.Add(hash, struct{}{})
}

// addKnownBlock adds a block hash to the known block cache.
func (c *conn) addKnownBlock(hash common.Hash) {
	if !c.shouldBroadcastBlocks && !c.shouldBroadcastBlockHashes {
		return
	}

	c.knownBlocks.Add(hash, struct{}{})
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

	txs := decodeTxs(rawTxs)
	tfs := time.Now()

	c.countMsgReceived((&eth.TransactionsPacket{}).Name(), float64(len(txs)))

	// Mark transactions as known from this peer
	for _, tx := range txs {
		c.addKnownTx(tx.Hash())
	}

	if len(txs) > 0 {
		c.db.WriteTransactions(ctx, c.node, txs, tfs)
	}

	// Cache transactions for duplicate detection and serving to peers
	hashes := make([]common.Hash, len(txs))
	for i, tx := range txs {
		c.conns.AddTx(tx.Hash(), tx)
		hashes[i] = tx.Hash()
	}

	// Broadcast transactions or hashes to other peers
	c.conns.BroadcastTxs(types.Transactions(txs))
	c.conns.BroadcastTxHashes(hashes)

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

	response := &eth.BlockHeadersPacket{
		RequestId:           request.RequestId,
		BlockHeadersRequest: headers,
	}
	c.countMsgSent(response.Name(), float64(len(headers)))
	return ethp2p.Send(c.rw, eth.BlockHeadersMsg, response)
}

func (c *conn) handleBlockHeaders(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.BlockHeadersPacket
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	tfs := time.Now()

	headers := packet.BlockHeadersRequest
	if len(headers) == 0 {
		return nil
	}

	c.countMsgReceived(packet.Name(), float64(len(headers)))

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

	response := &eth.BlockBodiesPacket{
		RequestId:           request.RequestId,
		BlockBodiesResponse: bodies,
	}
	c.countMsgSent(response.Name(), float64(len(bodies)))
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

	c.countMsgReceived((&eth.BlockBodiesPacket{}).Name(), float64(len(packet.BlockBodiesRLPResponse)))

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

	body := &eth.BlockBody{
		Transactions: decodeTxs(decoded.Transactions),
		Uncles:       decoded.Uncles,
		Withdrawals:  decoded.Withdrawals,
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
	if err := rlp.DecodeBytes(payload, &raw); err != nil {
		c.logger.Warn().Err(err).Msg("Failed to decode new block")
		return nil
	}

	block := types.NewBlockWithHeader(raw.Block.Header).WithBody(types.Body{
		Transactions: decodeTxs(raw.Block.Txs),
		Uncles:       raw.Block.Uncles,
		Withdrawals:  raw.Block.Withdrawals,
	})
	packet := &eth.NewBlockPacket{Block: block, TD: raw.TD}

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

	if err := c.getParentBlock(ctx, packet.Block.Header()); err != nil {
		return err
	}

	// Check if we already have the full block in the cache
	if cache, ok := c.conns.Blocks().Peek(hash); ok && cache.TD != nil {
		return nil
	}

	c.db.WriteBlock(ctx, c.node, packet.Block, packet.TD, tfs)

	// Update cache to store the full block
	c.conns.Blocks().Add(hash, BlockCache{
		Header: packet.Block.Header(),
		Body: &eth.BlockBody{
			Transactions: packet.Block.Transactions(),
			Uncles:       packet.Block.Uncles(),
			Withdrawals:  packet.Block.Withdrawals(),
		},
		TD: packet.TD,
	})

	// Broadcast block or block hash to other peers
	c.conns.BroadcastBlock(packet.Block, packet.TD)
	c.conns.BroadcastBlockHashes(
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

	// Try to serve from cache
	var txs []*types.Transaction
	for _, hash := range request.GetPooledTransactionsRequest {
		if tx, ok := c.conns.GetTx(hash); ok {
			txs = append(txs, tx)
		}
	}

	response := &eth.PooledTransactionsPacket{
		RequestId:                  request.RequestId,
		PooledTransactionsResponse: txs,
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

	packet := &eth.PooledTransactionsPacket{
		PooledTransactionsResponse: decodeTxs(raw.Txs),
	}

	tfs := time.Now()

	c.countMsgReceived(packet.Name(), float64(len(packet.PooledTransactionsResponse)))

	// Mark transactions as known from this peer
	for _, tx := range packet.PooledTransactionsResponse {
		c.addKnownTx(tx.Hash())
	}

	if len(packet.PooledTransactionsResponse) > 0 {
		c.db.WriteTransactions(ctx, c.node, packet.PooledTransactionsResponse, tfs)
	}

	// Cache transactions for duplicate detection and serving to peers
	hashes := make([]common.Hash, len(packet.PooledTransactionsResponse))
	for i, tx := range packet.PooledTransactionsResponse {
		c.conns.AddTx(tx.Hash(), tx)
		hashes[i] = tx.Hash()
	}

	// Broadcast transactions or hashes to other peers
	c.conns.BroadcastTxs(types.Transactions(packet.PooledTransactionsResponse))
	c.conns.BroadcastTxHashes(hashes)

	return nil
}

func (c *conn) handleGetReceipts(msg ethp2p.Msg) error {
	var request eth.GetReceiptsPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), float64(len(request.GetReceiptsRequest)))

	response := &eth.ReceiptsRLPPacket{RequestId: request.RequestId}
	c.countMsgSent((&eth.ReceiptsRLPResponse{}).Name(), 0)
	return ethp2p.Send(c.rw, eth.ReceiptsMsg, response)
}
