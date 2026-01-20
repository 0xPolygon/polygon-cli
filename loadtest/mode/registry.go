package mode

import (
	"fmt"
	"strings"
	"sync"
)

var (
	registry     = make(map[string]Runner)
	registryLock sync.RWMutex
)

// Register adds a mode to the registry.
func Register(mode Runner) {
	registryLock.Lock()
	defer registryLock.Unlock()

	// Register by canonical name (lowercase for case-insensitive lookup)
	name := strings.ToLower(mode.Name())
	registry[name] = mode

	// Register by aliases (preserve case for case-sensitive aliases like "r" vs "R")
	for _, alias := range mode.Aliases() {
		registry[alias] = mode
	}
}

// Get retrieves a mode by name or alias.
// First tries exact match (for case-sensitive aliases like "r" vs "R"),
// then falls back to lowercase match (for canonical names).
func Get(name string) (Runner, error) {
	registryLock.RLock()
	defer registryLock.RUnlock()

	// Try exact match first (for case-sensitive aliases)
	if mode, found := registry[name]; found {
		return mode, nil
	}
	// Fall back to lowercase match (for canonical names)
	if mode, found := registry[strings.ToLower(name)]; found {
		return mode, nil
	}
	return nil, fmt.Errorf("unrecognized load test mode: %s", name)
}

// GetAll returns all registered modes (by canonical name only).
func GetAll() map[string]Runner {
	registryLock.RLock()
	defer registryLock.RUnlock()

	// Return only canonical names to avoid duplicates
	result := make(map[string]Runner)
	seen := make(map[Runner]bool)
	for _, mode := range registry {
		if !seen[mode] {
			seen[mode] = true
			result[mode.Name()] = mode
		}
	}
	return result
}
