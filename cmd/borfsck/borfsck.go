package borfsck

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/rawdb"
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
		return nil
	},
}

//func openLevelDB() (*leveldb.DB, error) {
//	db, err := ethleveldb.New(fsckParams.dbPath, &opt.Options{
//		Filter:                 filter.NewBloomFilter(10),
//		DisableSeeksCompaction: true,
//		OpenFilesCacheCapacity: *fsckParams.openFilesCacheCapacity,
//		BlockCacheCapacity:     *fsckParams.cacheSize / 2 * opt.MiB,
//		// This tool should not be doing writes
//		ReadOnly: true,
//	})
//	if err != nil {
//		return nil, err
//	}
//	return db, nil
//
//}

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
