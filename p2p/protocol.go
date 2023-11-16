package p2p

import (
	"container/list"
	"context"
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/maticnetwork/polygon-cli/p2p/database"
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
	count     *MessageCount

	// messages keeps track of when a peer sends a message that the sensor cares
	// about. If the sensor doesn't receive a message sent to this channel in a
	// certain amount of time, it will disconnect with the peer.
	messages chan struct{}
	done     chan error

	// requests is used to store the request ID and the block hash. This is used
	// when fetching block bodies because the eth protocol block bodies do not
	// contain information about the block hash.
	requests   *list.List
	requestNum uint64

	// oldestBlock stores the first block the sensor has seen so when fetching
	// parent blocks, it does not request blocks older than this.
	oldestBlock *types.Header
}

// EthProtocolOptions is the options used when creating a new eth protocol.
type EthProtocolOptions struct {
	Context     context.Context
	Database    database.Database
	Genesis     *core.Genesis
	GenesisHash common.Hash
	RPC         string
	SensorID    string
	NetworkID   uint64
	Peers       chan *enode.Node
	Count       *MessageCount

	// Head keeps track of the current head block of the chain. This is required
	// when doing the status exchange.
	Head      *HeadBlock
	HeadMutex *sync.RWMutex
}

// HeadBlock contains the necessary head block data for the status message.
type HeadBlock struct {
	Hash            common.Hash
	TotalDifficulty *big.Int
	Number          uint64
}

// NewEthProctocol creates the new eth protocol. This will handle writing the
// status exchange, message handling, and writing blocks/txs to the database.
func NewEthProtocol(version uint, opts EthProtocolOptions) ethp2p.Protocol {
	return ethp2p.Protocol{
		Name:    "eth",
		Version: version,
		Length:  17,
		Run: func(p *ethp2p.Peer, rw ethp2p.MsgReadWriter) error {
			c := conn{
				sensorID:   opts.SensorID,
				node:       p.Node(),
				logger:     log.With().Str("peer", p.Node().URLv4()).Logger(),
				rw:         rw,
				db:         opts.Database,
				requests:   list.New(),
				requestNum: 0,
				head:       opts.Head,
				headMutex:  opts.HeadMutex,
				count:      opts.Count,
			}

			c.headMutex.RLock()
			status := eth.StatusPacket{
				ProtocolVersion: uint32(version),
				NetworkID:       opts.NetworkID,
				Genesis:         opts.GenesisHash,
				ForkID:          forkid.NewID(opts.Genesis.Config, opts.Genesis.ToBlock(), opts.Head.Number, uint64(time.Now().Unix())),
				Head:            opts.Head.Hash,
				TD:              opts.Head.TotalDifficulty,
			}
			err := c.statusExchange(&status)
			c.headMutex.RUnlock()
			if err != nil {
				return err
			}

			// Send the node to the peers channel. This allows the peers to be captured
			// across all connections and written to the nodes.json file.
			opts.Peers <- p.Node()
			ctx := opts.Context

			timeout := 10 * time.Minute
			timer := time.NewTimer(timeout)

			// Handle all the of the messages here.
			go c.handleMessages(ctx, rw, version)

			for {
				select {
				case <-c.messages:
					timer.Reset(timeout)
				case <-timer.C:
					p.Disconnect(p2p.DiscUselessPeer)
					return errors.New("useless peer")
				case err := <-c.done:
					p.Disconnect(p2p.DiscQuitting)
					return err
				}
			}
		},
	}
}

