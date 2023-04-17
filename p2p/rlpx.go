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
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
func Dial(n *enode.Node) (*Conn, error) {
	fd, err := net.Dial("tcp", fmt.Sprintf("%v:%d", n.IP(), n.TCP()))
	if err != nil {
		return nil, err
	}

	conn := Conn{
		Conn:   rlpx.NewConn(fd, n.Pubkey()),
		node:   n,
		logger: log.With().Interface("node", n).Logger(),
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
		// {Name: "eth", Version: 67},
		// {Name: "eth", Version: 68},
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
func (c *Conn) ReadAndServe() *Error {
	headers := make(map[common.Hash]*types.Header)
	for {
		start := time.Now()
		for time.Since(start) < timeout {
			if err := c.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
				c.logger.Error().Err(err).Msg("Failed to set read deadline")
			}

			msg := c.Read()
			switch msg := msg.(type) {
			case *Ping:
				if err := c.Write(&Pong{}); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write Pong response")
				}
			case *BlockHeaders:
				c.logger.Info().Msgf("Received %v block headers", len(msg.BlockHeadersPacket))
				for _, header := range msg.BlockHeadersPacket {
					headers[header.Hash()] = header
				}
			case *GetBlockHeaders:
				c.logger.Info().Interface("msg", msg).Msg("Received GetBlockHeaders request")
			case *GetBlockBodies:
				c.logger.Info().Msg("Received GetBlockBodies request")
			case *NewBlockHashes:
				c.logger.Info().Msgf("Received %v new block hashes", len(*msg))

				hashes := []common.Hash{}
				for _, hash := range *msg {
					hashes = append(hashes, hash.Hash)

					req := &GetBlockHeaders{
						GetBlockHeadersPacket: &eth.GetBlockHeadersPacket{
							// Providing both the hash and number will result in a `both origin
							// hash and number` error.
							Origin: eth.HashOrNumber{Hash: hash.Hash},
							Amount: 1,
						},
					}
					if err := c.Write(req); err != nil {
						c.logger.Error().Err(err).Msg("Failed to write GetBlockHeaders request")
					}
				}

				req := &GetBlockBodies{
					GetBlockBodiesPacket: hashes,
				}
				if err := c.Write(req); err != nil {
					c.logger.Error().Err(err).Msg("Failed to write GetBlockBodies request")
				}
			case *NewBlock:
				c.logger.Info().Interface("block", msg).Msg("Received new block")
			case *Transactions:
				c.logger.Info().Msgf("Received %v transaction(s)", len(*msg))
			case *NewPooledTransactionHashes:
				c.logger.Info().Msgf("Received %v pooled transaction(s)", len(msg.Hashes))
			case *NewPooledTransactionHashes66:
				c.logger.Info().Msgf("Received %v pooled transaction(s)", len(*msg))
			case *Error:
				c.logger.Debug().Err(msg.err).Msg("Received error")
				if !strings.Contains(msg.Error(), "timeout") {
					return msg
				}
			case *Disconnect:
				c.logger.Debug().Msgf("Disconnect received: %v", msg)
			case *Disconnects:
				c.logger.Debug().Msgf("Disconnect received: %v", msg)
			default:
				c.logger.Info().Interface("msg", msg).Int("code", msg.Code()).Msg("Received message")
			}
		}
	}
}
