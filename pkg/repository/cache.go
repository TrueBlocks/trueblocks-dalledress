package repository

import (
	"sync"
)

type Cache[T any] struct {
	data  []T
	mutex sync.RWMutex
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		data: make([]T, 0),
	}
}

func (c *Cache[T]) Add(item T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = append(c.data, item)
}

func (c *Cache[T]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = c.data[:0]
}

func (c *Cache[T]) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.data)
}
