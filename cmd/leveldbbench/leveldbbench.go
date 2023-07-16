package leveldbbench

import (
	"context"
	"crypto/sha1"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rs/zerolog/log"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	leveldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math/rand"
	"os"
	"regexp"
	"sort"
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
	writeLimit          *uint64
	noWriteMerge        *bool
	syncWrites          *bool
	dontFillCache       *bool
	readStrict          *bool
	keySize             *uint64
	degreeOfParallelism *uint8
	readLimit           *uint64
	rawSizeDistribution *string
	sizeDistribution    *IODistribution
	overwriteCount      *uint64
	sequentialReads     *bool
	sequentialWrites    *bool
	nilReadOptions      *bool
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
		db, err := leveldb.OpenFile("_benchmark_db", &opt.Options{
			Filter:                 filter.NewBloomFilter(10),
			DisableSeeksCompaction: true,
		})
		if err != nil {
			return err
		}
		ctx := context.Background()
		wo := opt.WriteOptions{
			NoWriteMerge: *noWriteMerge,
			Sync:         *syncWrites,
		}
		ro := &opt.ReadOptions{
			DontFillCache: *dontFillCache,
		}
		if *readStrict {
			ro.Strict = opt.StrictAll
		} else {
			ro.Strict = opt.DefaultStrict
		}
		if *nilReadOptions {
			ro = nil
		}
		var start time.Time
		trs := make([]*TestResult, 0)

		sequentialWritesDesc := "random"
		if *sequentialWrites {
			sequentialWritesDesc = "sequential"
		}
		sequentialReadsDesc := "random"
		if *sequentialReads {
			sequentialReadsDesc = "sequential"
		}

		start = time.Now()
		writeData(ctx, db, &wo, 0, *writeLimit, *sequentialWrites)
		trs = append(trs, NewTestResult(start, time.Now(), fmt.Sprintf("initial %s write", sequentialWritesDesc), *writeLimit, db))

		for i := 0; i < int(*overwriteCount); i += 1 {
			start = time.Now()
			writeData(ctx, db, &wo, 0, *writeLimit, *sequentialWrites)
			trs = append(trs, NewTestResult(start, time.Now(), fmt.Sprintf("%s overwrite %d", sequentialWritesDesc, i), *writeLimit, db))
		}

		if *sequentialReads {
			start = time.Now()
			readSeq(ctx, db, &wo, *readLimit)
			trs = append(trs, NewTestResult(start, time.Now(), fmt.Sprintf("%s read", sequentialReadsDesc), *readLimit, db))
		} else {
			start = time.Now()
			readRandom(ctx, db, ro, *readLimit)
			trs = append(trs, NewTestResult(start, time.Now(), fmt.Sprintf("%s read", sequentialWritesDesc), *readLimit, db))
		}

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
func writeData(ctx context.Context, db *leveldb.DB, wo *opt.WriteOptions, startIndex, writeLimit uint64, sequential bool) {
	var i uint64 = startIndex
	var wg sync.WaitGroup
	pool := make(chan bool, *degreeOfParallelism)
	bar := getNewProgressBar(int64(writeLimit), "Writing data")
	lim := writeLimit + startIndex
	for ; i < lim; i = i + 1 {
		pool <- true
		wg.Add(1)
		go func(i uint64) {
			_ = bar.Add(1)
			k, v := makeKV(i, sizeDistribution.GetSizeSample(), sequential)
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
			_ = pb.Add(1)
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
				_ = pb.Add(1)

				_, err := db.Get(randKey, ro)
				if err != nil {
					log.Error().Err(err).Msg("level db random read error")
				}
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
	tmpKey := make([]byte, *keySize)
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

	tmpValue := make([]byte, valueSize)
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
func NewIODistribution(ranges []IORange) (*IODistribution, error) {
	iod := new(IODistribution)
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].StartRange < ranges[j].StartRange
	})

	for i := 0; i < len(ranges)-1; i++ {
		if ranges[i].EndRange >= ranges[i+1].StartRange {
			return nil, fmt.Errorf("overlap found between ranges: %v and %v", ranges[i], ranges[i+1])
		}
	}

	iod.ranges = ranges
	f := 0
	for _, v := range ranges {
		f += v.Frequency
	}
	iod.totalFrequency = f
	return iod, nil
}

// GetSizeSample will return an IO size in accordance with the probability distribution
func (i *IODistribution) GetSizeSample() uint64 {
	randSrcMutex.Lock()
	randFreq := randSrc.Intn(i.totalFrequency)
	randSrcMutex.Unlock()

	log.Trace().Int("randFreq", randFreq).Int("totalFreq", i.totalFrequency).Msg("Getting Size Sample")
	var selectedRange *IORange
	currentFreq := 0
	for k, v := range i.ranges {
		currentFreq += v.Frequency
		if randFreq <= currentFreq {
			selectedRange = &i.ranges[k]
			break
		}
	}
	if selectedRange == nil {
		log.Fatal().Int("randFreq", randFreq).Int("totalFreq", i.totalFrequency).Msg("Potential off by 1 error in random sample")
		return 0 // lint
	}
	randRange := selectedRange.EndRange - selectedRange.StartRange
	randSrcMutex.Lock()
	randSize := randSrc.Intn(randRange)
	randSrcMutex.Unlock()
	return uint64(randSize+selectedRange.StartRange) * 1024
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
	return NewIODistribution(ioDist)
}

func init() {
	flagSet := LevelDBBenchCmd.PersistentFlags()
	writeLimit = flagSet.Uint64("write-limit", 1000000, "The number of entries to write in the db")
	readLimit = flagSet.Uint64("read-limit", 10000000, "the number of reads will attempt to complete in a given test")
	overwriteCount = flagSet.Uint64("overwrite-count", 5, "the number of times to overwrite the data")
	sequentialReads = flagSet.Bool("sequential-reads", false, "if true we'll perform reads sequentially")
	sequentialWrites = flagSet.Bool("sequential-writes", false, "if true we'll perform writes in somewhat sequential manner")
	keySize = flagSet.Uint64("key-size", 8, "The byte length of the keys that we'll use")
	degreeOfParallelism = flagSet.Uint8("degree-of-parallelism", 1, "The number of concurrent iops we'll perform")
	rawSizeDistribution = flagSet.String("size-kb-distribution", "4-7:23089,8-15:70350,16-31:11790,32-63:1193,64-127:204,128-255:271,256-511:1381", "the size distribution to use while testing")
	nilReadOptions = flagSet.Bool("nil-read-opts", false, "if true we'll use nil read opt (this is what geth/bor does)")
	dontFillCache = flagSet.Bool("dont-fill-read-cache", false, "if false, then random reads will be cached")
	readStrict = flagSet.Bool("read-strict", false, "if true the rand reads will be made in strict mode")
	noWriteMerge = flagSet.Bool("no-merge-write", false, "allows disabling write merge")
	syncWrites = flagSet.Bool("sync-writes", false, "sync each write")

	randSrc = rand.New(rand.NewSource(1))
}
