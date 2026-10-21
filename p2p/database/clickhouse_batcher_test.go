package database

import (
	"context"
	"sync"
	"testing"
	"time"
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

	start := time.Now()
	for i := range 3 {
		b.add(i)
	}

	timer := time.NewTimer(5 * interval)
	defer timer.Stop()
	select {
	case rows := <-flushed:
		if len(rows) != 3 {
			t.Fatalf("flushed %d rows, want 3", len(rows))
		}
		if elapsed := time.Since(start); elapsed > 3*interval {
			t.Fatalf("partial batch flushed after %v, want about %v", elapsed, interval)
		}
	case <-timer.C:
		t.Fatalf("partial batch not flushed within %v", 5*interval)
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
