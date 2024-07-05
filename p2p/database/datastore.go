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
	"google.golang.org/api/iterator"
)

const (
	// Kinds are the datastore equivalent of tables.
	BlocksKind            = "blocks"
	BlockEventsKind       = "block_events"
	TransactionsKind      = "transactions"
	TransactionEventsKind = "transaction_events"
	MaxAttempts           = 10
)

// Datastore wraps the datastore client, stores the sensorID, and other
// information needed when writing blocks and transactions.
type Datastore struct {
	client                       *datastore.Client
	sensorID                     string
	maxConcurrency               int
	shouldWriteBlocks            bool
	shouldWriteBlockEvents       bool
	shouldWriteTransactions      bool
	shouldWriteTransactionEvents bool
	jobs                         chan struct{}
	ttl                          time.Duration
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
	TTL      time.Time
}

// DatastoreHeader stores the data in manner that can be easily written without
// loss of precision.
type DatastoreHeader struct {
	ParentHash    *datastore.Key
	UncleHash     string
	Coinbase      string
	Root          string
	TxHash        string
	ReceiptHash   string
	Bloom         []byte `datastore:",noindex"`
	Difficulty    string
	Number        string
	GasLimit      string
	GasUsed       string
	Time          time.Time
	Extra         []byte `datastore:",noindex"`
	MixDigest     string
	Nonce         string
	BaseFee       string
	TimeFirstSeen time.Time
	TTL           time.Time
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
	Data          []byte `datastore:",noindex"`
	From          string
	Gas           string
	GasFeeCap     string
	GasPrice      string
	GasTipCap     string
	Nonce         string
	To            string
	Value         string
	V, R, S       string
	Time          time.Time
	TimeFirstSeen time.Time
	TTL           time.Time
	Type          int16
}

// DatastoreOptions is used when creating a NewDatastore.
type DatastoreOptions struct {
	ProjectID                    string
	DatabaseID                   string
	SensorID                     string
	MaxConcurrency               int
	ShouldWriteBlocks            bool
	ShouldWriteBlockEvents       bool
	ShouldWriteTransactions      bool
	ShouldWriteTransactionEvents bool
	TTL                          time.Duration
}

// NewDatastore connects to datastore and creates the client. This should
// only be called once unless trying to write to different databases.
func NewDatastore(ctx context.Context, opts DatastoreOptions) Database {
	client, err := datastore.NewClientWithDatabase(ctx, opts.ProjectID, opts.DatabaseID)
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to Datastore")
	}

	return &Datastore{
		client:                       client,
		sensorID:                     opts.SensorID,
		maxConcurrency:               opts.MaxConcurrency,
		shouldWriteBlocks:            opts.ShouldWriteBlocks,
		shouldWriteBlockEvents:       opts.ShouldWriteBlockEvents,
		shouldWriteTransactions:      opts.ShouldWriteTransactions,
		shouldWriteTransactionEvents: opts.ShouldWriteTransactionEvents,
		jobs:                         make(chan struct{}, opts.MaxConcurrency),
		ttl:                          opts.TTL,
	}
}

// WriteBlock writes the block and the block event to datastore.
func (d *Datastore) WriteBlock(ctx context.Context, peer *enode.Node, block *types.Block, td *big.Int) {
	if d.client == nil {
		return
	}

	if d.ShouldWriteBlockEvents() {
		d.jobs <- struct{}{}
		go func() {
			d.writeEvent(peer, BlockEventsKind, block.Hash(), BlocksKind)
			<-d.jobs
		}()
	}

	if d.ShouldWriteBlocks() {
		d.jobs <- struct{}{}
		go func() {
			d.writeBlock(ctx, block, td)
			<-d.jobs
		}()
	}
}

