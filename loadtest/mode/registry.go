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

	// Register by canonical name
	name := strings.ToLower(mode.Name())
	registry[name] = mode

	// Register by aliases
	for _, alias := range mode.Aliases() {
		registry[strings.ToLower(alias)] = mode
	}
}

// Get retrieves a mode by name or alias.
func Get(name string) (Runner, error) {
	registryLock.RLock()
	defer registryLock.RUnlock()

	mode, found := registry[strings.ToLower(name)]
	if !found {
		return nil, fmt.Errorf("unrecognized load test mode: %s", name)
	}
	return mode, nil
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
