// Copyright 2024 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package archivehelpers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestArchive(t *testing.T) {
	c := qt.New(t)

	_, err := New(32)
	c.Assert(err, qt.IsNotNil)

	c.Assert(TypeTarGz.String(), qt.Equals, "tar.gz")
	c.Assert(TypeUnknown.String(), qt.Equals, "unknown")

	tempDir := t.TempDir()

	for _, tp := range []Type{TypeTarGz} {

		archiveFilename := filepath.Join(tempDir, "myarchive1."+tp.String())
		f, err := os.Create(archiveFilename)
		c.Assert(err, qt.IsNil)

		a, err := New(tp)
		c.Assert(err, qt.IsNil)

		sourceDir := t.TempDir()
		subDir := filepath.Join(sourceDir, "subdir")
		c.Assert(os.MkdirAll(subDir, 0o755), qt.IsNil)
		c.Assert(os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("hello"), 0o644), qt.IsNil)
		c.Assert(os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("world"), 0o643), qt.IsNil)

		matchAll := func(string) bool { return true }
		c.Assert(a.ArchiveDirectory(sourceDir, matchAll, f), qt.IsNil)

		assertArchive := func(ararchiveFilename string, predicate func(string) bool) {
			resultDir := t.TempDir()
			f, err = os.Open(archiveFilename)
			c.Assert(err, qt.IsNil)
			c.Assert(a.Extract(f, resultDir), qt.IsNil)

			dirList := func(dirname string) string {
				var sb strings.Builder
				err := filepath.WalkDir(dirname, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.IsDir() {
						return nil
					}

					if !predicate(path) {
						return nil
					}

					fi, err := d.Info()
					if err != nil {
						return err
					}

					sb.WriteString(fmt.Sprintf("%s %04o %s\n", fi.Mode(), fi.Mode().Perm(), fi.Name()))
					return nil
				})

				c.Assert(err, qt.IsNil)
				return sb.String()
			}
			dirList1 := dirList(sourceDir)
			dirList2 := dirList(resultDir)
			c.Assert(dirList1, qt.Equals, dirList2)
		}

		assertArchive(archiveFilename, matchAll)

		archiveFilename = filepath.Join(tempDir, "myarchive2."+tp.String())
		f, err = os.Create(archiveFilename)
		matchSome := func(s string) bool {
			return filepath.Base(s) == "file1.txt"
		}
		c.Assert(a.ArchiveDirectory(sourceDir, matchSome, f), qt.IsNil)
		assertArchive(archiveFilename, matchSome)
	}
}
