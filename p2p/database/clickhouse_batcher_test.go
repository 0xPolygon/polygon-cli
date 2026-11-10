package database

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
)

// These exercise rowBatcher on its own, so unlike the integration tests above
// they run without a ClickHouse.

// newCollectingBatcher returns a batcher whose flushes are sent on the returned
// channel, and a stop func that cancels it and waits for the final drain.
func newCollectingBatcher(t *testing.T, maxRows int, interval time.Duration) (*rowBatcher[int], <-chan []int, func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	flushed := make(chan []int, 16)
	b := newRowBatcher(ctx, &wg, "test", maxRows, interval, func(rows []int) error {
		flushed <- append([]int(nil), rows...)
		return nil
	})
	return b, flushed, func() {
		cancel()
		wg.Wait()
	}
}

func TestRowBatcherFlushesPartialBatchOnInterval(t *testing.T) {
	const interval = 200 * time.Millisecond
	b, flushed, stop := newCollectingBatcher(t, 100, interval)
	defer stop()

	for i := range 3 {
		b.add(i)
	}

	// The batch is far below maxRows, so only the ticker can flush it. The bounds
	// are loose on purpose: a busy CI runner under -race can delay the tick.
	early := time.NewTimer(interval / 2)
	defer early.Stop()
	select {
	case <-flushed:
		t.Fatal("partial batch flushed before the interval")
	case <-early.C:
	}

	timer := time.NewTimer(10 * interval)
	defer timer.Stop()
	select {
	case rows := <-flushed:
		if len(rows) != 3 {
			t.Fatalf("flushed %d rows, want 3", len(rows))
		}
	case <-timer.C:
		t.Fatalf("partial batch not flushed within %v", 10*interval)
	}
}

func TestRowBatcherFlushesFullBatchBeforeInterval(t *testing.T) {
	// An interval far longer than the test: only reaching maxRows can flush.
	b, flushed, stop := newCollectingBatcher(t, 5, time.Hour)
	defer stop()

	for i := range 5 {
		b.add(i)
	}

	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	select {
	case rows := <-flushed:
		if len(rows) != 5 {
			t.Fatalf("flushed %d rows, want 5", len(rows))
		}
	case <-timer.C:
		t.Fatal("full batch not flushed before the interval")
	}
}

func TestRowBatcherDrainsOnStop(t *testing.T) {
	b, flushed, stop := newCollectingBatcher(t, 100, time.Hour)
	for i := range 3 {
		b.add(i)
	}
	stop()

	select {
	case rows := <-flushed:
		if len(rows) != 3 {
			t.Fatalf("drained %d rows, want 3", len(rows))
		}
	default:
		t.Fatal("buffered rows not flushed on stop")
	}
}

func TestNewClickHouseFlushInterval(t *testing.T) {
	tests := []struct {
		name string
		in   time.Duration
		want time.Duration
	}{
		{"zero uses default", 0, DefaultClickHouseFlushInterval},
		{"negative uses default", -time.Second, DefaultClickHouseFlushInterval},
		{"explicit", 10 * time.Second, 10 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// A malformed DSN never connects, so no batchers start and nothing is
			// dialled; the interval is resolved before that.
			db := NewClickHouse(context.Background(), ClickHouseOptions{DSN: "://bad", FlushInterval: tt.in})
			defer func() { _ = db.Close() }()
			c, ok := db.(*ClickHouse)
			if !ok {
				t.Fatalf("NewClickHouse returned %T", db)
			}
			if c.flushInterval != tt.want {
				t.Fatalf("flushInterval = %v, want %v", c.flushInterval, tt.want)
			}
		})
	}
}

// panicConn satisfies driver.Conn but panics if used, so a test can assert that a
// code path never reaches the database.
type panicConn struct{ driver.Conn }

func TestHasBlockSeesPendingBlock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	defer func() {
		cancel()
		wg.Wait()
	}()

	pending, err := lru.New(chPendingBlocks)
	if err != nil {
		t.Fatal(err)
	}
	c := &ClickHouse{
		conn:          panicConn{},
		pendingBlocks: pending,
		// Never flushes during the test, like a row still waiting for its tick.
		blocks: newRowBatcher(ctx, &wg, "blocks", 100, time.Hour, func([]chBlock) error { return nil }),
	}

	h := &types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(1)}
	c.addBlock(h)
	if !c.HasBlock(ctx, h.Hash()) {
		t.Fatal("HasBlock = false for an enqueued block, want true without querying")
	}
}

// TestRowBatcherDropsWhenFull checks that add never blocks: a row that does not
// fit in the buffer is counted as dropped.
func TestRowBatcherDropsWhenFull(t *testing.T) {
	b := &rowBatcher[int]{in: make(chan int, 1)}
	b.add(1)
	b.add(2)
	if got := b.dropped.Load(); got != 1 {
		t.Fatalf("dropped = %d, want 1", got)
	}
}

// TestRowBatcherSurvivesFlushErrorsAndReportsDrops checks that a failed flush
// does not stop the batcher, and that the dropped count is reported and reset
// on the next tick.
func TestRowBatcherSurvivesFlushErrorsAndReportsDrops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	defer func() { cancel(); wg.Wait() }()
	var mu sync.Mutex
	flushes := 0
	b := newRowBatcher(ctx, &wg, "test", 100, 10*time.Millisecond, func([]int) error {
		mu.Lock()
		defer mu.Unlock()
		flushes++
		return errors.New("insert failed")
	})
	b.dropped.Store(3)
	b.add(1)

	deadline := time.Now().Add(5 * time.Second)
	for {
		mu.Lock()
		n := flushes
		mu.Unlock()
		if n > 0 && b.dropped.Load() == 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("after 5s: %d flushes, dropped = %d; want a flush and the count reset", n, b.dropped.Load())
		}
		time.Sleep(5 * time.Millisecond)
	}
	b.add(2)
	deadline = time.Now().Add(5 * time.Second)
	for {
		mu.Lock()
		n := flushes
		mu.Unlock()
		if n > 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("batcher stopped flushing after a failed flush")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// TestRowBatcherDrainFlushesInFullBatches checks the shutdown drain: rows
// already queued when the context is cancelled are all flushed, in batches of
// at most maxRows.
func TestRowBatcherDrainFlushesInFullBatches(t *testing.T) {
	const maxRows = 2
	var batches [][]int
	b := &rowBatcher[int]{
		name:     "test",
		in:       make(chan int, 10),
		maxRows:  maxRows,
		interval: time.Hour,
		flush: func(rows []int) error {
			batches = append(batches, append([]int(nil), rows...))
			return nil
		},
	}
	for i := range 5 {
		b.in <- i
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	b.loop(ctx)

	total := 0
	for _, batch := range batches {
		if len(batch) > maxRows {
			t.Fatalf("flushed a batch of %d rows, more than maxRows %d", len(batch), maxRows)
		}
		total += len(batch)
	}
	if total != 5 {
		t.Fatalf("flushed %d rows, want all 5", total)
	}
}

func TestCopyBig(t *testing.T) {
	if copyBig(nil) != nil {
		t.Fatal("copyBig(nil) != nil")
	}
	v := big.NewInt(5)
	c := copyBig(v)
	v.SetInt64(6)
	if c.Int64() != 5 {
		t.Fatal("copyBig returned an alias, not a copy")
	}
}
