package dbbench

import (
	"context"
	"crypto/sha512"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var (
	//go:embed usage.md
	usage string

	randSrc                *rand.Rand
	randSrcMutex           sync.Mutex
	writeLimit             *uint64
	noWriteMerge           *bool
	syncWrites             *bool
	dontFillCache          *bool
	readStrict             *bool
	keySize                *uint64
	degreeOfParallelism    *uint8
	readLimit              *uint64
	rawSizeDistribution    *string
	sizeDistribution       *IODistribution
	overwriteCount         *uint64
	sequentialReads        *bool
	sequentialWrites       *bool
	nilReadOptions         *bool
	cacheSize              *int
	openFilesCacheCapacity *int
	writeZero              *bool
	readOnly               *bool
	dbPath                 *string
	fullScan               *bool
	dbMode                 *string
)

const (
	// This data was obtained by running a full scan on bor level db to get a sense how the key values are distributed
	// | Bucket | Min Size  | Max        | Count         |
	// |--------+-----------+------------+---------------|
	// |      0 | 0         | 1          | 2,347,864     |
	// |      1 | 2         | 3          | 804,394,856   |
	// |      2 | 4         | 7          | 541,267,689   |
	// |      3 | 8         | 15         | 738,828,593   |
	// |      4 | 16        | 31         | 261,122,372   |
	// |      5 | 32        | 63         | 1,063,470,933 |
	// |      6 | 64        | 127        | 3,584,745,195 |
	// |      7 | 128       | 255        | 1,605,760,137 |
	// |      8 | 256       | 511        | 316,074,206   |
	// |      9 | 512       | 1,023      | 312,887,514   |
	// |     10 | 1,024     | 2,047      | 328,894,149   |
	// |     11 | 2,048     | 4,095      | 141,180       |
	// |     12 | 4,096     | 8,191      | 92,789        |
	// |     13 | 8,192     | 16,383     | 256,060       |
	// |     14 | 16,384    | 32,767     | 261,806       |
	// |     15 | 32,768    | 65,535     | 191,032       |
	// |     16 | 65,536    | 131,071    | 99,715        |
	// |     17 | 131,072   | 262,143    | 73,782        |
	// |     18 | 262,144   | 524,287    | 17,552        |
	// |     19 | 524,288   | 1,048,575  | 717           |
	// |     20 | 1,048,576 | 2,097,151  | 995           |
	// |     21 | 2,097,152 | 4,194,303  | 1             |
	// |     22 | 4,194,304 | 8,388,607  | 0             |
	// |     23 | 8,388,608 | 16,777,215 | 1             |
	borDistribution = "0-1:2347864,2-3:804394856,4-7:541267689,8-15:738828593,16-31:261122372,32-63:1063470933,64-127:3584745195,128-255:1605760137,256-511:316074206,512-1023:312887514,1024-2047:328894149,2048-4095:141180,4096-8191:92789,8192-16383:256060,16384-32767:261806,32768-65535:191032,65536-131071:99715,131072-262143:73782,262144-524287:17552,524288-1048575:717,1048576-2097151:995,2097152-4194303:1,8388608-16777215:1"
)

type (
	LoadTestOperation int
	TestResult        struct {
		StartTime    time.Time
		EndTime      time.Time
		TestDuration time.Duration
		Description  string
		OpCount      uint64
		OpRate       float64
		ValueDist    []uint64
	}
	RandomKeySeeker struct {
		db            KeyValueDB
		iterator      iterator.Iterator
		iteratorMutex sync.Mutex
		firstKey      []byte
	}
	IORange struct {
		StartRange int
		EndRange   int
		Frequency  int
	}
	IODistribution struct {
		ranges         []IORange
		totalFrequency int
	}
	// KeyValueDB directly exposes the necessary methods of leveldb.DB that we need to run the test so that they can be
	// implemented by other KV stores
	KeyValueDB interface {
		Close() error
		Compact() error
		NewIterator() iterator.Iterator
		Get([]byte) ([]byte, error)
		Put([]byte, []byte) error
	}
)

func NewTestResult(startTime, endTime time.Time, desc string, opCount uint64) *TestResult {
	tr := new(TestResult)
	tr.StartTime = startTime
	tr.EndTime = endTime
	tr.TestDuration = endTime.Sub(startTime)
	tr.Description = desc
	tr.OpCount = opCount
	tr.OpRate = float64(opCount) / tr.TestDuration.Seconds()

	log.Info().Dur("testDuration", tr.TestDuration).Str("desc", tr.Description).Msg("Recorded result")
	log.Debug().Interface("result", tr).Msg("Recorded result")
	return tr
}

var DBBenchCmd = &cobra.Command{
	Use:   "dbbench [flags]",
	Short: "Perform a level/pebble db benchmark",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msg("Starting db test")
		var kvdb KeyValueDB
		var err error
		switch *dbMode {
		case "leveldb":
			kvdb, err = NewWrappedLevelDB()
			if err != nil {
				return err
			}
		case "pebbledb":
			kvdb, err = NewWrappedPebbleDB()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("the mode %s is not recognized", *dbMode)
		}

		ctx := context.Background()

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

		if *fullScan {
			start = time.Now()
			opCount, valueDist := runFullScan(ctx, kvdb)
			tr := NewTestResult(start, time.Now(), "full scan", opCount)
			tr.ValueDist = valueDist
			trs = append(trs, tr)
			return printSummary(trs)
		}

		// in no write mode, we assume the database as already been populated in a previous run or we're using some other database
		if !*readOnly {
			start = time.Now()
			writeData(ctx, kvdb, 0, *writeLimit, *sequentialWrites)
			trs = append(trs, NewTestResult(start, time.Now(), fmt.Sprintf("initial %s write", sequentialWritesDesc), *writeLimit))

			for i := 0; i < int(*overwriteCount); i += 1 {
				start = time.Now()
				writeData(ctx, kvdb, 0, *writeLimit, *sequentialWrites)
				trs = append(trs, NewTestResult(start, time.Now(), fmt.Sprintf("%s overwrite %d", sequentialWritesDesc, i), *writeLimit))
			}

			start = time.Now()
			runFullCompact(ctx, kvdb)
			trs = append(trs, NewTestResult(start, time.Now(), "compaction", 1))
		}

		if *sequentialReads {
			start = time.Now()
			readSeq(ctx, kvdb, *readLimit)
			trs = append(trs, NewTestResult(start, time.Now(), fmt.Sprintf("%s read", sequentialReadsDesc), *readLimit))
		} else {
			start = time.Now()
			readRandom(ctx, kvdb, *readLimit)
			trs = append(trs, NewTestResult(start, time.Now(), fmt.Sprintf("%s read", sequentialWritesDesc), *readLimit))
		}

		log.Info().Msg("Close DB")
		err = kvdb.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error while closing db")
		}

		return printSummary(trs)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		var err error
		sizeDistribution, err = parseRawSizeDistribution(*rawSizeDistribution)
		if err != nil {
			return err
		}
		if *keySize > 64 {
			return fmt.Errorf(" max supported key size is 64 bytes. %d is too big", *keySize)
		}
		return nil
	},
}

