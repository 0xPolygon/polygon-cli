package mode

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewChunkFeederRejectsBadArgs(t *testing.T) {
	if _, err := NewChunkFeeder(t.Context(), bytes.NewReader(nil), 0, 1); err == nil {
		t.Fatal("expected error for zero chunk size")
	}
	if _, err := NewChunkFeeder(t.Context(), bytes.NewReader(nil), 4, -1); err == nil {
		t.Fatal("expected error for negative buffer")
	}
}

func TestChunkFeederExactMultiple(t *testing.T) {
	input := []byte("aaaabbbbcccc")
	f, err := NewChunkFeeder(t.Context(), bytes.NewReader(input), 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"aaaa", "bbbb", "cccc"} {
		got, nextErr := f.Next(t.Context())
		if nextErr != nil {
			t.Fatalf("Next() error: %v", nextErr)
		}
		if string(got) != want {
			t.Fatalf("Next() = %q, want %q", got, want)
		}
	}
	if _, err = f.Next(t.Context()); !errors.Is(err, ErrInputExhausted) {
		t.Fatalf("Next() after EOF = %v, want ErrInputExhausted", err)
	}
	// Exhaustion is sticky.
	if _, err = f.Next(t.Context()); !errors.Is(err, ErrInputExhausted) {
		t.Fatalf("second Next() after EOF = %v, want ErrInputExhausted", err)
	}
}

func TestChunkFeederDropsPartialTail(t *testing.T) {
	f, err := NewChunkFeeder(t.Context(), bytes.NewReader([]byte("aaaabb")), 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	got, err := f.Next(t.Context())
	if err != nil || string(got) != "aaaa" {
		t.Fatalf("Next() = %q, %v; want \"aaaa\", nil", got, err)
	}
	if _, err = f.Next(t.Context()); !errors.Is(err, ErrInputExhausted) {
		t.Fatalf("Next() = %v, want ErrInputExhausted", err)
	}
}

func TestChunkFeederEmptyInput(t *testing.T) {
	f, err := NewChunkFeeder(t.Context(), bytes.NewReader(nil), 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.Next(t.Context()); !errors.Is(err, ErrInputExhausted) {
		t.Fatalf("Next() = %v, want ErrInputExhausted", err)
	}
}

func TestChunkFeederReadError(t *testing.T) {
	r := io.MultiReader(bytes.NewReader([]byte("aaaa")), &failingReader{err: errors.New("boom")})
	f, err := NewChunkFeeder(t.Context(), r, 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got, nextErr := f.Next(t.Context()); nextErr != nil || string(got) != "aaaa" {
		t.Fatalf("Next() = %q, %v; want \"aaaa\", nil", got, nextErr)
	}
	_, err = f.Next(t.Context())
	if err == nil || errors.Is(err, ErrInputExhausted) {
		t.Fatalf("Next() after read error = %v, want the wrapped read error, not a clean stop", err)
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("Next() error %q does not wrap the read error", err)
	}
	// The failure is sticky too.
	if _, err = f.Next(t.Context()); err == nil || errors.Is(err, ErrInputExhausted) {
		t.Fatalf("second Next() after read error = %v, want the wrapped read error", err)
	}
}

type failingReader struct{ err error }

func (r *failingReader) Read([]byte) (int, error) { return 0, r.err }

func TestChunkFeederNextRespectsCancellation(t *testing.T) {
	pr, pw := io.Pipe()
	defer func() { _ = pw.Close() }()

	ctx, cancel := context.WithCancel(t.Context())
	f, err := NewChunkFeeder(ctx, pr, 4, 1)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		_, nextErr := f.Next(ctx)
		done <- nextErr
	}()

	cancel()
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	select {
	case nextErr := <-done:
		if !errors.Is(nextErr, context.Canceled) {
			t.Fatalf("Next() = %v, want context.Canceled", nextErr)
		}
	case <-timer.C:
		t.Fatal("Next() did not return after cancellation")
	}
}

func TestChunkFeederAlreadyCancelledContext(t *testing.T) {
	f, err := NewChunkFeeder(t.Context(), bytes.NewReader([]byte("aaaa")), 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err = f.Next(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Next() = %v, want context.Canceled", err)
	}
}

func TestChunkFeederReaderExitsOnCancelWhenBufferFull(t *testing.T) {
	// Enough input to fill the buffer and block the reader on send, with no
	// consumer draining it. Cancelling must let the reader goroutine exit,
	// which we observe as the channel closing.
	input := bytes.Repeat([]byte("x"), 4*10)
	ctx, cancel := context.WithCancel(t.Context())
	f, err := NewChunkFeeder(ctx, bytes.NewReader(input), 4, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Wait until the buffer is full so the reader is parked on the send.
	deadline := time.Now().Add(5 * time.Second)
	for len(f.ch) < cap(f.ch) {
		if time.Now().After(deadline) {
			t.Fatal("buffer never filled")
		}
		time.Sleep(time.Millisecond)
	}
	cancel()

	// Drain whatever was buffered; after that the channel must close.
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-f.ch:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("reader goroutine did not exit after cancellation")
		}
	}
}

func TestChunkFeederConcurrentConsumers(t *testing.T) {
	const size = 8
	const chunks = 500
	const consumers = 16

	input := make([]byte, 0, size*chunks)
	for i := range chunks {
		chunk := make([]byte, size)
		for j := range chunk {
			chunk[j] = byte(i>>uint(8*(j%4))) + byte(j)
		}
		input = append(input, chunk...)
	}

	f, err := NewChunkFeeder(t.Context(), bytes.NewReader(input), size, consumers*2)
	if err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	var got [][]byte
	var wg sync.WaitGroup
	for range consumers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				chunk, nextErr := f.Next(t.Context())
				if errors.Is(nextErr, ErrInputExhausted) {
					return
				}
				if nextErr != nil {
					t.Errorf("Next() error: %v", nextErr)
					return
				}
				mu.Lock()
				got = append(got, chunk)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(got) != chunks {
		t.Fatalf("received %d chunks, want %d", len(got), chunks)
	}
	// Every chunk must appear exactly once regardless of which consumer got it.
	sort.Slice(got, func(i, j int) bool { return bytes.Compare(got[i], got[j]) < 0 })
	var want [][]byte
	for i := 0; i < len(input); i += size {
		want = append(want, input[i:i+size])
	}
	sort.Slice(want, func(i, j int) bool { return bytes.Compare(want[i], want[j]) < 0 })
	for i := range want {
		if !bytes.Equal(got[i], want[i]) {
			t.Fatalf("chunk %d mismatch: got %x want %x", i, got[i], want[i])
		}
	}
}
