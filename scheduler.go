package lru

import "time"

// startCleaner starts a background goroutine to clean expired items from the LRU cache.
// If the cache was initialized with expiry support, this function will periodically check
// for items with expired TTL (Time To Live) and remove them from the cache.
// The cleaner runs asynchronously and is meant to be started once when the cache is created.
//
// It is safe to call this function even if the cache was not initialized with expiry support.
// In that case, this function will have no effect.
func (l *lru[K, V]) startCleaner() {
	if !l.withExpiry {
		return
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				h := l.head
				for h != nil {
					if h.ttl.Before(time.Now()) {
						l.del(h.key)
					}

					h = h.next
				}
			}
		}
	}()
}
