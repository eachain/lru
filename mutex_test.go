package lru

import "testing"

func TestMutexOnEvicted(t *testing.T) {
	lru := NewWithMutex[string, int](3)

	keys := []string{"a", "b", "c", "d"}
	lru.OnEvicted(func(s string, i int) {
		if s != keys[0] {
			t.Fatalf("evicted: %q %v", s, i)
		}
	})

	for i, k := range keys {
		lru.Set(k, i)
	}
}

func TestMutexOnEvictedSet(t *testing.T) {
	lru := NewWithMutex[string, int](3)

	keys := []string{"a", "b", "c", "d"}
	var evicted int
	lru.OnEvicted(func(key string, value int) {
		if value != evicted {
			t.Fatalf("evicted: %q %v", key, value)
		}
		evicted++
		if key == "a" {
			lru.Set(key, 999)
		}
	})

	for i, k := range keys {
		lru.Set(k, i)
	}
}
