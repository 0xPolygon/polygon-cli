package render

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
)

// WatchFn is the data source a watcher polls. It should return an
// already-rendered output string for the current snapshot. Any error
// is logged to the watcher's Err writer but does not terminate the
// loop unless the context has been cancelled.
type WatchFn func(ctx context.Context) (string, error)

// Watch polls fn at interval, printing the output to out whenever it
// changes. Always exits cleanly on ctx cancellation; guarantees timer
// cleanup per CLAUDE.md.
//
// The first snapshot is always printed. Subsequent snapshots are
// printed only when they differ from the previous output.
func Watch(ctx context.Context, out, errOut io.Writer, interval time.Duration, fn WatchFn) error {
	if interval <= 0 {
		interval = 2 * time.Second
	}

	var prev string
	// Execute an initial snapshot synchronously so the caller sees
	// output immediately.
	snap, err := fn(ctx)
	if err != nil {
		fmt.Fprintf(errOut, "watch: %v\n", err)
	} else {
		writeSnapshot(out, prev, snap, time.Now())
		prev = snap
	}

	timer := time.NewTimer(interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			snap, err := fn(ctx)
			if err != nil {
				fmt.Fprintf(errOut, "watch: %v\n", err)
			} else if snap != prev {
				writeSnapshot(out, prev, snap, time.Now())
				prev = snap
			}
			timer.Reset(interval)
		}
	}
}

func writeSnapshot(w io.Writer, prev, snap string, at time.Time) {
	fmt.Fprintf(w, "--- %s ---\n", at.UTC().Format(time.RFC3339))
	if prev == "" {
		_, _ = io.Copy(w, bytesBuf(snap))
		if !hasTrailingNewline(snap) {
			fmt.Fprintln(w)
		}
		return
	}
	// Simple line-level diff. Cheap and good enough for dense cast-style
	// output.
	writeDiff(w, prev, snap)
}

func writeDiff(w io.Writer, before, after string) {
	beforeLines := splitLines(before)
	afterLines := splitLines(after)
	beforeSet := make(map[string]struct{}, len(beforeLines))
	for _, l := range beforeLines {
		beforeSet[l] = struct{}{}
	}
	afterSet := make(map[string]struct{}, len(afterLines))
	for _, l := range afterLines {
		afterSet[l] = struct{}{}
	}
	for _, l := range beforeLines {
		if _, ok := afterSet[l]; !ok {
			fmt.Fprintf(w, "- %s\n", l)
		}
	}
	for _, l := range afterLines {
		if _, ok := beforeSet[l]; !ok {
			fmt.Fprintf(w, "+ %s\n", l)
		}
	}
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	out := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func hasTrailingNewline(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '\n'
}

func bytesBuf(s string) *bytes.Buffer {
	return bytes.NewBufferString(s)
}
