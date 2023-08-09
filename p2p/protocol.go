package p2p

import (
	"container/list"
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/maticnetwork/polygon-cli/p2p/database"
	"github.com/maticnetwork/polygon-cli/rpctypes"
)

// conn represents an individual connection with a peer.
type conn struct {
	sensorID string
	node     *enode.Node
	logger   zerolog.Logger
	rw       ethp2p.MsgReadWriter
	db       database.Database

	// requests is used to store the request ID and the block hash. This is used
	// when fetching block bodies because the eth protocol block bodies do not
	// contain information about the block hash.
	requests   *list.List
	requestNum uint64

	// oldestBlock stores the first block the sensor has seen so when fetching
	// parent blocks, it does not request blocks older than this.
	oldestBlock *types.Header
}

type Eth66ProtocolOptions struct {
	Context     context.Context
	Database    database.Database
	Genesis     *core.Genesis
	GenesisHash common.Hash
	RPC         string
	SensorID    string
	NetworkID   uint64
	Peers       chan *enode.Node
}

func NewEth66Protocol(opts Eth66ProtocolOptions) ethp2p.Protocol {
	return ethp2p.Protocol{
		Name:    "eth",
		Version: 66,
		Length:  17,
		Run: func(p *ethp2p.Peer, rw ethp2p.MsgReadWriter) error {
			block, err := getLatestBlock(opts.RPC)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get latest block")
				return err
			}

			c := conn{
				sensorID:   opts.SensorID,
				node:       p.Node(),
				logger:     log.With().Str("peer", p.Node().URLv4()).Logger(),
				rw:         rw,
				db:         opts.Database,
				requests:   list.New(),
				requestNum: 0,
			}

			status := eth.StatusPacket{
				ProtocolVersion: 66,
				NetworkID:       opts.NetworkID,
				Genesis:         opts.GenesisHash,
				ForkID:          forkid.NewID(opts.Genesis.Config, opts.GenesisHash, block.Number.ToUint64()),
				Head:            block.Hash.ToHash(),
				TD:              block.TotalDifficulty.ToBigInt(),
			}
			if err = c.statusExchange(&status); err != nil {
				return err
			}

			opts.Peers <- p.Node()
			ctx := opts.Context

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
					err = c.handleNewPooledTransactionHashes(ctx, msg)
				case eth.GetPooledTransactionsMsg:
					err = c.handleGetPooledTransactions(msg)
				case eth.PooledTransactionsMsg:
					err = c.handlePooledTransactions(ctx, msg)
				case eth.GetReceiptsMsg:
					err = c.handleGetReceipts(msg)
				default:
					log.Trace().Interface("msg", msg).Send()
				}

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

func getLatestBlock(url string) (*rpctypes.RawBlockResponse, error) {
	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var block rpctypes.RawBlockResponse
	err = client.Call(&block, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (c *conn) statusExchange(packet *eth.StatusPacket) error {
	err := ethp2p.Send(c.rw, eth.StatusMsg, &packet)
	if err != nil {
		return err
	}

	msg, err := c.rw.ReadMsg()
	if err != nil {
		return err
	}

	var status eth.StatusPacket
	err = msg.Decode(&status)
	if err != nil {
		return err
	}

	if status.NetworkID != packet.NetworkID {
		return errors.New("network IDs mismatch")
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

	c.db.WriteTransactions(ctx, c.node, txs)

	return nil
}

func (c *conn) handleGetBlockHeaders(msg ethp2p.Msg) error {
	var request eth.GetBlockHeadersPacket66
	if err := msg.Decode(&request); err != nil {
		return err
	}
	return ethp2p.Send(
		c.rw,
		eth.GetBlockHeadersMsg,
		&eth.BlockHeadersPacket66{RequestId: request.RequestId},
	)
}

func (c *conn) handleBlockHeaders(ctx context.Context, msg ethp2p.Msg) error {
	var packet eth.BlockHeadersPacket66
	if err := msg.Decode(&packet); err != nil {
		return err
	}

	headers := packet.BlockHeadersPacket
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
	return ethp2p.Send(
		c.rw,
		eth.GetBlockHeadersMsg,
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
	return ethp2p.Send(c.rw, eth.GetPooledTransactionsMsg, &eth.PooledTransactionsPacket66{
		RequestId: request.RequestId,
	})
}

func (c *conn) handleNewPooledTransactionHashes(ctx context.Context, msg ethp2p.Msg) error {
	var txs eth.NewPooledTransactionHashesPacket
	if err := msg.Decode(&txs); err != nil {
		return err
	}

	if !c.db.ShouldWriteTransactions() || !c.db.ShouldWriteTransactionEvents() {
		return nil
	}

	var hashes []common.Hash = txs

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

	c.db.WriteTransactions(ctx, c.node, packet.PooledTransactionsPacket)

	return nil
}

func (c *conn) handleGetReceipts(msg ethp2p.Msg) error {
	var request eth.GetReceiptsPacket66
	if err := msg.Decode(&request); err != nil {
		return err
	}
	return ethp2p.Send(c.rw, eth.GetBlockHeadersMsg, &eth.ReceiptsPacket66{
		RequestId: request.RequestId,
	})
}
