package cache

import (
	"time"
)

// Cache is a storage for entries that need to be kept alive for a certain time period
type Cache struct {
	entries         map[string]entry
	TTL             time.Duration
	CleanupInterval time.Duration
}

// New will instantiate a new Cache with a TTL and cleanupInterval.
// Each of these can be disabled by passing `-1` in their place
func New(TTL time.Duration, cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		entries:         make(map[string]entry),
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

	c.entries[key] = entry{
		value:      value,
		expiryTime: expiryTime,
	}
}

// Get will retrieve the value from the cache if it exists.
// If TTL is enabled, it will lazy delete expired entries on lookup.
func (c *Cache) Get(key string) (interface{}, bool) {
	e := c.entries[key]
	if e == (entry{}) {
		return nil, false
	}

	if e.isExpired() && c.TTL >= 0 {
		delete(c.entries, key)
		return nil, false
	}

	return e.value, true

}

// Delete will remove an item from the Cache idempotently.
func (c *Cache) Delete(key string) {
	delete(c.entries, key)
}

// Cache Helper functions

func (c *Cache) processCleanupInterval() {
	time.Sleep(c.CleanupInterval)

	for k, e := range c.entries {
		if e.isExpired() {
			delete(c.entries, k)
		}
	}

	go c.processCleanupInterval()
}

type entry struct {
	value      interface{}
	expiryTime int64
}

func (e *entry) isExpired() bool {
	if e.expiryTime >= 0 {
		return time.Now().Unix() >= e.expiryTime
	}

	return false
}
