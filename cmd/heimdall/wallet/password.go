package wallet

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

// promptPassword reads a password interactively. If in is an *os.File
// pointing at a terminal we use term.ReadPassword so keystrokes are
// not echoed; otherwise we fall back to reading a single line so
// tests and piped input still work.
//
// When confirm is true the function asks for the same password twice
// and returns an error on mismatch.
func promptPassword(in io.Reader, errW io.Writer, label string, confirm bool) (string, error) {
	if label == "" {
		label = "password"
	}

	readOnce := func(prompt string) (string, error) {
		fmt.Fprintf(errW, "%s: ", prompt)
		if f, ok := in.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
			raw, err := term.ReadPassword(int(f.Fd()))
			fmt.Fprintln(errW)
			if err != nil {
				return "", fmt.Errorf("reading password: %w", err)
			}
			return string(raw), nil
		}
		scanner := bufio.NewScanner(in)
		// Allow up to 1MiB passwords — absurd, but pathological inputs
		// from CI pipes should not silently truncate.
		scanner.Buffer(make([]byte, 4096), 1<<20)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("reading password: %w", err)
			}
			return "", nil
		}
		return scanner.Text(), nil
	}

	pw, err := readOnce(label)
	if err != nil {
		return "", err
	}
	if !confirm {
		return pw, nil
	}
	pw2, err := readOnce(label + " (confirm)")
	if err != nil {
		return "", err
	}
	if pw != pw2 {
		return "", fmt.Errorf("passwords do not match")
	}
	return pw, nil
}
