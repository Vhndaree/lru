package lru

import "github.com/vhndaree/lru/lru"

// LRU is a generic interface representing a Least Recently Used (LRU) cache.
type LRU[K comparable, V any] interface {
	// ListAll returns a map containing all key-value pairs in the LRU cache.
	ListAll() map[K]V

	// Contains checks if a key is present in the LRU cache.
	Contains(key K) bool

	// Set adds or updates a key-value pair in the LRU cache.
	Set(key K, value V)

	// Get retrieves the value associated with the given key from the LRU cache.
	// If the key is found, the associated value and true are returned; otherwise,
	// a default value and false are returned.
	Get(key K) (value V, found bool)

	// Del removes a key-value pair from the LRU cache and returns true if successful.
	// If the key is not found, it returns false.
	Del(key K) bool
}

// New creates a new instance of a Least Recently Used (LRU) cache with the specified size.
// It returns an instance of the LRU[K, V] interface.
func New[K comparable, V any](size int) LRU[K, V] {
	return lru.New[K, V](size)
}
