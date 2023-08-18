package lru

import (
	"math"
	"testing"
	"time"
)

var n = int(math.Pow(2, 16))

func BenchmarkMYLRU(b *testing.B) {
	lr := New[int, string](n)

	b.Run("SET", func(b *testing.B) {
		for i := 0; i < n; i++ {
			lr.Set(i, time.Now().GoString())
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < n; i++ {
			lr.Get(i)
		}
	})

	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < n; i++ {
			lr.Contains(i)
		}
	})

	b.Run("Delete", func(b *testing.B) {
		for i := 0; i < n; i++ {
			lr.Del(i)
		}
	})
}
