package lru

import "container/list"

type item[K comparable, V any] struct {
	key   K
	value V
}

// LRU is a fixed size LRU cache.
type LRU[K comparable, V any] struct {
	elem    map[K]*list.Element // *item[K, V]
	items   *list.List          // *item[K, V]
	size    int
	evicted func(K, V)
}

// New creates a new LRU cache.
// If size is zero, the LRU has no limit
// and it's assumed that eviction is done by the caller.
func New[K comparable, V any](size int) *LRU[K, V] {
	return &LRU[K, V]{
		elem:  make(map[K]*list.Element),
		items: list.New(),
		size:  size,
	}
}

// OnEvicted optionally specifies a callback function to be
// executed when an entry is purged from the lru cache.
func (lru *LRU[K, V]) OnEvicted(cb func(K, V)) {
	lru.evicted = cb
}

// Set sets a value to the lru cache.
func (lru *LRU[K, V]) Set(key K, value V) {
	elem := lru.elem[key]
	if elem != nil {
		lru.items.MoveToFront(elem)
		elem.Value.(*item[K, V]).value = value
	} else {
		lru.elem[key] = lru.items.PushFront(&item[K, V]{key: key, value: value})
		if lru.size > 0 && lru.items.Len() > lru.size {
			lru.RemoveOldest()
		}
	}
}

// Get looks up a key's value from the lru cache.
func (lru *LRU[K, V]) Get(key K) (value V, ok bool) {
	elem := lru.elem[key]
	if elem != nil {
		lru.items.MoveToFront(elem)
		return elem.Value.(*item[K, V]).value, true
	}
	return
}

// Peek returns the key value (or undefined if not found)
// without updating the "recently used"-ness of the key.
func (lru *LRU[K, V]) Pick(key K) (value V, ok bool) {
	elem := lru.elem[key]
	if elem != nil {
		return elem.Value.(*item[K, V]).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (lru *LRU[K, V]) Remove(key K) (value V, ok bool) {
	elem := lru.elem[key]
	if elem != nil {
		item := elem.Value.(*item[K, V])
		elem.Value = nil

		delete(lru.elem, item.key)
		lru.items.Remove(elem)

		if lru.evicted != nil {
			lru.evicted(item.key, item.value)
		}
		return item.value, true
	}
	return
}

// RemoveOldest removes the oldest item from the cache.
func (lru *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	elem := lru.items.Back()
	if elem == nil {
		return
	}

	item := elem.Value.(*item[K, V])
	elem.Value = nil

	delete(lru.elem, item.key)
	lru.items.Remove(elem)

	if lru.evicted != nil {
		lru.evicted(item.key, item.value)
	}
	return item.key, item.value, true
}

// Resize changes the lru cache size.
func (lru *LRU[K, V]) Resize(size int) (evicted int) {
	lru.size = size
	for lru.Len() > size {
		lru.RemoveOldest()
		evicted++
	}
	return
}

// Len returns the number of items in the lru cache.
func (lru *LRU[K, V]) Len() int {
	return lru.items.Len()
}

// Clear purges all stored items from the lru cache.
func (lru *LRU[K, V]) Clear() {
	for lru.Len() > 0 {
		lru.RemoveOldest()
	}
}

// Backward returns an iterator over key-value pairs in the lru cache,
// traversing it from the newest item.
func (lru *LRU[K, V]) All() func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		var next *list.Element
		for elem := lru.items.Front(); elem != nil; elem = next {
			next = elem.Next()
			item := elem.Value.(*item[K, V])
			if !yield(item.key, item.value) {
				return
			}
		}
	}
}

// Backward returns an iterator over key-value pairs in the lru cache,
// traversing it from the oldest item.
func (lru *LRU[K, V]) Backward() func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		var prev *list.Element
		for elem := lru.items.Back(); elem != nil; elem = prev {
			prev = elem.Prev()
			item := elem.Value.(*item[K, V])
			if !yield(item.key, item.value) {
				return
			}
		}
	}
}
