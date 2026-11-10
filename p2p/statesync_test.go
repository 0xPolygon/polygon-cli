package p2p

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	ethp2p "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/0xPolygon/polygon-cli/p2p/bortx"
	"github.com/0xPolygon/polygon-cli/p2p/database"
	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

// stateSyncFixture is mainnet block 97,029,696, whose last transaction (index 48
// of 49) is a Bor state-sync transaction (type 0x7f).
const stateSyncFixture = "bortx/testdata/block_97029696.rlp.hex"

func loadStateSyncBlock(t *testing.T) rawBlock {
	t.Helper()
	data, err := os.ReadFile(stateSyncFixture)
	if err != nil {
		t.Fatal(err)
	}
	enc, err := hex.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatal(err)
	}
	var block rawBlock
	if err := rlp.DecodeBytes(enc, &block); err != nil {
		t.Fatal(err)
	}
	return block
}

// blockWriteDB records the block and transactions passed to WriteBlock.
type blockWriteDB struct {
	database.Database
	block      *types.Block
	txs        *bortx.Body
	bodyWrites int
	bodyTfs    time.Time
}

func (d *blockWriteDB) WriteBlock(_ context.Context, _ *enode.Node, block *types.Block, txs *bortx.Body, _ *big.Int, _ time.Time) {
	d.block = block
	d.txs = txs
}

// ShouldWriteBlocks reports true: blockWriteDB stands in for a database that
// stores blocks, so the handlers hold and write bodies for it.
func (d *blockWriteDB) ShouldWriteBlocks() bool { return true }

func (d *blockWriteDB) WriteBlockBody(_ context.Context, _ *eth.BlockBody, txs *bortx.Body, _ database.BlockAnnouncement, tfs time.Time) {
	d.txs = txs
	d.bodyWrites++
	d.bodyTfs = tfs
}

// TestBuildBlockBodyKeepsStateSyncTx is the regression test for BlockBodies
// responses carrying a Bor state-sync transaction: the body used to be rejected
// whole, and must now be kept with every transaction, byte for byte.
func TestBuildBlockBodyKeepsStateSyncTx(t *testing.T) {
	block := loadStateSyncBlock(t)
	raw, err := rlp.EncodeToBytes(rawBlockBody{Transactions: block.Txs, Uncles: block.Uncles})
	if err != nil {
		t.Fatal(err)
	}

	c := &conn{}
	built, txs, err := c.buildBlockBody(raw)
	if err != nil {
		t.Fatalf("buildBlockBody: %v", err)
	}

	want, err := rlp.EncodeToBytes(block.Txs)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(built.Transactions.Bytes(), want) {
		t.Fatal("built body transactions differ from the ones received")
	}
	if len(txs.Hashes) != len(block.Txs) || txs.StateSyncCount() != 1 {
		t.Fatalf("decoded %d transactions with %d state-sync, want %d with 1", len(txs.Hashes), txs.StateSyncCount(), len(block.Txs))
	}
	if txs.Root() != block.Header.TxHash {
		t.Fatalf("built body root %s does not match header TxHash %s", txs.Root().Hex(), block.Header.TxHash.Hex())
	}
}

// TestHandleNewBlockKeepsStateSyncTx is the regression test for NewBlock messages
// carrying a Bor state-sync transaction, which used to be dropped as undecodable.
// The block must be accepted, and the database and serving cache must both get
// the whole body.
func TestHandleNewBlockKeepsStateSyncTx(t *testing.T) {
	block := loadStateSyncBlock(t)
	conns := sharedTestConns(t, false)
	c := newTestConn(&recordRW{}, conns)
	db := &blockWriteDB{Database: database.NoDatabase()}
	c.db = db

	msg := encodeMsg(t, eth.NewBlockMsg, &rawNewBlockPacket{Block: block, TD: big.NewInt(1)})
	if err := c.handleNewBlock(context.Background(), msg); err != nil {
		t.Fatalf("handleNewBlock: %v", err)
	}

	hash := block.Header.Hash()
	if db.txs == nil {
		t.Fatal("block with a state-sync transaction was not written")
	}
	if got := len(db.txs.Hashes); got != len(block.Txs) {
		t.Fatalf("written body has %d transactions, want %d", got, len(block.Txs))
	}
	if got := len(db.block.Transactions()); got != len(block.Txs)-1 {
		t.Fatalf("written block holds %d transactions, want %d (all but the state-sync one)", got, len(block.Txs)-1)
	}
	cache, ok := conns.Blocks().Peek(hash)
	if !ok || cache.Body == nil {
		t.Fatal("block with a state-sync transaction was not cached")
	}
	if head := conns.HeadBlock(); head.Block.Hash() == hash && (head.Txs == nil || len(head.Txs.Hashes) != len(block.Txs)) {
		t.Fatal("head block does not carry the whole decoded body")
	}
	cached, err := bortx.DecodeList(&cache.Body.Transactions)
	if err != nil {
		t.Fatalf("DecodeList cached body: %v", err)
	}
	if cached.Root() != block.Header.TxHash {
		t.Fatalf("cached body root %s does not match header TxHash %s", cached.Root().Hex(), block.Header.TxHash.Hex())
	}
}

