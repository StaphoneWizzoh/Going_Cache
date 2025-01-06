package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"
)

// CacheEntry represents an individual entry in the cache.
type CacheEntry struct {
	Key        string
	Value      interface{}
	ExpiresAt  time.Time
}

// Cache represents the cache structure.
type Cache struct {
	capacity    int
	entries     map[string]*list.Element
	evictionList *list.List
	filePath    string
	mu          sync.Mutex
}

// NewCache creates a new cache with a specified capacity and persistence file path.
func NewCache(capacity int, filePath string) *Cache {
	cache := &Cache{
		capacity:    capacity,
		entries:     make(map[string]*list.Element),
		evictionList: list.New(),
		filePath:    filePath,
	}
	cache.loadFromDisk()
	return cache
}

// Get retrieves a value from the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.entries[key]; exists {
		entry := element.Value.(*CacheEntry)

		// Check if the entry has expired
		if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
			c.removeElement(element)
			return nil, false
		}

		// Move to the front of the eviction list
		c.evictionList.MoveToFront(element)
		return entry.Value, true
	}
	return nil, false
}

// Put adds a new key-value pair to the cache with an optional TTL.
func (c *Cache) Put(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the key already exists
	if element, exists := c.entries[key]; exists {
		c.evictionList.MoveToFront(element)
		entry := element.Value.(*CacheEntry)
		entry.Value = value
		if ttl > 0 {
			entry.ExpiresAt = time.Now().Add(ttl)
		} else {
			entry.ExpiresAt = time.Time{}
		}
		c.saveToDisk()
		return
	}

	// Add a new entry
	entry := &CacheEntry{
		Key:   key,
		Value: value,
	}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}

	element := c.evictionList.PushFront(entry)
	c.entries[key] = element

	// Check if we need to evict an entry
	if len(c.entries) > c.capacity {
		c.evict()
	}

	c.saveToDisk()
}

// evict removes the least recently used item from the cache.
func (c *Cache) evict() {
	if element := c.evictionList.Back(); element != nil {
		c.removeElement(element)
	}
}

// removeElement removes an element from the cache.
func (c *Cache) removeElement(element *list.Element) {
	c.evictionList.Remove(element)
	entry := element.Value.(*CacheEntry)
	delete(c.entries, entry.Key)
}

// Clear removes all entries from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*list.Element)
	c.evictionList.Init()
	os.Remove(c.filePath)
}

// saveToDisk saves the cache content to disk.
func (c *Cache) saveToDisk() {
	file, err := os.Create(c.filePath)
	if err != nil {
		fmt.Println("Error saving cache to disk:", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	data := make([]*CacheEntry, 0, len(c.entries))
	for _, element := range c.entries {
		entry := element.Value.(*CacheEntry)
		data = append(data, entry)
	}
	encoder.Encode(data)
}

// loadFromDisk loads the cache content from disk.
func (c *Cache) loadFromDisk() {
	file, err := os.Open(c.filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error loading cache from disk:", err)
		}
		return
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var data []*CacheEntry
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding cache data:", err)
		return
	}

	for _, entry := range data {
		element := c.evictionList.PushFront(entry)
		c.entries[entry.Key] = element
	}
}

// DisplayCache prints the current cache content (for debugging).
func (c *Cache) DisplayCache() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for element := c.evictionList.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*CacheEntry)
		fmt.Printf("Key: %s, Value: %v, ExpiresAt: %v\n", entry.Key, entry.Value, entry.ExpiresAt)
	}
}

func main() {
	cache := NewCache(3, "cache_data.gob")

	// Add items to cache
	cache.Put("key1", "value1", 5*time.Second)
	cache.Put("key2", "value2", 10*time.Second)
	cache.Put("key3", "value3", 0) // No expiration

	// Display cache
	fmt.Println("Initial cache:")
	cache.DisplayCache()

	// Access items
	value, found := cache.Get("key1")
	if found {
		fmt.Println("Retrieved key1:", value)
	}

	// Wait for TTL expiration
	fmt.Println("Waiting 6 seconds...")
	time.Sleep(6 * time.Second)

	// Clean expired items by accessing
	_, _ = cache.Get("key1")
	_, _ = cache.Get("key2")

	// Add a new item to trigger eviction
	cache.Put("key4", "value4", 0)

	fmt.Println("Cache after eviction and TTL expiration:")
	cache.DisplayCache()

	// Clear the cache
	cache.Clear()
	fmt.Println("Cache after clearing:")
	cache.DisplayCache()
}
