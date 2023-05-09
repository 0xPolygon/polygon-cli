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
	closed    chan struct{}
	db        *datastore.Client
	peers     uint32

	// settings
	revalidateInterval time.Duration
	mu                 sync.RWMutex
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
		closed:    make(chan struct{}),
		db:        db,
	}
	c.iters = append(c.iters, c.inputIter)
	// Copy input to output initially. Any nodes that fail validation
	// will be dropped from output during the run.
	for id, n := range input {
		c.output[id] = n
	}
	return c
}

func (c *client) run(timeout time.Duration, nthreads int) p2p.NodeSet {
	var (
		timeoutTimer = time.NewTimer(timeout)
		timeoutCh    <-chan time.Time
		statusTicker = time.NewTicker(time.Second * 8)
		doneCh       = make(chan enode.Iterator, len(c.iters))
		liveIters    = len(c.iters)
	)
	if nthreads < 1 {
		nthreads = 1
	}
	defer timeoutTimer.Stop()
	defer statusTicker.Stop()
	for _, it := range c.iters {
		go c.runIterator(doneCh, it)
	}
	var (
		added   uint64
		updated uint64
		skipped uint64
		recent  uint64
		removed uint64
		wg      sync.WaitGroup
	)
	wg.Add(nthreads)
	for i := 0; i < nthreads; i++ {
		go func() {
			defer wg.Done()
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
				case <-c.closed:
					return
				}
			}
		}()
	}

loop:
	for {
		select {
		case it := <-doneCh:
			if it == c.inputIter {
				// Enable timeout when we're done revalidating the input nodes.
				log.Info().Int("len", len(c.input)).Msg("Revalidation of input set is done")
				if timeout > 0 {
					timeoutCh = timeoutTimer.C
				}
			}
			if liveIters--; liveIters == 0 {
				break loop
			}
		case <-timeoutCh:
			break loop
		case <-statusTicker.C:
			log.Info().
				Uint64("added", atomic.LoadUint64(&added)).
				Uint64("updated", atomic.LoadUint64(&updated)).
				Uint64("removed", atomic.LoadUint64(&removed)).
				Uint64("ignored(recent)", atomic.LoadUint64(&removed)).
				Uint64("ignored(incompatible)", atomic.LoadUint64(&skipped)).
				Uint32("peers", atomic.LoadUint32(&c.peers)).
				Msg("Discovery in progress")
		}
	}

	close(c.closed)
	for _, it := range c.iters {
		it.Close()
	}
	for ; liveIters > 0; liveIters-- {
		<-doneCh
	}
	wg.Wait()
	return c.output
}

func (c *client) runIterator(done chan<- enode.Iterator, it enode.Iterator) {
	defer func() { done <- it }()
	for it.Next() {
		select {
		case c.ch <- it.Node():
		case <-c.closed:
			return
		}
	}
}

// shouldSkipNode filters out nodes by their network id. If there is a status
// message, skip nodes that don't have the correct network id. Otherwise, skip
// nodes that are unable to peer. If a peer is not being skipped and the client
// is not in crawler mode, then a goroutine will be spawned and read messages
// from the new peer.
func (c *client) shouldSkipNode(n *enode.Node) bool {
	// Exit early since crawling doesn't need to dial and peer.
	if inputClientParams.NetworkID <= 0 && inputClientParams.IsCrawler {
		return false
	}

	conn, err := p2p.Dial(n)
	if err != nil {
		log.Debug().Err(err).Msg("Dial failed")
		return true
	}
	conn.Sensor = inputClientParams.SensorID

	hello, message, err := conn.Peer()
	if err != nil {
		log.Debug().Err(err).Msg("Peer failed")
		conn.Close()
		return true
	}

	log.Debug().Interface("hello", hello).Interface("status", message).Msg("Peering messages received")

	skip := inputClientParams.NetworkID != int(message.NetworkID)
	if !skip && !inputClientParams.IsCrawler {
		go func() {
			defer conn.Close()
			atomic.AddUint32(&c.peers, 1)
			if err := conn.ReadAndServe(c.db); err != nil {
				log.Debug().Err(err.Unwrap()).Msg("Error received")
			}
			atomic.AddUint32(&c.peers, ^uint32(0))
		}()
	} else {
		conn.Close()
	}

	return skip
}

// updateNode updates the info about the given node, and returns a status about
// what changed.
func (c *client) updateNode(n *enode.Node) int {
	c.mu.RLock()
	node, ok := c.output[n.ID()]
	c.mu.RUnlock()

	// Skip validation of recently-seen nodes.
	if ok && time.Since(node.LastCheck) < c.revalidateInterval {
		log.Debug().Str("node", n.String()).Msg("Skipping node")
		return nodeSkipRecent
	}

	// Filter out incompatible nodes.
	if c.shouldSkipNode(n) {
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
	c.mu.Lock()
	defer c.mu.Unlock()

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
