package lru

import (
	"testing"
)

func TestOnEvicted(t *testing.T) {
	lru := New[string, int](3)

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

func TestGet(t *testing.T) {
	lru := New[string, int](3)
	lru.OnEvicted(func(s string, i int) {
		if s != "b" {
			t.Fatalf("evicted: %q %v", s, i)
		}
	})
	lru.Set("a", 1)
	lru.Set("b", 2)
	lru.Set("c", 3)
	lru.Get("a")
	lru.Set("c", 33)
	lru.Set("d", 4)
}

func TestPick(t *testing.T) {
	lru := New[string, int](3)
	lru.OnEvicted(func(s string, i int) {
		if s != "a" {
			t.Fatalf("evicted: %q %v", s, i)
		}
	})
	lru.Set("a", 1)
	lru.Set("b", 2)
	lru.Set("c", 3)
	lru.Pick("a")
	lru.Set("d", 4)
}

func TestRemove(t *testing.T) {
	lru := New[string, int](3)
	lru.Set("a", 1)
	lru.Set("b", 2)
	lru.Remove("a")
	if lru.Len() != 1 {
		t.Fatalf("len: %v", lru.Len())
	}
	if b, ok := lru.Get("b"); !ok || b != 2 {
		t.Fatalf("lru get b: %v %v", b, ok)
	}
}

func TestRemoveOldest(t *testing.T) {
	lru := New[string, int](3)
	lru.Set("a", 1)
	lru.Set("b", 2)
	key, value, ok := lru.RemoveOldest()
	if !ok {
		t.Fatalf("remove oldest failed")
	}
	if key != "a" || value != 1 {
		t.Fatalf("remove oldest: %q %v", key, value)
	}
}

func TestResize(t *testing.T) {
	lru := New[int, int](10)
	for i := 1; i <= 10; i++ {
		lru.Set(i, i*10+i)
	}

	evicted := lru.Resize(3)
	if evicted != 7 {
		t.Fatalf("resize evicted: %v", evicted)
	}
}

func TestClear(t *testing.T) {
	lru := New[string, int](3)
	lru.Set("a", 1)
	lru.Set("b", 2)
	lru.Set("c", 3)
	lru.Set("d", 4)

	lru.Clear()

	if lru.Len() != 0 {
		t.Fatalf("len after clear: %v", lru.Len())
	}
	c, ok := lru.Get("c")
	if ok || c != 0 {
		t.Fatalf("get c result: %v", c)
	}
}

func TestAll(t *testing.T) {
	lru := New[int, int](10)
	for i := 1; i <= 10; i++ {
		lru.Set(i, i)
	}

	n := 10
	lru.All()(func(key, value int) bool {
		if key != n || value != n {
			t.Fatalf("key value: %v %v", key, value)
		}
		n--
		return true
	})
}

func TestBackward(t *testing.T) {
	lru := New[int, int](10)
	for i := 1; i <= 10; i++ {
		lru.Set(i, i)
	}

	n := 1
	lru.Backward()(func(key, value int) bool {
		if key != n || value != n {
			t.Fatalf("key value: %v %v", key, value)
		}
		n++
		return true
	})
}
