// This file is a combination of these two files:
// https://github.com/maticnetwork/bor/blob/master/core/rawdb/bor_receipt.go
// https://github.com/maticnetwork/bor/blob/master/core/types/bor_receipt.go
//
// Ideally we could import bor and geth separately, but it's a little complicated. So in this case we're manually vendoring the bor module
package borfsck

import (
	"bytes"
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

var (
	// borTxLookupPrefix + hash -> transaction/receipt lookup metadata
	borTxLookupPrefix = []byte(borTxLookupPrefixStr)

	borReceiptPrefix = []byte("matic-bor-receipt-") // borReceiptPrefix + number + block hash -> bor block receipt

)

// BorReceiptKey = borReceiptPrefix + num (uint64 big endian) + hash
func borReceiptKey(number uint64, hash common.Hash) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)

	return append(append(borReceiptPrefix, enc...), hash.Bytes()...)
}

// isCanon is an internal utility method, to check whether the given number/hash
// is part of the ancient (canon) set.
func isCanon(reader ethdb.AncientReaderOp, number uint64, hash common.Hash) bool {
	h, err := reader.Ancient(rawdb.ChainFreezerHashTable, number)
	if err != nil {
		return false
	}

	return bytes.Equal(h, hash[:])
}

const (
	borTxLookupPrefixStr = "matic-bor-tx-lookup-"

	// freezerBorReceiptTable indicates the name of the freezer bor receipts table.
	freezerBorReceiptTable = "matic-bor-receipts"
)

// borTxLookupKey = borTxLookupPrefix + bor tx hash
func borTxLookupKey(hash common.Hash) []byte {
	return append(borTxLookupPrefix, hash.Bytes()...)
}

func ReadBorReceiptRLP(db ethdb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	var data []byte

	err := db.ReadAncients(func(reader ethdb.AncientReaderOp) error {
		// Check if the data is in ancients
		if isCanon(reader, number, hash) {
			data, _ = reader.Ancient(freezerBorReceiptTable, number)

			return nil
		}

		// If not, try reading from leveldb
		data, _ = db.Get(borReceiptKey(number, hash))

		return nil
	})

	if err != nil {
		log.Warn("during ReadBorReceiptRLP", "number", number, "hash", hash, "err", err)
	}

	return data
}

// ReadRawBorReceipt retrieves the block receipt belonging to a block.
// The receipt metadata fields are not guaranteed to be populated, so they
// should not be used. Use ReadBorReceipt instead if the metadata is needed.
func ReadRawBorReceipt(db ethdb.Reader, hash common.Hash, number uint64) *types.Receipt {
	// Retrieve the flattened receipt slice
	data := ReadBorReceiptRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}

	// Convert the receipts from their storage form to their internal representation
	var storageReceipt types.ReceiptForStorage
	if err := rlp.DecodeBytes(data, &storageReceipt); err != nil {
		log.Error("Invalid receipt array RLP", "hash", hash, "err", err)
		return nil
	}

	return (*types.Receipt)(&storageReceipt)
}

// ReadBorReceipt retrieves all the bor block receipts belonging to a block, including
// its correspoinding metadata fields. If it is unable to populate these metadata
// fields then nil is returned.
func ReadBorReceipt(db ethdb.Reader, hash common.Hash, number uint64, config *params.ChainConfig) *types.Receipt {

	// We're deriving many fields from the block body, retrieve beside the receipt
	borReceipt := ReadRawBorReceipt(db, hash, number)
	if borReceipt == nil {
		return nil
	}

	// We're deriving many fields from the block body, retrieve beside the receipt
	receipts := rawdb.ReadRawReceipts(db, hash, number)
	if receipts == nil {
		return nil
	}

	body := rawdb.ReadBody(db, hash, number)
	if body == nil {
		log.Error("Missing body but have bor receipt", "hash", hash, "number", number)
		return nil
	}

	if err := DeriveFieldsForBorReceipt(borReceipt, hash, number, receipts); err != nil {
		log.Error("Failed to derive bor receipt fields", "hash", hash, "number", number, "err", err)
		return nil
	}

	return borReceipt
}

func DeriveFieldsForBorReceipt(receipt *types.Receipt, hash common.Hash, number uint64, receipts types.Receipts) error {
	// get derived tx hash
	txHash := GetDerivedBorTxHash(BorReceiptKey(number, hash))
	txIndex := uint(len(receipts))

	// set tx hash and tx index
	receipt.TxHash = txHash
	receipt.TransactionIndex = txIndex
	receipt.BlockHash = hash
	receipt.BlockNumber = big.NewInt(0).SetUint64(number)

	logIndex := 0
	for i := 0; i < len(receipts); i++ {
		logIndex += len(receipts[i].Logs)
	}

	// The derived log fields can simply be set from the block and transaction
	for j := 0; j < len(receipt.Logs); j++ {
		receipt.Logs[j].BlockNumber = number
		receipt.Logs[j].BlockHash = hash
		receipt.Logs[j].TxHash = txHash
		receipt.Logs[j].TxIndex = txIndex
		receipt.Logs[j].Index = uint(logIndex)
		logIndex++
	}

	return nil
}

