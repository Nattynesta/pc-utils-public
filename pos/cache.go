package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

var appCache = &Cache{
	items: make(map[string]CacheItem),
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.items[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}
