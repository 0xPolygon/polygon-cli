package borfsck

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

type fsckParamsType struct {
	dbPath                 string
	cacheSize              *int
	openFilesCacheCapacity *int
	startBlock             *uint64
	txLookup               *bool
}

type blockStats struct {
	transactions uint64
	receipts     uint64
	logs         uint64
	borReceipts  uint64
}

func (bs *blockStats) Add(stats *blockStats) *blockStats {
	bs.transactions += stats.transactions
	bs.receipts += stats.receipts
	bs.logs += stats.logs
	bs.borReceipts += stats.borReceipts
	return bs
}

var fsckParams = fsckParamsType{}

var BorFsckCmd = &cobra.Command{
	Use:   "bor-fsck /path/to/chaindata",
	Short: "bor-fsck /path/to/chaindata",
	Long:  "TODO",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("bor-fsck expects exactly one argument: a path to leveldb")
		}

		fsckParams.dbPath = strings.TrimSuffix(args[0], "/")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Str("dbPath", fsckParams.dbPath).Msg("Attempting to open database")
		db, err := openDB()
		if err != nil {
			return err
		}
		defer db.Close()
		dbVersion := rawdb.ReadDatabaseVersion(db)
		log.Info().Uint64("dbVersion", *dbVersion).Send()

		hb := rawdb.ReadHeadBlock(db)
		log.Info().Uint64("headBlockNumber", hb.NumberU64()).Send()
		if *fsckParams.startBlock != 0 {
			startBlockHash := rawdb.ReadCanonicalHash(db, *fsckParams.startBlock)
			hb = rawdb.ReadBlock(db, startBlockHash, *fsckParams.startBlock)
		}
		log.Info().
			Uint64("blockNumber", hb.NumberU64()).
			Str("blockHash", hb.Hash().String()).
			Str("stateRoot", hb.Root().String()).
			Str("receiptRoot", hb.ReceiptHash().String()).
			Str("transactionRoot", hb.TxHash().String()).
			Str("unclesHash", hb.UncleHash().String()).
			Msg("starting check")
		return checkBlocks(db, hb.Hash(), hb.NumberU64())
	},
}

func checkBlocks(db ethdb.Database, hash common.Hash, number uint64) error {
	totalStats := new(blockStats)
	var blockCount uint64 = 0
	for {
		if number == 0 {
			log.Info().Str("hash", hash.String()).Msg("reached genesis")
			break
		}
		// TODO concurrency?
		b := rawdb.ReadBlock(db, hash, number)
		bs, err := checkBlock(db, b)
		if err != nil {
			return err
		}

		blockCount += 1
		totalStats = totalStats.Add(bs)
		hash = b.ParentHash()
		number = number - 1
	}
	log.Info().
		Uint64("transactions", totalStats.transactions).
		Uint64("receipts", totalStats.receipts).
		Uint64("logs", totalStats.logs).
		Uint64("borReceipts", totalStats.borReceipts).
		Uint64("blockCount", blockCount).
		Msg("done")
	return nil
}

