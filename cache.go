package lru

import (
	"sync"
	"time"
)

// cache represents an item in the cache.
type cache[K comparable, V any] struct {
	key   K            // Key associated with the cache item.
	value V            // Value associated with the cache item.
	prev  *cache[K, V] // Pointer to the previous cache item.
	next  *cache[K, V] // Pointer to the next cache item.
	ttl   *time.Time   // Cache expiry time.
}

// lru represents a Least Recently Used (LRU) cache.
type lru[K comparable, V any] struct {
	cache      map[K]*cache[K, V] // Map storing cached items.
	size       int                // Maximum number of items the cache can hold.
	withExpiry bool               // Flag to enable/disable LRU with expiry.
	head       *cache[K, V]       // Head of the linked list representing the LRU order.
	tail       *cache[K, V]       // Tail of the linked list representing the LRU order.
	length     int                // Current number of items in the cache.
	sync.Mutex                    // Mutex for concurrent access.
}

// Contains checks if the provided key is present in the LRU cache.
// It returns true if the key is found in the cache, and false otherwise.
// The function does not affect the cache's state or modify any data.
func (l *lru[K, V]) Contains(key K) bool {
	_, ok := l.cache[key]
	return ok
}

// Set adds or updates a key-value pair in the LRU cache with the provided key and value.
// If the key already exists in the cache, its corresponding value will be updated.
// If the key is new, a new entry will be created with the provided value.
//
// This function is thread-safe and utilizes a read-write lock to ensure concurrent access
// to the cache's internal data structures.
//
// Example usage:
//
//	cache.Set("myKey", "myValue")
func (l *lru[K, V]) Set(key K, value V) {
	l.Mutex.Lock()
	defer l.Unlock()

	var expiry time.Time
	l.set(key, value, expiry)
}

// SetWithExpiry adds or updates a key-value pair in the LRU cache with the provided key, value, and time-to-live (TTL).
// If the key already exists in the cache, its corresponding value and TTL will be updated.
// If the key is new, a new entry will be created with the provided value and TTL.
//
// The TTL parameter represents the time duration in milliseconds for which the key-value pair will be valid in the cache.
//
// This function is thread-safe and utilizes a read-write lock to ensure concurrent access
// to the cache's internal data structures.
//
// Example usage:
//
//	cache.SetWithExpiry("myKey", "myValue", 5000) // Sets the value with a TTL of 5 seconds
func (l *lru[K, V]) SetWithExpiry(key K, value V, ttl int) {
	l.Mutex.Lock()
	defer l.Unlock()

	l.set(key, value, time.Now().Add(time.Duration(ttl)*time.Millisecond))
}

func (l *lru[K, V]) set(key K, value V, expiry time.Time) {
	// if the key value already present in the lru
	// Linked list should be re-ordered
	// Cache value also should be updated in case of change
	if c, ok := l.cache[key]; ok {
		if c != l.head {
			if c.prev == nil {
				c.next.prev = nil
				l.head = c.next
			} else if c.next == nil {
				c.prev.next = nil
				l.tail = c.prev
			} else {
				c.next.prev = c.prev
				c.prev.next = c.next
			}
		}

		c.prev = nil
		c.next = l.head
		c.value = value
		c.ttl = &expiry

		l.head = c
		l.cache[key] = c
		return
	}

	// if lru length tries to exceed the capacity
	// drop last list/ which is least used cache
	if l.length >= l.size {
		l.del(l.tail.key)
	}

	c := &cache[K, V]{key: key, value: value, ttl: &expiry, next: l.head, prev: nil}

	if l.head == nil {
		l.tail = c
	} else {
		l.head.prev = c
	}

	l.head = c
	l.cache[key] = c
	l.length++
}

// Get retrieves the value associated with the provided key from the LRU cache.
// If the key exists in the cache, its corresponding value is returned along with a boolean true.
// If the key is not found in the cache, an empty value and boolean false are returned.
//
// The Get operation updates the order of items in the cache to reflect the most recently accessed item.
// If the item exists, it is moved to the head of the cache to prioritize recently accessed items.
func (l *lru[K, V]) Get(key K) (V, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if c, ok := l.cache[key]; ok {
		// if it was head do nothing just return value
		if c.prev == nil {
			return c.value, true
		}

		// if it was tail assign last item as tail and move it to head and link to previous head node
		if c.next == nil {
			c.prev.next = nil
			l.tail = c.prev
			c.prev = nil
			c.next = l.head
			l.head.prev = c
			l.head = c

			return c.value, true
		}

		// if it was neither head not tail then link its prev node to next node and move found value to head
		c.prev.next = c.next
		c.next.prev = c.prev
		c.next = l.head
		c.prev = nil
		l.head = c

		return c.value, true
	}

	var emptyVal V
	return emptyVal, false
}

// Del removes the key-value pair associated with the provided key from the LRU cache.
// If the key is found and the removal is successful, the function returns true.
// If the key is not found, it returns false.
//
// The function also adjusts the internal linked list of cache items to maintain its order.
// If the removed item was the head or tail of the list, appropriate adjustments are made.
// The deleted item's memory is released for garbage collection.
func (l *lru[K, V]) Del(key K) bool {
	l.Mutex.Lock()
	defer l.Unlock()

	return l.del(key)
}

func (l *lru[K, V]) del(key K) bool {
	if !l.Contains(key) {
		return false
	}

	c := l.cache[key]
	if c.prev == nil && c.next == nil {
		l.head = nil
		l.tail = nil
	} else {
		if c.prev == nil {
			c.next.prev = nil
			l.head = c.next
		} else if c.next == nil {
			c.prev.next = nil
			l.tail = c.prev
		} else {
			c.next.prev = c.prev
			c.prev.next = c.next
		}
	}

	delete(l.cache, key)
	l.length--
	c = nil

	return true
}
