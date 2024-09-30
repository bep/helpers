// Copyright 2024 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package slicehelpers

import "sync"

// Chunk splits the slice s into n number of chunks.
func Chunk[T any](s []T, n int) [][]T {
	if len(s) == 0 {
		return nil
	}
	var partitions [][]T
	sizeDefault := len(s) / n
	sizeBig := len(s) - sizeDefault*n
	size := sizeDefault + 1
	for i, idx := 0, 0; i < n; i++ {
		if i == sizeBig {
			size--
			if size == 0 {
				break
			}
		}
		partitions = append(partitions, s[idx:idx+size])
		idx += size
	}
	return partitions
}

// Partition partitions s into slices of size size.
func Partition[T any](s []T, size int) [][]T {
	if len(s) == 0 {
		return nil
	}
	if size <= 0 {
		return nil
	}
	var partitions [][]T
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		partitions = append(partitions, s[i:end])
	}
	return partitions
}

// stack represents a stack data structure.
type stack[T any] struct {
	zero  T
	items []T
}

type StackConfig struct {
	// ThreadSafe indicates if the stack should be thread safe.
	ThreadSafe bool
}

// NewStack returns a new Stack.
func NewStack[T any](conf StackConfig) Stack[T] {
	s := &stack[T]{}
	if !conf.ThreadSafe {
		return s
	}
	return &threadSafeStack[T]{stack: s}
}

// Push adds an element to the stack.
func (s *stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Peek returns the top element of the stack.
func (s *stack[T]) Peek() T {
	if len(s.items) == 0 {
		return s.zero
	}
	return s.items[len(s.items)-1]
}

// Pop removes and returns the top element of the stack.
func (s *stack[T]) Pop() T {
	if len(s.items) == 0 {
		return s.zero
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v
}

// Drain removes all elements from the stack and returns them.
func (s *stack[T]) Drain() []T {
	items := s.items
	s.items = nil
	return items
}

// Len returns the number of elements in the stack.
func (s *stack[T]) Len() int {
	return len(s.items)
}

type Stack[T any] interface {
	// Push adds an element to the stack.
	Push(v T)
	// Pop removes and returns the top element of the stack.
	Pop() T
	// Peek returns the top element of the stack.
	Peek() T
	// Drain removes all elements from the stack and returns them.
	Drain() []T
	// Len returns the number of elements in the stack.
	Len() int
}

type threadSafeStack[T any] struct {
	*stack[T]
	mu sync.RWMutex
}

func (s *threadSafeStack[T]) Push(v T) {
	s.mu.Lock()
	s.stack.Push(v)
	s.mu.Unlock()
}

func (s *threadSafeStack[T]) Pop() T {
	s.mu.Lock()
	v := s.stack.Pop()
	s.mu.Unlock()
	return v
}

func (s *threadSafeStack[T]) Peek() T {
	s.mu.RLock()
	v := s.stack.Peek()
	s.mu.RUnlock()
	return v
}

func (s *threadSafeStack[T]) Drain() []T {
	s.mu.Lock()
	items := s.stack.Drain()
	s.mu.Unlock()
	return items
}

func (s *threadSafeStack[T]) Len() int {
	s.mu.RLock()
	l := s.stack.Len()
	s.mu.RUnlock()
	return l
}
