package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value string
	ExpiryTime time.Time
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

func (c *Cache) Set(key, value string, ttl time.Duration) {
	// Lock the cache to allow only one instance of writing
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value: value,
		ExpiryTime: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	// Only allows for reading operations
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]

	// If the item is not found or has expired, return false
	if !found || time.Now().After(item.ExpiryTime){
		return "", false
	}

	return item.Value, true
}

func (c *Cache) startEvictionTicker(d time.Duration){
	ticker := time.NewTicker(d)
	go func ()  {
		for range ticker.C{
			c.evictExpiredItems()
		}	
	}()
}

func (c *Cache) evictExpiredItems(){
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items{
		if now.After(item.ExpiryTime){
			delete(c.items, key)
		}
	}
}