// WriteBlockHeaders will write the block headers to datastore. It will not
// write block events because headers will only be sent to the sensor when
// requested. The block events will be written when the hash is received
// instead.
func (d *Datastore) WriteBlockHeaders(ctx context.Context, headers []*types.Header) {
	if d.client == nil || !d.ShouldWriteBlocks() {
		return
	}

	for _, h := range headers {
		d.jobs <- struct{}{}
		go func(header *types.Header) {
			d.writeBlockHeader(ctx, header)
			<-d.jobs
		}(h)
	}
}

// WriteBlockHeaders will write the block bodies to datastore. It will not
// write block events because bodies will only be sent to the sensor when
// requested. The block events will be written when the hash is received
// instead. It will write the uncles and transactions to datastore if they
// don't already exist.
func (d *Datastore) WriteBlockBody(ctx context.Context, body *eth.BlockBody, hash common.Hash) {
	if d.client == nil || !d.ShouldWriteBlocks() {
		return
	}

	d.jobs <- struct{}{}
	go func() {
		d.writeBlockBody(ctx, body, hash)
		<-d.jobs
	}()
}

// WriteBlockHashes will write the block events to datastore.
func (d *Datastore) WriteBlockHashes(ctx context.Context, peer *enode.Node, hashes []common.Hash) {
	if d.client == nil || !d.ShouldWriteBlockEvents() || len(hashes) == 0 {
		return
	}

	d.jobs <- struct{}{}
	go func() {
		d.writeEvents(ctx, peer, BlockEventsKind, hashes, BlocksKind)
		<-d.jobs
	}()
}

// WriteTransactions will write the transactions and transaction events to datastore.
func (d *Datastore) WriteTransactions(ctx context.Context, peer *enode.Node, txs []*types.Transaction) {
	if d.client == nil {
		return
	}

	if d.ShouldWriteTransactions() {
		d.jobs <- struct{}{}
		go func() {
			d.writeTransactions(ctx, txs)
			<-d.jobs
		}()
	}

	if d.ShouldWriteTransactionEvents() {
		hashes := make([]common.Hash, 0, len(txs))
		for _, tx := range txs {
			hashes = append(hashes, tx.Hash())
		}

		d.jobs <- struct{}{}
		go func() {
			d.writeEvents(ctx, peer, TransactionEventsKind, hashes, TransactionsKind)
			<-d.jobs
		}()
	}
}

func (d *Datastore) MaxConcurrentWrites() int {
	return d.maxConcurrency
}

func (d *Datastore) ShouldWriteBlocks() bool {
	return d.shouldWriteBlocks
}

func (d *Datastore) ShouldWriteBlockEvents() bool {
	return d.shouldWriteBlockEvents
}

func (d *Datastore) ShouldWriteTransactions() bool {
	return d.shouldWriteTransactions
}

func (d *Datastore) ShouldWriteTransactionEvents() bool {
	return d.shouldWriteTransactionEvents
}

func (d *Datastore) HasBlock(ctx context.Context, hash common.Hash) bool {
	if d.client == nil {
		return true
	}

	key := datastore.NameKey(BlocksKind, hash.Hex(), nil)
	var block DatastoreBlock
	err := d.client.Get(ctx, key, &block)

	return err == nil && block.DatastoreHeader != nil
}

// newDatastoreHeader creates a DatastoreHeader from a types.Header. Some
// values are converted into strings to prevent a loss of precision.
func (d *Datastore) newDatastoreHeader(header *types.Header) *DatastoreHeader {
	now := time.Now()

	return &DatastoreHeader{
		ParentHash:    datastore.NameKey(BlocksKind, header.ParentHash.Hex(), nil),
		UncleHash:     header.UncleHash.Hex(),
		Coinbase:      header.Coinbase.Hex(),
		Root:          header.Root.Hex(),
		TxHash:        header.TxHash.Hex(),
		ReceiptHash:   header.ReceiptHash.Hex(),
		Bloom:         header.Bloom.Bytes(),
		Difficulty:    header.Difficulty.String(),
		Number:        header.Number.String(),
		GasLimit:      fmt.Sprint(header.GasLimit),
		GasUsed:       fmt.Sprint(header.GasUsed),
		Time:          time.Unix(int64(header.Time), 0),
		Extra:         header.Extra,
		MixDigest:     header.MixDigest.String(),
		Nonce:         fmt.Sprint(header.Nonce.Uint64()),
		BaseFee:       header.BaseFee.String(),
		TimeFirstSeen: now,
		TTL:           now.Add(d.ttl),
	}
}

