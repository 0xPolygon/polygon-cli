package p2p

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"golang.org/x/exp/slices"
)

type DatastoreEvent struct {
	SensorId string
	PeerId   string
	Hash     *datastore.Key
	Time     time.Time
}

// DatastoreHeader stores the data in manner that can be easily written without
// loss of precision.
type DatastoreHeader struct {
	ParentHash  *datastore.Key
	UncleHash   string
	Coinbase    string
	Root        string
	TxHash      string
	ReceiptHash string
	Bloom       []byte
	Difficulty  string
	Number      string
	GasLimit    string
	GasUsed     string
	Time        time.Time
	Extra       []byte
	MixDigest   string
	Nonce       string
	BaseFee     string
}

type DatastoreBlock struct {
	*DatastoreHeader
	Transactions []*datastore.Key
	Uncles       []*datastore.Key
}

type DatastoreTransaction struct {
	BlockHashes []*datastore.Key
	Data        []byte `datastore:",noindex"`
	From        string
	Gas         string
	GasFeeCap   string
	GasPrice    string
	GasTipCap   string
	Nonce       string
	To          string
	Value       string
	V, R, S     string
	Time        time.Time
	Type        int16
}

func NewDatastoreHeader(header *types.Header) *DatastoreHeader {
	return &DatastoreHeader{
		ParentHash:  datastore.NameKey("blocks", header.ParentHash.Hex(), nil),
		UncleHash:   header.UncleHash.Hex(),
		Coinbase:    header.Coinbase.Hex(),
		Root:        header.Root.Hex(),
		TxHash:      header.TxHash.Hex(),
		ReceiptHash: header.ReceiptHash.Hex(),
		Bloom:       header.Bloom.Bytes(),
		Difficulty:  header.Difficulty.String(),
		Number:      header.Number.String(),
		GasLimit:    fmt.Sprint(header.GasLimit),
		GasUsed:     fmt.Sprint(header.GasUsed),
		Time:        time.Unix(int64(header.Time), 0),
		Extra:       header.Extra,
		MixDigest:   header.MixDigest.String(),
		Nonce:       fmt.Sprint(header.Nonce.Uint64()),
		BaseFee:     header.BaseFee.String(),
	}
}

func NewDatastoreTransaction(tx *types.Transaction) *DatastoreTransaction {
	v, r, s := tx.RawSignatureValues()
	var from, to string

	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
	if err == nil {
		from = msg.From().Hex()
	}

	if tx.To() != nil {
		to = tx.To().Hex()
	}

	return &DatastoreTransaction{
		Data:      tx.Data(),
		From:      from,
		Gas:       fmt.Sprint(tx.Gas()),
		GasFeeCap: tx.GasFeeCap().String(),
		GasPrice:  tx.GasPrice().String(),
		GasTipCap: tx.GasTipCap().String(),
		Nonce:     fmt.Sprint(tx.Nonce()),
		To:        to,
		Value:     tx.Value().String(),
		V:         v.String(),
		R:         r.String(),
		S:         s.String(),
		Time:      time.Now(),
		Type:      int16(tx.Type()),
	}
}

func (c *Conn) writeEvent(ctx context.Context, client *datastore.Client, eventKind string, hash common.Hash, hashKind string) {
	key := datastore.IncompleteKey(eventKind, nil)
	event := DatastoreEvent{
		SensorId: c.Sensor,
		PeerId:   c.node.URLv4(),
		Hash:     datastore.NameKey(hashKind, hash.Hex(), nil),
		Time:     time.Now(),
	}
	if _, err := client.Put(ctx, key, &event); err != nil {
		c.logger.Error().Err(err).Msgf("Failed to write to %v", eventKind)
	}
}

func (c *Conn) writeEvents(ctx context.Context, client *datastore.Client, eventKind string, hashes []common.Hash, hashKind string) {
	keys := make([]*datastore.Key, 0, len(hashes))
	events := make([]*DatastoreEvent, 0, len(hashes))
	now := time.Now()

	for _, hash := range hashes {
		key := datastore.IncompleteKey(eventKind, nil)
		keys = append(keys, key)

		event := DatastoreEvent{
			SensorId: c.Sensor,
			PeerId:   c.node.URLv4(),
			Hash:     datastore.NameKey(hashKind, hash.Hex(), nil),
			Time:     now,
		}
		events = append(events, &event)
	}

	if _, err := client.PutMulti(ctx, keys, events); err != nil {
		c.logger.Error().Err(err).Msgf("Failed to write to %v", eventKind)
	}
}