func (c *conn) handleMessages(ctx context.Context, rw ethp2p.MsgReadWriter, version uint) {
	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			c.done <- err
			return
		}

		switch msg.Code {
		case eth.NewBlockHashesMsg:
			err = c.handleNewBlockHashes(ctx, msg)
			c.messages <- struct{}{}
		case eth.TransactionsMsg:
			err = c.handleTransactions(ctx, msg)
		case eth.GetBlockHeadersMsg:
			err = c.handleGetBlockHeaders(msg)
		case eth.BlockHeadersMsg:
			err = c.handleBlockHeaders(ctx, msg)
			c.messages <- struct{}{}
		case eth.GetBlockBodiesMsg:
			err = c.handleGetBlockBodies(msg)
		case eth.BlockBodiesMsg:
			err = c.handleBlockBodies(ctx, msg)
			c.messages <- struct{}{}
		case eth.NewBlockMsg:
			err = c.handleNewBlock(ctx, msg)
			c.messages <- struct{}{}
		case eth.NewPooledTransactionHashesMsg:
			err = c.handleNewPooledTransactionHashes(ctx, version, msg)
		case eth.GetPooledTransactionsMsg:
			err = c.handleGetPooledTransactions(msg)
		case eth.PooledTransactionsMsg:
			err = c.handlePooledTransactions(ctx, msg)
		case eth.GetReceiptsMsg:
			err = c.handleGetReceipts(msg)
		default:
			log.Trace().Interface("msg", msg).Send()
		}

		// All the handler functions are built in a way where returning an error
		// should drop the connection. If the connection shouldn't be dropped,
		// then return nil and log the error instead.
		if err != nil {
			c.logger.Error().Err(err).Send()
			c.done <- err
			return
		}

		if err = msg.Discard(); err != nil {
			c.done <- err
			return
		}
	}
}

// statusExchange will exchange status message between the nodes. It will return
// and error if the nodes are incompatible.
func (c *conn) statusExchange(packet *eth.StatusPacket) error {
	err := ethp2p.Send(c.rw, eth.StatusMsg, &packet)
	if err != nil {
		return err
	}

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
		return ethp2p.DiscUselessPeer
	}

	c.logger.Info().Interface("status", status).Msg("New peer")

	return nil
}

// getBlockData will send a GetBlockHeaders and GetBlockBodies request to the
// peer. It will return an error if the sending either of the requests failed.
func (c *conn) getBlockData(hash common.Hash) error {
	headersRequest := &GetBlockHeaders{
		GetBlockHeadersPacket: &eth.GetBlockHeadersPacket{
			// Providing both the hash and number will result in a `both origin
			// hash and number` error.
			Origin: eth.HashOrNumber{Hash: hash},
			Amount: 1,
		},
	}

	if err := ethp2p.Send(c.rw, eth.GetBlockHeadersMsg, headersRequest); err != nil {
		return err
	}

	c.requestNum++
	c.requests.PushBack(request{
		requestID: c.requestNum,
		hash:      hash,
	})
	bodiesRequest := &GetBlockBodies{
		RequestId:            c.requestNum,
		GetBlockBodiesPacket: []common.Hash{hash},
	}

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

	atomic.AddInt32(&c.count.BlockHashes, int32(len(packet)))

	hashes := make([]common.Hash, 0, len(packet))
	for _, hash := range packet {
		hashes = append(hashes, hash.Hash)
		if err := c.getBlockData(hash.Hash); err != nil {
			return err
		}
	}

	c.db.WriteBlockHashes(ctx, c.node, hashes)

	return nil
}

func (c *conn) handleTransactions(ctx context.Context, msg ethp2p.Msg) error {
	var txs eth.TransactionsPacket
	if err := msg.Decode(&txs); err != nil {
		return err
	}

	atomic.AddInt32(&c.count.Transactions, int32(len(txs)))

	c.db.WriteTransactions(ctx, c.node, txs)

	return nil
}

func (c *conn) handleGetBlockHeaders(msg ethp2p.Msg) error {
	var request eth.GetBlockHeadersPacket66
	if err := msg.Decode(&request); err != nil {
		return err
	}

	atomic.AddInt32(&c.count.BlockHeaderRequests, 1)

	return ethp2p.Send(
		c.rw,
		eth.BlockHeadersMsg,
		&eth.BlockHeadersPacket66{RequestId: request.RequestId},
	)
}

func (c *conn) handleBlockHeaders(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.BlockHeadersPacket66
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	headers := packet.BlockHeadersPacket
	atomic.AddInt32(&c.count.BlockHeaders, int32(len(headers)))

	for _, header := range headers {
		if err := c.getParentBlock(ctx, header); err != nil {
			return err
		}
	}

	c.db.WriteBlockHeaders(ctx, headers)

	return nil
}

