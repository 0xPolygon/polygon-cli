package datastructures

import mapset "github.com/deckarep/golang-set/v2"

// BoundedSet is a simple set-based collection with a maximum size.
// When the set reaches capacity, the oldest element is evicted via Pop().
// This provides lower memory overhead compared to a full LRU cache when
// only membership tracking is needed without value storage.
type BoundedSet[T comparable] struct {
	set mapset.Set[T]
	max int
}

// NewBoundedSet creates a new BoundedSet with the specified maximum size.
func NewBoundedSet[T comparable](max int) *BoundedSet[T] {
	return &BoundedSet[T]{
		max: max,
		set: mapset.NewSet[T](),
	}
}

// Add adds an element to the set, evicting the oldest element if at capacity.
func (b *BoundedSet[T]) Add(elem T) {
	for b.set.Cardinality() >= b.max {
		b.set.Pop()
	}
	b.set.Add(elem)
}

// Contains returns true if the element exists in the set.
func (b *BoundedSet[T]) Contains(elem T) bool {
	return b.set.Contains(elem)
}

// Len returns the number of elements in the set.
func (b *BoundedSet[T]) Len() int {
	return b.set.Cardinality()
}

// Clear removes all elements from the set.
func (b *BoundedSet[T]) Clear() {
	b.set.Clear()
}
