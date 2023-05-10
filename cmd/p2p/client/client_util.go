// Copyright 2019 The go-ethereum Authors
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

package client

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type client struct {
	input     p2p.NodeSet
	output    p2p.NodeSet
	disc      *discover.UDPv4
	iters     []enode.Iterator
	inputIter enode.Iterator
	ch        chan *enode.Node
	db        *datastore.Client
	peers     map[*enode.Node]struct{}

	// settings
	revalidateInterval time.Duration
	outputMutex        sync.RWMutex
	peersMutex         sync.RWMutex
}

const (
	nodeRemoved = iota
	nodeSkipRecent
	nodeSkipIncompat
	nodeAdded
	nodeUpdated
)

type resolver interface {
	RequestENR(*enode.Node) (*enode.Node, error)
}

func newClient(input p2p.NodeSet, disc *discover.UDPv4, iters ...enode.Iterator) *client {
	db, err := datastore.NewClient(context.Background(), inputClientParams.ProjectID)
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to Datastore")
	}

	c := &client{
		input:     input,
		output:    make(p2p.NodeSet, len(input)),
		disc:      disc,
		iters:     iters,
		inputIter: enode.IterNodes(input.Nodes()),
		ch:        make(chan *enode.Node),
		db:        db,
		peers:     make(map[*enode.Node]struct{}),
	}
	c.iters = append(c.iters, c.inputIter)
	// Copy input to output initially. Any nodes that fail validation
	// will be dropped from output during the run.
	for id, n := range input {
		c.output[id] = n
	}
	return c
}

func (c *client) run(nthreads int) {
	statusTicker := time.NewTicker(time.Second * 8)
	defer statusTicker.Stop()

	if nthreads < 1 {
		nthreads = 1
	}

	for _, it := range c.iters {
		go c.runIterator(it)
	}

	var (
		added   uint64
		updated uint64
		skipped uint64
		recent  uint64
		removed uint64
	)

	for i := 0; i < nthreads; i++ {
		go func() {
			for {
				select {
				case n := <-c.ch:
					switch c.updateNode(n) {
					case nodeSkipIncompat:
						atomic.AddUint64(&skipped, 1)
					case nodeSkipRecent:
						atomic.AddUint64(&recent, 1)
					case nodeRemoved:
						atomic.AddUint64(&removed, 1)
					case nodeAdded:
						atomic.AddUint64(&added, 1)
					default:
						atomic.AddUint64(&updated, 1)
					}
				}
			}
		}()
	}

	for {
		select {
		case <-statusTicker.C:
			c.peersMutex.RLock()
			log.Info().
				Uint64("added", atomic.LoadUint64(&added)).
				Uint64("updated", atomic.LoadUint64(&updated)).
				Uint64("removed", atomic.LoadUint64(&removed)).
				Uint64("ignored(recent)", atomic.LoadUint64(&removed)).
				Uint64("ignored(incompatible)", atomic.LoadUint64(&skipped)).
				Int("peers", len(c.peers)).
				Msg("Discovery in progress")
			c.peersMutex.RUnlock()
		}
	}
}

func (c *client) runIterator(it enode.Iterator) {
	for it.Next() {
		select {
		case c.ch <- it.Node():
		}
	}
}

// peerNode peers with nodes with matching network IDs.
//
// Nodes will exchange hello and status messages to determine the network ID.
// If the network IDs match then a connection will be opened with the node to
// receive blocks and transactions.
func (c *client) peerNode(n *enode.Node) bool {
	conn, err := p2p.Dial(n)
	if err != nil {
		log.Debug().Err(err).Msg("Dial failed")
		return true
	}
	conn.Sensor = inputClientParams.SensorID

	hello, status, err := conn.Peer()
	if err != nil {
		log.Debug().Err(err).Msg("Peer failed")
		conn.Close()
		return true
	}

	log.Debug().Interface("hello", hello).Interface("status", status).Msg("Peering messages received")

	skip := inputClientParams.NetworkID != int(status.NetworkID)
	if !skip {
		c.peersMutex.Lock()
		defer c.peersMutex.Unlock()

		// Don't open duplicate connections.
		if _, ok := c.peers[n]; !ok {
			c.peers[n] = struct{}{}

			go func() {
				defer conn.Close()
				if err := conn.ReadAndServe(c.db); err != nil {
					log.Debug().Err(err.Unwrap()).Msg("Error received")
				}

				c.peersMutex.Lock()
				defer c.peersMutex.Unlock()
				delete(c.peers, n)
			}()
		}
	}

	return skip
}

// updateNode updates the info about the given node, and returns a status about
// what changed.
func (c *client) updateNode(n *enode.Node) int {
	c.outputMutex.RLock()
	node, ok := c.output[n.ID()]
	c.outputMutex.RUnlock()

	// Skip validation of recently-seen nodes.
	if ok && time.Since(node.LastCheck) < c.revalidateInterval {
		log.Debug().Str("node", n.String()).Msg("Skipping node")
		return nodeSkipRecent
	}

	// Skip the node if unable to peer.
	if c.peerNode(n) {
		return nodeSkipIncompat
	}

	// Request the node record.
	status := nodeUpdated
	node.LastCheck = truncNow()

	if nn, err := c.disc.RequestENR(n); err != nil {
		if node.Score == 0 {
			// Node doesn't implement EIP-868.
			log.Debug().Str("node", n.String()).Msg("Skipping node")
			return nodeSkipIncompat
		}
		node.Score /= 2
	} else {
		node.N = nn
		node.Seq = nn.Seq()
		node.Score++
		if node.FirstResponse.IsZero() {
			node.FirstResponse = node.LastCheck
			status = nodeAdded
		}
		node.LastResponse = node.LastCheck
	}

	// Store/update node in output set.
	c.outputMutex.Lock()
	defer c.outputMutex.Unlock()

	if node.Score <= 0 {
		log.Debug().Str("node", n.String()).Msg("Removing node")
		delete(c.output, n.ID())
		return nodeRemoved
	}

	log.Debug().Str("node", n.String()).Uint64("seq", n.Seq()).Int("score", node.Score).Msg("Updating node")
	c.output[n.ID()] = node

	if err := p2p.WriteNodesJSON(inputClientParams.NodesFile, c.output); err != nil {
		log.Error().Err(err).Msg("Failed to write nodes json")
	}

	return status
}

func truncNow() time.Time {
	return time.Now().UTC().Truncate(1 * time.Second)
}
