package cache

import (
	"time"
)

// Cache is a storage for entries that need to be kept alive for a certain time period
type Cache struct {
	entries         []entry
	TTL             time.Duration
	CleanupInterval time.Duration
}

// New will instantiate a new Cache
func New(TTL time.Duration, cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		TTL:             TTL,
		CleanupInterval: cleanupInterval,
	}

	if cache.CleanupInterval >= 0 {
		go cache.processCleanupInterval()
	}

	return cache
}

// Put will upsert a key/value pair to the cache
func (c *Cache) Put(key string, value interface{}) {
	var expiryTime int64

	if c.TTL >= 0 {
		expiryTime = time.Now().Add(c.TTL).Unix()
	} else {
		expiryTime = -1
	}

	e := entry{
		key:        key,
		value:      value,
		expiryTime: expiryTime,
	}

	for i, e := range c.entries {
		if e.key == key {
			c.removeEntry(i)
		}
	}

	c.entries = append(c.entries, e)
}

// Get will retrieve the value from the cache if it exists
func (c *Cache) Get(key string) (interface{}, bool) {
	for i, entry := range c.entries {
		if entry.key == key {
			if entry.IsExpired() && c.TTL >= 0 {
				c.removeEntry(i)
				return nil, false
			}
			return entry.value, true
		}
	}

	return nil, false
}

// Delete will remove an item from the Cache idempotently
func (c *Cache) Delete(key string) {
	for i, e := range c.entries {
		if e.key == key {
			c.removeEntry(i)
		}
	}
}

func (c *Cache) processCleanupInterval() {
	time.Sleep(c.CleanupInterval)

	for i, e := range c.entries {
		if e.IsExpired() {
			c.removeEntry(i)
		}
	}

	go c.processCleanupInterval()
}

func (c *Cache) removeEntry(index int) {
	c.entries[index] = c.entries[len(c.entries)-1] // Copy last element to index i
	c.entries[len(c.entries)-1] = entry{}
	c.entries = c.entries[:len(c.entries)-1]
}

type entry struct {
	key        string
	value      interface{}
	expiryTime int64
}

func (e *entry) IsExpired() bool {
	if e.expiryTime >= 0 {
		return time.Now().Unix() >= e.expiryTime
	}

	return false
}
