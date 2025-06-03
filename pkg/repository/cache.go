package repository

import (
	"sync"
)

type Cache[T any] struct {
	data  []T
	mutex sync.RWMutex
}

// NewCache creates a new cache for type T
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{data: make([]T, 0)}
}

// Add adds an item to the cache
func (c *Cache[T]) Add(item T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = append(c.data, item)
}

// Clear empties the cache
func (c *Cache[T]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = c.data[:0]
}

// Count returns the number of items in the cache
func (c *Cache[T]) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.data)
}