// TestDecodeNewBlockKeepsStateSyncTx is the regression test for the rlpx reader
// used by `p2p ping`: a NewBlock carrying a Bor state-sync transaction decoded to
// an *Error, and ReadAndServe ends the connection on one.
func TestDecodeNewBlockKeepsStateSyncTx(t *testing.T) {
	block := loadStateSyncBlock(t)
	raw, err := rlp.EncodeToBytes(&rawNewBlockPacket{Block: block, TD: big.NewInt(1)})
	if err != nil {
		t.Fatal(err)
	}

	msg, ok := decodeNewBlock(raw).(*NewBlock)
	if !ok {
		t.Fatalf("decodeNewBlock returned %T, want *NewBlock", decodeNewBlock(raw))
	}
	if msg.Block.Hash() != block.Header.Hash() {
		t.Fatalf("block hash = %s, want %s", msg.Block.Hash().Hex(), block.Header.Hash().Hex())
	}
	if got := len(msg.Block.Transactions()); got != len(block.Txs)-1 {
		t.Fatalf("block holds %d transactions, want %d (all but the state-sync one)", got, len(block.Txs)-1)
	}
}

// blockBodiesMsg wraps one raw body in a BlockBodies response.
func blockBodiesMsg(t *testing.T, reqID uint64, txs []rlp.RawValue, uncles []*types.Header) ethp2p.Msg {
	t.Helper()
	bodyRLP, err := rlp.EncodeToBytes(rawBlockBody{Transactions: txs, Uncles: uncles})
	if err != nil {
		t.Fatal(err)
	}
	return encodeMsg(t, eth.BlockBodiesMsg, &eth.BlockBodiesRLPPacket{
		RequestId:              reqID,
		BlockBodiesRLPResponse: eth.BlockBodiesRLPResponse{rlp.RawValue(bodyRLP)},
	})
}

// TestHandleBlockBodiesChecksRootAgainstCachedHeader checks that a body which
// does not match its already-cached header is neither stored nor cached, while
// the real body is.
func TestHandleBlockBodiesChecksRootAgainstCachedHeader(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	conns := sharedTestConns(t, false)
	// The test Conns is shared, and other tests cache this block whole; start
	// from the header alone so the body path runs.
	conns.Blocks().Update(hash, func(BlockCache) BlockCache {
		return BlockCache{Header: block.Header}
	})

	ann := database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()}

	c := newTestConn(&recordRW{}, conns)
	db := &blockWriteDB{Database: database.NoDatabase()}
	c.db = db
	c.requests.Add(1, ann)
	// Well formed, but missing a transaction, so its root does not match.
	tampered := append([]rlp.RawValue{}, block.Txs[1:]...)
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 1, tampered, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	if db.txs != nil {
		t.Fatalf("stored a body that does not match its header (%d transactions)", len(db.txs.Hashes))
	}
	if cache, _ := conns.Blocks().Peek(hash); cache.Body != nil {
		t.Fatal("cached a body that does not match its header")
	}

	c.requests.Add(2, ann)
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 2, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	if db.txs == nil || len(db.txs.Hashes) != len(block.Txs) {
		t.Fatal("did not store the body that matches its header")
	}
}

// TestNewBlockFromBodyMatchesNetworkEncoding checks that a NewBlock rebuilt from
// the header and cached raw body is byte for byte the message the network sent,
// state-sync transaction included, so rebroadcasting it does not alter the body.
func TestNewBlockFromBodyMatchesNetworkEncoding(t *testing.T) {
	block := loadStateSyncBlock(t)
	td := big.NewInt(12345)
	want, err := rlp.EncodeToBytes(&rawNewBlockPacket{Block: block, TD: td})
	if err != nil {
		t.Fatal(err)
	}

	txs, err := bortx.Decode(block.Txs)
	if err != nil {
		t.Fatal(err)
	}
	body, err := newRawBlockBody(txs, block.Uncles, block.Withdrawals)
	if err != nil {
		t.Fatal(err)
	}
	got, err := rlp.EncodeToBytes(newBlockFromBody(block.Header, body, td))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("NewBlock rebuilt from the cached body differs from the network encoding")
	}
}

