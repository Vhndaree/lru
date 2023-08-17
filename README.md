# LRU Cache Library

The LRU Cache Library is a Go package that provides a generic implementation of a Least Recently Used (LRU) cache. This cache is designed to store key-value pairs, automatically evicting the least recently used items when the cache reaches its capacity.

## Features

- Simple and easy-to-use API.
- Thread-safe implementation with read-write locks for concurrent access.
- Support for regular LRU caching as well as LRU caching with item expiry.

## Installation

Install the package using `go get`:

```sh
go get github.com/vhndaree/lru
```

## Usage

Import the package into your Go code:
```Go
import "github.com/vhndaree/lru"
```

### Basic usage
```Go
// Create a new basic LRU cache with a specified size.
cache := lru.New[int, string](cacheSize)

// Set key-value pairs in the cache.
cache.Set(1, "value1")
cache.Set(2, "value2")

// Retrieve values from the cache.
value, found := cache.Get(1)
if found {
    fmt.Println("Value:", value)
} else {
    fmt.Println("Key not found")
}

// Check if a key exists in the cache.
if cache.Contains(2) {
    fmt.Println("Key exists")
} else {
    fmt.Println("Key not found")
}

// Delete a key from the cache.
if cache.Del(2) {
    fmt.Println("Key deleted")
} else {
    fmt.Println("Key not found")
}
```

### LRU Cache with Expiry
```Go
// Create a new LRU cache with item expiry and a specified size.
cacheWithExpiry := lru.NewWithExpiry[int, string](cacheSize)

// Set key-value pairs with expiry in the cache.
cacheWithExpiry.SetWithExpiry(1, "value1", ttlMilliseconds)

// ... (Same retrieval and deletion operations as in basic cache)

// Note: The cache will automatically expire items after the specified TTL.
```

## Contribution
Contributions are welcome! If you find a bug or have suggestions for improvements, please open an issue or submit a pull request.
