package leveldbbench

import (
	"context"
	"crypto/sha1"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rs/zerolog/log"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	leveldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	//go:embed usage.md
	usage string
	// memory leak?
	knownKeys           map[string][]byte
	knownKeysMutex      sync.RWMutex
	randSrc             *rand.Rand
	randSrcMutex        sync.Mutex
	smallFillLimit      *uint64
	largeFillLimit      *uint64
	noWriteMerge        *bool
	syncWrites          *bool
	dontFillCache       *bool
	readStrict          *bool
	keySize             *uint64
	smallValueSize      *uint64
	largeValueSize      *uint64
	degreeOfParallelism *uint8
	readLimit           *uint64
	rawSizeDistribution *string
	sizeDistribution    *IODistribution
)

type (
	LoadTestOperation int
	TestResult        struct {
		StartTime    time.Time
		EndTime      time.Time
		TestDuration time.Duration
		Description  string
		OpCount      uint64
		Stats        *leveldb.DBStats
		OpRate       float64
	}
)

func NewTestResult(startTime, endTime time.Time, desc string, opCount uint64, db *leveldb.DB) *TestResult {
	tr := new(TestResult)
	s := new(leveldb.DBStats)
	err := db.Stats(s)
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve db stats")
	}
	tr.Stats = s
	tr.StartTime = startTime
	tr.EndTime = endTime
	tr.TestDuration = endTime.Sub(startTime)
	tr.Description = desc
	tr.OpCount = opCount
	tr.OpRate = float64(opCount) / tr.TestDuration.Seconds()

	log.Info().Dur("testDuration", tr.TestDuration).Str("desc", tr.Description).Msg("recorded result")
	log.Debug().Interface("result", tr).Msg("recorded result")
	return tr
}

var LevelDBBenchCmd = &cobra.Command{
	Use:   "leveldbbench [flags]",
	Short: "Perform a level db benchmark",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msg("Starting level db test")
		knownKeys = make(map[string][]byte, 0)
		db, err := leveldb.OpenFile("_benchmark_db", nil)
		if err != nil {
			return err
		}
		ctx := context.Background()
		wo := opt.WriteOptions{
			NoWriteMerge: *noWriteMerge,
			Sync:         *syncWrites,
		}
		ro := opt.ReadOptions{
			DontFillCache: *dontFillCache,
		}
		if *readStrict {
			ro.Strict = opt.StrictAll
		} else {
			ro.Strict = opt.DefaultStrict
		}
		var start time.Time
		trs := make([]*TestResult, 0)

		start = time.Now()
		writeData(ctx, db, &wo, *smallValueSize, 0, *smallFillLimit, true)
		trs = append(trs, NewTestResult(start, time.Now(), "small seq fill", *smallFillLimit, db))

		start = time.Now()
		writeData(ctx, db, &wo, *smallValueSize, 0, *smallFillLimit, true)
		trs = append(trs, NewTestResult(start, time.Now(), "small seq overwrite", *smallFillLimit, db))

		start = time.Now()
		writeData(ctx, db, &wo, *smallValueSize, 0, *smallFillLimit, false)
		trs = append(trs, NewTestResult(start, time.Now(), "small rand fill", *smallFillLimit, db))

		start = time.Now()
		writeData(ctx, db, &wo, *smallValueSize, 0, *smallFillLimit, false)
		trs = append(trs, NewTestResult(start, time.Now(), "small rand overwrite", *smallFillLimit, db))

		start = time.Now()
		writeData(ctx, db, &wo, *smallValueSize, 0, *smallFillLimit, false)
		trs = append(trs, NewTestResult(start, time.Now(), "small rand overwrite", *smallFillLimit, db))

		start = time.Now()
		writeData(ctx, db, &wo, *smallValueSize, 0, *smallFillLimit, false)
		trs = append(trs, NewTestResult(start, time.Now(), "small rand overwrite", *smallFillLimit, db))

		start = time.Now()
		readSeq(ctx, db, &wo, *readLimit)
		trs = append(trs, NewTestResult(start, time.Now(), "sequential read", *readLimit, db))

		start = time.Now()
		writeData(ctx, db, &wo, *largeValueSize, *smallFillLimit*2, *largeFillLimit, false)
		trs = append(trs, NewTestResult(start, time.Now(), "large rand fill", *largeFillLimit, db))

		start = time.Now()
		writeData(ctx, db, &wo, *largeValueSize, *smallFillLimit*2, *largeFillLimit, false)
		trs = append(trs, NewTestResult(start, time.Now(), "large rand overwrite", *largeFillLimit, db))

		start = time.Now()
		readSeq(ctx, db, &wo, *readLimit)
		trs = append(trs, NewTestResult(start, time.Now(), "sequential read", *readLimit, db))

		start = time.Now()
		readRandom(ctx, db, &ro, *readLimit)
		trs = append(trs, NewTestResult(start, time.Now(), "random read", *readLimit, db))

		start = time.Now()
		runFullCompact(ctx, db, &wo)
		trs = append(trs, NewTestResult(start, time.Now(), "compaction", 1, db))

		log.Info().Msg("Close DB")
		err = db.Close()
		if err != nil {
			log.Error().Err(err).Msg("error while closing db")
		}

		jsonResults, err := json.Marshal(trs)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonResults))
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		var err error
		sizeDistribution, err = parseRawSizeDistribution(*rawSizeDistribution)
		if err != nil {
			return err
		}
		return nil
	},
}

