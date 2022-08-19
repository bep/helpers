package envhelpers

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestExpand(t *testing.T) {
	c := qt.New(t)

	m := func(s string) string {
		return s + "-expanded"
	}

	c.Assert(Expand("", m), qt.Equals, "")
	c.Assert(Expand("no vars", m), qt.Equals, "no vars")
	c.Assert(Expand("one var: ${myvar}.", m), qt.Equals, "one var: myvar-expanded.")
	c.Assert(Expand("two vars: first:  ${first}, second ${second}.", m), qt.Equals, "two vars: first:  first-expanded, second second-expanded.")
	c.Assert(Expand("with space: ${ myvar }.", m), qt.Equals, "with space: ${ myvar }.")
	c.Assert(Expand("with special char: ${myvar&}.", m), qt.Equals, "with special char: ${myvar&}.")
	c.Assert(Expand("not without brackets: $myvar", m), qt.Equals, "not without brackets: $myvar")
	c.Assert(Expand("multiline: ${first}\n\nanother: ${second}", m), qt.Equals, "multiline: first-expanded\n\nanother: second-expanded")

	c.Assert(Expand("unquoute single: '${myvar@U}'.", m), qt.Equals, "unquoute single: myvar-expanded.")
	c.Assert(Expand("unquoute double: \"${myvar@U}\".", m), qt.Equals, "unquoute double: myvar-expanded.")
}
