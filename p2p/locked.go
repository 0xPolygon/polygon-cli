package p2p

import "sync"

// Locked wraps a value with a RWMutex for thread-safe access.
type Locked[T any] struct {
	value T
	mu    sync.RWMutex
}

// Get returns the current value (thread-safe read).
func (l *Locked[T]) Get() T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.value
}

// Set updates the value (thread-safe write).
func (l *Locked[T]) Set(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.value = value
}

// Update atomically updates the value using a function.
// The function receives the current value and returns the new value and a result.
// The result is returned to the caller.
func (l *Locked[T]) Update(fn func(T) (T, bool)) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	newValue, changed := fn(l.value)
	l.value = newValue
	return changed
}
