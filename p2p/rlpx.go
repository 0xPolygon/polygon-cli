// Copyright 2021 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package p2p

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/rs/zerolog/log"
)

const (
	MaxNumRequests = 1000
)

var (
	timeout = 20 * time.Second
)

type MessageCounter struct {
	BlockHeaders        int `json:",omitempty"`
	BlockBodies         int `json:",omitempty"`
	Blocks              int `json:",omitempty"`
	BlockHashes         int `json:",omitempty"`
	BlockHeaderRequest  int `json:",omitempty"`
	BlockBodiesRequests int `json:",omitempty"`
	Transactions        int `json:",omitempty"`
	TransactionHashes   int `json:",omitempty"`
	TransactionRequests int `json:",omitempty"`
	Pings               int `json:",omitempty"`
	Errors              int `json:",omitempty"`
	Disconnects         int `json:",omitempty"`
}

// Dial attempts to Dial the given node and perform a handshake,
// returning the created Conn if successful.
func Dial(n *enode.Node) (*Conn, error) {
	fd, err := net.Dial("tcp", fmt.Sprintf("%v:%d", n.IP(), n.TCP()))
	if err != nil {
		return nil, err
	}

	conn := Conn{
		Conn:   rlpx.NewConn(fd, n.Pubkey()),
		node:   n,
		logger: log.With().Str("node", n.String()).Logger(),
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

// ReadAndServe reads messages from peers.
func (c *Conn) ReadAndServe(client *datastore.Client) *Error {
	ctx := context.Background()
	requests := make(map[uint64]common.Hash)
	var count uint64 = 0

	done := make(chan struct{})

	counter := MessageCounter{}
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.logger.Info().Interface("messages", counter).Send()
				counter = MessageCounter{}
			case <-done:
				return
			}
		}
	}()

	type BodyMessage struct {
		body *eth.BlockBody
		hash common.Hash
	}
	blockChan := make(chan *types.Block)
	blockHeaderChan := make(chan *types.Header)
	blockBodyChan := make(chan BodyMessage)
	blockHashesChan := make(chan []common.Hash)
	transactionsChan := make(chan []*types.Transaction)
	if client != nil {
		go func() {
			for {
				select {
				case block := <-blockChan:
					c.writeEvent(ctx, client, "block_events", block.Hash(), "blocks")
					c.writeBlockHeader(ctx, client, block.Header())
					c.writeBlockBody(ctx, client, block.Hash().Hex(),
						&eth.BlockBody{
							Transactions: block.Transactions(),
							Uncles:       block.Uncles(),
						},
					)
				case header := <-blockHeaderChan:
					c.writeBlockHeader(ctx, client, header)
				case blockBody := <-blockBodyChan:
					c.writeBlockBody(ctx, client, blockBody.hash.Hex(), blockBody.body)
				case hashes := <-blockHashesChan:
					c.writeEvents(ctx, client, "block_events", hashes, "blocks")
				case transactions := <-transactionsChan:
					c.writeTransactions(ctx, client, transactions)
				case <-done:
					return
				}
			}
		}()
	}

	for {
		start := time.Now()

	loop:
		for time.Since(start) < timeout {
			if err := c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
				c.logger.Error().Err(err).Msg("Failed to set read deadline")
			}

			switch msg := c.Read().(type) {
			case *Ping:
				counter.Pings++
				if err := c.Write(&Pong{}); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write Pong response")
				}
			case *BlockHeaders:
				counter.BlockHeaders++
				for _, header := range msg.BlockHeadersPacket {
					blockHeaderChan <- header
				}
			case *GetBlockHeaders:
				counter.BlockHeaderRequest++
				res := &BlockHeaders{
					RequestId: msg.RequestId,
				}
				if err := c.Write(res); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write BlockHeaders response")
				}
			case *BlockBodies:
				counter.BlockBodies += len(msg.BlockBodiesPacket)
				if hash, ok := requests[msg.RequestId]; ok {
					if len(msg.BlockBodiesPacket) > 0 {
						blockBodyChan <- BodyMessage{
							hash: hash,
							body: msg.BlockBodiesPacket[0],
						}
					}
					delete(requests, msg.RequestId)
				}
			case *GetBlockBodies:
				counter.BlockBodiesRequests++
				res := &BlockBodies{
					RequestId: msg.RequestId,
				}
				if err := c.Write(res); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write BlockBodies response")
				}
			case *NewBlockHashes:
				counter.BlockHashes += len(*msg)

				hashes := make([]common.Hash, 0, len(*msg))
				for _, hash := range *msg {
					hashes = append(hashes, hash.Hash)

					headersRequest := &GetBlockHeaders{
						GetBlockHeadersPacket: &eth.GetBlockHeadersPacket{
							// Providing both the hash and number will result in a `both origin
							// hash and number` error.
							Origin: eth.HashOrNumber{Hash: hash.Hash},
							Amount: 1,
						},
					}
					if err := c.Write(headersRequest); err != nil {
						c.logger.Error().Err(err).Msg("Failed to write GetBlockHeaders request")
						break loop
					}

					count++
					if count > MaxNumRequests {
						count = 0
					}
					requests[count] = hash.Hash
					bodiesRequest := &GetBlockBodies{
						RequestId:            count,
						GetBlockBodiesPacket: []common.Hash{hash.Hash},
					}
					if err := c.Write(bodiesRequest); err != nil {
						c.logger.Error().Err(err).Msg("Failed to write GetBlockBodies request")
						break loop
					}
				}

				blockHashesChan <- hashes
			case *NewBlock:
				counter.Blocks++
				blockChan <- msg.Block
			case *Transactions:
				counter.Transactions += len(*msg)
				transactionsChan <- *msg
			case *PooledTransactions:
				counter.Transactions += len(msg.PooledTransactionsPacket)
				transactionsChan <- msg.PooledTransactionsPacket
			case *NewPooledTransactionHashes:
				c.processNewPooledTransactions(ctx, client, &counter, msg.Hashes)
			case *NewPooledTransactionHashes66:
				c.processNewPooledTransactions(ctx, client, &counter, *msg)
			case *GetPooledTransactions:
				c.logger.Info().Interface("msg", msg).Msg("Received GetPooledTransactions request")
				res := &PooledTransactions{
					RequestId: msg.RequestId,
				}
				if err := c.Write(res); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write PooledTransactions response")
				}
			case *Error:
				counter.Errors++
				if !strings.Contains(msg.Error(), "timeout") {
					ticker.Stop()
					close(done)
					return msg
				}
			case *Disconnect:
				counter.Disconnects++
				c.logger.Debug().Msgf("Disconnect received: %v", msg)
			case *Disconnects:
				counter.Disconnects++
				c.logger.Debug().Msgf("Disconnect received: %v", msg)
			default:
				c.logger.Info().Interface("msg", msg).Int("code", msg.Code()).Msg("Received message")
			}
		}
	}
}

func (c *Conn) processNewPooledTransactions(ctx context.Context, client *datastore.Client, counter *MessageCounter, hashes []common.Hash) {
	counter.TransactionHashes += len(hashes)
	req := &GetPooledTransactions{
		RequestId:                   0,
		GetPooledTransactionsPacket: hashes,
	}
	if err := c.Write(req); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write GetPooledTransactions request")
	}
}
