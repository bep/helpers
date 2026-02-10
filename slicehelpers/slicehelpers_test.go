// Copyright 2024 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package slicehelpers

import (
	"sync"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestChunk(t *testing.T) {
	c := qt.New(t)
	c.Assert(Chunk(
		[]int{1, 2, 3, 4, 5}, 2),
		qt.DeepEquals,
		[][]int{
			{1, 2, 3},
			{4, 5},
		},
	)

	c.Assert(Chunk(
		[]int{1, 2}, 3),
		qt.DeepEquals,
		[][]int{
			{1},
			{2},
		},
	)

	c.Assert(Chunk(
		[]int{1}, 2),
		qt.DeepEquals,
		[][]int{
			{1},
		},
	)

	c.Assert(Chunk(
		[]int{}, 2),
		qt.IsNil,
	)

	c.Assert(
		Chunk([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}, 3),
		qt.DeepEquals,
		[][]string{
			{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
			{"j", "k", "l", "m", "n", "o", "p", "q", "r"},
			{"s", "t", "u", "v", "w", "x", "y", "z"},
		},
	)

	c.Assert(
		Chunk([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}, 7),
		qt.DeepEquals,
		[][]string{
			{"a", "b", "c", "d"},
			{"e", "f", "g", "h"},
			{"i", "j", "k", "l"},
			{"m", "n", "o", "p"},
			{"q", "r", "s", "t"},
			{"u", "v", "w"},
			{"x", "y", "z"},
		},
	)
}

func TestPartition(t *testing.T) {
	c := qt.New(t)

	c.Assert(
		Partition([]int{1, 2, 3, 4, 5}, 2),
		qt.DeepEquals,
		[][]int{
			{1, 2},
			{3, 4},
			{5},
		},
	)

	c.Assert(
		Partition([]int{1, 2, 3, 4, 5}, 3),
		qt.DeepEquals,
		[][]int{
			{1, 2, 3},
			{4, 5},
		},
	)

	c.Assert(
		Partition([]int{1, 2, 3, 4, 5}, 1),
		qt.DeepEquals,
		[][]int{
			{1},
			{2},
			{3},
			{4},
			{5},
		},
	)

	c.Assert(
		Partition([]int{1, 2, 3, 4, 5}, 5),
		qt.DeepEquals,
		[][]int{
			{1, 2, 3, 4, 5},
		},
	)

	c.Assert(
		Partition([]int{1, 2, 3, 4, 5}, 6),
		qt.DeepEquals,
		[][]int{
			{1, 2, 3, 4, 5},
		},
	)

	c.Assert(
		Partition([]int{1, 2, 3, 4, 5}, 7),
		qt.DeepEquals,
		[][]int{
			{1, 2, 3, 4, 5},
		},
	)

	c.Assert(
		Partition([]int{1, 2, 3, 4, 5}, 0),
		qt.IsNil,
	)

	c.Assert(
		Partition([]int{}, 2),
		qt.IsNil,
	)
}

func TestStack(t *testing.T) {
	c := qt.New(t)

	s := NewStack[int](StackConfig{})

	c.Assert(s.Peek(), qt.Equals, 0)
	c.Assert(s.Pop(), qt.Equals, 0)
	s.Push(1)
	c.Assert(s.Peek(), qt.Equals, 1)
	c.Assert(s.Pop(), qt.Equals, 1)
	c.Assert(s.Pop(), qt.Equals, 0)
	c.Assert(s.Peek(), qt.Equals, 0)

	s.Push(2)
	s.Push(3)
	c.Assert(s.Len(), qt.Equals, 2)
	c.Assert(s.Peek(), qt.Equals, 3)
	c.Assert(s.Pop(), qt.Equals, 3)
	c.Assert(s.Pop(), qt.Equals, 2)
	c.Assert(s.Pop(), qt.Equals, 0)
	c.Assert(s.Peek(), qt.Equals, 0)

	s.Push(4)
	s.Push(5)
	c.Assert(s.Drain(), qt.DeepEquals, []int{4, 5})
	c.Assert(s.Len(), qt.Equals, 0)
}

func TestStackThreadSafe(t *testing.T) {
	s := NewStack[int](StackConfig{ThreadSafe: true})

	var wg sync.WaitGroup

	for range 20 {
		wg.Add(3)
		go func() {
			defer wg.Done()
			for i := range 100 {
				s.Push(i)
				s.Len()
				s.Peek()
			}
		}()

		go func() {
			defer wg.Done()
			for range 50 {
				s.Pop()
			}
		}()

		go func() {
			defer wg.Done()
			for range 50 {
				s.Drain()
			}
		}()

	}

	wg.Wait()
}
