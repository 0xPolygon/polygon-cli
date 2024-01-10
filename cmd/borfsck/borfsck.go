package borfsck

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type fsckParamsType struct {
	dbPath                 string
	cacheSize              *int
	openFilesCacheCapacity *int
}

var fsckParams = fsckParamsType{}

var BorFsckCmd = &cobra.Command{
	Use:   "bor-fsck /path/to/leveldb",
	Short: "bor-fsck /path/to/leveldb",
	Long:  "TODO",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("bor-fsck expects exactly one argument: a path to leveldb")
		}

		fsckParams.dbPath = args[0]
		db, err := openLevelDB()
		if err != nil {
			return err
		}
		_ = db
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func openLevelDB() (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(fsckParams.dbPath, &opt.Options{
		Filter:                 filter.NewBloomFilter(10),
		DisableSeeksCompaction: true,
		OpenFilesCacheCapacity: *fsckParams.openFilesCacheCapacity,
		BlockCacheCapacity:     *fsckParams.cacheSize / 2 * opt.MiB,
		// This tool should not be doing writes
		ReadOnly: true,
	})
	if err != nil {
		return nil, err
	}
	return db, nil

}

func init() {
	flagSet := BorFsckCmd.PersistentFlags()
	fsckParams.cacheSize = flagSet.Int("cache-size", 512, "the number of megabytes to use as our internal cache size")
	fsckParams.openFilesCacheCapacity = flagSet.Int("handles", 500, "defines the capacity of the open files caching. Use -1 for zero, this has same effect as specifying NoCacher to OpenFilesCacher.")

}