// BorReceiptKey = borReceiptPrefix + num (uint64 big endian) + hash
func BorReceiptKey(number uint64, hash common.Hash) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)

	return append(append(borReceiptPrefix, enc...), hash.Bytes()...)
}

// GetDerivedBorTxHash get derived tx hash from receipt key
func GetDerivedBorTxHash(receiptKey []byte) common.Hash {
	return common.BytesToHash(crypto.Keccak256(receiptKey))
}

// WriteBorReceipt stores all the bor receipt belonging to a block.
func WriteBorReceipt(db ethdb.KeyValueWriter, hash common.Hash, number uint64, borReceipt *types.ReceiptForStorage) {
	// Convert the bor receipt into their storage form and serialize them
	bytes, err := rlp.EncodeToBytes(borReceipt)
	if err != nil {
		log.Crit("Failed to encode bor receipt", "err", err)
	}

	// Store the flattened receipt slice
	if err := db.Put(borReceiptKey(number, hash), bytes); err != nil {
		log.Crit("Failed to store bor receipt", "err", err)
	}
}

// DeleteBorReceipt removes receipt data associated with a block hash.
func DeleteBorReceipt(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	key := borReceiptKey(number, hash)

	if err := db.Delete(key); err != nil {
		log.Crit("Failed to delete bor receipt", "err", err)
	}
}

// ReadBorTransactionWithBlockHash retrieves a specific bor (fake) transaction by tx hash and block hash, along with
// its added positional metadata.
func ReadBorTransactionWithBlockHash(db ethdb.Reader, txHash common.Hash, blockHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	blockNumber := ReadBorTxLookupEntry(db, txHash)
	if blockNumber == nil {
		return nil, common.Hash{}, 0, 0
	}

	body := rawdb.ReadBody(db, blockHash, *blockNumber)
	if body == nil {
		log.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash)
		return nil, common.Hash{}, 0, 0
	}

	// fetch receipt and return it
	return NewBorTransaction(), blockHash, *blockNumber, uint64(len(body.Transactions))
}

// NewBorTransaction create new bor transaction for bor receipt
func NewBorTransaction() *types.Transaction {
	return types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), make([]byte, 0))
}

// ReadBorTransaction retrieves a specific bor (fake) transaction by hash, along with
// its added positional metadata.
func ReadBorTransaction(db ethdb.Reader, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	blockNumber := ReadBorTxLookupEntry(db, hash)
	if blockNumber == nil {
		return nil, common.Hash{}, 0, 0
	}

	blockHash := rawdb.ReadCanonicalHash(db, *blockNumber)
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}

	body := rawdb.ReadBody(db, blockHash, *blockNumber)
	if body == nil {
		log.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash)
		return nil, common.Hash{}, 0, 0
	}

	// fetch receipt and return it
	return NewBorTransaction(), blockHash, *blockNumber, uint64(len(body.Transactions))
}

//
// Indexes for reverse lookup
//

// ReadBorTxLookupEntry retrieves the positional metadata associated with a transaction
// hash to allow retrieving the bor transaction or bor receipt using tx hash.
func ReadBorTxLookupEntry(db ethdb.Reader, txHash common.Hash) *uint64 {
	data, _ := db.Get(borTxLookupKey(txHash))
	if len(data) == 0 {
		return nil
	}

	number := new(big.Int).SetBytes(data).Uint64()

	return &number
}

// WriteBorTxLookupEntry stores a positional metadata for bor transaction using block hash and block number
func WriteBorTxLookupEntry(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	txHash := GetDerivedBorTxHash(borReceiptKey(number, hash))
	if err := db.Put(borTxLookupKey(txHash), big.NewInt(0).SetUint64(number).Bytes()); err != nil {
		log.Crit("Failed to store bor transaction lookup entry", "err", err)
	}
}

// DeleteBorTxLookupEntry removes bor transaction data associated with block hash and block number
func DeleteBorTxLookupEntry(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	txHash := GetDerivedBorTxHash(borReceiptKey(number, hash))
	DeleteBorTxLookupEntryByTxHash(db, txHash)
}

// DeleteBorTxLookupEntryByTxHash removes bor transaction data associated with a bor tx hash.
func DeleteBorTxLookupEntryByTxHash(db ethdb.KeyValueWriter, txHash common.Hash) {
	if err := db.Delete(borTxLookupKey(txHash)); err != nil {
		log.Crit("Failed to delete bor transaction lookup entry", "err", err)
	}
}
