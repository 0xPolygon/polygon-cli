package leveldbbench

import (
	"context"
	_ "embed"
	"encoding/binary"
	"github.com/rs/zerolog/log"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	leveldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"math/rand"
	"sync"
	"time"
)

var (
	//go:embed usage.md
	usage string

	randSrc             *rand.Rand
	randSrcMutex        sync.Mutex
	fillLimit           *uint64
	noWriteMerge        *bool
	syncWrites          *bool
	keySize             *uint64
	smallValueSize      *uint64
	largeValueSize      *uint64
	degreeOfParallelism *uint8
)
var LevelDBBenchCmd = &cobra.Command{
	Use:   "leveldbbench [flags]",
	Short: "Perform a level db benchmark",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msg("hi")
		db, err := leveldb.OpenFile("_benchmark_db", nil)
		if err != nil {
			return err
		}
		ctx := context.Background()
		wo := opt.WriteOptions{
			NoWriteMerge: *noWriteMerge,
			Sync:         *syncWrites,
		}
		performFill(ctx, db, &wo)
		defer db.Close()
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func performFill(ctx context.Context, db *leveldb.DB, wo *opt.WriteOptions) {
	var i uint64 = 0
	var wg sync.WaitGroup
	pool := make(chan bool, *degreeOfParallelism)
	bar := getNewProgessBar(int64(*fillLimit), "performing initial fill")
	for i = 0; i < *fillLimit; i = i + 1 {
		pool <- true
		wg.Add(1)
		go func() {
			k, v := makeKV(i, *smallValueSize)
			bar.Add(1)
			err := db.Put(k, v, wo)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to put value")
			}
			wg.Done()
			<-pool
		}()
	}
	wg.Wait()
}

func getNewProgessBar(max int64, description string) *progressbar.ProgressBar {
	pb := progressbar.NewOptions64(max,
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetItsString("iop"),
		progressbar.OptionSetRenderBlankState(false),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionThrottle(1*time.Second),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	//return progressbar.Default(max, description)
	return pb
}

func makeKV(seed, valueSize uint64) ([]byte, []byte) {
	tmpKey := make([]byte, *keySize, *keySize)
	binary.PutUvarint(tmpKey, seed)
	tmpValue := make([]byte, valueSize, valueSize)
	randSrcMutex.Lock()
	randSrc.Read(tmpValue)
	randSrcMutex.Unlock()
	return tmpKey, tmpValue
}

func init() {
	flagSet := LevelDBBenchCmd.PersistentFlags()
	fillLimit = flagSet.Uint64("fill-limit", 1000000, "The number of unique entries to set in the db")
	smallValueSize = flagSet.Uint64("small-value-size", 32, "the number of random bytes to store")
	largeValueSize = flagSet.Uint64("large-value-size", 102400, "the number of random bytes to store for large tests")
	noWriteMerge = flagSet.Bool("no-merge-write", false, "allows disabling write merge")
	syncWrites = flagSet.Bool("sync-writes", false, "sync each write")
	keySize = flagSet.Uint64("key-size", 16, "The byte length of the keys that we'll use")
	degreeOfParallelism = flagSet.Uint8("degree-of-parallelism", 1, "The number of concurrent iops we'll perform")

	randSrc = rand.New(rand.NewSource(1))
}
