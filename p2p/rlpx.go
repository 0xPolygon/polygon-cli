package p2p

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/rs/zerolog/log"
)

var (
	timeout = 20 * time.Second
)

// Dial attempts to Dial the given node and perform a handshake,
// returning the created Conn if successful.
func Dial(n *enode.Node) (*rlpxConn, error) {
	fd, err := net.Dial("tcp", fmt.Sprintf("%v:%d", n.IP(), n.TCP()))
	if err != nil {
		return nil, err
	}

	conn := rlpxConn{
		Conn:   rlpx.NewConn(fd, n.Pubkey()),
		node:   n,
		logger: log.With().Str("peer", n.URLv4()).Logger(),
		caps: []p2p.Cap{
			{Name: "eth", Version: 66},
			{Name: "eth", Version: 67},
			{Name: "eth", Version: 68},
		},
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

	return &conn, nil
}

// Peer performs both the protocol handshake and the status message
// exchange with the node in order to Peer with it.
func (c *rlpxConn) Peer() (*Hello, *Status, error) {
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
func (c *rlpxConn) handshake() (*Hello, error) {
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
func (c *rlpxConn) statusExchange() (*Status, error) {
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
	time      time.Time
}

// ReadAndServe reads messages from peers and writes it to a database.
func (c *rlpxConn) ReadAndServe(count *MessageCount) error {
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
				atomic.AddInt32(&count.BlockHeaders, int32(len(msg.BlockHeadersRequest)))
				c.logger.Trace().Msgf("Received %v BlockHeaders", len(msg.BlockHeadersRequest))
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
				atomic.AddInt32(&count.BlockBodies, int32(len(msg.BlockBodiesResponse)))
				c.logger.Trace().Msgf("Received %v BlockBodies", len(msg.BlockBodiesResponse))
			case *GetBlockBodies:
				atomic.AddInt32(&count.BlockBodiesRequests, int32(len(msg.GetBlockBodiesRequest)))
				c.logger.Trace().Msgf("Received %v GetBlockBodies request", len(msg.GetBlockBodiesRequest))

				res := &BlockBodies{
					RequestId: msg.RequestId,
				}
				if err := c.Write(res); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write BlockBodies response")
				}
			case *NewBlockHashes:
				atomic.AddInt32(&count.BlockHashes, int32(len(*msg)))
				c.logger.Trace().Msgf("Received %v NewBlockHashes", len(*msg))

				for _, hash := range *msg {
					headersRequest := &GetBlockHeaders{
						GetBlockHeadersRequest: &eth.GetBlockHeadersRequest{
							// Providing both the hash and number will result in a `both origin
							// hash and number` error.
							Origin: eth.HashOrNumber{Hash: hash.Hash},
							Amount: 1,
						},
					}

					if err := c.Write(headersRequest); err != nil {
						c.logger.Error().Err(err).Msg("Failed to write GetBlockHeaders request")
					}

					bodiesRequest := &GetBlockBodies{
						GetBlockBodiesRequest: []common.Hash{hash.Hash},
					}

					if err := c.Write(bodiesRequest); err != nil {
						c.logger.Error().Err(err).Msg("Failed to write GetBlockBodies request")
					}
				}

			case *NewBlock:
				atomic.AddInt32(&count.Blocks, 1)
				c.logger.Trace().Str("hash", msg.Block.Hash().Hex()).Msg("Received NewBlock")
			case *Transactions:
				atomic.AddInt32(&count.Transactions, int32(len(*msg)))
				c.logger.Trace().Msgf("Received %v Transactions", len(*msg))
			case *PooledTransactions:
				atomic.AddInt32(&count.Transactions, int32(len(msg.PooledTransactionsResponse)))
				c.logger.Trace().Msgf("Received %v PooledTransactions", len(msg.PooledTransactionsResponse))
			case *NewPooledTransactionHashes:
				if err := c.processNewPooledTransactionHashes(count, msg.Hashes); err != nil {
					return err
				}
			case *NewPooledTransactionHashes66:
				if err := c.processNewPooledTransactionHashes(count, *msg); err != nil {
					return err
				}
			case *GetPooledTransactions:
				atomic.AddInt32(&count.TransactionRequests, int32(len(msg.GetPooledTransactionsRequest)))
				c.logger.Trace().Msgf("Received %v GetPooledTransactions request", len(msg.GetPooledTransactionsRequest))

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
func (c *rlpxConn) processNewPooledTransactionHashes(count *MessageCount, hashes []common.Hash) error {
	atomic.AddInt32(&count.TransactionHashes, int32(len(hashes)))
	c.logger.Trace().Msgf("Received %v NewPooledTransactionHashes", len(hashes))

	req := &GetPooledTransactions{
		RequestId:                    rand.Uint64(),
		GetPooledTransactionsRequest: hashes,
	}
	if err := c.Write(req); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write GetPooledTransactions request")
		return err
	}

	return nil
}

// QueryHeaders requests block headers given start and count
func (c *rlpxConn) QueryHeaders(start, count uint64) error {
	c.logger.Trace().Msgf("Querying headers from %v for %v blocks", start, count)

	// Prepare the `GetBlockHeaders` request.
	req := &GetBlockHeaders{
		RequestId: rand.Uint64(),
		GetBlockHeadersRequest: &eth.GetBlockHeadersRequest{
			Origin: eth.HashOrNumber{Number: start},
			Amount: count,
		},
	}
	if err := c.Write(req); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write GetBlockHeaders request")
		return err
	}

	return nil
}

// ListenHeaders keeps listening for requested headers from the p2p connection.
func (c *rlpxConn) ListenHeaders() (eth.BlockHeadersRequest, error) {
	for {
		start := time.Now()

		for time.Since(start) < timeout {
			if err := c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
				c.logger.Error().Err(err).Msg("Failed to set read deadline")
			}

			switch msg := c.Read().(type) {
			case *BlockHeaders:
				c.logger.Trace().Msgf("Received %v BlockHeaders", len(msg.BlockHeadersRequest))
				return msg.BlockHeadersRequest, nil
			case *Error:
				c.logger.Trace().Err(msg.Unwrap()).Msg("Received Error")

				if !strings.Contains(msg.Error(), "timeout") {
					return nil, msg.Unwrap()
				}
			case *Disconnect:
				c.logger.Debug().Msgf("Disconnect received: %v", msg)
				return nil, fmt.Errorf("disconnect received: %v", msg)
			case *Disconnects:
				c.logger.Debug().Msgf("Disconnect received: %v", msg)
				return nil, fmt.Errorf("disconnect received: %v", msg)
			default:
				c.logger.Trace().Interface("msg", msg).Int("code", msg.Code()).Msg("Received message")
			}
		}
	}
}
