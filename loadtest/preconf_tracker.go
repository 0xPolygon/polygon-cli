package loadtest

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type PreconfTracker struct {
	client *ethclient.Client

	// metrics
	preconfSuccess     atomic.Uint64
	preconfFail        atomic.Uint64
	totalTasks         atomic.Uint64
	bothFailedCount    atomic.Uint64
	uneffectivePreconf atomic.Uint64
	falsePositiveCount atomic.Uint64
	confidence         atomic.Uint64

	mu               sync.Mutex
	preconfDurations []time.Duration
	blockDiffs       []uint64
}

func NewPreconfTracker(client *ethclient.Client) *PreconfTracker {
	return &PreconfTracker{
		client:           client,
		preconfDurations: make([]time.Duration, 0, 1024),
	}
}

func (pt *PreconfTracker) Track(txHash common.Hash) {
	currentBlock, err := pt.client.BlockNumber(context.Background())
	if err != nil {
		return
	}

	// wait for preconf
	var wg sync.WaitGroup
	var preconfStatus bool
	var preconfError error
	var preconfDuration time.Duration
	wg.Add(1)
	go func() {
		defer wg.Done()

		preconfStartTime := time.Now()
		defer func() {
			preconfDuration = time.Since(preconfStartTime)
		}()

		preconfStatus, preconfError = util.WaitPreconf(context.Background(), pt.client, txHash, time.Minute)
	}()

	// wait for receipt
	var receipt *types.Receipt
	var receiptError error
	var receiptDuration time.Duration
	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(100 * time.Millisecond)

		receiptTime := time.Now()
		defer func() {
			receiptDuration = time.Since(receiptTime)
		}()

		receipt, receiptError = util.WaitReceiptWithTimeout(context.Background(), pt.client, txHash, time.Minute)
	}()

	wg.Wait()

	pt.totalTasks.Add(1)
	if preconfStatus {
		pt.preconfSuccess.Add(1)
		pt.mu.Lock()
		pt.preconfDurations = append(pt.preconfDurations, preconfDuration)
		pt.mu.Unlock()
	} else {
		pt.preconfFail.Add(1)
	}

	switch {
	case preconfError != nil && receiptError != nil:
		// Both failed: no tx inclusion in txpool or block
		pt.bothFailedCount.Add(1)

	case preconfError == nil && receiptError != nil:
		// False positive: preconf said tx is included but never got executed
		pt.falsePositiveCount.Add(1)

	case preconfError != nil && receiptError == nil:
		// Receipt arrived but preconf failed: preconf wasn't effective
		pt.uneffectivePreconf.Add(1)

	case preconfError == nil && receiptError == nil:
		// Both succeeded
		if preconfDuration > receiptDuration {
			// Receipt arrived before preconf: preconf wasn't effective
			pt.uneffectivePreconf.Add(1)
		}
		// Track block diff for confidence
		blockDiff := receipt.BlockNumber.Uint64() - currentBlock
		if blockDiff < 10 {
			pt.mu.Lock()
			pt.blockDiffs = append(pt.blockDiffs, blockDiff)
			pt.mu.Unlock()
			pt.confidence.Add(1)
		}
	}
}

func (pt *PreconfTracker) Stats() {
	log.Info().Uint64("total_tasks", pt.totalTasks.Load()).
		Uint64("preconf_success", pt.preconfSuccess.Load()).
		Uint64("preconf_fail", pt.preconfFail.Load()).
		Uint64("both_failed", pt.bothFailedCount.Load()).
		Uint64("uneffective_preconf", pt.uneffectivePreconf.Load()).
		Uint64("false_positives", pt.falsePositiveCount.Load()).
		Uint64("confidence", pt.confidence.Load()).
		Msg("Preconf Tracker Stats")

	// Copy data under lock, then write files without holding lock
	pt.mu.Lock()
	durations := make([]time.Duration, len(pt.preconfDurations))
	copy(durations, pt.preconfDurations)
	blockDiffs := make([]uint64, len(pt.blockDiffs))
	copy(blockDiffs, pt.blockDiffs)
	pt.mu.Unlock()

	timestamp := time.Now().Format("20060102_150405")
	path := "preconf_durations_" + timestamp + ".csv"
	if err := dumpDurationsCSV(path, durations); err != nil {
		log.Error().Err(err).Msg("Error dumping preconf durations")
	} else {
		log.Info().Str("path", path).Msg("Dumped preconf durations into file")
	}

	path = "preconf_block_diffs_" + timestamp + ".csv"
	if err := dumpBlockDiff(path, blockDiffs); err != nil {
		log.Error().Err(err).Msg("Error dumping preconf block diffs")
	} else {
		log.Info().Str("path", path).Msg("Dumped preconf block diffs into file")
	}
}

func dumpBlockDiff(path string, diffs []uint64) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// header
	if err := w.Write([]string{"idx", "diff"}); err != nil {
		return err
	}

	for i, d := range diffs {
		row := []string{
			strconv.Itoa(i),
			strconv.FormatUint(d, 10),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

func dumpDurationsCSV(path string, durations []time.Duration) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// header
	if err := w.Write([]string{"idx", "duration_ns", "duration_ms"}); err != nil {
		return err
	}

	for i, d := range durations {
		row := []string{
			strconv.Itoa(i),
			strconv.FormatInt(d.Nanoseconds(), 10),
			strconv.FormatInt(d.Milliseconds(), 10),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}
