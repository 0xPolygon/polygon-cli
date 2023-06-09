package database

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
)

const (
	// Kinds are the datastore equivalent of tables.
	blocksKind            = "blocks"
	blockEventsKind       = "block_events"
	transactionsKind      = "transactions"
	transactionEventsKind = "transaction_events"
)

// datastoreWrapper wraps the datastore client and stores the sensorID so
// writing block and transaction events possible.
type datastoreWrapper struct {
	client                  *datastore.Client
	sensorID                string
	maxConcurrentWrites     int
	shouldWriteBlocks       bool
	shouldWriteTransactions bool
}

// DatastoreEvent can represent a peer sending the sensor a transaction hash or
// a block hash. In this implementation, the block and transactions are written
// to different tables by specifying a kind during key creation see writeEvents
// for more.
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

// DatastoreBlock represents a block stored in datastore.
type DatastoreBlock struct {
	*DatastoreHeader
	TotalDifficulty string
	Transactions    []*datastore.Key
	Uncles          []*datastore.Key
}

// DatastoreTransaction represents a transaction stored in datastore. Data is
// not indexed because there is a max sized for indexed byte slices, which Data
// will occasionally exceed.
type DatastoreTransaction struct {
	Data      []byte `datastore:",noindex"`
	From      string
	Gas       string
	GasFeeCap string
	GasPrice  string
	GasTipCap string
	Nonce     string
	To        string
	Value     string
	V, R, S   string
	Time      time.Time
	Type      int16
}

// NewDatastore connects to datastore and creates the client. This should
// only be called once unless trying to write to different databases.
func NewDatastore(ctx context.Context, projectID string, sensorID string, maxConcurrentWrites int, shouldWriteBlocks bool, shouldWriteTransactions bool) Database {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to Datastore")
		return nil
	}

	return &datastoreWrapper{
		client:                  client,
		sensorID:                sensorID,
		maxConcurrentWrites:     maxConcurrentWrites,
		shouldWriteBlocks:       shouldWriteBlocks,
		shouldWriteTransactions: shouldWriteTransactions,
	}
}

// WriteBlock writes the block and the block event to datastore.
func (d *datastoreWrapper) WriteBlock(ctx context.Context, peer *enode.Node, block *types.Block, td *big.Int) {
	d.writeEvent(peer, blockEventsKind, block.Hash(), blocksKind)

	key := datastore.NameKey(blocksKind, block.Hash().Hex(), nil)
	var dsBlock DatastoreBlock
	// Fetch the block. We don't check the error because if some of the fields
	// are nil we will just set them.
	_ = d.client.Get(ctx, key, &dsBlock)

	shouldWrite := false

	if dsBlock.DatastoreHeader == nil {
		shouldWrite = true
		dsBlock.DatastoreHeader = newDatastoreHeader(block.Header())
	}

	if len(dsBlock.TotalDifficulty) == 0 {
		shouldWrite = true
		dsBlock.TotalDifficulty = td.String()
	}

	if dsBlock.Transactions == nil && len(block.Transactions()) > 0 {
		shouldWrite = true
		if d.shouldWriteTransactions {
			d.writeTransactions(ctx, block.Transactions())
		}

		dsBlock.Transactions = make([]*datastore.Key, 0, len(block.Transactions()))
		for _, tx := range block.Transactions() {
			dsBlock.Transactions = append(dsBlock.Transactions, datastore.NameKey(transactionsKind, tx.Hash().Hex(), nil))
		}
	}

	if dsBlock.Uncles == nil && len(block.Uncles()) > 0 {
		shouldWrite = true
		dsBlock.Uncles = make([]*datastore.Key, 0, len(block.Uncles()))
		for _, uncle := range block.Uncles() {
			d.writeBlockHeader(ctx, uncle)
			dsBlock.Uncles = append(dsBlock.Uncles, datastore.NameKey(blocksKind, uncle.Hash().Hex(), nil))
		}
	}

	if !shouldWrite {
		return
	}

	if _, err := d.client.Put(ctx, key, &dsBlock); err != nil {
		log.Error().Err(err).Msg("Failed to write new block")
	}
}

// WriteBlockHeaders will write the block headers to datastore. It will not
// write block events because headers will only be sent to the sensor when
// requested. The block events will be written when the hash is received
// instead.
func (d *datastoreWrapper) WriteBlockHeaders(ctx context.Context, headers []*types.Header) {
	for _, header := range headers {
		d.writeBlockHeader(ctx, header)
	}
}

// WriteBlockHeaders will write the block bodies to datastore. It will not
// write block events because bodies will only be sent to the sensor when
// requested. The block events will be written when the hash is received
// instead. It will write the uncles and transactions to datastore if they
// don't already exist.
func (d *datastoreWrapper) WriteBlockBody(ctx context.Context, body *eth.BlockBody, hash common.Hash) {
	key := datastore.NameKey(blocksKind, hash.Hex(), nil)
	var block DatastoreBlock

	if err := d.client.Get(ctx, key, &block); err != nil {
		log.Debug().Err(err).Str("hash", hash.Hex()).Msg("Failed to fetch block when writing block body")
	}

	if block.Transactions == nil && len(body.Transactions) > 0 {
		if d.shouldWriteTransactions {
			d.writeTransactions(ctx, body.Transactions)
		}

		block.Transactions = make([]*datastore.Key, 0, len(body.Transactions))
		for _, tx := range body.Transactions {
			block.Transactions = append(block.Transactions, datastore.NameKey(transactionsKind, tx.Hash().Hex(), nil))
		}
	}

	if block.Uncles == nil && len(body.Uncles) > 0 {
		block.Uncles = make([]*datastore.Key, 0, len(body.Uncles))
		for _, uncle := range body.Uncles {
			d.writeBlockHeader(ctx, uncle)
			block.Uncles = append(block.Uncles, datastore.NameKey(blocksKind, uncle.Hash().Hex(), nil))
		}
	}

	if _, err := d.client.Put(ctx, key, &block); err != nil {
		log.Error().Err(err).Msg("Failed to write block header")
	}
}

