package main

import (
	"fmt"
	"sync"
	"time"
)

//CacheEntry represents an entry in the cache
type CacheEntry struct {
	value   string
	expires time.Time
}

//Cache represents a map[string]string that is safe
//for concurrent access
type Cache struct {
	//TODO: protect this map with a RWMutex
	entries map[string]*CacheEntry // key: string, value: *CacheEntry
	mu      sync.RWMutex
	quit    chan bool // common pattern in GoLang, not exported
}

//NewCache creates and returns a new Cache
func NewCache() *Cache {
	c := &Cache{
		entries: make(map[string]*CacheEntry),
		mu:      sync.RWMutex{},
		quit:    make(chan bool),
	}
	go c.startJanitor()
	return c
}

// Close allows external packages to quit
func (c *Cache) Close() {
	c.quit <- true
}

// clear's memstore of expired entries
// called as a goroutine
func (c *Cache) startJanitor() {
	ticker := time.NewTicker(time.Second)
	for { // runs indefinitely
		select {
		case <-ticker.C: // read from the ticker channel
			c.purgeExpired() // called every second
		case <-c.quit:
			return
		}
	}
}

func (c *Cache) purgeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	nPurged := 0
	for key, entry := range c.entries {
		if now.After(entry.expires) {
			delete(c.entries, key) // map delete function
			nPurged++
		}
	}

	fmt.Printf("purge %d entries\n", nPurged)
}

//Get returns the value associated with the requested key.
//The returned boolean will be false if the key was not
//in the cache.
func (c *Cache) Get(key string) (string, bool) {
	//TODO: implement this method and
	//replace the return statement below

	c.mu.RLock() //
	defer c.mu.RUnlock()

	entry := c.entries[key]
	if entry == nil {
		return "", false
	}

	return entry.value, true
}

//Set sets the value associated with the given key.
//If the key is not yet in the cache, it will be added.
func (c *Cache) Set(key string, value string, ttl time.Duration) {
	//TODO: implement this method

	c.mu.Lock()
	defer c.mu.Unlock()

	entry := c.entries[key]
	if entry == nil {
		entry = &CacheEntry{}  // create new cacheentry struct
		c.entries[key] = entry // insert struct into map
	}

	entry.value = value
	entry.expires = time.Now().Add(ttl)

}
