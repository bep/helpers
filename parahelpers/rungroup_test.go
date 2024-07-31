package parahelpers

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestNew(t *testing.T) {
	c := qt.New(t)

	var result int
	adder := func(ctx context.Context, i int) error {
		result += i
		return nil
	}

	g := RunGroup[int](
		context.Background(),
		GroupConfig[int]{
			Handle: adder,
		},
	)

	c.Assert(g, qt.IsNotNil)
	g.Enqueue(32)
	g.Enqueue(33)
	c.Assert(g.Wait(), qt.IsNil)
	c.Assert(result, qt.Equals, 65)
}
