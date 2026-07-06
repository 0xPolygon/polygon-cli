package p2p

import (
	"bytes"
	"context"
	"io"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog"

	"github.com/0xPolygon/polygon-cli/p2p/database"
	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

// recordRW is a minimal ethp2p.MsgReadWriter that records written messages
// (code + payload), so tests can assert which requests/responses a handler sent.
type recordRW struct {
	codes    []uint64
	payloads [][]byte
}

func (r *recordRW) ReadMsg() (ethp2p.Msg, error) { return ethp2p.Msg{}, io.EOF }

func (r *recordRW) WriteMsg(m ethp2p.Msg) error {
	var buf []byte
	if m.Payload != nil {
		buf, _ = io.ReadAll(m.Payload)
	}
	r.codes = append(r.codes, m.Code)
	r.payloads = append(r.payloads, buf)
	return nil
}

func (r *recordRW) reset() {
	r.codes = nil
	r.payloads = nil
}

// encodeMsg RLP-encodes val and wraps it in an ethp2p.Msg with the given code.
func encodeMsg(t *testing.T, code uint64, val any) ethp2p.Msg {
	t.Helper()
	enc, err := rlp.EncodeToBytes(val)
	if err != nil {
		t.Fatalf("encode msg (code %d): %v", code, err)
	}
	return ethp2p.Msg{Code: code, Size: uint32(len(enc)), Payload: bytes.NewReader(enc)}
}

func makeAnnounce(t *testing.T, hash common.Hash, number uint64) ethp2p.Msg {
	t.Helper()
	return encodeMsg(t, eth.NewBlockHashesMsg, NewBlockHashesPacket{{Hash: hash, Number: number}})
}

// makeBodies builds a BlockBodies response for one (empty) body under the given
// request id.
func makeBodies(t *testing.T, reqID uint64) ethp2p.Msg {
	t.Helper()
	bodyRLP, err := rlp.EncodeToBytes(rawBlockBody{})
	if err != nil {
		t.Fatalf("encode body: %v", err)
	}
	return encodeMsg(t, eth.BlockBodiesMsg, &eth.BlockBodiesRLPPacket{
		RequestId:              reqID,
		BlockBodiesRLPResponse: eth.BlockBodiesRLPResponse{rlp.RawValue(bodyRLP)},
	})
}

// newTestConn builds a minimal conn wired to a recording MsgReadWriter, a no-op
// database, and the given connection manager, sufficient to drive the block
// message handlers in unit tests.
func newTestConn(rw ethp2p.MsgReadWriter, conns *Conns) *conn {
	return &conn{
		sensorID:    "test",
		logger:      zerolog.Nop(),
		rw:          rw,
		db:          database.NoDatabase(),
		requests:    ds.NewLRU[uint64, common.Hash](ds.LRUOptions{MaxSize: 1024}),
		parents:     ds.NewLRU[common.Hash, struct{}](ds.LRUOptions{MaxSize: 1024}),
		conns:       conns,
		messages:    NewPeerMessages(),
		latestBlock: &ds.Locked[latestBlock]{},
		version:     68,
	}
}

// sharedTestConns returns a single Conns for the whole test binary. NewConns
// registers Prometheus collectors on the default registry, so it can only be
// called once per process; tests set cacheOnlyValidated directly and use
// distinct block hashes to stay independent.
var (
	testConnsOnce sync.Once
	testConns     *Conns
)

func sharedTestConns(t *testing.T, cacheOnlyValidated bool) *Conns {
	t.Helper()
	testConnsOnce.Do(func() {
		testConns = NewConns(ConnsOptions{
			Head:        NewBlockPacket{Block: types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)}), TD: big.NewInt(1)},
			BlocksCache: ds.LRUOptions{MaxSize: 1024},
		})
	})
	testConns.cacheOnlyValidated = cacheOnlyValidated
	return testConns
}

// TestHandleNewBlockHashesRefetch verifies that an announced hash is (re-)fetched
// until we hold the full block: the first announcer triggers a fetch, later
// announcers still fetch while the block is incomplete, and once the full block
// is cached further announcements do not fetch.
func TestHandleNewBlockHashesRefetch(t *testing.T) {
	rw := &recordRW{}
	conns := sharedTestConns(t, false)
	c := newTestConn(rw, conns)

	ctx := context.Background()
	hash := common.HexToHash("0xabc123")

	// 1. First announcement of an unknown hash -> fetch header + body.
	if err := c.handleNewBlockHashes(ctx, makeAnnounce(t, hash, 100)); err != nil {
		t.Fatalf("first announce: %v", err)
	}
	if got := len(rw.codes); got != 2 {
		t.Fatalf("first announce: want 2 requests (headers+bodies), got %d %v", got, rw.codes)
	}
	assertContains(t, rw.codes, eth.GetBlockHeadersMsg, eth.GetBlockBodiesMsg)

	// 2. Re-announcement (e.g. from another peer) while the block is still only a
	//    hash marker -> fetch AGAIN. This is the fix: previously the marker caused
	//    an unconditional skip, stranding the block if the first peer never served it.
	rw.reset()
	if err := c.handleNewBlockHashes(ctx, makeAnnounce(t, hash, 100)); err != nil {
		t.Fatalf("re-announce (incomplete): %v", err)
	}
	if got := len(rw.codes); got != 2 {
		t.Fatalf("re-announce while incomplete: want 2 re-fetch requests, got %d %v", got, rw.codes)
	}

	// 3. Once the full block is cached, a re-announcement must NOT fetch.
	conns.Blocks().Update(hash, func(bc BlockCache) BlockCache {
		bc.Header = &types.Header{Number: big.NewInt(100)}
		bc.Body = &eth.BlockBody{}
		bc.TD = big.NewInt(1)
		return bc
	})
	rw.reset()
	if err := c.handleNewBlockHashes(ctx, makeAnnounce(t, hash, 100)); err != nil {
		t.Fatalf("re-announce (complete): %v", err)
	}
	if got := len(rw.codes); got != 0 {
		t.Fatalf("re-announce with full block cached: want 0 requests, got %d %v", got, rw.codes)
	}
}

