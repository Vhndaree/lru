package lru

import (
	"reflect"
	"testing"
)

func TestLRU(t *testing.T) {
	t.Run("should return list with number of size items", func(t *testing.T) {
		l := New[int, int](3)

		l.Set(1, 1)
		l.Set(2, 2)
		l.Set(3, 3)
		l.Set(4, 4)

		expected := map[int]int{4: 4, 3: 3, 2: 2}
		actual := l.ListAll()

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
		l := New[int, int](3)

		l.Set(1, 1)
		l.Set(2, 2)
		l.Set(3, 3)
		l.Get(1)
		l.Set(4, 4)

		expected := map[int]int{4: 4, 3: 3, 1: 1}
		actual := l.ListAll()

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %v; Actual = %v", expected, actual)
		}
	})
}
