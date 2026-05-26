package msgs

import (
	"sort"
	"sync"

	"github.com/spf13/cobra"

	"github.com/0xPolygon/polygon-cli/internal/heimdall/config"
)

// Factory builds one instance of a single msg subcommand. Every msg
// supported by mktx/send/estimate registers one Factory; the umbrella
// commands call each factory with their mode so each umbrella owns its
// own command tree (cobra commands can only have one parent).
//
// The factory returns a *cobra.Command that is fully wired: shared tx
// flags are attached via RegisterFlags, per-msg flags bound, and the
// RunE closure assembles and executes the builder using Execute.
type Factory func(mode Mode, globalFlags *config.Flags) *cobra.Command

// registryEntry binds a factory to its canonical name (the subcommand
// verb, e.g. "withdraw"). Additional aliases are declared on the
// returned *cobra.Command via the Aliases field in the factory itself.
type registryEntry struct {
	Name    string
	Factory Factory
}

var (
	registryMu sync.RWMutex
	registry   = map[string]registryEntry{}
)

// RegisterFactory adds a msg factory to the package-level registry.
// W3 registers `withdraw` via an init() in msgs/withdraw.go; W4
// registers additional msg subcommands the same way. Duplicate names
// panic during init so the error is caught at startup.
//
// name is the cobra subcommand verb (e.g. "withdraw"); it must be
// non-empty and unique across the registry.
func RegisterFactory(name string, factory Factory) {
	if name == "" {
		panic("msgs.RegisterFactory: name is empty")
	}
	if factory == nil {
		panic("msgs.RegisterFactory: factory is nil")
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, ok := registry[name]; ok {
		panic("msgs.RegisterFactory: duplicate name " + name)
	}
	registry[name] = registryEntry{Name: name, Factory: factory}
}

// Names returns the sorted list of registered msg subcommands. Used
// by tests to assert the registry shape and by the umbrella commands
// to build their `Long` usage hints.
func Names() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]string, 0, len(registry))
	for n := range registry {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// BuildChildren invokes every registered factory for the given mode
// and returns the resulting cobra subcommands in registry order. The
// umbrella command then AddCommand's them onto its own cobra tree.
//
// A fresh slice of *cobra.Command is returned on every call so each
// umbrella owns independent children. Cobra's command tree is not
// thread-safe and a single *cobra.Command can only have one parent.
func BuildChildren(mode Mode, globalFlags *config.Flags) []*cobra.Command {
	registryMu.RLock()
	defer registryMu.RUnlock()
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	out := make([]*cobra.Command, 0, len(names))
	for _, n := range names {
		entry := registry[n]
		cmd := entry.Factory(mode, globalFlags)
		if cmd == nil {
			continue
		}
		out = append(out, cmd)
	}
	return out
}

// resetRegistryForTest clears the registry. Intended for tests only.
// Kept unexported and short-named because it is not part of the
// public surface; test files in this package can call it via
// internal access.
func resetRegistryForTest() {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = map[string]registryEntry{}
}
