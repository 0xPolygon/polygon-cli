package crawl

import (
	"errors"
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
	mu                 sync.Mutex
}

const (
	nodeAdded = iota
	nodeDialErr
	nodePeerErr
	nodeIncompatible
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
		added        uint64
		dialErr      uint64
		peerErr      uint64
		incompatible uint64
		wg           sync.WaitGroup
	)
	wg.Add(nthreads)
	for i := 0; i < nthreads; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case n := <-c.ch:
					switch c.updateNode(n) {
					case nodeAdded:
						atomic.AddUint64(&added, 1)
					case nodeDialErr:
						atomic.AddUint64(&dialErr, 1)
					case nodePeerErr:
						atomic.AddUint64(&peerErr, 1)
					case nodeIncompatible:
						atomic.AddUint64(&incompatible, 1)
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
				Uint64("dial_err", atomic.LoadUint64(&dialErr)).
				Uint64("peer_err", atomic.LoadUint64(&peerErr)).
				Uint64("incompatible", atomic.LoadUint64(&incompatible)).
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

// updateNode updates the info about the given node, and returns a status about
// what changed.
func (c *crawler) updateNode(n *enode.Node) int {
	var (
		hello  *p2p.Hello
		status *p2p.Status
		err    error
		nodes  int
	)

	c.mu.Lock()
	if _, ok := c.output[n.ID()]; !ok {
		c.output[n.ID()] = []p2p.NodeJSON{}
	}
	nodes = len(c.output)
	c.mu.Unlock()

	// Add result to output node set.
	defer func() {
		c.mu.Lock()
		result := p2p.NodeJSON{
			URL:    n.URLv4(),
			Hello:  hello,
			Status: status,
			Time:   time.Now().Unix(),
			Nodes:  nodes,
		}
		if err != nil {
			result.Error = err.Error()
		}

		c.output[n.ID()] = append(c.output[n.ID()], result)
		c.mu.Unlock()
	}()

	conn, dialErr := p2p.Dial(n)
	if err = dialErr; err != nil {
		log.Error().Err(err).Msg("Dial failed")
		return nodeDialErr
	}
	defer conn.Close()

	hello, status, err = conn.Peer()
	if err != nil {
		log.Error().Err(err).Msg("Peer failed")
		return nodePeerErr
	}

	log.Debug().Interface("hello", hello).Interface("status", status).Msg("Message received")
	if inputCrawlParams.NetworkID != 0 && inputCrawlParams.NetworkID != status.NetworkID {
		err = errors.New("network ID mismatch")
		return nodeIncompatible
	}

	return nodeAdded
}