// TestBodyBeforeHeaderIsCheckedWhenHeaderArrives covers a body that arrives
// before its header: it cannot be checked then, so it must be checked, and
// dropped if it does not match, when the header is cached.
func TestBodyBeforeHeaderIsCheckedWhenHeaderArrives(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)

	c := newTestConn(&recordRW{}, conns)
	c.requests.Add(3, database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()})
	tampered := append([]rlp.RawValue{}, block.Txs[1:]...)
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 3, tampered, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	if cache, _ := conns.Blocks().Peek(hash); cache.Body == nil {
		t.Fatal("body arriving before its header was not held provisionally")
	}

	c.cacheAndAnnounceHeader(block.Header, false, false)
	cache, ok := conns.Blocks().Peek(hash)
	if !ok || cache.Header == nil {
		t.Fatal("header was not cached")
	}
	if cache.Body != nil {
		t.Fatal("kept a provisional body that does not match the header that arrived")
	}
}

// blockHeadersMsg wraps headers in a BlockHeaders response.
func blockHeadersMsg(t *testing.T, headers ...*types.Header) ethp2p.Msg {
	t.Helper()
	list, err := rlp.EncodeToRawList(headers)
	if err != nil {
		t.Fatal(err)
	}
	return encodeMsg(t, eth.BlockHeadersMsg, &eth.BlockHeadersPacket{RequestId: 0, List: list})
}

// TestBodyIsWrittenOnlyOnceCheckedAgainstHeader is the regression test for a body
// that arrives before its header being written to the database unchecked. It
// must be held, then written with its arrival time if the header matches, and
// never written if it does not.
func TestBodyIsWrittenOnlyOnceCheckedAgainstHeader(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	ann := database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()}
	conns := sharedTestConns(t, false)

	cases := []struct {
		name      string
		txs       []rlp.RawValue
		wantWrite bool
	}{
		{"matching body", block.Txs, true},
		{"body missing a transaction", block.Txs[1:], false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			useFreshBodyGate(t, conns)
			conns.Blocks().Remove(hash)

			c := newTestConn(&recordRW{}, conns)
			db := &blockWriteDB{Database: database.NoDatabase()}
			c.db = db

			c.requests.Add(7, ann)
			if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 7, tc.txs, block.Uncles)); err != nil {
				t.Fatalf("handleBlockBodies: %v", err)
			}
			if db.bodyWrites != 0 {
				t.Fatal("wrote a body before its header arrived")
			}
			arrived := time.Now()

			if err := c.handleBlockHeaders(context.Background(), blockHeadersMsg(t, block.Header)); err != nil {
				t.Fatalf("handleBlockHeaders: %v", err)
			}
			if tc.wantWrite {
				if db.bodyWrites != 1 || len(db.txs.Hashes) != len(block.Txs) {
					t.Fatalf("body writes = %d, want 1 with %d transactions", db.bodyWrites, len(block.Txs))
				}
				if db.bodyTfs.After(arrived) {
					t.Fatalf("body written with time %v, want its arrival time (before %v)", db.bodyTfs, arrived)
				}
				return
			}
			if db.bodyWrites != 0 {
				t.Fatal("wrote a body that does not match its header")
			}
		})
	}
}

// TestBodyAfterHeaderIsWrittenImmediately covers the usual order: the header is
// known, so the body is checked and written as it arrives.
func TestBodyAfterHeaderIsWrittenImmediately(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)

	c := newTestConn(&recordRW{}, conns)
	db := &blockWriteDB{Database: database.NoDatabase()}
	c.db = db

	if err := c.handleBlockHeaders(context.Background(), blockHeadersMsg(t, block.Header)); err != nil {
		t.Fatalf("handleBlockHeaders: %v", err)
	}
	c.requests.Add(8, database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()})
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 8, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	if db.bodyWrites != 1 {
		t.Fatalf("body writes = %d, want 1", db.bodyWrites)
	}
}

// useFreshBodyGate gives a test an empty body gate on the shared test Conns and
// restores the previous one when the test ends.
func useFreshBodyGate(t *testing.T, conns *Conns) {
	t.Helper()
	prev := conns.bodies
	conns.bodies = newBodyGate(ds.LRUOptions{MaxSize: 64})
	t.Cleanup(func() { conns.bodies = prev })
}

// failRW fails every write, so a handler's parent-block request fails.
type failRW struct{ recordRW }

func (*failRW) WriteMsg(ethp2p.Msg) error { return errors.New("write failed") }

// parentFetchDB makes getParentBlock request the parent: block writes are on and
// no parent is stored.
type parentFetchDB struct{ *blockWriteDB }

func (parentFetchDB) ShouldWriteBlocks() bool                    { return true }
func (parentFetchDB) HasBlock(context.Context, common.Hash) bool { return false }

