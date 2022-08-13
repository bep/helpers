package filehelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()
	abs := func(s string) string {
		return filepath.Join(tempDir, s)
	}

	c := qt.New(t)
	_, err := os.Create(abs("f1.txt"))
	c.Assert(err, qt.IsNil)
	c.Assert(CopyFile(abs("f1.txt"), abs("f2.txt")), qt.IsNil)
	fi1, err := os.Stat(abs("f1.txt"))
	c.Assert(err, qt.IsNil)
	fi2, err := os.Stat(abs("f2.txt"))
	c.Assert(err, qt.IsNil)
	c.Assert(fi1.Mode(), qt.Equals, fi2.Mode())
}

func TestCopyDir(t *testing.T) {
	tempDir := t.TempDir()
	abs := func(s string) string {
		return filepath.Join(tempDir, s)
	}
	c := qt.New(t)
	c.Assert(os.MkdirAll(abs("a/b/c"), 0755), qt.IsNil)
	_, err := os.Create(abs("a/b/c/f1.txt"))
	c.Assert(err, qt.IsNil)
	_, err = os.Create(abs("a/b/c/f2.txt"))
	c.Assert(err, qt.IsNil)
	_, err = os.Create(abs("a/b/c/f3.txt"))
	c.Assert(err, qt.IsNil)

	c.Assert(CopyDir(abs("a"), abs("b"), func(filename string) bool { return !strings.Contains(filename, "f3") }), qt.IsNil)

	for i := 1; i <= 2; i++ {
		fi1, err := os.Stat(abs(fmt.Sprintf("a/b/c/f%d.txt", i)))
		c.Assert(err, qt.IsNil)
		fi2, err := os.Stat(abs(fmt.Sprintf("b/b/c/f%d.txt", i)))
		c.Assert(err, qt.IsNil)
		c.Assert(fi1.Mode(), qt.Equals, fi2.Mode())
	}

	rdir2, _ := os.ReadDir(abs("b/b/c"))
	c.Assert(len(rdir2), qt.Equals, 2)

}
