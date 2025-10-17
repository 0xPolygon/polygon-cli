package p2p

import (
	"container/list"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0xPolygon/polygon-cli/p2p/database"
)

// conn represents an individual connection with a peer.
type conn struct {
	sensorID  string
	node      *enode.Node
	logger    zerolog.Logger
	rw        ethp2p.MsgReadWriter
	db        database.Database
	head      *HeadBlock
	headMutex *sync.RWMutex
	counter   *prometheus.CounterVec
	peer      *ethp2p.Peer

	// requests is used to store the request ID and the block hash. This is used
	// when fetching block bodies because the eth protocol block bodies do not
	// contain information about the block hash.
	requests   *Cache[uint64, common.Hash]
	requestNum uint64

	// Linked list of seen block hashes with timestamps.
	blockHashes *list.List

	// oldestBlock stores the first block the sensor has seen so when fetching
	// parent blocks, it does not request blocks older than this.
	oldestBlock *types.Header

	// connectedAt stores when this peer connection was established.
	connectedAt time.Time

	// activeGoroutines tracks the number of active goroutines spawned by this peer
	// for database write operations
	activeGoroutines atomic.Int64
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
	MsgCounter  *prometheus.CounterVec

	// Head keeps track of the current head block of the chain. This is required
	// when doing the status exchange.
	Head      *HeadBlock
	HeadMutex *sync.RWMutex

	// Requests cache configuration
	MaxRequests      int
	RequestsCacheTTL time.Duration
}

// HeadBlock contains the necessary head block data for the status message.
type HeadBlock struct {
	Hash            common.Hash
	TotalDifficulty *big.Int
	Number          uint64
	Time            uint64
}

type BlockHashEntry struct {
	hash common.Hash
	time time.Time
}

// blockHashTTL defines the time-to-live for block hash entries in blockHashes list.
var blockHashTTL = 10 * time.Minute