// TestCorrectBodyIsNotShutOutByUncheckedCachedBody is the regression test for a
// wrong body cached before its header making handleBlockBodies skip the correct
// body from another peer, so that no body was ever written for the block.
func TestCorrectBodyIsNotShutOutByUncheckedCachedBody(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	ann := database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()}
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)
	t.Cleanup(func() { conns.Blocks().Remove(hash) })

	db := &blockWriteDB{Database: database.NoDatabase()}
	bad := newTestConn(&recordRW{}, conns)
	bad.db = db
	bad.requests.Add(1, ann)
	if err := bad.handleBlockBodies(context.Background(), blockBodiesMsg(t, 1, block.Txs[1:], block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies (wrong body): %v", err)
	}
	good := newTestConn(&recordRW{}, conns)
	good.db = db
	good.requests.Add(2, ann)
	if err := good.handleBlockBodies(context.Background(), blockBodiesMsg(t, 2, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies (correct body): %v", err)
	}
	if err := good.handleBlockHeaders(context.Background(), blockHeadersMsg(t, block.Header)); err != nil {
		t.Fatalf("handleBlockHeaders: %v", err)
	}
	if db.bodyWrites != 1 || len(db.txs.Hashes) != len(block.Txs) {
		t.Fatalf("body writes = %d, want the correct body written once", db.bodyWrites)
	}
}

// TestHeldBodyIsWrittenWhenHeaderHandlerFailsEarly is the regression test for
// handleBlockHeaders returning on a failed parent request before releasing the
// bodies held for its headers.
func TestHeldBodyIsWrittenWhenHeaderHandlerFailsEarly(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)
	t.Cleanup(func() { conns.Blocks().Remove(hash) })

	db := &blockWriteDB{Database: database.NoDatabase()}
	c := newTestConn(&failRW{}, conns)
	c.db = parentFetchDB{db}
	c.requests.Add(3, database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()})
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 3, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	if err := c.handleBlockHeaders(context.Background(), blockHeadersMsg(t, block.Header)); err == nil {
		t.Fatal("handleBlockHeaders succeeded; the test needs the parent request to fail")
	}
	if db.bodyWrites != 1 {
		t.Fatalf("body writes = %d, want the held body written once", db.bodyWrites)
	}
}

// TestHeldBodyIsWrittenWhenNewBlockHandlerFailsEarly is the regression test for
// handleNewBlock discarding a held body and then returning on a failed parent
// request before writing its own.
func TestHeldBodyIsWrittenWhenNewBlockHandlerFailsEarly(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)
	t.Cleanup(func() { conns.Blocks().Remove(hash) })

	db := &blockWriteDB{Database: database.NoDatabase()}
	c := newTestConn(&failRW{}, conns)
	c.db = parentFetchDB{db}
	c.requests.Add(4, database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()})
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 4, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	msg := encodeMsg(t, eth.NewBlockMsg, &rawNewBlockPacket{Block: block, TD: big.NewInt(1)})
	if err := c.handleNewBlock(context.Background(), msg); err == nil {
		t.Fatal("handleNewBlock succeeded; the test needs the parent request to fail")
	}
	if db.bodyWrites != 1 {
		t.Fatalf("body writes = %d, want the held body written once", db.bodyWrites)
	}
}

// TestBodyIsNotHeldWhenNothingWritesIt checks that a body arriving before its
// header is not held when the database writes neither blocks nor transactions,
// since holding it would only use memory.
func TestBodyIsNotHeldWhenNothingWritesIt(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)
	t.Cleanup(func() { conns.Blocks().Remove(hash) })

	c := newTestConn(&recordRW{}, conns)
	c.requests.Add(5, database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()})
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 5, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	if _, held := conns.bodies.pending.Peek(hash); held {
		t.Fatal("held a body for a database that does not write it")
	}
}

// TestNewBlockDoesNotRewriteHeldBody checks that when a NewBlock arrives for a
// block whose body is held, only the NewBlock's write happens: the held body is
// not written as well.
func TestNewBlockDoesNotRewriteHeldBody(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)
	t.Cleanup(func() { conns.Blocks().Remove(hash) })

	c := newTestConn(&recordRW{}, conns)
	db := &blockWriteDB{Database: database.NoDatabase()}
	c.db = db
	c.requests.Add(6, database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()})
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 6, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	msg := encodeMsg(t, eth.NewBlockMsg, &rawNewBlockPacket{Block: block, TD: big.NewInt(1)})
	if err := c.handleNewBlock(context.Background(), msg); err != nil {
		t.Fatalf("handleNewBlock: %v", err)
	}
	if db.block == nil {
		t.Fatal("the NewBlock was not written")
	}
	if db.bodyWrites != 0 {
		t.Fatalf("held body written %d times as well as the NewBlock", db.bodyWrites)
	}
	if _, held := conns.bodies.pending.Peek(hash); held {
		t.Fatal("the held body was left in the gate")
	}
}
