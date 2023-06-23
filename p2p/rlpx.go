package p2p

import (
	"container/list"
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/maticnetwork/polygon-cli/p2p/database"
	"github.com/rs/zerolog/log"
)

var (
	timeout = 20 * time.Second
)

// Dial attempts to Dial the given node and perform a handshake,
// returning the created Conn if successful.
func Dial(n *enode.Node) (*Conn, error) {
	fd, err := net.Dial("tcp", fmt.Sprintf("%v:%d", n.IP(), n.TCP()))
	if err != nil {
		return nil, err
	}

	conn := Conn{
		Conn:       rlpx.NewConn(fd, n.Pubkey()),
		node:       n,
		logger:     log.With().Str("peer", n.URLv4()).Logger(),
		requests:   list.New(),
		requestNum: 0,
	}

	if conn.ourKey, err = crypto.GenerateKey(); err != nil {
		return nil, err
	}

	defer func() { _ = conn.SetDeadline(time.Time{}) }()
	if err = conn.SetDeadline(time.Now().Add(20 * time.Second)); err != nil {
		return nil, err
	}
	if _, err = conn.Handshake(conn.ourKey); err != nil {
		conn.Close()
		return nil, err
	}

	conn.caps = []p2p.Cap{
		{Name: "eth", Version: 66},
	}

	return &conn, nil
}

// Peer performs both the protocol handshake and the status message
// exchange with the node in order to Peer with it.
func (c *Conn) Peer() (*Hello, *Status, error) {
	hello, err := c.handshake()
	if err != nil {
		return nil, nil, fmt.Errorf("handshake failed: %v", err)
	}
	status, err := c.statusExchange()
	if err != nil {
		return hello, nil, fmt.Errorf("status exchange failed: %v", err)
	}
	return hello, status, nil
}

// handshake performs a protocol handshake with the node.
func (c *Conn) handshake() (*Hello, error) {
	defer func() { _ = c.SetDeadline(time.Time{}) }()
	if err := c.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, err
	}

	// write hello to client
	pub0 := crypto.FromECDSAPub(&c.ourKey.PublicKey)[1:]
	ourHandshake := &Hello{
		Version: 5,
		Caps:    c.caps,
		ID:      pub0,
	}
	if err := c.Write(ourHandshake); err != nil {
		return nil, fmt.Errorf("write to connection failed: %v", err)
	}

	// read hello from client
	switch msg := c.Read().(type) {
	case *Hello:
		if msg.Version >= 5 {
			c.SetSnappy(true)
		}
		return msg, nil
	case *Disconnect:
		return nil, fmt.Errorf("disconnect received: %v", msg)
	case *Disconnects:
		return nil, fmt.Errorf("disconnect received: %v", msg)
	default:
		return nil, fmt.Errorf("bad handshake: %v", msg)
	}
}

// statusExchange gets the Status message from the given node.
func (c *Conn) statusExchange() (*Status, error) {
	defer func() { _ = c.SetDeadline(time.Time{}) }()
	if err := c.SetDeadline(time.Now().Add(20 * time.Second)); err != nil {
		return nil, err
	}

	var status *Status
loop:
	for {
		switch msg := c.Read().(type) {
		case *Status:
			status = msg
			break loop
		case *Disconnect:
			return nil, fmt.Errorf("disconnect received: %v", msg)
		case *Disconnects:
			return nil, fmt.Errorf("disconnect received: %v", msg)
		case *Ping:
			if err := c.Write(&Pong{}); err != nil {
				c.logger.Error().Err(err).Msg("Write pong failed")
			}
		default:
			return nil, fmt.Errorf("bad status message: %v", msg)
		}
	}

	if err := c.Write(status); err != nil {
		return nil, fmt.Errorf("write to connection failed: %v", err)
	}

	return status, nil
}

// request stores the request ID and the block's hash.
type request struct {
	requestID uint64
	hash      common.Hash
}

