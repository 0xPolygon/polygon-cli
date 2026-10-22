package database

import (
	"context"
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
