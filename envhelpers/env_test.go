// Copyright 2024 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package envhelpers

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestSetEnvVars(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	vars := []string{"FOO=bar", "HUGO=cool", "BAR=foo"}
	SetEnvVars(&vars, "HUGO", "rocking!", "NEW", "bar")
	c.Assert(vars, qt.DeepEquals, []string{"FOO=bar", "HUGO=rocking!", "BAR=foo", "NEW=bar"})

	key, val := SplitEnvVar("HUGO=rocks")
	c.Assert(key, qt.Equals, "HUGO")
	c.Assert(val, qt.Equals, "rocks")
}