// NewEthProtocol creates the new eth protocol. This will handle writing the
// status exchange, message handling, and writing blocks/txs to the database.
func NewEthProtocol(version uint, opts EthProtocolOptions) ethp2p.Protocol {
	return ethp2p.Protocol{
		Name:    "eth",
		Version: version,
		Length:  17,
		Run: func(p *ethp2p.Peer, rw ethp2p.MsgReadWriter) error {
			c := &conn{
				sensorID:    opts.SensorID,
				node:        p.Node(),
				logger:      log.With().Str("peer", p.Node().URLv4()).Logger(),
				rw:          rw,
				db:          opts.Database,
				requests:    NewCache[uint64, common.Hash](opts.MaxRequests, opts.RequestsCacheTTL),
				requestNum:  0,
				head:        opts.Head,
				headMutex:   opts.HeadMutex,
				counter:     opts.MsgCounter,
				peer:        p,
				blockHashes: list.New(),
				connectedAt: time.Now(),
			}

			c.headMutex.RLock()
			status := eth.StatusPacket{
				ProtocolVersion: uint32(version),
				NetworkID:       opts.NetworkID,
				Genesis:         opts.GenesisHash,
				ForkID:          opts.ForkID,
				Head:            opts.Head.Hash,
				TD:              opts.Head.TotalDifficulty,
			}
			err := c.statusExchange(&status)
			c.headMutex.RUnlock()
			if err != nil {
				return err
			}

			// Send the connection object to the conns manager for RPC broadcasting
			opts.Conns.Add(c)
			defer opts.Conns.Remove(c)

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
func (c *conn) statusExchange(packet *eth.StatusPacket) error {
	errc := make(chan error, 2)

	go func() {
		c.countMsgSent((&eth.StatusPacket{}).Name(), 1)
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

// countMsg increments the prometheus counter for this connection with the given direction, message name, and count.
func (c *conn) countMsg(direction Direction, messageName string, count float64) {
	c.counter.WithLabelValues(messageName, c.node.URLv4(), c.peer.Fullname(), string(direction)).Add(count)
}

// countMsgReceived increments the prometheus counter for received messages.
func (c *conn) countMsgReceived(messageName string, count float64) {
	c.countMsg(MsgReceived, messageName, count)
	c.countMsg(MsgReceived, messageName+PacketSuffix, 1)
}

// countMsgSent increments the prometheus counter for sent messages.
func (c *conn) countMsgSent(messageName string, count float64) {
	c.countMsg(MsgSent, messageName, count)
	c.countMsg(MsgSent, messageName+PacketSuffix, 1)
}

// TrackGoroutine runs a function in a goroutine and tracks the active goroutine count.
// This implements the database.GoroutineTracker interface.
func (c *conn) TrackGoroutine(f func()) {
	c.activeGoroutines.Add(1)
	go func() {
		defer c.activeGoroutines.Add(-1)
		f()
	}()
}

// ActiveGoroutines returns the number of active goroutines for this peer.
func (c *conn) ActiveGoroutines() int64 {
	return c.activeGoroutines.Load()
}

func (c *conn) readStatus(packet *eth.StatusPacket) error {
	msg, err := c.rw.ReadMsg()
	if err != nil {
		return err
	}

	if msg.Code != eth.StatusMsg {
		return errors.New("expected status message code")
	}

	var status eth.StatusPacket
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

// getBlockData will send a GetBlockHeaders and GetBlockBodies request to the
// peer. It will return an error if the sending either of the requests failed.
func (c *conn) getBlockData(hash common.Hash) error {
	headersRequest := &GetBlockHeaders{
		GetBlockHeadersRequest: &eth.GetBlockHeadersRequest{
			// Providing both the hash and number will result in a `both origin
			// hash and number` error.
			Origin: eth.HashOrNumber{Hash: hash},
			Amount: 1,
		},
	}

	c.countMsgSent(headersRequest.Name(), 1)
	if err := ethp2p.Send(c.rw, eth.GetBlockHeadersMsg, headersRequest); err != nil {
		return err
	}

	c.requestNum++
	c.requests.Add(c.requestNum, hash)

	bodiesRequest := &GetBlockBodies{
		RequestId:             c.requestNum,
		GetBlockBodiesRequest: []common.Hash{hash},
	}

	c.countMsgSent(bodiesRequest.Name(), 1)
	return ethp2p.Send(c.rw, eth.GetBlockBodiesMsg, bodiesRequest)
}

// getParentBlock will send a request to the peer if the parent of the header
// does not exist in the database.
func (c *conn) getParentBlock(ctx context.Context, header *types.Header) error {
	if !c.db.ShouldWriteBlocks() || !c.db.ShouldWriteBlockEvents() {
		return nil
	}

	if c.oldestBlock == nil {
		c.logger.Info().Interface("block", header).Msg("Setting oldest block")
		c.oldestBlock = header
		return nil
	}

	if c.db.HasBlock(ctx, header.ParentHash) || header.Number.Cmp(c.oldestBlock.Number) != 1 {
		return nil
	}

	c.logger.Info().
		Str("hash", header.ParentHash.Hex()).
		Str("number", new(big.Int).Sub(header.Number, big.NewInt(1)).String()).
		Msg("Fetching missing parent block")

	return c.getBlockData(header.ParentHash)
}

func (c *conn) handleNewBlockHashes(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.NewBlockHashesPacket
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	tfs := time.Now()

	c.countMsgReceived(packet.Name(), float64(len(packet)))

	// Collect unique hashes for database write.
	uniqueHashes := make([]common.Hash, 0, len(packet))

	for _, entry := range packet {
		hash := entry.Hash

		// Check if we've seen the hash and remove old entries
		if c.hasSeenBlockHash(hash) {
			continue
		}

		// Attempt to fetch block data first
		if err := c.getBlockData(hash); err != nil {
			return err
		}

		// Now that we've successfully fetched, record the new block hash
		c.addBlockHash(hash)
		uniqueHashes = append(uniqueHashes, hash)
	}

	// Write only unique hashes to the database.
	if len(uniqueHashes) > 0 {
		c.db.WriteBlockHashes(ctx, c, c.node, uniqueHashes, tfs)
	}

	return nil
}

// addBlockHash adds a new block hash with a timestamp to the blockHashes list.
func (c *conn) addBlockHash(hash common.Hash) {
	now := time.Now()

	// Add the new block hash entry to the list.
	c.blockHashes.PushBack(BlockHashEntry{
		hash: hash,
		time: now,
	})
}

// Helper method to check if a block hash is already in blockHashes.
func (c *conn) hasSeenBlockHash(hash common.Hash) bool {
	now := time.Now()
	for e := c.blockHashes.Front(); e != nil; e = e.Next() {
		entry := e.Value.(BlockHashEntry)
		// Check if the hash matches. We can short circuit here because there will
		// be block hashes that we haven't seen before, which will make a full
		// iteration of the blockHashes linked list.
		if entry.hash.Cmp(hash) == 0 {
			return true
		}
		// Remove entries older than blockHashTTL.
		if now.Sub(entry.time) > blockHashTTL {
			c.blockHashes.Remove(e)
		}
	}
	return false
}

func (c *conn) handleTransactions(ctx context.Context, msg ethp2p.Msg) error {
	var txs eth.TransactionsPacket
	if err := msg.Decode(&txs); err != nil {
		return err
	}

	tfs := time.Now()

	c.countMsgReceived(txs.Name(), float64(len(txs)))

	c.db.WriteTransactions(ctx, c, c.node, txs, tfs)

	return nil
}

func (c *conn) handleGetBlockHeaders(msg ethp2p.Msg) error {
	var request eth.GetBlockHeadersPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), 1)

	response := &eth.BlockHeadersPacket{RequestId: request.RequestId}
	c.countMsgSent(response.Name(), 0)
	return ethp2p.Send(c.rw, eth.BlockHeadersMsg, response)
}

func (c *conn) handleBlockHeaders(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.BlockHeadersPacket
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	tfs := time.Now()

	headers := packet.BlockHeadersRequest
	c.countMsgReceived(packet.Name(), float64(len(headers)))

	for _, header := range headers {
		if err := c.getParentBlock(ctx, header); err != nil {
			return err
		}
	}

	c.db.WriteBlockHeaders(ctx, c, headers, tfs)
	return nil
}

func (c *conn) handleGetBlockBodies(msg ethp2p.Msg) error {
	var request eth.GetBlockBodiesPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), float64(len(request.GetBlockBodiesRequest)))

	response := &eth.BlockBodiesPacket{RequestId: request.RequestId}
	c.countMsgSent(response.Name(), 0)
	return ethp2p.Send(c.rw, eth.BlockBodiesMsg, response)
}

func (c *conn) handleBlockBodies(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.BlockBodiesPacket
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	tfs := time.Now()

	if len(packet.BlockBodiesResponse) == 0 {
		return nil
	}

	c.countMsgReceived(packet.Name(), float64(len(packet.BlockBodiesResponse)))

	hash, ok := c.requests.Get(packet.RequestId)
	if !ok {
		c.logger.Warn().Msg("No block hash found for block body")
		return nil
	}
	c.requests.Remove(packet.RequestId)

	c.db.WriteBlockBody(ctx, c, packet.BlockBodiesResponse[0], hash, tfs)

	return nil
}

func (c *conn) handleNewBlock(ctx context.Context, msg ethp2p.Msg) error {
	var block eth.NewBlockPacket
	if err := msg.Decode(&block); err != nil {
		return err
	}

	tfs := time.Now()

	c.countMsgReceived(block.Name(), 1)

	// Set the head block if newer.
	c.headMutex.Lock()
	if block.Block.Number().Uint64() > c.head.Number && block.TD.Cmp(c.head.TotalDifficulty) == 1 {
		*c.head = HeadBlock{
			Hash:            block.Block.Hash(),
			TotalDifficulty: block.TD,
			Number:          block.Block.Number().Uint64(),
			Time:            block.Block.Time(),
		}
		c.logger.Info().Interface("head", c.head).Msg("Setting head block")
	}
	c.headMutex.Unlock()

	if err := c.getParentBlock(ctx, block.Block.Header()); err != nil {
		return err
	}

	c.db.WriteBlock(ctx, c, c.node, block.Block, block.TD, tfs)

	return nil
}

func (c *conn) handleGetPooledTransactions(msg ethp2p.Msg) error {
	var request eth.GetPooledTransactionsPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), float64(len(request.GetPooledTransactionsRequest)))

	response := &eth.PooledTransactionsPacket{RequestId: request.RequestId}
	c.countMsgSent(response.Name(), 0)
	return ethp2p.Send(c.rw, eth.PooledTransactionsMsg, response)
}

func (c *conn) handleNewPooledTransactionHashes(version uint, msg ethp2p.Msg) error {
	var hashes []common.Hash
	var name string

	switch version {
	case 67, 68:
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
	var packet eth.PooledTransactionsPacket
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	tfs := time.Now()

	c.countMsgReceived(packet.Name(), float64(len(packet.PooledTransactionsResponse)))

	c.db.WriteTransactions(ctx, c, c.node, packet.PooledTransactionsResponse, tfs)

	return nil
}

func (c *conn) handleGetReceipts(msg ethp2p.Msg) error {
	var request eth.GetReceiptsPacket
	if err := msg.Decode(&request); err != nil {
		return err
	}

	c.countMsgReceived(request.Name(), float64(len(request.GetReceiptsRequest)))

	response := &eth.ReceiptsPacket{RequestId: request.RequestId}
	c.countMsgSent(response.Name(), 0)
	return ethp2p.Send(c.rw, eth.ReceiptsMsg, response)
}
