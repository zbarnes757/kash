package cache

import (
	"time"
)

type Cache struct {
	entries            []Entry
	expirationDuration time.Duration
}

type Entry struct {
	key       string
	value     interface{}
	entryTime time.Time
}

func (c *Cache) Cleanup() {
	time.Sleep(1 * time.Second)

	for i, entry := range c.entries {
		if time.Now().Unix() >= entry.entryTime.Add(c.expirationDuration).Unix() {
			c.entries[i] = c.entries[len(c.entries)-1] // Copy last element to index i
			c.entries[len(c.entries)-1] = Entry{}
			c.entries = c.entries[:len(c.entries)-1]
		}
	}

	go c.Cleanup()
}

func (c *Cache) Put(key string, value interface{}) {
	e := Entry{key, value, time.Now()}
	c.entries = append(c.entries, e)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	for _, entry := range c.entries {
		if entry.key == key {
			return entry.value, true
		}
	}

	return nil, false
}