// newDatastoreTransaction creates a DatastoreTransaction from a types.Transaction. Some
// values are converted into strings to prevent a loss of precision.
func (d *Datastore) newDatastoreTransaction(tx *types.Transaction) *DatastoreTransaction {
	v, r, s := tx.RawSignatureValues()
	var from, to string

	address, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err == nil {
		from = address.String()
	}

	if tx.To() != nil {
		to = tx.To().Hex()
	}

	now := time.Now()

	return &DatastoreTransaction{
		Data:          tx.Data(),
		From:          from,
		Gas:           fmt.Sprint(tx.Gas()),
		GasFeeCap:     tx.GasFeeCap().String(),
		GasPrice:      tx.GasPrice().String(),
		GasTipCap:     tx.GasTipCap().String(),
		Nonce:         fmt.Sprint(tx.Nonce()),
		To:            to,
		Value:         tx.Value().String(),
		V:             v.String(),
		R:             r.String(),
		S:             s.String(),
		Time:          tx.Time(),
		TimeFirstSeen: now,
		TTL:           now.Add(d.ttl),
		Type:          int16(tx.Type()),
	}
}

func (d *Datastore) writeBlock(ctx context.Context, block *types.Block, td *big.Int) {
	key := datastore.NameKey(BlocksKind, block.Hash().Hex(), nil)

	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var dsBlock DatastoreBlock
		// Fetch the block. We don't check the error because if some of the fields
		// are nil we will just set them.
		_ = tx.Get(key, &dsBlock)

		shouldWrite := false

		if dsBlock.DatastoreHeader == nil {
			shouldWrite = true
			dsBlock.DatastoreHeader = d.newDatastoreHeader(block.Header())
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
				dsBlock.Transactions = append(dsBlock.Transactions, datastore.NameKey(TransactionsKind, tx.Hash().Hex(), nil))
			}
		}

		if dsBlock.Uncles == nil && len(block.Uncles()) > 0 {
			shouldWrite = true
			dsBlock.Uncles = make([]*datastore.Key, 0, len(block.Uncles()))
			for _, uncle := range block.Uncles() {
				d.writeBlockHeader(ctx, uncle)
				dsBlock.Uncles = append(dsBlock.Uncles, datastore.NameKey(BlocksKind, uncle.Hash().Hex(), nil))
			}
		}

		if shouldWrite {
			_, err := tx.Put(key, &dsBlock)
			return err
		}

		return nil
	}, datastore.MaxAttempts(MaxAttempts))

	if err != nil {
		log.Error().Err(err).Str("hash", block.Hash().Hex()).Msg("Failed to write new block")
	}
}

// writeEvent writes either a block or transaction event to datastore depending
// on the provided eventKind and hashKind.
func (d *Datastore) writeEvent(peer *enode.Node, eventKind string, hash common.Hash, hashKind string) {
	key := datastore.IncompleteKey(eventKind, nil)
	now := time.Now()

	event := DatastoreEvent{
		SensorId: d.sensorID,
		PeerId:   peer.URLv4(),
		Hash:     datastore.NameKey(hashKind, hash.Hex(), nil),
		Time:     now,
		TTL:      now.Add(d.ttl),
	}
	if _, err := d.client.Put(context.Background(), key, &event); err != nil {
		log.Error().Err(err).Msgf("Failed to write to %v", eventKind)
	}
}

