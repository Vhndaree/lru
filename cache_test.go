package lru

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func listAll[K comparable, V any](in *cache[K, V]) map[K]V {
	out := map[K]V{}

	h := in
	for h.next != nil {
		out[h.key] = h.value
		h = h.next
	}
	out[h.key] = h.value

	return out
}

func TestLRU(t *testing.T) {
	t.Run("LRU without expiry", func(t *testing.T) {
		t.Run("should return list with number of size items", func(t *testing.T) {
			l := &lru[int, int]{
				cache: map[int]*cache[int, int]{},
				size:  3,
			}

			l.Set(1, 1)
			l.Set(2, 2)
			l.Set(3, 3)
			l.Set(4, 4)

			expected := map[int]int{4: 4, 3: 3, 2: 2}
			actual := listAll[int, int](l.head)

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("Expected %v; Actual = %v", expected, actual)
			}
		})

		t.Run("should return value for key", func(t *testing.T) {
			l := New[int, int](3)

			l.Set(1, 1)
			l.Set(2, 2)
			l.Set(3, 3)
			actual, ok := l.Get(2)

			if !reflect.DeepEqual(true, ok) {
				t.Errorf("Expected true; Actual = %v", ok)
			}

			if !reflect.DeepEqual(2, actual) {
				t.Errorf("Expected 2; Actual = %v", actual)
			}
		})

		t.Run("should return nil and false for undefined key", func(t *testing.T) {
			l := New[int, int](3)

			l.Set(1, 1)
			l.Set(2, 2)
			l.Set(3, 3)
			l.Del(2)
			actual, ok := l.Get(2)

			if !reflect.DeepEqual(false, ok) {
				t.Errorf("Expected false; Actual = %v", ok)
			}

			if !reflect.DeepEqual(0, actual) {
				t.Errorf("Expected 0; Actual = %v", actual)
			}
		})

		t.Run("should delete lRU item", func(t *testing.T) {
			l := &lru[int, int]{
				cache: map[int]*cache[int, int]{},
				size:  3,
			}

			l.Set(1, 1)
			l.Set(2, 2)
			l.Set(3, 3)
			l.Get(1)
			l.Set(4, 4)

			expected := map[int]int{4: 4, 3: 3, 1: 1}
			actual := listAll[int, int](l.head)

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("Expected %v; Actual = %v", expected, actual)
			}
		})
	})

	t.Run("LRU with expiry", func(t *testing.T) {
		t.Run("should clean up expired items", func(t *testing.T) {
			l := NewWithExpiry[int, int](3)

			l.SetWithExpiry(1, 1, 20000)
			l.SetWithExpiry(2, 2, 2000)
			l.SetWithExpiry(3, 3, 2000)

			time.Sleep(10 * time.Second)

			_, ok := l.Get(1)
			if !reflect.DeepEqual(true, ok) {
				t.Errorf("Expected true; Actual = %v", ok)
			}

			_, ok = l.Get(2)
			if !reflect.DeepEqual(false, ok) {
				t.Errorf("Expected false; Actual = %v", ok)
			}
		})
	})

	t.Run("should handle concurrency", func(t *testing.T) {
		count := 5
		cache := New[string, int](count)

		var wg sync.WaitGroup

		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				cache.Set(fmt.Sprintf("%d", i), i)
			}(i)
		}

		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				cache.Get(fmt.Sprintf("%d", i))
			}(i)
		}

		wg.Wait()
	})
}
