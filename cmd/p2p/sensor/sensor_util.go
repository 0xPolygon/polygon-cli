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

package sensor

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"

	"github.com/maticnetwork/polygon-cli/p2p"
	"github.com/maticnetwork/polygon-cli/p2p/database"
)

type sensor struct {
	input     p2p.NodeSet
	output    p2p.NodeSet
	disc      resolver
	iters     []enode.Iterator
	inputIter enode.Iterator
	nodeCh    chan *enode.Node
	db        database.Database
	peers     map[string]struct{}
	count     *p2p.MessageCount

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

// newSensor creates the new sensor and establishes the connection. If you want
// to change the database, this is where you should do it.
func newSensor(input p2p.NodeSet, disc resolver, iters ...enode.Iterator) *sensor {
	s := &sensor{
		input:     input,
		output:    make(p2p.NodeSet, len(input)),
		disc:      disc,
		iters:     iters,
		inputIter: enode.IterNodes(input.Nodes()),
		nodeCh:    make(chan *enode.Node),
		db: database.NewDatastore(
			inputSensorParams.ProjectID,
			inputSensorParams.SensorID,
			inputSensorParams.MaxConcurrentDatabaseWrites,
			inputSensorParams.ShouldWriteBlocks,
			inputSensorParams.ShouldWriteTransactions,
		),
		peers: make(map[string]struct{}),
		count: &p2p.MessageCount{},
	}
	s.iters = append(s.iters, s.inputIter)
	// Copy input to output initially. Any nodes that fail validation
	// will be dropped from output during the run.
	for id, n := range input {
		s.output[id] = n
	}
	return s
}

// run will start nthreads number of goroutines for discovery and a goroutine
// for logging.
func (s *sensor) run(nthreads int) {
	statusTicker := time.NewTicker(time.Second * 8)

	if nthreads < 1 {
		nthreads = 1
	}

	for _, it := range s.iters {
		go s.runIterator(it)
	}

	var (
		added   uint64
		updated uint64
		skipped uint64
		recent  uint64
		removed uint64
	)

	// This will start the goroutines responsible for discovery.
	for i := 0; i < nthreads; i++ {
		go func() {
			for {
				switch s.updateNode(<-s.nodeCh) {
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
		}()
	}

	// Start logging message counts and peer status.
	go p2p.LogMessageCount(s.count, time.NewTicker(time.Second))

	for {
		<-statusTicker.C
		s.peersMutex.RLock()
		log.Info().
			Uint64("added", atomic.LoadUint64(&added)).
			Uint64("updated", atomic.LoadUint64(&updated)).
			Uint64("removed", atomic.LoadUint64(&removed)).
			Uint64("ignored(recent)", atomic.LoadUint64(&removed)).
			Uint64("ignored(incompatible)", atomic.LoadUint64(&skipped)).
			Int("peers", len(s.peers)).
			Msg("Discovery in progress")
		s.peersMutex.RUnlock()
	}
}

func (s *sensor) runIterator(it enode.Iterator) {
	for it.Next() {
		s.nodeCh <- it.Node()
	}
}

// peerNode peers with nodes with matching network IDs. Nodes will exchange
// hello and status messages to determine the network ID. If the network IDs
// match then a connection will be opened with the node to receive blocks and
// transactions.
func (s *sensor) peerNode(n *enode.Node) bool {
	conn, err := p2p.Dial(n)
	if err != nil {
		log.Debug().Err(err).Msg("Dial failed")
		return true
	}
	conn.SensorID = inputSensorParams.SensorID

	hello, status, err := conn.Peer()
	if err != nil {
		log.Debug().Err(err).Msg("Peer failed")
		conn.Close()
		return true
	}

	log.Debug().Interface("hello", hello).Interface("status", status).Msg("Peering messages received")

	skip := inputSensorParams.NetworkID != int(status.NetworkID)
	if !skip {
		s.peersMutex.Lock()
		defer s.peersMutex.Unlock()

		// Don't open duplicate connections with peers. The enode is used here
		// rather than the enr because the enr can change. Ensure the maximum
		// number of peers is not exceeded to prevent memory overuse.
		peer := n.URLv4()
		if _, ok := s.peers[peer]; !ok && len(s.peers) < inputSensorParams.MaxPeers {
			s.peers[peer] = struct{}{}

			go func() {
				defer conn.Close()
				if err := conn.ReadAndServe(s.db, s.count); err != nil {
					log.Debug().Err(err).Msg("Received error")
				}

				s.peersMutex.Lock()
				defer s.peersMutex.Unlock()
				delete(s.peers, peer)
			}()
		}
	}

	return skip
}

// updateNode updates the info about the given node, and returns a status about
// what changed. If the node is compatible, then it will peer with the node and
// start receiving block and transaction data.
func (s *sensor) updateNode(n *enode.Node) int {
	s.outputMutex.RLock()
	node, ok := s.output[n.ID()]
	s.outputMutex.RUnlock()

	// Skip validation of recently-seen nodes.
	if ok && time.Since(node.LastCheck) < s.revalidateInterval {
		log.Debug().Str("node", n.String()).Msg("Skipping node")
		return nodeSkipRecent
	}

	// Skip the node if unable to peer.
	if s.peerNode(n) {
		return nodeSkipIncompat
	}

	// Request the node record.
	status := nodeUpdated
	node.LastCheck = truncNow()

	if nn, err := s.disc.RequestENR(n); err != nil {
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
	s.outputMutex.Lock()
	defer s.outputMutex.Unlock()

	if node.Score <= 0 {
		log.Debug().Str("node", n.String()).Msg("Removing node")
		delete(s.output, n.ID())
		return nodeRemoved
	}

	log.Debug().Str("node", n.String()).Uint64("seq", n.Seq()).Int("score", node.Score).Msg("Updating node")
	s.output[n.ID()] = node

	// Update the nodes file at the end of each iteration.
	if err := p2p.WriteNodesJSON(inputSensorParams.NodesFile, s.output); err != nil {
		log.Error().Err(err).Msg("Failed to write nodes json")
	}

	return status
}

func truncNow() time.Time {
	return time.Now().UTC().Truncate(1 * time.Second)
}
