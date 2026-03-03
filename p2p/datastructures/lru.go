package datastructures

import (
	"container/list"
	"sync"
	"time"
)

// LRUOptions contains configuration for LRU caches with TTL.
type LRUOptions struct {
	MaxSize int
	TTL     time.Duration
}

// LRU is a thread-safe LRU cache with optional TTL-based expiration.
type LRU[K comparable, V any] struct {
	mu      sync.RWMutex
	maxSize int
	ttl     time.Duration
	items   map[K]*list.Element
	list    *list.List
}

type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt *time.Time
}

// NewLRU creates a new LRU cache with the given options.
// If opts.MaxSize <= 0, the cache has no size limit.
// If opts.TTL is 0, entries never expire based on time.
func NewLRU[K comparable, V any](opts LRUOptions) *LRU[K, V] {
	return &LRU[K, V]{
		maxSize: opts.MaxSize,
		ttl:     opts.TTL,
		items:   make(map[K]*list.Element),
		list:    list.New(),
	}
}

// Add adds or updates a value in the cache.
func (c *LRU[K, V]) Add(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt *time.Time
	if c.ttl > 0 {
		t := time.Now().Add(c.ttl)
		expiresAt = &t
	}

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		e := elem.Value.(*entry[K, V])
		e.value = value
		e.expiresAt = expiresAt
		return
	}

	e := &entry[K, V]{
		key:       key,
		value:     value,
		expiresAt: expiresAt,
	}
	elem := c.list.PushFront(e)
	c.items[key] = elem

	if c.maxSize > 0 && c.list.Len() > c.maxSize {
		back := c.list.Back()
		if back != nil {
			c.list.Remove(back)
			e := back.Value.(*entry[K, V])
			delete(c.items, e.key)
		}
	}
}

// Get retrieves a value from the cache and updates LRU ordering.
func (c *LRU[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	e := elem.Value.(*entry[K, V])

	if e.expiresAt != nil && time.Now().After(*e.expiresAt) {
		c.list.Remove(elem)
		delete(c.items, key)
		var zero V
		return zero, false
	}

	c.list.MoveToFront(elem)
	return e.value, true
}

// Peek retrieves a value from the cache without updating LRU ordering.
// Uses a read lock for better concurrency.
func (c *LRU[K, V]) Peek(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	e := elem.Value.(*entry[K, V])

	if e.expiresAt != nil && time.Now().After(*e.expiresAt) {
		var zero V
		return zero, false
	}

	return e.value, true
}

// PeekMany retrieves multiple values from the cache without updating LRU ordering.
// Uses a single read lock for all lookups, providing better concurrency than GetMany
// when LRU ordering updates are not needed.
func (c *LRU[K, V]) PeekMany(keys []K) []V {
	if len(keys) == 0 {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now()
	result := make([]V, 0, len(keys))

	for _, key := range keys {
		elem, ok := c.items[key]
		if !ok {
			continue
		}

		e := elem.Value.(*entry[K, V])

		if e.expiresAt != nil && now.After(*e.expiresAt) {
			continue
		}

		result = append(result, e.value)
	}

	return result
}

// Update atomically updates a value in the cache using the provided update function.
// The update function receives the current value (or zero value if not found) and
// returns the new value to store. This is thread-safe and prevents race conditions
// in get-modify-add patterns.
func (c *LRU[K, V]) Update(key K, updateFn func(V) V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	var expiresAt *time.Time
	if c.ttl > 0 {
		t := now.Add(c.ttl)
		expiresAt = &t
	}

	var currentVal V
	if elem, ok := c.items[key]; ok {
		e := elem.Value.(*entry[K, V])
		if e.expiresAt == nil || !now.After(*e.expiresAt) {
			currentVal = e.value
			// Update existing entry
			c.list.MoveToFront(elem)
			e.value = updateFn(currentVal)
			e.expiresAt = expiresAt
			return
		}
		// Entry expired, remove it
		c.list.Remove(elem)
		delete(c.items, key)
	}

	// Create new entry
	newVal := updateFn(currentVal)
	e := &entry[K, V]{
		key:       key,
		value:     newVal,
		expiresAt: expiresAt,
	}
	elem := c.list.PushFront(e)
	c.items[key] = elem

	// Enforce size limit
	if c.maxSize > 0 && c.list.Len() > c.maxSize {
		back := c.list.Back()
		if back != nil {
			c.list.Remove(back)
			e := back.Value.(*entry[K, V])
			delete(c.items, e.key)
		}
	}
}

// Remove removes a key from the cache and returns the value if it existed.
func (c *LRU[K, V]) Remove(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		e := elem.Value.(*entry[K, V])
		c.list.Remove(elem)
		delete(c.items, key)
		return e.value, true
	}

	var zero V
	return zero, false
}

// AddBatch adds multiple key-value pairs to the cache.
// Uses a single write lock for all additions, reducing lock contention
// compared to calling Add in a loop. Keys and values must have the same length.
func (c *LRU[K, V]) AddBatch(keys []K, values []V) {
	if len(keys) == 0 || len(keys) != len(values) {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt *time.Time
	if c.ttl > 0 {
		t := time.Now().Add(c.ttl)
		expiresAt = &t
	}

	for i, key := range keys {
		value := values[i]

		if elem, ok := c.items[key]; ok {
			c.list.MoveToFront(elem)
			e := elem.Value.(*entry[K, V])
			e.value = value
			e.expiresAt = expiresAt
			continue
		}

		e := &entry[K, V]{
			key:       key,
			value:     value,
			expiresAt: expiresAt,
		}
		elem := c.list.PushFront(e)
		c.items[key] = elem
	}

	// Enforce size limit after all additions
	for c.maxSize > 0 && c.list.Len() > c.maxSize {
		back := c.list.Back()
		if back == nil {
			break
		}
		c.list.Remove(back)
		e := back.Value.(*entry[K, V])
		delete(c.items, e.key)
	}
}