// writeEvents writes either block or transaction events to datastore depending
// on the provided eventKind and hashKind. This is similar to writeEvent but
// batches the request.
func (d *Datastore) writeEvents(ctx context.Context, peer *enode.Node, eventKind string, hashes []common.Hash, hashKind string) {
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
			TTL:      now.Add(d.ttl),
		}
		events = append(events, &event)
	}

	if _, err := d.client.PutMulti(ctx, keys, events); err != nil {
		log.Error().Err(err).Msgf("Failed to write to %v", eventKind)
	}
}

// writeBlockHeader will write the block header to datastore if it doesn't
// exist.
func (d *Datastore) writeBlockHeader(ctx context.Context, header *types.Header) {
	key := datastore.NameKey(BlocksKind, header.Hash().Hex(), nil)

	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var block DatastoreBlock
		if err := tx.Get(key, &block); err == nil && block.DatastoreHeader != nil {
			return nil
		}

		block.DatastoreHeader = d.newDatastoreHeader(header)
		_, err := tx.Put(key, &block)
		return err
	}, datastore.MaxAttempts(MaxAttempts))

	if err != nil {
		log.Error().Err(err).Str("hash", header.Hash().Hex()).Msg("Failed to write block header")
	}
}

func (d *Datastore) writeBlockBody(ctx context.Context, body *eth.BlockBody, hash common.Hash) {
	key := datastore.NameKey(BlocksKind, hash.Hex(), nil)

	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var block DatastoreBlock
		if err := tx.Get(key, &block); err != nil {
			log.Debug().Err(err).Str("hash", hash.Hex()).Msg("Failed to fetch block when writing block body")
		}

		shouldWrite := false

		if block.Transactions == nil && len(body.Transactions) > 0 {
			shouldWrite = true
			if d.shouldWriteTransactions {
				d.writeTransactions(ctx, body.Transactions)
			}

			block.Transactions = make([]*datastore.Key, 0, len(body.Transactions))
			for _, tx := range body.Transactions {
				block.Transactions = append(block.Transactions, datastore.NameKey(TransactionsKind, tx.Hash().Hex(), nil))
			}
		}

		if block.Uncles == nil && len(body.Uncles) > 0 {
			shouldWrite = true
			block.Uncles = make([]*datastore.Key, 0, len(body.Uncles))
			for _, uncle := range body.Uncles {
				d.writeBlockHeader(ctx, uncle)
				block.Uncles = append(block.Uncles, datastore.NameKey(BlocksKind, uncle.Hash().Hex(), nil))
			}
		}

		if shouldWrite {
			_, err := tx.Put(key, &block)
			return err
		}

		return nil
	}, datastore.MaxAttempts(MaxAttempts))

	if err != nil {
		log.Error().Err(err).Str("hash", hash.Hex()).Msg("Failed to write block body")
	}
}

// writeTransactions will write the transactions to datastore and return the
// transaction hashes.
func (d *Datastore) writeTransactions(ctx context.Context, txs []*types.Transaction) {
	keys := make([]*datastore.Key, 0, len(txs))
	transactions := make([]*DatastoreTransaction, 0, len(txs))

	for _, tx := range txs {
		keys = append(keys, datastore.NameKey(TransactionsKind, tx.Hash().Hex(), nil))
		transactions = append(transactions, d.newDatastoreTransaction(tx))
	}

	if _, err := d.client.PutMulti(ctx, keys, transactions); err != nil {
		log.Error().Err(err).Msg("Failed to write transactions")
	}
}

func (d *Datastore) NodeList(ctx context.Context, limit int) ([]string, error) {
	query := datastore.NewQuery(BlockEventsKind).Order("-Time")
	iter := d.client.Run(ctx, query)

	enodes := make(map[string]struct{})
	for len(enodes) < limit {
		var event DatastoreEvent
		_, err := iter.Next(&event)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("Failed to get next block event")
			continue
		}

		enodes[event.PeerId] = struct{}{}
	}

	log.Info().Int("enodes", len(enodes)).Send()

	nodelist := []string{}
	for enode := range enodes {
		nodelist = append(nodelist, enode)
	}

	return nodelist, nil
}
