package main

import (
	"container/list"
	"sync"
	"time"
)

type CacheItem struct {
	Value      string
	ExpiryTime time.Time
}

type Cache struct {
	mu       sync.RWMutex
	items    map[string]*list.Element
	eviction *list.List
	capacity int
}

type entry struct {
	key   string
	value CacheItem
}

func NewCache(capacity int) *Cache {
	return &Cache{
		items:    make(map[string]*list.Element),
		eviction: list.New(),
		capacity: capacity,
	}
}

func (c *Cache) Set(key, value string, ttl time.Duration) {
	// Lock the cache to allow only one instance of writing
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove the old value if it exists
	if elem, found := c.items[key]; found {
		c.eviction.Remove(elem)
		delete(c.items, key)
	}

	// Evict the least recently used item if the cache is at capacity
	if c.eviction.Len() >= c.capacity {
		c.evictLRU()
	}

	item := CacheItem{
		Value:      value,
		ExpiryTime: time.Now().Add(ttl),
	}
	elem := c.eviction.PushFront(&entry{key, item})
	c.items[key] = elem
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.Unlock()

	elem, found := c.items[key]
	if !found || time.Now().After(elem.Value.(*entry).value.ExpiryTime) {
		// If the item is not found or has expired, return false
		if found {
			c.eviction.Remove(elem)
			delete(c.items, key)
		}
		return "", false
	}

	// Move the accessed element to the front of the eviction list
	c.eviction.MoveToFront(elem)
	return elem.Value.(*entry).value.Value, true
}

func (c *Cache) startEvictionTicker(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			c.evictExpiredItems()
		}
	}()
}

func (c *Cache) evictLRU() {
	elem := c.eviction.Back()
	if elem != nil {
		c.eviction.Remove(elem)
		kv := elem.Value.(*entry)
		delete(c.items, kv.key)
	}
}

func (c *Cache) evictExpiredItems() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, elem := range c.items {
		if now.After(elem.Value.(*entry).value.ExpiryTime) {
			c.eviction.Remove(elem)
			delete(c.items, key)
		}
	}
}
