package render

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

// TestEnableWatchPassThrough asserts that a command with the --watch
// flag unset behaves exactly like the original command: one invocation
// of RunE, no extra output.
func TestEnableWatchPassThrough(t *testing.T) {
	var calls int32
	cmd := &cobra.Command{
		Use: "stub",
		RunE: func(c *cobra.Command, _ []string) error {
			atomic.AddInt32(&calls, 1)
			fmt.Fprintln(c.OutOrStdout(), "hello")
			return nil
		},
	}
	EnableWatch(cmd)
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("RunE calls = %d, want 1", got)
	}
	if !strings.Contains(out.String(), "hello") {
		t.Fatalf("expected hello in output, got %q", out.String())
	}
}

// TestEnableWatchLoopsUntilCancel asserts that --watch causes the
// command to run repeatedly until the context is cancelled, and that
// the cancellation returns promptly.
func TestEnableWatchLoopsUntilCancel(t *testing.T) {
	var calls int32
	cmd := &cobra.Command{
		Use: "stub",
		RunE: func(c *cobra.Command, _ []string) error {
			atomic.AddInt32(&calls, 1)
			fmt.Fprintln(c.OutOrStdout(), "tick")
			return nil
		},
	}
	EnableWatch(cmd)
	cmd.SetArgs([]string{"--watch", "25ms"})

	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- cmd.ExecuteContext(ctx) }()

	select {
	case err := <-done:
		if err != context.DeadlineExceeded && err != nil {
			// context.Canceled is also acceptable; anything else is
			// unexpected.
			if !strings.Contains(err.Error(), "context") {
				t.Fatalf("Execute returned unexpected error: %v", err)
			}
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("watch loop failed to exit within 2s; collected %d calls", atomic.LoadInt32(&calls))
	}

	if got := atomic.LoadInt32(&calls); got < 2 {
		t.Fatalf("expected at least 2 iterations, got %d", got)
	}
	if !strings.Contains(out.String(), "tick") {
		t.Fatalf("expected tick in output, got %q", out.String())
	}
}

// TestEnableWatchTreeCoversDescendants asserts that EnableWatchTree
// recurses into descendants and registers --watch on each leaf with a
// RunE, but does not touch umbrellas that lack one.
func TestEnableWatchTreeCoversDescendants(t *testing.T) {
	umbrella := &cobra.Command{Use: "umbrella"}
	leaf := &cobra.Command{
		Use:  "leaf",
		RunE: func(c *cobra.Command, _ []string) error { return nil },
	}
	bareParent := &cobra.Command{Use: "bare"}
	umbrella.AddCommand(leaf)
	umbrella.AddCommand(bareParent)

	EnableWatchTree(umbrella)

	if leaf.Flags().Lookup(WatchFlag) == nil {
		t.Fatal("leaf did not get --watch flag")
	}
	if bareParent.Flags().Lookup(WatchFlag) != nil {
		t.Fatal("umbrella parent without RunE should not get --watch flag")
	}
}
