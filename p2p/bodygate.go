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

// maxPendingBodiesPerBlock bounds how many distinct bodies (by txRoot) are held
// for one block while its header is outstanding. Honest peers all send the same
// body, so it takes one slot; the rest only absorb bodies that will fail the
// check.
const maxPendingBodiesPerBlock = 4

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
// hold the body; record the header, release any held body -- happen under one
// lock. Otherwise a body held just after its header went by would never be
// written.
//
// Bodies are held per block by txRoot, first arrival kept, so a later body with
// different contents cannot displace one that will match the header.
type bodyGate struct {
	mu      sync.Mutex
	roots   *ds.LRU[common.Hash, common.Hash]
	pending *ds.LRU[common.Hash, map[common.Hash]pendingBody]
}

func newBodyGate(opts ds.LRUOptions) *bodyGate {
	return &bodyGate{
		roots:   ds.NewLRU[common.Hash, common.Hash](opts),
		pending: ds.NewLRU[common.Hash, map[common.Hash]pendingBody](opts),
	}
}

// bodyArrived returns the txRoot of the block's header if it has been seen.
// Otherwise it holds the body until headerSeen is called for the block. A body
// whose txRoot is already held, or that would exceed maxPendingBodiesPerBlock,
// is not held again.
func (g *bodyGate) bodyArrived(hash common.Hash, body pendingBody) (common.Hash, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if root, ok := g.roots.Get(hash); ok {
		return root, true
	}
	held, ok := g.pending.Peek(hash)
	if !ok {
		g.pending.Add(hash, map[common.Hash]pendingBody{body.root: body})
		return common.Hash{}, false
	}
	// Added to the block's existing entry in place, without re-adding it, so a
	// repeat or rejected body does not reset the entry's expiry or LRU position.
	if _, dup := held[body.root]; !dup && len(held) < maxPendingBodiesPerBlock {
		held[body.root] = body
	}
	return common.Hash{}, false
}

// headerSeen records a header's txRoot and removes every body held for the
// block. It returns the one matching root, if any, for the caller to write, and
// the roots of the others, which are dropped.
func (g *bodyGate) headerSeen(hash, root common.Hash) (match pendingBody, ok bool, dropped []common.Hash) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.roots.Add(hash, root)
	held, found := g.pending.Remove(hash)
	if !found {
		return pendingBody{}, false, nil
	}
	for r, b := range held {
		if r == root {
			match, ok = b, true
			continue
		}
		dropped = append(dropped, r)
	}
	return match, ok, dropped
}
