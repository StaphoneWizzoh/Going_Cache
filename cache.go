package main

import "sync"

type CacheItem struct {
	Value string
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

func (c *Cache) Set(key, value string) {
	// Lock the cache to allow only one instance of writing
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value: value,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	// Only allows for reading operations
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return "", false
	}

	return item.Value, true
}
