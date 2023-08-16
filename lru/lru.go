package lru

import "sync"

// cache is a generic struct representing an item in a cache.
type cache[K comparable, V any] struct {
	key   K            // The key associated with the cache item.
	value V            // The value associated with the cache item.
	prev  *cache[K, V] // Pointer to the previous cache item.
	next  *cache[K, V] // Pointer to the next cache item.
}

// lru is a generic struct representing a Least Recently Used (LRU) cache.
type lru[K comparable, V any] struct {
	cache        map[K]*cache[K, V] // The map storing cached items.
	size         int                // The maximum number of items the cache can hold.
	head         *cache[K, V]       // The head of the linked list representing the LRU order.
	tail         *cache[K, V]       // The tail of the linked list representing the LRU order.
	length       int                // The current number of items in the cache.
	sync.RWMutex                    // A mutex for concurrent access.
}

// New creates a new instance of a Least Recently Used (LRU) cache with the specified size.
// It returns a pointer to an lru[K, V] instance.
func New[K comparable, V any](size int) *lru[K, V] {
	return &lru[K, V]{
		cache:  map[K]*cache[K, V]{},
		size:   size,
		length: 0,
		head:   nil,
	}
}

func (l *lru[K, V]) Contains(key K) bool {
	_, ok := l.cache[key]
	return ok
}

func (l *lru[K, V]) Set(key K, value V) {
	l.RWMutex.RLock()
	defer l.RUnlock()
	// if the key value already present in the lru
	// Linked list should be re-ordered
	// Cache value also should be updated in case of change
	if l.Contains(key) {
		c := l.cache[key]
		// key is at head
		if c.prev == nil {
			c.next.prev = nil
			l.head = c.next
		} else if c.next == nil { // key is at tail
			c.prev.next = nil
			l.tail = c.prev
		} else {
			c.next.prev = c.prev
			c.prev.next = c.next
		}

		c.prev = nil
		c.next = l.head
		c.value = value

		l.head = c
		l.cache[key] = c
		return
	}

	// if lru length tries to exceed the capacity
	// drop last list/ which is least used cache
	if l.length >= l.size {
		l.del(l.tail.key)
	}

	if l.head == nil {
		c := &cache[K, V]{key: key, value: value, next: nil, prev: nil}
		l.head = c
		l.tail = c
		l.cache[key] = c
		l.length++
	} else {
		c := &cache[K, V]{key: key, value: value, next: l.head, prev: nil}
		l.head.prev = c
		l.head = c
		l.cache[key] = c
		l.length++
	}
}

func (l *lru[K, V]) Get(key K) (V, bool) {
	l.RWMutex.RLock()
	defer l.RWMutex.RUnlock()

	if l.Contains(key) {
		c := l.cache[key]

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

func (l *lru[K, V]) Del(key K) bool {
	l.RWMutex.RLock()
	defer l.RUnlock()

	return l.del(key)
}

func (l *lru[K, V]) del(key K) bool {
	if !l.Contains(key) {
		return false
	}

	c := l.cache[key]
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

	delete(l.cache, key)
	l.length--
	c = nil

	return true
}

func (l *lru[K, V]) ListAll() map[K]V {
	out := map[K]V{}

	h := l.head
	for h.next != nil {
		out[h.key] = h.value
		h = h.next
	}
	out[h.key] = h.value

	return out
}
