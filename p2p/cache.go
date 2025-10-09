package p2p

import (
	"container/list"
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	mu      sync.RWMutex
	maxSize int
	ttl     time.Duration
	list    *list.List
}

// entry represents an entry in the cache.
type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

// NewCache creates a new cache with the given max size and optional TTL.
// If maxSize is 0 or negative, the cache has no size limit (only TTL eviction).
// If ttl is 0, entries never expire based on time.
func NewCache[K comparable, V any](maxSize int, ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		maxSize: maxSize,
		ttl:     ttl,
		list:    list.New(),
	}
}

// Add adds or updates a value in the cache.
func (c *Cache[K, V]) Add(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	expiresAt := time.Time{}
	if c.ttl > 0 {
		expiresAt = now.Add(c.ttl)
	}

	// Check if key exists, update it and move to front
	for elem := c.list.Front(); elem != nil; elem = elem.Next() {
		e := elem.Value.(*entry[K, V])
		if e.key == key {
			c.list.MoveToFront(elem)
			e.value = value
			e.expiresAt = expiresAt
			return
		}
	}

	// Add new entry at front
	e := &entry[K, V]{
		key:       key,
		value:     value,
		expiresAt: expiresAt,
	}
	c.list.PushFront(e)

	// Evict oldest if over max size (only if maxSize is set)
	if c.maxSize > 0 && c.list.Len() > c.maxSize {
		c.list.Remove(c.list.Back())
	}
}

// Get retrieves a value from the cache.
// Returns the value and true if found and not expired, otherwise zero value and false.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for elem := c.list.Front(); elem != nil; elem = elem.Next() {
		e := elem.Value.(*entry[K, V])
		if e.key == key {
			// Check if expired
			if c.ttl > 0 && now.After(e.expiresAt) {
				c.list.Remove(elem)
				var zero V
				return zero, false
			}
			// Move to front (LRU)
			c.list.MoveToFront(elem)
			return e.value, true
		}
	}

	var zero V
	return zero, false
}

// Contains checks if a key exists in the cache and is not expired.
func (c *Cache[K, V]) Contains(key K) bool {
	_, ok := c.Get(key)
	return ok
}

// Remove removes a key from the cache.
func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for elem := c.list.Front(); elem != nil; elem = elem.Next() {
		e := elem.Value.(*entry[K, V])
		if e.key == key {
			c.list.Remove(elem)
			return
		}
	}
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list.Len()
}

// Purge clears all items from the cache.
func (c *Cache[K, V]) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list.Init()
}

// Keys returns all keys in the cache (including potentially expired ones).
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, c.list.Len())
	for elem := c.list.Front(); elem != nil; elem = elem.Next() {
		e := elem.Value.(*entry[K, V])
		keys = append(keys, e.key)
	}
	return keys
}
