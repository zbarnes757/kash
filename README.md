[![Build Status](https://travis-ci.org/zbarnes757/kash.svg?branch=master)](https://travis-ci.org/zbarnes757/kash) [![](https://godoc.org/github.com/zbarnes757/kash?status.svg)](http://godoc.org/github.com/zbarnes757/kash)

# Kash

Kash is a simple in-memory cache for Go that uses a simple TTL and cleanup system.

### Usage

```go
import (
  "time"

  "github.com/zbarnes757/kash"
)

func main() {
  // Pass in desired TTL and cleanup interval
  // This example has a TTL of 1 minute and cleans up every 2 minutes
  cache := kash.New(1 * time.Minute, 2 * time.Minute)

  // Add a new value to the cache. Will also update existing entry if the key exists
  cache.Put("myKey", 42)

  // Retrieve a value for a key
  value, found := cache.Get("myKey")

  // delete a key/value pair from the cache
  cache.Delete("myKey")

  // To disable automatic cleanup, set value to -1
  // This will still lazy delete an entry on lookup if it has expired
  cache = kash.New(1 * time.Minute, -1)

  // To prevent any automatic deletion of entries, set both values to -1
  cache = kash.New(-1, -1)
}
```
