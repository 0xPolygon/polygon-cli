package p2p

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"

	"github.com/0xPolygon/polygon-cli/p2p/bortx"
	"github.com/0xPolygon/polygon-cli/p2p/database"
	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

// pendingBody is a block body that arrived before its header, held until the
// header shows whether it belongs to the block.
type pendingBody struct {
	body *eth.BlockBody
	txs  *bortx.Body
	ann  database.BlockAnnouncement
	tfs  time.Time
	root common.Hash
}

// bodyGate keeps a block body out of the database until it has been checked
// against its header's txRoot. Headers and bodies arrive in either order and on
// different peers' goroutines, so the two lookups -- is the header known, else
// park the body; record the header, release any parked body -- happen under one
// lock. Otherwise a body parked just after its header went by would never be
// written.
type bodyGate struct {
	mu      sync.Mutex
	roots   *ds.LRU[common.Hash, common.Hash]
	pending *ds.LRU[common.Hash, pendingBody]
}

func newBodyGate(opts ds.LRUOptions) *bodyGate {
	return &bodyGate{
		roots:   ds.NewLRU[common.Hash, common.Hash](opts),
		pending: ds.NewLRU[common.Hash, pendingBody](opts),
	}
}

// bodyArrived returns the txRoot of the block's header if it has been seen.
// Otherwise it holds the body until headerSeen is called for the block.
func (g *bodyGate) bodyArrived(hash common.Hash, body pendingBody) (common.Hash, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if root, ok := g.roots.Get(hash); ok {
		return root, true
	}
	g.pending.Add(hash, body)
	return common.Hash{}, false
}

// headerSeen records a header's txRoot and returns the body that was waiting for
// it, if any, for the caller to check and write.
func (g *bodyGate) headerSeen(hash, root common.Hash) (pendingBody, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.roots.Add(hash, root)
	return g.pending.Remove(hash)
}