func runFullCompact(ctx context.Context, db *leveldb.DB, wo *opt.WriteOptions) {
	err := db.CompactRange(util.Range{Start: nil, Limit: nil})
	if err != nil {
		log.Fatal().Err(err).Msg("error compacting data")
	}
}
func writeData(ctx context.Context, db *leveldb.DB, wo *opt.WriteOptions, valueSize, startIndex, writeLimit uint64, sequential bool) {
	var i uint64 = startIndex
	var wg sync.WaitGroup
	pool := make(chan bool, *degreeOfParallelism)
	bar := getNewProgressBar(int64(writeLimit), fmt.Sprintf("Write: %d", valueSize))
	lim := writeLimit + startIndex
	for ; i < lim; i = i + 1 {
		pool <- true
		wg.Add(1)
		go func(i uint64) {
			_ = bar.Add(1)
			k, v := makeKV(i, valueSize, sequential)
			err := db.Put(k, v, wo)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to put value")
			}
			wg.Done()
			<-pool
		}(i)
	}
	wg.Wait()
	_ = bar.Finish()
}

func readSeq(ctx context.Context, db *leveldb.DB, wo *opt.WriteOptions, limit uint64) {
	pb := getNewProgressBar(int64(limit), "sequential reads")
	var rCount uint64 = 0
	pool := make(chan bool, *degreeOfParallelism)
	var wg sync.WaitGroup
benchLoop:
	for {
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			rCount += 1
			pb.Add(1)
			pool <- true
			wg.Add(1)
			go func(i iterator.Iterator) {
				_ = i.Key()
				_ = i.Value()
				wg.Done()
				<-pool
			}(iter)

			if rCount >= limit {
				iter.Release()
				break benchLoop
			}
		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			log.Fatal().Err(err).Msg("Error reading sequentially")
		}
	}
	wg.Wait()
	_ = pb.Finish()
}
func readRandom(ctx context.Context, db *leveldb.DB, ro *opt.ReadOptions, limit uint64) {
	pb := getNewProgressBar(int64(limit), "random reads")
	var rCount uint64 = 0
	pool := make(chan bool, *degreeOfParallelism)
	var wg sync.WaitGroup

benchLoop:
	for {
		for _, randKey := range knownKeys {
			pool <- true
			wg.Add(1)
			go func() {
				rCount += 1
				pb.Add(1)

				db.Get(randKey, ro)
				wg.Done()
				<-pool
			}()
			if rCount >= limit {
				break benchLoop
			}
		}
	}
	wg.Wait()
	_ = pb.Finish()
}