// TestHandleBlockBodiesRetainsBodyBeforeHeader is a regression test for the
// body-before-header race under cache-only-validated-blocks: the body response
// can arrive before the header response. Previously the body was retained only
// if a header was already cached, so a body-first response was dropped and the
// entry stayed header-only even for a valid block (unservable on GetBlockBodies).
// Now the body is retained whenever the block already has a cache entry (the
// announcement marker), and completed when the header later arrives.
func TestHandleBlockBodiesRetainsBodyBeforeHeader(t *testing.T) {
	rw := &recordRW{}
	conns := sharedTestConns(t, true) // cache-only-validated enabled
	c := newTestConn(rw, conns)

	ctx := context.Background()
	hash := common.HexToHash("0xdef456")

	// The announcement created a hash-only marker; the header has NOT arrived yet.
	conns.Blocks().Update(hash, func(bc BlockCache) BlockCache { return bc })
	const reqID = 7
	c.requests.Add(reqID, hash)

	// Body response arrives before the header.
	if err := c.handleBlockBodies(ctx, makeBodies(t, reqID)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}

	cache, ok := conns.Blocks().Peek(hash)
	if !ok {
		t.Fatal("cache entry vanished")
	}
	if cache.Body == nil {
		t.Fatal("body was dropped when it arrived before the header (regression)")
	}
	if cache.Header != nil {
		t.Fatal("header should still be absent at this point")
	}
}

// makeGetBodies builds a GetBlockBodies request for the given hashes.
func makeGetBodies(t *testing.T, reqID uint64, hashes ...common.Hash) ethp2p.Msg {
	t.Helper()
	return encodeMsg(t, eth.GetBlockBodiesMsg, &eth.GetBlockBodiesPacket{
		RequestId:             reqID,
		GetBlockBodiesRequest: eth.GetBlockBodiesRequest(hashes),
	})
}

// servedBodies decodes the BlockBodies response the handler sent and returns the
// number of bodies it served.
func servedBodies(t *testing.T, rw *recordRW) int {
	t.Helper()
	for i, code := range rw.codes {
		if code != eth.BlockBodiesMsg {
			continue
		}
		var pkt eth.BlockBodiesRLPPacket
		if err := rlp.DecodeBytes(rw.payloads[i], &pkt); err != nil {
			t.Fatalf("decode block bodies response: %v", err)
		}
		return len(pkt.BlockBodiesRLPResponse)
	}
	t.Fatal("no BlockBodies response was sent")
	return 0
}

// TestHandleGetBlockBodiesServesOnlyComplete verifies the serve-gate half of the
// reorder fix: a body is served to peers only when its (validated) header is also
// cached. A provisional body-only entry must not be served.
func TestHandleGetBlockBodiesServesOnlyComplete(t *testing.T) {
	rw := &recordRW{}
	conns := sharedTestConns(t, true)
	c := newTestConn(rw, conns)

	// An empty body built the same way the body handler would, so it encodes
	// cleanly on the serve path.
	emptyRawBody, err := rlp.EncodeToBytes(rawBlockBody{})
	if err != nil {
		t.Fatalf("encode raw body: %v", err)
	}
	body, err := c.buildBlockBody(emptyRawBody)
	if err != nil {
		t.Fatalf("build body: %v", err)
	}

	// Provisional: body present but no header -> must NOT be served.
	bodyOnly := common.HexToHash("0x1111")
	conns.Blocks().Update(bodyOnly, func(bc BlockCache) BlockCache {
		bc.Body = body
		return bc
	})
	// Complete: header + body -> must be served.
	complete := common.HexToHash("0x2222")
	conns.Blocks().Update(complete, func(bc BlockCache) BlockCache {
		bc.Header = &types.Header{Number: big.NewInt(5)}
		bc.Body = body
		return bc
	})

	if err := c.handleGetBlockBodies(makeGetBodies(t, 1, bodyOnly, complete)); err != nil {
		t.Fatalf("handleGetBlockBodies: %v", err)
	}
	if n := servedBodies(t, rw); n != 1 {
		t.Fatalf("want 1 served body (only the complete block), got %d", n)
	}
}

func assertContains(t *testing.T, codes []uint64, want ...uint64) {
	t.Helper()
	set := make(map[uint64]bool, len(codes))
	for _, c := range codes {
		set[c] = true
	}
	for _, w := range want {
		if !set[w] {
			t.Errorf("want a message with code %d in %v", w, codes)
		}
	}
}
