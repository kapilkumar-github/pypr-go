package cache

import (
	"sync"
	"time"
)

type item[T any] struct {
	value     T
	expiresAt time.Time
}

type Cache[T any] struct {
	mu    sync.RWMutex
	items map[string]item[T]
	ttl   time.Duration
}

func New[T any](ttl time.Duration) *Cache[T] {
	return &Cache[T]{
		items: make(map[string]item[T]),
		ttl:   ttl,
	}
}

func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item[T]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		var zero T
		return zero, false
	}

	if time.Now().After(item.expiresAt) {
		c.Delete(key)

		var zero T
		return zero, false
	}

	return item.value, true
}

func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]item[T])
}