func getNewProgressBar(max int64, description string) *progressbar.ProgressBar {
	pb := progressbar.NewOptions64(max,
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetItsString("iop"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionThrottle(1*time.Second),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprintln(os.Stderr)
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionSetWidth(10),
		progressbar.OptionFullWidth(),
	)
	return pb
}

func makeKV(seed, valueSize uint64, sequential bool) ([]byte, []byte) {
	tmpKey := make([]byte, *keySize, *keySize)
	if sequential {
		// We're going to hack sequential by counting in reverse
		binary.BigEndian.PutUint64(tmpKey, math.MaxUint64-seed)
	} else {
		// For random (non-sequential) we'll just hash the number so it's still deterministic
		binary.LittleEndian.PutUint64(tmpKey, seed)
		hashedKey := sha1.Sum(tmpKey)
		tmpKey = hashedKey[0:*keySize]
	}
	knownKeysMutex.Lock()
	knownKeys[string(tmpKey)] = tmpKey
	knownKeysMutex.Unlock()

	log.Trace().Str("tmpKey", hex.EncodeToString(tmpKey)).Msg("Generated key")

	tmpValue := make([]byte, valueSize, valueSize)
	randSrcMutex.Lock()
	randSrc.Read(tmpValue)
	randSrcMutex.Unlock()
	return tmpKey, tmpValue
}

type (
	IORange struct {
		StartRange int
		EndRange   int
		Frequency  int
	}
	IODistribution struct {
		ranges         []IORange
		totalFrequency int
	}
)

func (i *IORange) Validate() error {
	if i.EndRange < i.StartRange {
		return fmt.Errorf("the end of the range %d  is less than the start of the range %d", i.EndRange, i.StartRange)
	}
	if i.EndRange <= 0 {
		return fmt.Errorf("the provided end range %d is less than 0", i.EndRange)
	}
	if i.StartRange < 0 {
		return fmt.Errorf("the provided start range %d is less than 0", i.StartRange)
	}
	if i.Frequency <= 0 {
		return fmt.Errorf("the relative frequency must be greater than 0, but got %d", i.Frequency)
	}
	return nil
}
func NewIODistribution(ranges []IORange) *IODistribution {
	iod := new(IODistribution)
	iod.ranges = ranges
	f := 0
	for _, v := range ranges {
		f += v.Frequency
	}
	iod.totalFrequency = f
	return iod
}
func (i *IODistribution) GetSizeSample() int {
	randSrcMutex.Lock()
	randFreq := randSrc.Intn(i.totalFrequency)
	randSrcMutex.Unlock()

	var selectedRange *IORange
	currentFreq := 0
	for _, v := range i.ranges {
		currentFreq += v.Frequency
		if currentFreq <= randFreq {
			selectedRange = &v
		}
	}
	if selectedRange == nil {
		log.Fatal().Int("randFreq", randFreq).Int("totalFreq", i.totalFrequency).Msg("Potential off by 1 error in random sample")
	}
	randRange := selectedRange.EndRange - selectedRange.StartRange
	randSrcMutex.Lock()
	randSize := randSrc.Intn(randRange)
	randSrcMutex.Unlock()
	return randSize + selectedRange.StartRange
}

func parseRawSizeDistribution(dist string) (*IODistribution, error) {
	buckets := strings.Split(dist, ",")
	if len(buckets) == 0 {
		return nil, fmt.Errorf("at least one size bucket must be provided")
	}
	ioDist := make([]IORange, 0)
	bucketRegEx := regexp.MustCompile(`^(\d*)-(\d*):(\d*)$`)
	for _, r := range buckets {
		matches := bucketRegEx.FindAllStringSubmatch(r, -1)
		if len(matches) != 1 {
			return nil, fmt.Errorf("the bucket %s did not match expected format of start-end:ratio", r)
		}
		if len(matches[0]) != 4 {
			return nil, fmt.Errorf("the bucket %s didn't match expected number of sub groups", r)
		}
		startRange, err := strconv.Atoi(matches[0][1])
		if err != nil {
			return nil, err
		}
		endRange, err := strconv.Atoi(matches[0][2])
		if err != nil {
			return nil, err
		}
		frequency, err := strconv.Atoi(matches[0][3])
		if err != nil {
			return nil, err
		}
		ioRange := new(IORange)
		ioRange.StartRange = startRange
		ioRange.EndRange = endRange
		ioRange.Frequency = frequency
		err = ioRange.Validate()
		if err != nil {
			return nil, err
		}
		ioDist = append(ioDist, *ioRange)
	}
	iod := NewIODistribution(ioDist)

	return iod, nil
}

func init() {
	flagSet := LevelDBBenchCmd.PersistentFlags()
	smallFillLimit = flagSet.Uint64("small-fill-limit", 1000000, "The number of small entries to write in the db")
	largeFillLimit = flagSet.Uint64("large-fill-limit", 2000, "The number of large entries to write in the db")
	readLimit = flagSet.Uint64("read-limit", 10000000, "the number of reads will attempt to complete in a given test")
	smallValueSize = flagSet.Uint64("small-value-size", 32, "the number of random bytes to store")
	largeValueSize = flagSet.Uint64("large-value-size", 102400, "the number of random bytes to store for large tests")
	dontFillCache = flag.Bool("dont-fill-read-cache", false, "if false, then random reads will be cached")
	readStrict = flag.Bool("read-strict", false, "if true the rand reads will be made in strict mode")
	keySize = flagSet.Uint64("key-size", 8, "The byte length of the keys that we'll use")
	degreeOfParallelism = flagSet.Uint8("degree-of-parallelism", 1, "The number of concurrent iops we'll perform")
	noWriteMerge = flagSet.Bool("no-merge-write", false, "allows disabling write merge")
	syncWrites = flagSet.Bool("sync-writes", false, "sync each write")
	rawSizeDistribution = flagSet.String("size-kb-distribution", "4-7:23,8-15:57,16-61:16,32-63:4", "the size distribution to use while testing")

	randSrc = rand.New(rand.NewSource(1))
}