func (c *Conn) writeBlockHeader(ctx context.Context, client *datastore.Client, header *types.Header) {
	key := datastore.NameKey("blocks", header.Hash().Hex(), nil)
	var block DatastoreBlock

	if err := client.Get(ctx, key, &block); err == nil && block.DatastoreHeader != nil {
		return
	}

	block.DatastoreHeader = NewDatastoreHeader(header)

	if _, err := client.Put(ctx, key, &block); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write block header")
	}
}

func (c *Conn) writeBlockHeaders(ctx context.Context, client *datastore.Client, headers []*types.Header) {
	for _, header := range headers {
		c.writeBlockHeader(ctx, client, header)
	}
}

func (c *Conn) writeBlockBody(ctx context.Context, client *datastore.Client, hash string, body *eth.BlockBody) {
	key := datastore.NameKey("blocks", hash, nil)
	var block DatastoreBlock

	if err := client.Get(ctx, key, &block); err != nil {
		c.logger.Warn().Err(err).Str("hash", hash).Msg("Failed to fetch block when writing block body")
	}

	if block.Transactions == nil {
		block.Transactions = make([]*datastore.Key, 0, len(body.Transactions))
		for _, tx := range body.Transactions {
			c.writeTransaction(ctx, client, tx, hash)
			block.Transactions = append(block.Transactions, datastore.NameKey("transactions", tx.Hash().Hex(), nil))
		}
	}

	if block.Uncles == nil {
		c.writeBlockHeaders(ctx, client, body.Uncles)

		block.Uncles = make([]*datastore.Key, 0, len(body.Uncles))
		for _, uncle := range body.Uncles {
			block.Uncles = append(block.Uncles, datastore.NameKey("blocks", uncle.Hash().Hex(), nil))
		}
	}

	if _, err := client.Put(ctx, key, &block); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write block header")
	}
}

func (c *Conn) writeBlockBodies(ctx context.Context, client *datastore.Client, hashes []common.Hash, bodies []*eth.BlockBody) {
	if len(hashes) != len(bodies) {
		c.logger.Error().Msg("Mismatch hashes and bodies length")
		return
	}

	for i, body := range bodies {
		hash := hashes[i].Hex()
		c.writeBlockBody(ctx, client, hash, body)
	}
}

func (c *Conn) writeTransactions(ctx context.Context, client *datastore.Client, txs []*types.Transaction) {
	hashes := make([]common.Hash, 0, len(txs))
	for _, tx := range txs {
		hashes = append(hashes, tx.Hash())

		go func(tx *types.Transaction) {
			key := datastore.NameKey("transactions", tx.Hash().Hex(), nil)

			var transaction *DatastoreTransaction
			if err := client.Get(ctx, key, transaction); err == nil {
				return
			}

			transaction = NewDatastoreTransaction(tx)

			if _, err := client.Put(ctx, key, transaction); err != nil {
				c.logger.Error().Err(err).Msg("Failed to write transaction")
			}
		}(tx)
	}

	c.writeEvents(ctx, client, "transaction_events", hashes, "transactions")
}

func (c *Conn) writeTransaction(ctx context.Context, client *datastore.Client, tx *types.Transaction, blockHash string) {
	txKey := datastore.NameKey("transactions", tx.Hash().Hex(), nil)

	var transaction *DatastoreTransaction
	if err := client.Get(ctx, txKey, transaction); err != nil {
		transaction = NewDatastoreTransaction(tx)
	}

	blockKey := datastore.NameKey("blocks", blockHash, nil)
	if !slices.Contains(transaction.BlockHashes, blockKey) {
		transaction.BlockHashes = append(transaction.BlockHashes, blockKey)
	}

	if _, err := client.Put(ctx, txKey, transaction); err != nil {
		c.logger.Error().Err(err).Msg("Failed to write transaction")
	}
}

func unseenTransactions(ctx context.Context, client *datastore.Client, hashes []common.Hash) []common.Hash {
	unseen := make([]common.Hash, 0, len(hashes))

	for _, hash := range hashes {
		key := datastore.NameKey("transactions", hash.Hex(), nil)

		var transaction *DatastoreTransaction
		if err := client.Get(ctx, key, transaction); err != nil {
			unseen = append(unseen, hash)
		}
	}

	return unseen
}