func checkBlock(db ethdb.Database, block *types.Block) (*blockStats, error) {
	bStats := new(blockStats)
	if block == nil {
		return bStats, fmt.Errorf("nil block")
	}
	log.Debug().Uint64("bn", block.NumberU64()).Str("hash", block.Hash().String()).Msg("checking block")
	err := block.SanityCheck()
	if err != nil {
		return bStats, err
	}
	txs := block.Transactions()
	txHashes := make(map[common.Hash]struct{}, 0)
	for idx, tx := range txs {
		log.Trace().Str("txHash", tx.Hash().String()).Msg("checking tx")

		bStats.transactions += 1
		txHashes[tx.Hash()] = struct{}{}
		lErr := log.Error().Str("blockHash", block.Hash().String()).Uint64("blockNumber", block.NumberU64()).Str("txHash", tx.Hash().String())

		if *fsckParams.txLookup {
			rtx, blockHash, blockNumber, txIndex := rawdb.ReadTransaction(db, tx.Hash())
			if rtx == nil {
				lErr.Msg("tx lookup failed")
			}
			if rtx != nil && rtx.Hash() != tx.Hash() {
				lErr.Str("rTxHash", rtx.Hash().String()).Msg("hash mismatch")
			}
			if txIndex != uint64(idx) {
				lErr.Int("idx", idx).Uint64("txIndex", txIndex).Msg("tx indices do not match")
			}
			if blockHash != block.Hash() {
				lErr.Str("innerHash", blockHash.String()).Str("outerHash", block.Hash().String()).Msg("block hash mismatch")
			}
			if blockNumber != block.NumberU64() {
				lErr.Uint64("innerNumber", blockNumber).Uint64("outerNumber", block.NumberU64()).Msg("blokc number mismatch")
			}
		}
	}

	blockLogs := rawdb.ReadLogs(db, block.Hash(), block.NumberU64())
	for receiptIndex, logs := range blockLogs {
		for logIndex, l := range logs {
			bStats.logs += 1
			lErr := log.Error().Str("blockHash", block.Hash().String()).Uint64("blockNumber", block.NumberU64()).Int("receiptIndex", receiptIndex).Int("logIndex", logIndex)
			if l == nil {
				lErr.Msg("nil log entry")
			}
			// These are the only consensus fields. I guess we Could sum the topics / data lengths but there's no easy way to validate these
			// Address
			// Topics
			// Data
		}
	}

	// Reading the derived fields will be complicated, so for now we'll avoid that
	receipts := rawdb.ReadRawReceipts(db, block.Hash(), block.NumberU64())
	var gasDiff uint64 = 0
	for rIdx, receipt := range receipts {
		bStats.receipts += 1
		// These fields are derived, so there is no point checking them
		// Type
		// TxHash
		// EffectiveGasPrice
		// BlobGasUsed
		// BlobGasPrice
		// BlockHash
		// BlockNumber
		// TransactionIndex
		// ContractAddress
		// GasUsed
		// Logs

		// I guess we could check these? but I'm not exactly sure how to sanity check them?
		// receipt.PostState
		// receipt.Status
		// receipt.CumulativeGasUsed
		// receipt.Bloom

		if receipt.CumulativeGasUsed-gasDiff < 21000 {
			log.Error().Uint64("blockNumber", block.NumberU64()).Int("index", rIdx).Uint64("gasUsed", receipt.CumulativeGasUsed).Msg("gas used is less than 21000")
		}
		gasDiff = receipt.CumulativeGasUsed
	}

	borReceipt := ReadRawBorReceipt(db, block.Hash(), block.NumberU64())
	if borReceipt != nil {
		bStats.borReceipts += 1
		// It seems like this is always 0? It seems like it should actually = the CumulativeGasUsed of the last receipt... but I guess not
		if borReceipt.CumulativeGasUsed != 0 {
			log.Error().Uint64("blockNumber", block.NumberU64()).Uint64("gasUsed", borReceipt.CumulativeGasUsed).Msg("gas used is not 0")
		}
	}

	return bStats, nil
}

func openDB() (ethdb.Database, error) {
	oo := rawdb.OpenOptions{
		Type:              "leveldb",
		Directory:         fsckParams.dbPath,
		AncientsDirectory: fsckParams.dbPath + "/ancient",
		Namespace:         "",
		Cache:             *fsckParams.cacheSize,
		Handles:           *fsckParams.openFilesCacheCapacity,
		ReadOnly:          true, // Want to read only for this use case
		Ephemeral:         false,
	}
	return rawdb.Open(oo)
}

func init() {
	flagSet := BorFsckCmd.PersistentFlags()
	fsckParams.cacheSize = flagSet.Int("cache-size", 512, "the number of megabytes to use as our internal cache size")
	fsckParams.openFilesCacheCapacity = flagSet.Int("handles", 4096, "number of files to be open simultaneously")
	fsckParams.startBlock = flagSet.Uint64("start-block", 0, "The block to start from")
	fsckParams.txLookup = flagSet.Bool("tx-lookup", false, "attempt a tx lookup for each hash")

}
