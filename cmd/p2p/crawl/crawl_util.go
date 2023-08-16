package crawl

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type crawler struct {
	input     []*enode.Node
	output    p2p.NodeSet
	disc      resolver
	iters     []enode.Iterator
	inputIter enode.Iterator
	ch        chan *enode.Node
	closed    chan struct{}

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

func newCrawler(input []*enode.Node, disc resolver, iters ...enode.Iterator) *crawler {
	c := &crawler{
		input:     input,
		output:    make(p2p.NodeSet, len(input)),
		disc:      disc,
		iters:     iters,
		inputIter: enode.IterNodes(input),
		ch:        make(chan *enode.Node),
		closed:    make(chan struct{}),
	}
	c.iters = append(c.iters, c.inputIter)
	// Copy input to output initially. Any nodes that fail validation
	// will be dropped from output during the run.
	for _, n := range input {
		c.output[n.ID()] = n.URLv4()
	}
	return c
}

func (c *crawler) run(timeout time.Duration, nthreads int) p2p.NodeSet {
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
		skipped uint64
		recent  uint64
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
					case nodeAdded:
						atomic.AddUint64(&added, 1)
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
				Uint64("ignored(recent)", atomic.LoadUint64(&recent)).
				Uint64("ignored(incompatible)", atomic.LoadUint64(&skipped)).
				Msg("Crawling in progress")
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

// shouldSkipNode filters out nodes by their network id. If there is a status
// message, skip nodes that don't have the correct network id. Otherwise, skip
// nodes that are unable to peer.
func shouldSkipNode(n *enode.Node) bool {
	if inputCrawlParams.NetworkID == 0 {
		return false
	}

	conn, err := p2p.Dial(n)
	if err != nil {
		log.Error().Err(err).Msg("Dial failed")
		return true
	}
	defer conn.Close()

	hello, status, err := conn.Peer()
	if err != nil {
		log.Error().Err(err).Msg("Peer failed")
		return true
	}

	log.Debug().Interface("hello", hello).Interface("status", status).Msg("Message received")
	return inputCrawlParams.NetworkID != status.NetworkID
}

// updateNode updates the info about the given node, and returns a status about
// what changed.
func (c *crawler) updateNode(n *enode.Node) int {
	c.mu.RLock()
	_, ok := c.output[n.ID()]
	c.mu.RUnlock()

	if ok {
		return nodeSkipRecent
	}

	// Filter out incompatible nodes.
	if shouldSkipNode(n) {
		return nodeSkipIncompat
	}

	nn, err := c.disc.RequestENR(n)
	if err != nil {
		return nodeSkipIncompat
	}

	// Store/update node in output set.
	c.mu.Lock()
	c.output[nn.ID()] = nn.URLv4()
	c.mu.Unlock()

	return nodeAdded
}