// ReadAndServe reads messages from peers and writes it to a database.
func (c *Conn) ReadAndServe(db database.Database, count *MessageCount) error {
	// dbCh is used to limit the number of database goroutines running at one
	// time with a buffered channel. Without this, a large influx of messages can
	// bog down the system and leak memory.
	var dbCh chan struct{}
	if db != nil {
		dbCh = make(chan struct{}, db.MaxConcurrentWrites())
	}

	ctx := context.Background()

	for {
		start := time.Now()

		for time.Since(start) < timeout {
			if err := c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
				c.logger.Error().Err(err).Msg("Failed to set read deadline")
			}

			switch msg := c.Read().(type) {
			case *Ping:
				atomic.AddInt32(&count.Pings, 1)
				c.logger.Trace().Msg("Received Ping")

				if err := c.Write(&Pong{}); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write Pong response")
				}
			case *BlockHeaders:
				atomic.AddInt32(&count.BlockHeaders, int32(len(msg.BlockHeadersPacket)))
				c.logger.Trace().Msgf("Received %v BlockHeaders", len(msg.BlockHeadersPacket))

				if db != nil && db.ShouldWriteBlocks() {
					for _, header := range msg.BlockHeadersPacket {
						if err := c.getParentBlock(ctx, db, header); err != nil {
							return err
						}
					}

					dbCh <- struct{}{}
					go func() {
						db.WriteBlockHeaders(ctx, msg.BlockHeadersPacket)
						<-dbCh
					}()
				}
			case *GetBlockHeaders:
				atomic.AddInt32(&count.BlockHeaderRequests, 1)
				c.logger.Trace().Msgf("Received GetBlockHeaders request")

				res := &BlockHeaders{
					RequestId: msg.RequestId,
				}
				if err := c.Write(res); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write BlockHeaders response")
					return err
				}
			case *BlockBodies:
				atomic.AddInt32(&count.BlockBodies, int32(len(msg.BlockBodiesPacket)))
				c.logger.Trace().Msgf("Received %v BlockBodies", len(msg.BlockBodiesPacket))

				var hash *common.Hash
				for e := c.requests.Front(); e != nil; e = e.Next() {
					r, ok := e.Value.(request)
					if !ok {
						log.Error().Msg("Request type assertion failed")
						continue
					}

					if r.requestID == msg.ReqID() {
						hash = &r.hash
						c.requests.Remove(e)
						break
					}
				}

				if hash == nil {
					c.logger.Warn().Msg("No block hash found for block body")
					break
				}

				if db != nil && len(msg.BlockBodiesPacket) > 0 && db.ShouldWriteBlocks() {
					dbCh <- struct{}{}
					go func() {
						db.WriteBlockBody(ctx, msg.BlockBodiesPacket[0], *hash)
						<-dbCh
					}()
				}
			case *GetBlockBodies:
				atomic.AddInt32(&count.BlockBodiesRequests, int32(len(msg.GetBlockBodiesPacket)))
				c.logger.Trace().Msgf("Received %v GetBlockBodies request", len(msg.GetBlockBodiesPacket))

				res := &BlockBodies{
					RequestId: msg.RequestId,
				}
				if err := c.Write(res); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write BlockBodies response")
				}
			case *NewBlockHashes:
				atomic.AddInt32(&count.BlockHashes, int32(len(*msg)))
				c.logger.Trace().Msgf("Received %v NewBlockHashes", len(*msg))

				hashes := make([]common.Hash, 0, len(*msg))
				for _, hash := range *msg {
					hashes = append(hashes, hash.Hash)
					if err := c.getBlockData(hash.Hash); err != nil {
						return err
					}
				}

				if db != nil && db.ShouldWriteBlockEvents() && len(hashes) > 0 {
					dbCh <- struct{}{}
					go func() {
						db.WriteBlockHashes(ctx, c.node, hashes)
						<-dbCh
					}()
				}
			case *NewBlock:
				atomic.AddInt32(&count.Blocks, 1)
				c.logger.Trace().Str("hash", msg.Block.Hash().Hex()).Msg("Received NewBlock")

				if db != nil && (db.ShouldWriteBlocks() || db.ShouldWriteBlockEvents()) {
					if err := c.getParentBlock(ctx, db, msg.Block.Header()); err != nil {
						return err
					}

					dbCh <- struct{}{}
					go func() {
						db.WriteBlock(ctx, c.node, msg.Block, msg.TD)
						<-dbCh
					}()
				}
			case *Transactions:
				atomic.AddInt32(&count.Transactions, int32(len(*msg)))
				c.logger.Trace().Msgf("Received %v Transactions", len(*msg))

				if db != nil && (db.ShouldWriteTransactions() || db.ShouldWriteTransactionEvents()) {
					dbCh <- struct{}{}
					go func() {
						db.WriteTransactions(ctx, c.node, *msg)
						<-dbCh
					}()
				}
			case *PooledTransactions:
				atomic.AddInt32(&count.Transactions, int32(len(msg.PooledTransactionsPacket)))
				c.logger.Trace().Msgf("Received %v PooledTransactions", len(msg.PooledTransactionsPacket))

				if db != nil && (db.ShouldWriteTransactions() || db.ShouldWriteTransactionEvents()) {
					dbCh <- struct{}{}
					go func() {
						db.WriteTransactions(ctx, c.node, msg.PooledTransactionsPacket)
						<-dbCh
					}()
				}
			case *NewPooledTransactionHashes:
				if err := c.processNewPooledTransactionHashes(count, msg.Hashes); err != nil {
					return err
				}
			case *NewPooledTransactionHashes66:
				if err := c.processNewPooledTransactionHashes(count, *msg); err != nil {
					return err
				}
			case *GetPooledTransactions:
				atomic.AddInt32(&count.TransactionRequests, int32(len(msg.GetPooledTransactionsPacket)))
				c.logger.Trace().Msgf("Received %v GetPooledTransactions request", len(msg.GetPooledTransactionsPacket))

				res := &PooledTransactions{
					RequestId: msg.RequestId,
				}
				if err := c.Write(res); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write PooledTransactions response")
				}
			case *Error:
				atomic.AddInt32(&count.Errors, 1)
				c.logger.Trace().Err(msg.Unwrap()).Msg("Received Error")

				if !strings.Contains(msg.Error(), "timeout") {
					return msg.Unwrap()
				}
			case *Disconnect:
				atomic.AddInt32(&count.Disconnects, 1)
				c.logger.Debug().Msgf("Disconnect received: %v", msg)
			case *Disconnects:
				atomic.AddInt32(&count.Disconnects, 1)
				c.logger.Debug().Msgf("Disconnect received: %v", msg)
			default:
				c.logger.Info().Interface("msg", msg).Int("code", msg.Code()).Msg("Received message")
			}
		}
	}
}

