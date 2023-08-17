package lru

type Base[K comparable, V any] interface {
	// Contains checks if the provided key is present in the LRU cache.
	// It returns true if the key is found in the cache, and false otherwise.
	// The function does not affect the cache's state or modify any data.
	Contains(key K) bool

	// Get retrieves the value associated with the provided key from the LRU cache.
	// If the key exists in the cache, its corresponding value is returned along with a boolean true.
	// If the key is not found in the cache, an empty value and boolean false are returned.
	//
	// The Get operation updates the order of items in the cache to reflect the most recently accessed item.
	// If the item exists, it is moved to the head of the cache to prioritize recently accessed items.
	Get(key K) (value V, found bool)

	// Del removes the key-value pair associated with the provided key from the LRU cache.
	// If the key is found and the removal is successful, the function returns true.
	// If the key is not found, it returns false.
	//
	// The function also adjusts the internal linked list of cache items to maintain its order.
	// If the removed item was the head or tail of the list, appropriate adjustments are made.
	// The deleted item's memory is released for garbage collection.
	Del(key K) bool
}

// LRU is a generic interface representing a Least Recently Used (LRU) cache.
type LRU[K comparable, V any] interface {
	Base[K, V]

	// Set adds or updates a key-value pair in the LRU cache with the provided key and value.
	// If the key already exists in the cache, its corresponding value will be updated.
	// If the key is new, a new entry will be created with the provided value.
	//
	// This function is thread-safe and utilizes a read-write lock to ensure concurrent access
	// to the cache's internal data structures.
	Set(key K, value V)
}

// LRUWithExpiry is a generic interface representing a Least Recently Used (LRU) cache.
type LRUWithExpiry[K comparable, V any] interface {
	Base[K, V]

	// SetWithExpiry adds or updates a key-value pair in the LRU cache with the provided key, value, and time-to-live (TTL).
	// If the key already exists in the cache, its corresponding value and TTL will be updated.
	// If the key is new, a new entry will be created with the provided value and TTL.
	//
	// The TTL parameter represents the time duration in milliseconds for which the key-value pair will be valid in the cache.
	//
	// This function is thread-safe and utilizes a read-write lock to ensure concurrent access
	// to the cache's internal data structures.
	SetWithExpiry(key K, value V, ttl int)
}

// New creates a new instance of a Least Recently Used (LRU) cache with the specified size.
// It returns a pointer to an lru[K, V] instance.
func New[K comparable, V any](size int) LRU[K, V] {
	return &lru[K, V]{
		cache:  map[K]*cache[K, V]{},
		size:   size,
		length: 0,
		head:   nil,
	}
}

// New creates a new instance of a Least Recently Used (LRU) cache with the specified size.
// It returns a pointer to an lru[K, V] instance.
func NewWithExpiry[K comparable, V any](size int) LRUWithExpiry[K, V] {
	out := &lru[K, V]{
		cache:      map[K]*cache[K, V]{},
		size:       size,
		withExpiry: true,
		length:     0,
		head:       nil,
	}
	out.startCleaner()

	return out
}
