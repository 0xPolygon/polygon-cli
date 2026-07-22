package p2p

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/0xPolygon/polygon-cli/p2p/database"
)

// recordingDB is a Database that counts the hashes passed to the event writers
// and reports configurable Should* flags. Unwritten methods fall through to the
// no-op backend.
type recordingDB struct {
	database.Database
	blockEvents, txEvents                  int
	fullBlock, firstBlock, fullTx, firstTx bool
}

func (r *recordingDB) WriteBlockEvents(_ context.Context, _ *enode.Node, hashes []common.Hash, _ time.Time) {
	r.blockEvents += len(hashes)
}

func (r *recordingDB) WriteTransactionEvents(_ context.Context, _ *enode.Node, hashes []common.Hash, _ time.Time) {
	r.txEvents += len(hashes)
}

func (r *recordingDB) ShouldWriteBlockEvents() bool           { return r.fullBlock }
func (r *recordingDB) ShouldWriteFirstBlockEvent() bool       { return r.firstBlock }
func (r *recordingDB) ShouldWriteTransactionEvents() bool     { return r.fullTx }
func (r *recordingDB) ShouldWriteFirstTransactionEvent() bool { return r.firstTx }
func (r *recordingDB) ShouldWriteTransactions() bool          { return false }

func makeTxAnnounce(t *testing.T, hash common.Hash) ethp2p.Msg {
	t.Helper()
	return encodeMsg(t, eth.NewPooledTransactionHashesMsg, &eth.NewPooledTransactionHashesPacket{
		Types:  []byte{0},
		Sizes:  []uint32{100},
		Hashes: []common.Hash{hash},
	})
}

func TestEventHashes(t *testing.T) {
	all := []common.Hash{{1}, {2}}
	first := []common.Hash{{1}}
	if got := eventHashes(all, first, true, false); len(got) != 2 {
		t.Errorf("full: got %d hashes, want 2", len(got))
	}
	if got := eventHashes(all, first, false, true); len(got) != 1 {
		t.Errorf("first: got %d hashes, want 1", len(got))
	}
	if got := eventHashes(all, first, false, false); got != nil {
		t.Errorf("neither: got %v, want nil", got)
	}
	if got := eventHashes(all, first, true, true); len(got) != 2 {
		t.Errorf("full takes precedence: got %d hashes, want 2", len(got))
	}
}

// TestBlockEventFullVsFirst announces the same hash twice (two peers). The full
// per-peer stream records every announcement; first-seen mode records it once;
// neither records nothing.
func TestBlockEventFullVsFirst(t *testing.T) {
	conns := sharedTestConns(t, false)
	ctx := context.Background()
	cases := []struct {
		name        string
		full, first bool
		want        int
	}{
		{"full", true, false, 2},
		{"first", false, true, 1},
		{"neither", false, false, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := &recordingDB{Database: database.NoDatabase(), fullBlock: tc.full, firstBlock: tc.first}
			c := newTestConn(&recordRW{}, conns)
			c.db = rec
			hash := common.BytesToHash([]byte("blkevt-" + tc.name))
			for i := 0; i < 2; i++ {
				if err := c.handleNewBlockHashes(ctx, makeAnnounce(t, hash, 200)); err != nil {
					t.Fatalf("announce %d: %v", i, err)
				}
			}
			if rec.blockEvents != tc.want {
				t.Fatalf("%s: recorded %d block events, want %d", tc.name, rec.blockEvents, tc.want)
			}
		})
	}
}

// TestTransactionEventFullVsFirst is the tx mirror of TestBlockEventFullVsFirst.
func TestTransactionEventFullVsFirst(t *testing.T) {
	conns := sharedTestConns(t, false)
	ctx := context.Background()
	cases := []struct {
		name        string
		full, first bool
		want        int
	}{
		{"full", true, false, 2},
		{"first", false, true, 1},
		{"neither", false, false, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := &recordingDB{Database: database.NoDatabase(), fullTx: tc.full, firstTx: tc.first}
			c := newTestConn(&recordRW{}, conns)
			c.db = rec
			hash := common.BytesToHash([]byte("txevt-" + tc.name))
			for i := 0; i < 2; i++ {
				if err := c.handleNewPooledTransactionHashes(ctx, 68, makeTxAnnounce(t, hash)); err != nil {
					t.Fatalf("announce %d: %v", i, err)
				}
			}
			if rec.txEvents != tc.want {
				t.Fatalf("%s: recorded %d tx events, want %d", tc.name, rec.txEvents, tc.want)
			}
		})
	}
}
