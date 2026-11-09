package p2p

import (
	"bytes"
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/0xPolygon/polygon-cli/p2p/bortx"
	"github.com/0xPolygon/polygon-cli/p2p/database"
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

// blockWriteDB records the block and body passed to WriteBlock.
type blockWriteDB struct {
	database.Database
	block *types.Block
	body  *eth.BlockBody
}

func (d *blockWriteDB) WriteBlock(_ context.Context, _ *enode.Node, block *types.Block, body *eth.BlockBody, _ *big.Int, _ time.Time) {
	d.block = block
	d.body = body
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
	built, err := c.buildBlockBody(raw)
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
	txs, err := bortx.DecodeList(&built.Transactions)
	if err != nil {
		t.Fatalf("DecodeList: %v", err)
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
	if db.body == nil {
		t.Fatal("block with a state-sync transaction was not written")
	}
	if got := db.body.Transactions.Len(); got != len(block.Txs) {
		t.Fatalf("written body has %d transactions, want %d", got, len(block.Txs))
	}
	if got := len(db.block.Transactions()); got != len(block.Txs)-1 {
		t.Fatalf("written block holds %d transactions, want %d (all but the state-sync one)", got, len(block.Txs)-1)
	}
	cache, ok := conns.Blocks().Peek(hash)
	if !ok || cache.Body == nil {
		t.Fatal("block with a state-sync transaction was not cached")
	}
	cached, err := bortx.DecodeList(&cache.Body.Transactions)
	if err != nil {
		t.Fatalf("DecodeList cached body: %v", err)
	}
	if cached.Root() != block.Header.TxHash {
		t.Fatalf("cached body root %s does not match header TxHash %s", cached.Root().Hex(), block.Header.TxHash.Hex())
	}
}
