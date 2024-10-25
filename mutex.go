package lru

import (
	"sync"
)

// MutexLRU is a thread-safe fixed size LRU cache.
type MutexLRU[K comparable, V any] struct {
	mut sync.RWMutex
	lru *LRU[K, V]
}

// NewWithMutex creates a new thread-safe LRU cache.
// If size is zero, the LRU has no limit
// and it's assumed that eviction is done by the caller.
func NewWithMutex[K comparable, V any](size int) *MutexLRU[K, V] {
	return &MutexLRU[K, V]{
		lru: New[K, V](size),
	}
}

// OnEvicted optionally specifies a callback function to be
// executed when an entry is purged from the lru cache.
func (m *MutexLRU[K, V]) OnEvicted(cb func(K, V)) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.lru.OnEvicted(func(k K, v V) {
		m.mut.Unlock()
		defer m.mut.Lock()
		cb(k, v)
	})
}

// Set sets a value to the lru cache.
func (m *MutexLRU[K, V]) Set(key K, value V) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.lru.Set(key, value)
}

// Get looks up a key's value from the lru cache.
func (m *MutexLRU[K, V]) Get(key K) (value V, ok bool) {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.lru.Get(key)
}

// Peek returns the key value (or undefined if not found)
// without updating the "recently used"-ness of the key.
func (m *MutexLRU[K, V]) Pick(key K) (value V, ok bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.lru.Pick(key)
}

// Remove removes the provided key from the lru cache.
func (m *MutexLRU[K, V]) Remove(key K) (value V, ok bool) {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.lru.Remove(key)
}

// RemoveOldest removes the oldest item from the cache.
func (m *MutexLRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.lru.RemoveOldest()
}

// Resize changes the lru cache size.
func (m *MutexLRU[K, V]) Resize(size int) (evicted int) {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.lru.Resize(size)
}

// Len returns the number of items in the lru cache.
func (m *MutexLRU[K, V]) Len() int {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.lru.Len()
}

// Clear purges all stored items from the lru cache.
func (m *MutexLRU[K, V]) Clear() {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.lru.Clear()
}

// Backward returns an iterator over key-value pairs in the lru cache,
// traversing it from the newest item.
func (m *MutexLRU[K, V]) All() func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		m.mut.Lock()
		defer m.mut.Unlock()
		m.lru.All()(func(key K, value V) bool {
			m.mut.Unlock()
			defer m.mut.Lock()
			return yield(key, value)
		})
	}
}

// Backward returns an iterator over key-value pairs in the lru cache,
// traversing it from the oldest item.
func (m *MutexLRU[K, V]) Backward() func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		m.mut.Lock()
		defer m.mut.Unlock()
		m.lru.Backward()(func(key K, value V) bool {
			m.mut.Unlock()
			defer m.mut.Lock()
			return yield(key, value)
		})
	}
}
