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

package crawl

import (
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/p2p"

	// "github.com/ethereum/go-ethereum/rpc"

	"github.com/rs/zerolog/log"
)

type crawler struct {
	input     nodeSet
	output    nodeSet
	disc      resolver
	iters     []enode.Iterator
	inputIter enode.Iterator
	ch        chan *enode.Node
	closed    chan struct{}

	// settings
	revalidateInterval time.Duration
}

type resolver interface {
	RequestENR(*enode.Node) (*enode.Node, error)
}

func newCrawler(input nodeSet, disc resolver, iters ...enode.Iterator) *crawler {
	c := &crawler{
		input:     input,
		output:    make(nodeSet, len(input)),
		disc:      disc,
		iters:     iters,
		inputIter: enode.IterNodes(input.nodes()),
		ch:        make(chan *enode.Node),
		closed:    make(chan struct{}),
	}
	c.iters = append(c.iters, c.inputIter)
	// Copy input to output initially. Any nodes that fail validation
	// will be dropped from output during the run.
	for id, n := range input {
		c.output[id] = n
	}
	return c
}

func (c *crawler) run(timeout time.Duration, server p2p.Server) nodeSet {
	var (
		timeoutTimer = time.NewTimer(timeout)
		timeoutCh    <-chan time.Time
		doneCh       = make(chan enode.Iterator, len(c.iters))
		liveIters    = len(c.iters)
	)
	defer timeoutTimer.Stop()
	for _, it := range c.iters {
		go c.runIterator(doneCh, it)
	}

loop:
	for {
		select {
		case n := <-c.ch:
			c.updateNode(n, server, inputCrawlParams.Client)
		case it := <-doneCh:
			if it == c.inputIter {
				// Enable timeout when we're done revalidating the input nodes.
				log.Info().Msgf("Revalidation of input set is done %d", len(c.input))
				if timeout > 0 {
					timeoutCh = timeoutTimer.C
				}
			}
			if liveIters--; liveIters == 0 {
				break loop
			}
		case <-timeoutCh:
			break loop
		}
	}

	close(c.closed)
	for _, it := range c.iters {
		it.Close()
	}
	for ; liveIters > 0; liveIters-- {
		<-doneCh
	}
	return c.output
}

func (c *crawler) runIterator(done chan<- enode.Iterator, it enode.Iterator) {
	defer func() { done <- it }()
	for it.Next() {
		select {
		case c.ch <- it.Node():
		case <-c.closed:
			return
		}
	}
}

func (c *crawler) updateNode(n *enode.Node, server p2p.Server, clientName *string) {
	nodeItem, ok := c.output[n.ID()]

	// Skip validation of recently-seen nodes.
	if ok && time.Since(nodeItem.LastCheck) < c.revalidateInterval {
		return
	}

	// log.Info().Msgf("URL: %s", fmt.Sprintf("http://%s:%d", n.IP(), n.TCP()))
	// log.Info().Msgf("URL: %s", n.URLv4())

	// Request the node record.
	nn, err := c.disc.RequestENR(n)
	nodeItem.LastCheck = truncNow()
	if err != nil {
		if nodeItem.Score == 0 {
			// Node doesn't implement EIP-868.
			log.Debug().Msgf("Skipping node: id %s", n.ID())
			return
		}
		nodeItem.Score /= 2
	} else {
		nodeItem.N = nn
		nodeItem.Seq = nn.Seq()
		nodeItem.Score++
		if nodeItem.FirstResponse.IsZero() {
			nodeItem.FirstResponse = nodeItem.LastCheck
		}
		nodeItem.LastResponse = nodeItem.LastCheck
	}

	// Store/update node in output set.
	if nodeItem.Score <= 0 {
		log.Info().Msgf("Removing node id %s", n.ID())
		delete(c.output, n.ID())
	} else {
		log.Info().Msgf("Updating node id %s seq %d", n.ID(), n.Seq())
		c.output[n.ID()] = nodeItem
	}

}

func truncNow() time.Time {
	return time.Now().UTC().Truncate(1 * time.Second)
}
