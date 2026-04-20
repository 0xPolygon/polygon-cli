package heimdallutil

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// buildRootWithCompletions constructs a minimal cobra tree that mirrors
// the runtime wiring (root → util → completions) so the completion
// command can walk up to the root via cmd.Parent().
func buildRootWithCompletions() *cobra.Command {
	root := &cobra.Command{Use: "polycli-test", SilenceUsage: true}
	util := &cobra.Command{Use: "util", Args: cobra.NoArgs}
	util.AddCommand(newCompletionsCmd(util))
	// A second dummy subcommand so completion output has something to chew on.
	util.AddCommand(&cobra.Command{Use: "ping", Short: "test subcommand", Run: func(*cobra.Command, []string) {}})
	root.AddCommand(util)
	return root
}

func runCompletions(t *testing.T, shell string) string {
	t.Helper()
	root := buildRootWithCompletions()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"util", "completions", shell})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("completions %s: %v", shell, err)
	}
	return buf.String()
}

func TestCompletionsBashSmoke(t *testing.T) {
	out := runCompletions(t, "bash")
	if out == "" {
		t.Fatal("empty bash completion")
	}
	// Cobra's v2 bash completion declares __start_<root>() and invokes
	// `complete ... -F __start_polycli-test polycli-test` near the end.
	if !strings.Contains(out, "-F __start_polycli-test") {
		t.Errorf("bash completion missing `-F __start_polycli-test` marker:\n%s", out)
	}
	if !strings.Contains(out, "polycli-test") {
		t.Errorf("bash completion missing root command name:\n%s", out)
	}
}

func TestCompletionsZshSmoke(t *testing.T) {
	out := runCompletions(t, "zsh")
	if out == "" {
		t.Fatal("empty zsh completion")
	}
	if !strings.Contains(out, "#compdef") {
		t.Errorf("zsh completion missing #compdef directive:\n%s", out)
	}
}

func TestCompletionsFishSmoke(t *testing.T) {
	out := runCompletions(t, "fish")
	if out == "" {
		t.Fatal("empty fish completion")
	}
	if !strings.Contains(out, "complete -c polycli-test") {
		t.Errorf("fish completion missing `complete -c` directive:\n%s", out)
	}
}

func TestCompletionsPowerShellSmoke(t *testing.T) {
	out := runCompletions(t, "powershell")
	if out == "" {
		t.Fatal("empty powershell completion")
	}
	if !strings.Contains(out, "Register-ArgumentCompleter") {
		t.Errorf("powershell completion missing `Register-ArgumentCompleter`:\n%s", out)
	}
}

func TestCompletionsCaseInsensitive(t *testing.T) {
	// Users commonly type `BASH` or `Zsh`; we lower-case at the boundary.
	root := buildRootWithCompletions()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"util", "completions", "BASH"})
	if err := root.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("completions BASH: %v", err)
	}
	if !strings.Contains(buf.String(), "-F __start_polycli-test") {
		t.Errorf("uppercase shell name did not produce bash completion:\n%s", buf.String())
	}
}

func TestCompletionsRejectsUnknownShell(t *testing.T) {
	root := buildRootWithCompletions()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"util", "completions", "tcsh"})
	err := root.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
	if !strings.Contains(err.Error(), "tcsh") {
		t.Errorf("error does not mention offending shell: %v", err)
	}
}
