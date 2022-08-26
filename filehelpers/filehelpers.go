package filehelpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file.
func CopyFile(from, to string) error {
	sf, err := os.Open(from)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(to)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}
	si, err := os.Stat(from)
	if err == nil {
		err = os.Chmod(to, si.Mode())
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyDir copies a directory. Any directory or file matching the filter will be copied.
func CopyDir(from, to string, filter func(filename string) bool) error {
	fi, err := os.Stat(from)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%q is not a directory", from)
	}

	err = os.MkdirAll(to, 0777) // before umask
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(from)
	if err != nil {
		return err
	}

	if filter == nil {
		filter = func(filename string) bool {
			return true
		}
	}

	for _, entry := range entries {
		fromFilename := filepath.Join(from, entry.Name())
		toFilename := filepath.Join(to, entry.Name())
		if !filter(fromFilename) {
			continue
		}
		if entry.IsDir() {
			if err := CopyDir(fromFilename, toFilename, filter); err != nil {
				return err
			}
		} else {
			if err := CopyFile(fromFilename, toFilename); err != nil {
				return err
			}
		}

	}

	return nil
}