func (c *conn) handleGetBlockBodies(msg ethp2p.Msg) error {
	var request eth.GetBlockBodiesPacket66
	if err := msg.Decode(&request); err != nil {
		return err
	}

	atomic.AddInt32(&c.count.BlockBodiesRequests, int32(len(request.GetBlockBodiesPacket)))

	return ethp2p.Send(
		c.rw,
		eth.BlockBodiesMsg,
		&eth.BlockBodiesPacket66{RequestId: request.RequestId},
	)
}

func (c *conn) handleBlockBodies(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.BlockBodiesPacket66
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	if len(packet.BlockBodiesPacket) == 0 {
		return nil
	}

	atomic.AddInt32(&c.count.BlockBodies, int32(len(packet.BlockBodiesPacket)))

	var hash *common.Hash
	for e := c.requests.Front(); e != nil; e = e.Next() {
		r, ok := e.Value.(request)
		if !ok {
			log.Error().Msg("Request type assertion failed")
			continue
		}

		if r.requestID == packet.RequestId {
			hash = &r.hash
			c.requests.Remove(e)
			break
		}
	}

	if hash == nil {
		c.logger.Warn().Msg("No block hash found for block body")
		return nil
	}

	c.db.WriteBlockBody(ctx, packet.BlockBodiesPacket[0], *hash)

	return nil
}

func (c *conn) handleNewBlock(ctx context.Context, msg ethp2p.Msg) error {
	var block eth.NewBlockPacket
	if err := msg.Decode(&block); err != nil {
		return err
	}

	atomic.AddInt32(&c.count.Blocks, 1)

	// Set the head block if newer.
	c.headMutex.Lock()
	if block.Block.Number().Uint64() > c.head.Number && block.TD.Cmp(c.head.TotalDifficulty) == 1 {
		*c.head = HeadBlock{
			Hash:            block.Block.Hash(),
			TotalDifficulty: block.TD,
			Number:          block.Block.Number().Uint64(),
		}

		c.logger.Info().Interface("head", c.head).Msg("Setting head block")
	}
	c.headMutex.Unlock()

	if err := c.getParentBlock(ctx, block.Block.Header()); err != nil {
		return err
	}

	c.db.WriteBlock(ctx, c.node, block.Block, block.TD)

	return nil
}

func (c *conn) handleGetPooledTransactions(msg ethp2p.Msg) error {
	var request eth.GetPooledTransactionsPacket66
	if err := msg.Decode(&request); err != nil {
		return err
	}

	atomic.AddInt32(&c.count.TransactionRequests, int32(len(request.GetPooledTransactionsPacket)))

	return ethp2p.Send(
		c.rw,
		eth.PooledTransactionsMsg,
		&eth.PooledTransactionsPacket66{RequestId: request.RequestId})
}

func (c *conn) handleNewPooledTransactionHashes(ctx context.Context, version uint, msg ethp2p.Msg) error {
	var hashes []common.Hash

	switch version {
	case 66, 67:
		var txs eth.NewPooledTransactionHashesPacket66
		if err := msg.Decode(&txs); err != nil {
			return err
		}
		hashes = txs
	case 68:
		var txs eth.NewPooledTransactionHashesPacket68
		if err := msg.Decode(&txs); err != nil {
			return err
		}
		hashes = txs.Hashes
	default:
		return errors.New("protocol version not found")
	}

	atomic.AddInt32(&c.count.TransactionHashes, int32(len(hashes)))

	if !c.db.ShouldWriteTransactions() || !c.db.ShouldWriteTransactionEvents() {
		return nil
	}

	return ethp2p.Send(
		c.rw,
		eth.GetPooledTransactionsMsg,
		&eth.GetPooledTransactionsPacket66{GetPooledTransactionsPacket: hashes},
	)
}

func (c *conn) handlePooledTransactions(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.PooledTransactionsPacket66
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	atomic.AddInt32(&c.count.Transactions, int32(len(packet.PooledTransactionsPacket)))

	c.db.WriteTransactions(ctx, c.node, packet.PooledTransactionsPacket)

	return nil
}

func (c *conn) handleGetReceipts(msg ethp2p.Msg) error {
	var request eth.GetReceiptsPacket66
	if err := msg.Decode(&request); err != nil {
		return err
	}
	return ethp2p.Send(
		c.rw,
		eth.ReceiptsMsg,
		&eth.ReceiptsPacket66{RequestId: request.RequestId},
	)
}