// processNewPooledTransactionHashes processes NewPooledTransactionHashes
// messages by requesting the transaction bodies.
func (c *Conn) processNewPooledTransactionHashes(count *MessageCount, hashes []common.Hash) error {
	atomic.AddInt32(&count.TransactionHashes, int32(len(hashes)))
	c.logger.Trace().Msgf("Received %v NewPooledTransactionHashes", len(hashes))

	req := &GetPooledTransactions{
		RequestId:                   rand.Uint64(),
		GetPooledTransactionsPacket: hashes,
	}
	if err := c.Write(req); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write GetPooledTransactions request")
		return err
	}

	return nil
}

// getBlockData will send a GetBlockHeaders and GetBlockBodies request to the
// peer. It will return an error if the sending either of the requests failed.
func (c *Conn) getBlockData(hash common.Hash) error {
	headersRequest := &GetBlockHeaders{
		GetBlockHeadersPacket: &eth.GetBlockHeadersPacket{
			// Providing both the hash and number will result in a `both origin
			// hash and number` error.
			Origin: eth.HashOrNumber{Hash: hash},
			Amount: 1,
		},
	}

	if err := c.Write(headersRequest); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write GetBlockHeaders request")
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
	if err := c.Write(bodiesRequest); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write GetBlockBodies request")
		return err
	}

	return nil
}

// getParentBlock will send a request to the peer if the parent of the header
// does not exist in the database.
func (c *Conn) getParentBlock(ctx context.Context, db database.Database, header *types.Header) error {
	if c.oldestBlock == nil {
		c.logger.Info().Interface("block", header).Msg("Setting oldest block")
		c.oldestBlock = header
		return nil
	}

	if !db.HasParentBlock(ctx, header.ParentHash) && header.Number.Cmp(c.oldestBlock.Number) == 1 {
		c.logger.Info().
			Str("hash", header.ParentHash.Hex()).
			Str("number", new(big.Int).Sub(header.Number, big.NewInt(1)).String()).
			Msg("Fetching missing parent block")

		return c.getBlockData(header.ParentHash)
	}

	return nil
}
