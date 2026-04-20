package heimdall

import (
	"bytes"
	"strings"
	"testing"
)

func TestHeimdallCmdHelp(t *testing.T) {
	cmd := HeimdallCmd
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("--help exited non-zero: %v", err)
	}
	if !strings.Contains(out.String(), "heimdall") {
		t.Fatalf("help output missing 'heimdall': %q", out.String())
	}
}
