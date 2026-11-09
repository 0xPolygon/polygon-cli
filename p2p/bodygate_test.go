package p2p

import (
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

func TestBodyGateHeaderFirst(t *testing.T) {
	g := newBodyGate(ds.LRUOptions{MaxSize: 16})
	hash, root := common.Hash{1}, common.Hash{2}

	if _, ok := g.headerSeen(hash, root); ok {
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
	pending, ok := g.headerSeen(hash, root)
	if !ok || pending.root != root {
		t.Fatal("headerSeen did not release the waiting body")
	}
	if _, ok := g.headerSeen(hash, root); ok {
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
		wg.Go(func() { _, released = g.headerSeen(hash, common.Hash{}) })
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
