// Copyright 2026 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package maphelpers

import (
	"iter"
	"sync"
)

// NewConcurrentMap creates a new ConcurrentMap.
func NewConcurrentMap[K comparable, T any]() *ConcurrentMap[K, T] {
	return &ConcurrentMap[K, T]{
		m: make(map[K]T),
	}
}

// ConcurrentMap is a thread safe map backed by a Go map.
type ConcurrentMap[K comparable, T any] struct {
	m  map[K]T
	mu sync.RWMutex
}

// Get gets the value for the given key.
// It returns the zero value of T if the key is not found.
func (m *ConcurrentMap[K, T]) Get(key K) T {
	v, _ := m.Lookup(key)
	return v
}

// Lookup looks up the given key in the map.
// It returns the value and a boolean indicating whether the key was found.
func (m *ConcurrentMap[K, T]) Lookup(key K) (T, bool) {
	m.mu.RLock()
	v, found := m.m[key]
	m.mu.RUnlock()
	return v, found
}

// GetOrCreate gets the value for the given key if it exists, or creates it if not.
func (m *ConcurrentMap[K, T]) GetOrCreate(key K, create func() (T, error)) (T, error) {
	v, found := m.Lookup(key)
	if found {
		return v, nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	v, found = m.m[key]
	if found {
		return v, nil
	}
	v, err := create()
	if err != nil {
		return v, err
	}
	m.m[key] = v
	return v, nil
}

// Set sets the given key to the given value.
func (m *ConcurrentMap[K, T]) Set(key K, value T) {
	m.mu.Lock()
	m.m[key] = value
	m.mu.Unlock()
}

// SetIfAbsent sets the given key to the given value if the key is not already present in the map.
// It returns true if the key was set, false if the key was already present.
func (m *ConcurrentMap[K, T]) SetIfAbsent(key K, value T) bool {
	m.mu.RLock()
	if _, found := m.m[key]; found {
		m.mu.RUnlock()
		return false
	}
	m.mu.RUnlock()
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.m[key]; found {
		return false
	}
	m.m[key] = value
	return true
}

// Delete deletes the given key from the map.
// It returns true if the key was found and deleted, false otherwise.
func (m *ConcurrentMap[K, T]) Delete(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.m[key]; found {
		delete(m.m, key)
		return true
	}
	return false
}

// WithReadLock executes the given function with a read lock on the map.
// Note that the map m should only be accessed within the function f.
func (m *ConcurrentMap[K, T]) WithReadLock(f func(m map[K]T) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return f(m.m)
}

// WithWriteLock executes the given function with a write lock on the map.
// Note that the map m should only be accessed within the function f.
func (m *ConcurrentMap[K, T]) WithWriteLock(f func(m map[K]T) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return f(m.m)
}

// All returns an iterator over all key/value pairs in the map.
// A read lock is held during the iteration.
func (m *ConcurrentMap[K, T]) All() iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()
		for k, v := range m.m {
			if !yield(k, v) {
				return
			}
		}
	}
}