func printSummary(trs []*TestResult) error {
	jsonResults, err := json.Marshal(trs)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonResults))
	return nil
}

func runFullCompact(ctx context.Context, db KeyValueDB) {
	err := db.Compact()
	if err != nil {
		log.Fatal().Err(err).Msg("Error compacting data")
	}
}
func runFullScan(ctx context.Context, db KeyValueDB) (uint64, []uint64) {
	pool := make(chan bool, *degreeOfParallelism)
	var wg sync.WaitGroup
	// 32 should be safe here. That would correspond to a single value that's 4.2 GB
	buckets := make([]uint64, 32)
	var bucketsMutex sync.Mutex
	iter := db.NewIterator()
	var opCount uint64 = 0
	for iter.Next() {
		pool <- true
		wg.Add(1)
		go func(i iterator.Iterator) {
			opCount += 1
			k := i.Key()
			v := i.Value()

			bucket := bits.Len(uint(len(v)))
			bucketsMutex.Lock()
			buckets[bucket] += 1
			bucketsMutex.Unlock()

			if bucket >= 22 {
				// 9:19PM INF encountered giant value currentKey=536e617073686f744a6f75726e616c
				log.Info().Str("currentKey", hex.EncodeToString(k)).Int("bytes", len(v)).Msg("Encountered giant value")
			}

			if opCount%1000000 == 0 {
				log.Debug().Uint64("opCount", opCount).Str("currentKey", hex.EncodeToString(k)).Msg("Continuing full scan")
			}
			wg.Done()
			<-pool
		}(iter)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Fatal().Err(err).Msg("Error running full scan")
	}

	wg.Wait()

	for k, v := range buckets {
		if v == 0 {
			continue
		}
		start := math.Exp2(float64(k))
		end := math.Exp2(float64(k+1)) - 1
		if k == 0 {
			start = 0
		}
		log.Debug().
			Int("bucket", k).
			Float64("start", start).
			Float64("end", end).
			Uint64("count", v).Msg("Buckets")
	}
	return opCount, buckets
}
func writeData(ctx context.Context, db KeyValueDB, startIndex, writeLimit uint64, sequential bool) {
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
			err := db.Put(k, v)
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

func readSeq(ctx context.Context, db KeyValueDB, limit uint64) {
	pb := getNewProgressBar(int64(limit), "sequential reads")
	var rCount uint64 = 0
	pool := make(chan bool, *degreeOfParallelism)
	var wg sync.WaitGroup
benchLoop:
	for {
		iter := db.NewIterator()
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
func readRandom(ctx context.Context, db KeyValueDB, limit uint64) {
	pb := getNewProgressBar(int64(limit), "random reads")
	var rCount uint64 = 0
	pool := make(chan bool, *degreeOfParallelism)
	var wg sync.WaitGroup
	rks := NewRandomKeySeeker(db)
	defer rks.iterator.Release()

	var rCountLock sync.Mutex
	var keyLock sync.Mutex
benchLoop:
	for {
		for {
			pool <- true
			wg.Add(1)
			go func() {
				rCountLock.Lock()
				rCount += 1
				rCountLock.Unlock()
				_ = pb.Add(1)

				// It's not entirely obvious WHY this is needed, but without it, there are issues with the way that
				// pebble db manages it's iterators and internal state. Level db works fine though.
				keyLock.Lock()
				tmpKey := rks.Key()
				_, err := db.Get(tmpKey)
				keyLock.Unlock()
				if err != nil {
					log.Error().Str("key", hex.EncodeToString(tmpKey)).Err(err).Msg("db random read error")
				}
				wg.Done()
				<-pool
			}()
			rCountLock.Lock()
			if rCount >= limit {
				rCountLock.Unlock()
				break benchLoop
			}
			rCountLock.Unlock()

		}
	}
	wg.Wait()
	_ = pb.Finish()
}

func NewRandomKeySeeker(db KeyValueDB) *RandomKeySeeker {
	rks := new(RandomKeySeeker)
	rks.db = db
	rks.iterator = db.NewIterator()
	rks.firstKey = rks.iterator.Key()
	return rks
}
func (r *RandomKeySeeker) Key() []byte {
	seekKey := make([]byte, 8)
	randSrcMutex.Lock()
	randSrc.Read(seekKey)
	randSrcMutex.Unlock()

	log.Trace().Str("seekKey", hex.EncodeToString(seekKey)).Msg("Searching for key")

	r.iteratorMutex.Lock()

	// first try to just get a random key
	exists := r.iterator.Seek(seekKey)

	// if that key doesn't exist exactly advance to the next key
	if !exists {
		exists = r.iterator.Next()
	}

	// if there is no next key, to back to the beginning
	if !exists {
		exists = r.iterator.First()
	}

	// if there is no first key try advancing again
	if !exists {
		exists = r.iterator.Next()
	}

	// if after trying to all these ways to find a valid key... something must be very wrong
	if !exists {
		log.Fatal().Msg("Unable to select random key!?")
	}
	if err := r.iterator.Error(); err != nil {
		log.Error().Err(err).Msg("Issue getting random key")
	}
	resultKey := r.iterator.Key()
	r.iteratorMutex.Unlock()
	log.Trace().Str("seekKey", hex.EncodeToString(seekKey)).Str("resultKey", hex.EncodeToString(resultKey)).Msg("Found random key")
	return resultKey
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
	binary.LittleEndian.PutUint64(tmpKey, seed)
	hashedKey := sha512.Sum512(tmpKey)
	tmpKey = hashedKey[0:*keySize]
	if sequential {
		// binary.BigEndian.PutUint64(tmpKey, seed)
		binary.BigEndian.PutUint64(tmpKey, seed)
	}

	log.Trace().Str("tmpKey", hex.EncodeToString(tmpKey)).Uint64("valueSize", valueSize).Uint64("seed", seed).Msg("Generated key")

	tmpValue := make([]byte, valueSize)
	if !*writeZero {
		// Assuming we're not in zero mode, we'll fill the data with random data
		randSrcMutex.Lock()
		randSrc.Read(tmpValue)
		randSrcMutex.Unlock()
	}
	return tmpKey, tmpValue
}

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
	randSize := randSrc.Intn(randRange + 1)
	randSrcMutex.Unlock()
	return uint64(randSize + selectedRange.StartRange)
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
	flagSet := DBBenchCmd.PersistentFlags()
	writeLimit = flagSet.Uint64("write-limit", 1000000, "The number of entries to write in the db")
	readLimit = flagSet.Uint64("read-limit", 10000000, "the number of reads will attempt to complete in a given test")
	overwriteCount = flagSet.Uint64("overwrite-count", 5, "the number of times to overwrite the data")
	sequentialReads = flagSet.Bool("sequential-reads", false, "if true we'll perform reads sequentially")
	sequentialWrites = flagSet.Bool("sequential-writes", false, "if true we'll perform writes in somewhat sequential manner")
	keySize = flagSet.Uint64("key-size", 32, "The byte length of the keys that we'll use")
	degreeOfParallelism = flagSet.Uint8("degree-of-parallelism", 2, "The number of concurrent goroutines we'll use")
	rawSizeDistribution = flagSet.String("size-distribution", borDistribution, "the size distribution to use while testing")
	nilReadOptions = flagSet.Bool("nil-read-opts", false, "if true we'll use nil read opt (this is what geth/bor does)")
	dontFillCache = flagSet.Bool("dont-fill-read-cache", false, "if false, then random reads will be cached")
	readStrict = flagSet.Bool("read-strict", false, "if true the rand reads will be made in strict mode")
	noWriteMerge = flagSet.Bool("no-merge-write", false, "allows disabling write merge")
	syncWrites = flagSet.Bool("sync-writes", false, "sync each write")
	// https://github.com/maticnetwork/bor/blob/eedeaed1fb17d73dd46d8999644d5035e176e22a/eth/backend.go#L141
	// https://github.com/maticnetwork/bor/blob/eedeaed1fb17d73dd46d8999644d5035e176e22a/eth/ethconfig/config.go#L86C2-L86C15
	cacheSize = flagSet.Int("cache-size", 512, "the number of megabytes to use as our internal cache size")
	openFilesCacheCapacity = flagSet.Int("handles", 500, "defines the capacity of the open files caching. Use -1 for zero, this has same effect as specifying NoCacher to OpenFilesCacher.")
	writeZero = flagSet.Bool("write-zero", false, "if true, we'll write 0s rather than random data")
	readOnly = flagSet.Bool("read-only", false, "if true, we'll skip all the write operations and open the DB in read only mode")
	dbPath = flagSet.String("db-path", "_benchmark_db", "the path of the database that we'll use for testing")
	fullScan = flagSet.Bool("full-scan-mode", false, "if true, the application will scan the full database as fast as possible and print a summary")
	dbMode = flagSet.String("db-mode", "leveldb", "The mode to use: leveldb or pebbledb")

	randSrc = rand.New(rand.NewSource(1))
}
