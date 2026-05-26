package render

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// WatchFlag is the name of the per-command duration flag added by
// EnableWatch. Exported so tests can set it directly via cmd.Flags().
const WatchFlag = "watch"

// EnableWatch decorates cmd with a `--watch DURATION` flag and wraps
// cmd.RunE so the command repeats its output at the requested interval
// until the context is cancelled. The first iteration always runs; if
// --watch is zero (unset), the wrapper is a no-op and the original
// RunE is invoked verbatim.
//
// Between iterations the watcher clears the terminal with the standard
// VT100 sequence, matching `watch(1)` behaviour. When the command is
// not attached to a TTY the separator is a plain divider line, so
// piping `polycli heimdall <cmd> --watch 5s | cat` stays readable.
//
// The wrapper captures the original RunE once and is idempotent: a
// second call on the same command is a no-op. Usage strings document
// the flag so it surfaces in `--help` output.
func EnableWatch(cmd *cobra.Command) {
	if cmd == nil || cmd.RunE == nil {
		return
	}
	if cmd.Flags().Lookup(WatchFlag) != nil {
		return
	}
	var interval time.Duration
	cmd.Flags().DurationVar(&interval, WatchFlag, 0, "repeat every DURATION (e.g. 5s) until Ctrl-C; 0 disables")

	orig := cmd.RunE
	cmd.RunE = func(c *cobra.Command, args []string) error {
		if interval <= 0 {
			return orig(c, args)
		}
		ctx := c.Context()
		outWriter := c.OutOrStdout()
		// Each iteration writes into a buffer so the screen is cleared
		// between snapshots regardless of how many Fprintf calls the
		// underlying RunE makes.
		isTTY := writerIsTerminal(outWriter)

		timer := time.NewTimer(0)
		defer timer.Stop()
		// First tick fires immediately; subsequent ticks wait
		// `interval`.
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
			}
			var buf bytes.Buffer
			c.SetOut(&buf)
			err := orig(c, args)
			c.SetOut(outWriter)

			clearScreen(outWriter, isTTY)
			fmt.Fprintf(outWriter, "--- %s (every %s) ---\n", time.Now().UTC().Format(time.RFC3339), interval)
			_, _ = io.Copy(outWriter, &buf)
			if err != nil {
				// Stream the error alongside the snapshot but don't
				// abort the loop — a transient node blip shouldn't
				// kill a long-running watch.
				fmt.Fprintf(c.ErrOrStderr(), "watch: %v\n", err)
			}
			if !timerReset(timer, interval) {
				return nil
			}
		}
	}
}

// EnableWatchTree recursively applies EnableWatch to cmd and every
// descendant that has a RunE. Parents-only (pure umbrella) commands
// are skipped because they do not print anything loopable.
func EnableWatchTree(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	if cmd.RunE != nil {
		EnableWatch(cmd)
	}
	for _, sub := range cmd.Commands() {
		EnableWatchTree(sub)
	}
}

// timerReset guards against the documented "race between Stop and
// channel drain" pattern; we always Reset on a drained timer.
func timerReset(t *time.Timer, d time.Duration) bool {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	t.Reset(d)
	return true
}

// clearScreen emits the VT100 clear sequence when attached to a
// terminal. Otherwise it prints a short divider so pipes remain
// readable.
func clearScreen(w io.Writer, tty bool) {
	if tty {
		// ANSI: home cursor + clear screen.
		_, _ = io.WriteString(w, "\x1b[H\x1b[2J")
		return
	}
	_, _ = fmt.Fprintln(w, strings.Repeat("-", 40))
}

func writerIsTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
