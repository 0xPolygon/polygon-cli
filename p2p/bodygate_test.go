package p2p

import (
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

func TestBodyGateHeaderFirst(t *testing.T) {
	g := newBodyGate(ds.LRUOptions{MaxSize: 16})
	hash, root := common.Hash{1}, common.Hash{2}

	if _, ok, _ := g.headerSeen(hash, root); ok {
		t.Fatal("headerSeen released a body that never arrived")
	}
	got, known := g.bodyArrived(hash, pendingBody{root: root})
	if !known || got != root {
		t.Fatalf("bodyArrived = %s, %v, want the header's root", got.Hex(), known)
	}
}

func TestBodyGateBodyFirst(t *testing.T) {
	g := newBodyGate(ds.LRUOptions{MaxSize: 16})
	hash, root := common.Hash{1}, common.Hash{2}

	if _, known := g.bodyArrived(hash, pendingBody{root: root}); known {
		t.Fatal("bodyArrived reported a header that was never seen")
	}
	pending, ok, _ := g.headerSeen(hash, root)
	if !ok || pending.root != root {
		t.Fatal("headerSeen did not release the waiting body")
	}
	if _, ok, _ := g.headerSeen(hash, root); ok {
		t.Fatal("a body was released twice")
	}
}

// TestBodyGateRaceNeverLosesABody runs a header and its body concurrently many
// times. Exactly one side must see the other: either the body finds the header
// or the header releases the body. A body that does neither would never be
// written.
func TestBodyGateRaceNeverLosesABody(t *testing.T) {
	g := newBodyGate(ds.LRUOptions{MaxSize: 4096})
	const rounds = 2000
	var lost, twice atomic.Int32
	for i := range rounds {
		hash := common.BigToHash(big.NewInt(int64(i + 1)))
		var bodyKnown, released bool
		var wg sync.WaitGroup
		wg.Go(func() { _, bodyKnown = g.bodyArrived(hash, pendingBody{}) })
		wg.Go(func() { _, released, _ = g.headerSeen(hash, common.Hash{}) })
		wg.Wait()
		switch {
		case !bodyKnown && !released:
			lost.Add(1)
		case bodyKnown && released:
			twice.Add(1)
		}
	}
	if lost.Load() != 0 || twice.Load() != 0 {
		t.Fatalf("lost %d bodies and handled %d twice over %d rounds", lost.Load(), twice.Load(), rounds)
	}
}

// TestBodyGateKeepsFirstBodyPerRoot is the regression test for a later body
// replacing a held one: a wrong body arriving second must not displace the one
// that matches the header.
func TestBodyGateKeepsFirstBodyPerRoot(t *testing.T) {
	g := newBodyGate(ds.LRUOptions{MaxSize: 16})
	hash, good, bad := common.Hash{1}, common.Hash{2}, common.Hash{3}

	g.bodyArrived(hash, pendingBody{root: good})
	g.bodyArrived(hash, pendingBody{root: bad})
	pending, ok, dropped := g.headerSeen(hash, good)
	if !ok || pending.root != good {
		t.Fatal("the matching body was not released")
	}
	if len(dropped) != 1 || dropped[0] != bad {
		t.Fatalf("dropped = %v, want only the wrong body's root", dropped)
	}
}

// TestBodyGateBoundsBodiesPerBlock checks that at most maxPendingBodiesPerBlock
// distinct bodies are held for a block, and that a repeat of a held body is not
// counted again.
func TestBodyGateBoundsBodiesPerBlock(t *testing.T) {
	g := newBodyGate(ds.LRUOptions{MaxSize: 16})
	hash := common.Hash{1}
	for i := range maxPendingBodiesPerBlock + 3 {
		g.bodyArrived(hash, pendingBody{root: common.BigToHash(big.NewInt(int64(i + 1)))})
		g.bodyArrived(hash, pendingBody{root: common.BigToHash(big.NewInt(1))})
	}
	_, ok, dropped := g.headerSeen(hash, common.BigToHash(big.NewInt(1)))
	if !ok {
		t.Fatal("the first body was not kept")
	}
	if len(dropped) != maxPendingBodiesPerBlock-1 {
		t.Fatalf("held %d other bodies, want %d", len(dropped), maxPendingBodiesPerBlock-1)
	}
}

// TestBodyGateRepeatDoesNotExtendExpiry is the regression test for a repeat
// body re-adding its block's entry: that reset the entry's expiry, so peers
// resending bodies for a block whose header never arrived kept it alive.
func TestBodyGateRepeatDoesNotExtendExpiry(t *testing.T) {
	const ttl = 200 * time.Millisecond
	g := newBodyGate(ds.LRUOptions{MaxSize: 16, TTL: ttl})
	hash, root := common.Hash{1}, common.Hash{2}

	g.bodyArrived(hash, pendingBody{root: root})
	time.Sleep(ttl / 2)
	g.bodyArrived(hash, pendingBody{root: root})
	g.bodyArrived(hash, pendingBody{root: common.Hash{3}})
	time.Sleep(ttl/2 + ttl/4)
	if _, ok := g.pending.Peek(hash); ok {
		t.Fatal("repeat bodies extended the held entry past its expiry")
	}
}
