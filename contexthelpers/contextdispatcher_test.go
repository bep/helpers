// Copyright 2024 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package contexthelpers

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestContextDispatcher(t *testing.T) {
	c := qt.New(t)

	ctx := context.Background()
	dispatcher1 := NewContextDispatcher[string](1)
	dispatcher2 := NewContextDispatcher[string](2)
	ctx = dispatcher1.Set(ctx, "testValue")
	c.Assert(dispatcher1.Get(ctx), qt.Equals, "testValue")
	c.Assert(dispatcher2.Get(ctx), qt.Equals, "")
	c.Assert(dispatcher1.Get(context.Background()), qt.Equals, "")

	value, found := dispatcher1.Lookup(ctx)
	c.Assert(found, qt.IsTrue)
	c.Assert(value, qt.Equals, "testValue")
	_, found = dispatcher1.Lookup(context.Background())
	c.Assert(found, qt.IsFalse)
	value, found = dispatcher2.Lookup(ctx)
	c.Assert(found, qt.IsFalse)
}
