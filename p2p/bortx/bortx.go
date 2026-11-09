// Package bortx handles the block body transactions Bor produces that upstream
// go-ethereum cannot decode.
//
// Since PIP-74 Bor appends a state-sync transaction (type 0x7f) to the body of
// any block that carries bridge events, and the header's transaction root
// commits to it. go-ethereum's types.Transaction has no such type and cannot be
// given one from outside its package, so these transactions are kept as their
// raw encoding: enough to hash them, to verify a body against its header, and to
// serve the body back to peers byte for byte.
package bortx

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

// StateSyncTxType is Bor's state-sync transaction type.
const StateSyncTxType = 0x7f

// stateSyncData mirrors Bor's types.StateSyncData, the payload of a state-sync
// transaction is an RLP list of these.
type stateSyncData struct {
	ID       uint64
	Contract common.Address
	Data     []byte
	TxHash   common.Hash
}

// Body is the decoded transactions of one block body.
type Body struct {
	// Txs holds every transaction go-ethereum can decode, in body order. Bor
	// state-sync transactions are not in it.
	Txs []*types.Transaction
	// Hashes holds the hash of every transaction in body order, state-sync
	// transactions included, so len(Hashes) is the body's transaction count.
	Hashes []common.Hash

	byIndex  []*types.Transaction
	binaries [][]byte
}

// Tx returns the transaction at body index i, or nil if it is a Bor state-sync
// transaction, which has no go-ethereum form.
func (b *Body) Tx(i int) *types.Transaction {
	return b.byIndex[i]
}

// DecodeError reports the body transaction that failed to decode.
type DecodeError struct {
	Index int
	Count int
	// Type is the transaction type byte, or -1 if the item is not a typed
	// transaction.
	Type int
	Size int
	Hash common.Hash
	Err  error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("transaction %d of %d failed to decode: %v", e.Index, e.Count, e.Err)
}

func (e *DecodeError) Unwrap() error { return e.Err }

func newDecodeError(index, count int, item []byte, err error) *DecodeError {
	encoded := item
	if binary, berr := Binary(item); berr == nil {
		encoded = binary
	}
	txType := -1
	if len(encoded) > 0 && encoded[0] <= StateSyncTxType {
		txType = int(encoded[0])
	}
	return &DecodeError{
		Index: index,
		Count: count,
		Type:  txType,
		Size:  len(encoded),
		Hash:  crypto.Keccak256Hash(encoded),
		Err:   err,
	}
}

// StateSyncCount returns how many of the body's transactions are Bor state-sync
// transactions.
func (b *Body) StateSyncCount() int {
	return len(b.Hashes) - len(b.Txs)
}

// Root returns the transaction root of the body, the value a header's TxHash
// must equal.
func (b *Body) Root() common.Hash {
	return types.DeriveSha(binaryList(b.binaries), trie.NewStackTrie(nil))
}

// RawList returns the body's transactions as the raw list eth.BlockBody holds,
// each in its canonical body encoding: a legacy transaction as its RLP list,
// a typed one (state-sync included) wrapped in an RLP string. This is how
// go-ethereum encodes a transaction, so for a body without state-sync
// transactions the list is the same as encoding the decoded transactions.
func (b *Body) RawList() (rlp.RawList[*types.Transaction], error) {
	var list rlp.RawList[*types.Transaction]
	for _, binary := range b.binaries {
		item := binary
		if binary[0] < 0xc0 {
			enc, err := rlp.EncodeToBytes(binary)
			if err != nil {
				return list, err
			}
			item = enc
		}
		if err := list.AppendRaw(item); err != nil {
			return list, err
		}
	}
	return list, nil
}

// binaryList implements types.DerivableList over consensus-encoded transactions.
type binaryList [][]byte

func (l binaryList) Len() int { return len(l) }

func (l binaryList) EncodeIndex(i int, w *bytes.Buffer) {
	w.Write(l[i])
}

// Binary returns the consensus encoding of one block body transaction item: the
// RLP list itself for a legacy transaction, or type byte and payload for a typed
// one, which a body carries wrapped in an RLP string.
func Binary(item []byte) ([]byte, error) {
	kind, content, _, err := rlp.Split(item)
	if err != nil {
		return nil, err
	}
	switch kind {
	case rlp.List:
		return item, nil
	case rlp.String:
		if len(content) == 0 {
			return nil, errors.New("empty typed transaction")
		}
		return content, nil
	default:
		return nil, errors.New("transaction is a single byte")
	}
}

// IsStateSync reports whether a consensus-encoded transaction is a well-formed
// Bor state-sync transaction.
func IsStateSync(binary []byte) bool {
	if len(binary) == 0 || binary[0] != StateSyncTxType {
		return false
	}
	var data []stateSyncData
	return rlp.DecodeBytes(binary[1:], &data) == nil
}

// Decode decodes the transaction items of a block body. It fails if any item is
// neither a transaction go-ethereum can decode nor a Bor state-sync
// transaction, because a body missing a transaction is not the body its block
// hash commits to.
func Decode(items []rlp.RawValue) (*Body, error) {
	body := newBody(len(items))
	for i, item := range items {
		if err := body.add(item); err != nil {
			return nil, newDecodeError(i, len(items), item, err)
		}
	}
	return body, nil
}

// DecodeList decodes the transactions of a block body held as a raw list, the
// form eth.BlockBody carries them in.
func DecodeList(list *rlp.RawList[*types.Transaction]) (*Body, error) {
	body := newBody(list.Len())
	it := list.ContentIterator()
	for i := 0; it.Next(); i++ {
		if err := body.add(it.Value()); err != nil {
			return nil, newDecodeError(i, list.Len(), it.Value(), err)
		}
	}
	if err := it.Err(); err != nil {
		return nil, err
	}
	return body, nil
}

func newBody(n int) *Body {
	return &Body{
		Txs:      make([]*types.Transaction, 0, n),
		Hashes:   make([]common.Hash, 0, n),
		byIndex:  make([]*types.Transaction, 0, n),
		binaries: make([][]byte, 0, n),
	}
}

func (b *Body) add(item []byte) error {
	binary, err := Binary(item)
	if err != nil {
		return err
	}
	if IsStateSync(binary) {
		b.Hashes = append(b.Hashes, crypto.Keccak256Hash(binary))
		b.byIndex = append(b.byIndex, nil)
		b.binaries = append(b.binaries, binary)
		return nil
	}
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(binary); err != nil {
		return err
	}
	b.Txs = append(b.Txs, tx)
	b.Hashes = append(b.Hashes, tx.Hash())
	b.byIndex = append(b.byIndex, tx)
	b.binaries = append(b.binaries, binary)
	return nil
}