// WriteBlockHashes will write the block events to datastore.
func (d *datastoreWrapper) WriteBlockHashes(ctx context.Context, peer *enode.Node, hashes []common.Hash) {
	d.writeEvents(ctx, peer, blockEventsKind, hashes, blocksKind)
}

// WriteTransactions will write the transactions and transaction events to datastore.
func (d *datastoreWrapper) WriteTransactions(ctx context.Context, peer *enode.Node, txs []*types.Transaction) {
	hashes := d.writeTransactions(ctx, txs)
	d.writeEvents(ctx, peer, transactionEventsKind, hashes, transactionsKind)
}

func (d *datastoreWrapper) MaxConcurrentWrites() int {
	return d.maxConcurrentWrites
}

func (d *datastoreWrapper) ShouldWriteBlocks() bool {
	return d.shouldWriteBlocks
}

func (d *datastoreWrapper) ShouldWriteTransactions() bool {
	return d.shouldWriteTransactions
}

func (d *datastoreWrapper) HasParentBlock(ctx context.Context, hash common.Hash) bool {
	key := datastore.NameKey(blocksKind, hash.Hex(), nil)
	var block DatastoreBlock
	err := d.client.Get(ctx, key, &block)

	return err == nil && block.DatastoreHeader != nil
}

// newDatastoreHeader creates a DatastoreHeader from a types.Header. Some
// values are converted into strings to prevent a loss of precision.
func newDatastoreHeader(header *types.Header) *DatastoreHeader {
	return &DatastoreHeader{
		ParentHash:  datastore.NameKey(blocksKind, header.ParentHash.Hex(), nil),
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

// newDatastoreTransaction creates a DatastoreTransaction from a types.Transaction. Some
// values are converted into strings to prevent a loss of precision.
func newDatastoreTransaction(tx *types.Transaction) *DatastoreTransaction {
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

// writeEvent writes either a block or transaction event to datastore depending
// on the provided eventKind and hashKind.
func (d *datastoreWrapper) writeEvent(peer *enode.Node, eventKind string, hash common.Hash, hashKind string) {
	key := datastore.IncompleteKey(eventKind, nil)
	event := DatastoreEvent{
		SensorId: d.sensorID,
		PeerId:   peer.URLv4(),
		Hash:     datastore.NameKey(hashKind, hash.Hex(), nil),
		Time:     time.Now(),
	}
	if _, err := d.client.Put(context.Background(), key, &event); err != nil {
		log.Error().Err(err).Msgf("Failed to write to %v", eventKind)
	}
}

// writeEvents writes either block or transaction events to datastore depending
// on the provided eventKind and hashKind. This is similar to writeEvent but
// batches the request.
func (d *datastoreWrapper) writeEvents(ctx context.Context, peer *enode.Node, eventKind string, hashes []common.Hash, hashKind string) {
	keys := make([]*datastore.Key, 0, len(hashes))
	events := make([]*DatastoreEvent, 0, len(hashes))
	now := time.Now()

	for _, hash := range hashes {
		keys = append(keys, datastore.IncompleteKey(eventKind, nil))

		event := DatastoreEvent{
			SensorId: d.sensorID,
			PeerId:   peer.URLv4(),
			Hash:     datastore.NameKey(hashKind, hash.Hex(), nil),
			Time:     now,
		}
		events = append(events, &event)
	}

	if _, err := d.client.PutMulti(ctx, keys, events); err != nil {
		log.Error().Err(err).Msgf("Failed to write to %v", eventKind)
	}
}

// writeBlockHeader will write the block header to datastore if it doesn't
// exist.
func (d *datastoreWrapper) writeBlockHeader(ctx context.Context, header *types.Header) {
	key := datastore.NameKey(blocksKind, header.Hash().Hex(), nil)
	var block DatastoreBlock

	if err := d.client.Get(ctx, key, &block); err == nil && block.DatastoreHeader != nil {
		return
	}

	block.DatastoreHeader = newDatastoreHeader(header)

	if _, err := d.client.Put(ctx, key, &block); err != nil {
		log.Error().Err(err).Msg("Failed to write block header")
	}
}

// writeTransactions will write the transactions to datastore and return the
// transaction hashes.
func (d *datastoreWrapper) writeTransactions(ctx context.Context, txs []*types.Transaction) []common.Hash {
	hashes := make([]common.Hash, 0, len(txs))
	keys := make([]*datastore.Key, 0, len(txs))
	transactions := make([]*DatastoreTransaction, 0, len(txs))

	for _, tx := range txs {
		hashes = append(hashes, tx.Hash())
		keys = append(keys, datastore.NameKey(transactionsKind, tx.Hash().Hex(), nil))
		transactions = append(transactions, newDatastoreTransaction(tx))
	}

	if _, err := d.client.PutMulti(ctx, keys, transactions); err != nil {
		log.Error().Err(err).Msg("Failed to write transactions")
	}

	return hashes
}
