package p2p

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/0xPolygon/polygon-cli/p2p/database"
)

// txFetchDB is a Database that asks for transaction bodies, so announcements
// trigger GetPooledTransactions requests.
type txFetchDB struct {
	database.Database
}

func (txFetchDB) ShouldWriteTransactions() bool { return true }

// TestHandleNewPooledTransactionHashesChunksRequests verifies that a large
// announcement is fetched in requests of at most maxPooledTxsPerRequest
// hashes, in order, with every hash requested exactly once.
func TestHandleNewPooledTransactionHashesChunksRequests(t *testing.T) {
	conns := sharedTestConns(t, false)
	ctx := context.Background()

	for _, n := range []int{1, maxPooledTxsPerRequest, maxPooledTxsPerRequest + 1, 16384} {
		rw := &recordRW{}
		c := newTestConn(rw, conns)
		c.db = txFetchDB{Database: database.NoDatabase()}

		hashes := make([]common.Hash, n)
		types := make([]byte, n)
		sizes := make([]uint32, n)
		for i := range hashes {
			// Distinct per case so the shared Conns tx cache never filters them.
			hashes[i] = common.BigToHash(big.NewInt(int64(n)<<32 | int64(i)))
			sizes[i] = 100
		}
		msg := encodeMsg(t, eth.NewPooledTransactionHashesMsg, &eth.NewPooledTransactionHashesPacket71{
			Types:  types,
			Sizes:  sizes,
			Hashes: hashes,
		})
		if err := c.handleNewPooledTransactionHashes(ctx, 68, msg); err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}

		wantReqs := (n + maxPooledTxsPerRequest - 1) / maxPooledTxsPerRequest
		if len(rw.codes) != wantReqs {
			t.Fatalf("n=%d: sent %d messages, want %d", n, len(rw.codes), wantReqs)
		}
		var got []common.Hash
		for i, code := range rw.codes {
			if code != eth.GetPooledTransactionsMsg {
				t.Fatalf("n=%d: message %d has code %#x, want GetPooledTransactionsMsg", n, i, code)
			}
			var req eth.GetPooledTransactionsPacket
			if err := rlp.DecodeBytes(rw.payloads[i], &req); err != nil {
				t.Fatalf("n=%d: decode request %d: %v", n, i, err)
			}
			if len(req.GetPooledTransactionsRequest) > maxPooledTxsPerRequest {
				t.Fatalf("n=%d: request %d has %d hashes, want <= %d", n, i, len(req.GetPooledTransactionsRequest), maxPooledTxsPerRequest)
			}
			got = append(got, req.GetPooledTransactionsRequest...)
		}
		if len(got) != n {
			t.Fatalf("n=%d: requested %d hashes, want %d", n, len(got), n)
		}
		for i := range got {
			if got[i] != hashes[i] {
				t.Fatalf("n=%d: hash %d is %x, want %x", n, i, got[i], hashes[i])
			}
		}
	}
}
