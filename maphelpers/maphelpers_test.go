// Copyright 2026 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package maphelpers

import (
	"errors"
	"fmt"
	"maps"
	"sync"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestConcurrentMap(t *testing.T) {
	c := qt.New(t)

	m := NewConcurrentMap[string, int]()

	m.Set("b", 42)
	v, found := m.Lookup("b")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 42)
	v = m.Get("b")
	c.Assert(v, qt.Equals, 42)
	v, found = m.Lookup("c")
	c.Assert(found, qt.Equals, false)
	c.Assert(v, qt.Equals, 0)
	v = m.Get("c")
	c.Assert(v, qt.Equals, 0)
	v, err := m.GetOrCreate("d", func() (int, error) {
		return 100, nil
	})
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.Equals, 100)
	v, found = m.Lookup("d")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 100)

	v, err = m.GetOrCreate("d", func() (int, error) {
		return 200, nil
	})
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.Equals, 100)

	m.WithWriteLock(func(m map[string]int) error {
		m["f"] = 500
		return nil
	})
	v, found = m.Lookup("f")
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, 500)

	// Delete existing key.
	deleted := m.Delete("b")
	c.Assert(deleted, qt.Equals, true)
	_, found = m.Lookup("b")
	c.Assert(found, qt.Equals, false)

	// Delete non-existing key.
	deleted = m.Delete("nonexistent")
	c.Assert(deleted, qt.Equals, false)

	// GetOrCreate with error.
	testErr := errors.New("create failed")
	_, err = m.GetOrCreate("g", func() (int, error) {
		return 0, testErr
	})
	c.Assert(err, qt.Equals, testErr)
	_, found = m.Lookup("g")
	c.Assert(found, qt.Equals, false)

	// WithWriteLock returning error.
	err = m.WithWriteLock(func(m map[string]int) error {
		return testErr
	})
	c.Assert(err, qt.Equals, testErr)

	// All iterator.
	all := NewConcurrentMap[string, int]()
	all.Set("a", 1)
	all.Set("b", 2)
	all.Set("c", 3)
	collected := maps.Collect(all.All())
	c.Assert(collected, qt.DeepEquals, map[string]int{"a": 1, "b": 2, "c": 3})

	// All iterator with early break.
	count := 0
	for range all.All() {
		count++
		break
	}
	c.Assert(count, qt.Equals, 1)
}

func TestConcurrentMapConcurrency(t *testing.T) {
	m := NewConcurrentMap[string, int]()

	var wg sync.WaitGroup
	n := 100

	// Concurrent writes.
	for i := range n {
		wg.Go(func() {
			m.Set(fmt.Sprintf("key%d", i), i)
		})
	}

	// Concurrent reads.
	for i := range n {
		wg.Go(func() {
			m.Get(fmt.Sprintf("key%d", i))
		})
	}

	// Concurrent GetOrCreate.
	for i := range n {
		wg.Go(func() {
			m.GetOrCreate(fmt.Sprintf("shared%d", i), func() (int, error) {
				return i, nil
			})
		})
	}

	// Concurrent Delete.
	for i := range n {
		wg.Go(func() {
			m.Delete(fmt.Sprintf("key%d", i))
		})
	}

	wg.Wait()
}
