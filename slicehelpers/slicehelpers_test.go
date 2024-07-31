package slicehelpers

import (
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
