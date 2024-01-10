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
		return checkBlocks(db, hb.Hash(), hb.NumberU64())
	},
}

func checkBlocks(db ethdb.Database, hash common.Hash, number uint64) error {
	for {
		if number == 0 {
			log.Info().Str("hash", hash.String()).Msg("reached genesis")
			break
		}
		// TODO concurrency?
		b := rawdb.ReadBlock(db, hash, number)
		err := checkBlock(db, b)
		if err != nil {
			return err
		}
		hash = b.ParentHash()
		number = number - 1
	}
	return nil
}

func checkBlock(db ethdb.Database, block *types.Block) error {
	if block == nil {
		return fmt.Errorf("nil block")
	}
	log.Debug().Uint64("bn", block.NumberU64()).Str("hash", block.Hash().String()).Msg("checking block")
	err := block.SanityCheck()
	if err != nil {
		return err
	}
	txs := block.Transactions()
	for idx, tx := range txs {
		log.Trace().Str("txHash", tx.Hash().String()).Msg("checking tx")
		rtx, blockHash, blockNumber, txIndex := rawdb.ReadTransaction(db, tx.Hash())
		lErr := log.Error().Str("txHash", tx.Hash().String())
		if rtx == nil {
			lErr.Msg("tx lookup failed")
			continue
		}
		if rtx.Hash() != tx.Hash() {
			lErr.Str("rTxHash", rtx.Hash().String()).Msg("hash mismatch")
			continue
		}
		if txIndex != uint64(idx) {
			lErr.Int("idx", idx).Uint64("txIndex", txIndex).Msg("tx indices do not match")
			continue
		}
		if blockHash != block.Hash() {
			lErr.Str("innerHash", blockHash.String()).Str("outerHash", block.Hash().String()).Msg("block hash mismatch")
			continue
		}
		if blockNumber != block.NumberU64() {
			lErr.Uint64("innerNumber", blockNumber).Uint64("outerNumber", block.NumberU64()).Msg("blokc number mismatch")
		}
	}

	return nil
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

}
