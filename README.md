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
}
```
