package p2p

import (
	"container/list"
	"sync"
	"time"
)

// Cache is a thread-safe LRU cache with optional TTL-based expiration.
type Cache[K comparable, V any] struct {
	mu      sync.RWMutex
	maxSize int
	ttl     time.Duration
	items   map[K]*list.Element
	list    *list.List
}

type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

// NewCache creates a new cache with the given max size and optional TTL.
// If maxSize <= 0, the cache has no size limit.
// If ttl is 0, entries never expire based on time.
func NewCache[K comparable, V any](maxSize int, ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		maxSize: maxSize,
		ttl:     ttl,
		items:   make(map[K]*list.Element),
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
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	e := elem.Value.(*entry[K, V])

	if c.ttl > 0 && time.Now().After(e.expiresAt) {
		c.list.Remove(elem)
		delete(c.items, key)
		var zero V
		return zero, false
	}

	c.list.MoveToFront(elem)
	return e.value, true
}

// Contains checks if a key exists in the cache and is not expired.
// Uses a read lock and doesn't update LRU ordering.
func (c *Cache[K, V]) Contains(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.items[key]
	if !ok {
		return false
	}

	e := elem.Value.(*entry[K, V])

	if c.ttl > 0 && time.Now().After(e.expiresAt) {
		return false
	}

	return true
}

// Remove removes a key from the cache.
func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.Remove(elem)
		delete(c.items, key)
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

	c.items = make(map[K]*list.Element)
	c.list.Init()
}

// Keys returns all keys in the cache.
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
