package p2p

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/0xPolygon/polygon-cli/p2p/bortx"
	"github.com/0xPolygon/polygon-cli/p2p/database"
)

// fakeUncle is a header a peer could add to a real block's body.
var fakeUncle = &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(1)}

func TestBodyRootMatchesHeader(t *testing.T) {
	block := loadStateSyncBlock(t)
	txs, err := bortx.Decode(block.Txs)
	if err != nil {
		t.Fatal(err)
	}
	if decodedBodyRoot(txs.Root(), block.Uncles, block.Withdrawals) != headerBodyRoot(block.Header) {
		t.Fatal("the real body does not match its header")
	}
	if decodedBodyRoot(txs.Root(), []*types.Header{fakeUncle}, block.Withdrawals) == headerBodyRoot(block.Header) {
		t.Fatal("a body with an added uncle matches the header")
	}
	withdrawals := []*types.Withdrawal{{Index: 1, Validator: 2, Address: common.Address{3}, Amount: 4}}
	if decodedBodyRoot(txs.Root(), block.Uncles, withdrawals) == headerBodyRoot(block.Header) {
		t.Fatal("a body with withdrawals matches a header without a withdrawals root")
	}

	// No withdrawals, an empty list and the empty withdrawals root are alike.
	if decodedBodyRoot(txs.Root(), block.Uncles, []*types.Withdrawal{}) != headerBodyRoot(block.Header) {
		t.Fatal("an empty withdrawals list does not match a header without a withdrawals root")
	}
	emptyRoot := types.CopyHeader(block.Header)
	empty := types.EmptyWithdrawalsHash
	emptyRoot.WithdrawalsHash = &empty
	if decodedBodyRoot(txs.Root(), block.Uncles, []*types.Withdrawal{}) != headerBodyRoot(emptyRoot) {
		t.Fatal("an empty withdrawals list does not match the empty withdrawals root")
	}

	// Real withdrawals are checked exactly.
	withRoot := types.CopyHeader(block.Header)
	root := types.DeriveSha(types.Withdrawals(withdrawals), trie.NewStackTrie(nil))
	withRoot.WithdrawalsHash = &root
	if decodedBodyRoot(txs.Root(), block.Uncles, withdrawals) != headerBodyRoot(withRoot) {
		t.Fatal("matching withdrawals do not match the header")
	}
	if decodedBodyRoot(txs.Root(), block.Uncles, nil) == headerBodyRoot(withRoot) {
		t.Fatal("a body without withdrawals matches a header with a withdrawals root")
	}
	other := []*types.Withdrawal{{Index: 1, Validator: 2, Address: common.Address{3}, Amount: 5}}
	if decodedBodyRoot(txs.Root(), block.Uncles, other) == headerBodyRoot(withRoot) {
		t.Fatal("different withdrawals match the header")
	}
}

// TestBodyWithFakeUnclesIsNotWritten is the regression test for a body with the
// right transactions but added uncles passing the check, which compared only
// the transaction root, and having its uncles written under the real block.
func TestBodyWithFakeUnclesIsNotWritten(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)
	t.Cleanup(func() { conns.Blocks().Remove(hash) })

	c := newTestConn(&recordRW{}, conns)
	db := &blockWriteDB{Database: database.NoDatabase()}
	c.db = db
	if err := c.handleBlockHeaders(context.Background(), blockHeadersMsg(t, block.Header)); err != nil {
		t.Fatalf("handleBlockHeaders: %v", err)
	}
	c.requests.Add(10, database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()})
	uncles := []*types.Header{fakeUncle}
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 10, block.Txs, uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	if db.bodyWrites != 0 {
		t.Fatal("wrote a body whose uncles do not match its header")
	}
	if cache, _ := conns.Blocks().Peek(hash); cache.Body != nil {
		t.Fatal("cached a body whose uncles do not match its header")
	}
}

// TestNewBlockWithFakeUnclesIsDropped checks the NewBlock path: a block whose
// transactions match the header but whose uncles do not is not accepted.
func TestNewBlockWithFakeUnclesIsDropped(t *testing.T) {
	block := loadStateSyncBlock(t)
	txs, err := bortx.Decode(block.Txs)
	if err != nil {
		t.Fatal(err)
	}
	tampered := block
	tampered.Uncles = []*types.Header{fakeUncle}
	if _, err := blockFromRaw(tampered, txs); err == nil {
		t.Fatal("blockFromRaw accepted uncles that do not match the header")
	}
	if _, err := blockFromRaw(block, txs); err != nil {
		t.Fatalf("blockFromRaw rejected the real block: %v", err)
	}
}

// TestHeldBodyAndArrivingBodyAreWrittenOnce is the regression test for a body
// arriving after its header was cached by a path that did not release the
// held body: the held body was written, then the arriving one, which is the
// same body, was written again.
func TestHeldBodyAndArrivingBodyAreWrittenOnce(t *testing.T) {
	block := loadStateSyncBlock(t)
	hash := block.Header.Hash()
	ann := database.BlockAnnouncement{Hash: hash, Number: block.Header.Number.Uint64()}
	conns := sharedTestConns(t, false)
	useFreshBodyGate(t, conns)
	conns.Blocks().Remove(hash)
	t.Cleanup(func() { conns.Blocks().Remove(hash) })

	c := newTestConn(&recordRW{}, conns)
	db := &blockWriteDB{Database: database.NoDatabase()}
	c.db = db
	c.requests.Add(11, ann)
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 11, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	// The header reaches the cache without going through the gate.
	conns.Blocks().Update(hash, func(BlockCache) BlockCache { return BlockCache{Header: block.Header} })

	c.requests.Add(12, ann)
	if err := c.handleBlockBodies(context.Background(), blockBodiesMsg(t, 12, block.Txs, block.Uncles)); err != nil {
		t.Fatalf("handleBlockBodies: %v", err)
	}
	if db.bodyWrites != 1 {
		t.Fatalf("body writes = %d, want 1", db.bodyWrites)
	}
}
