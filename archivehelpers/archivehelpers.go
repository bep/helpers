// Copyright 2024 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package archivehelpers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// TypeUnknown is the default value for an unknown type.
	TypeUnknown Type = iota
	// TypeTarGz is a tar.gz archive.
	TypeTarGz
)

// New returns a new Archiver for the given type.
func New(typ Type) (Archiver, error) {
	switch typ {
	case TypeTarGz:
		return &archivist{archiver: &tarGzExtractor{}}, nil
	default:
		return nil, fmt.Errorf("unknown type %d", typ)
	}
}

// Archiver is an interface for archiving files and directories.
type Archiver interface {
	Extracter

	// ArchiveDirectory archives the given directory into the given output stream
	// for all files matching predicate.
	// out is closed by this method.
	ArchiveDirectory(directory string, predicate func(string) bool, out io.WriteCloser) error
}

type Extracter interface {
	// Extract extracts the given archive into the given directory.
	Extract(in io.ReadCloser, targetDir string) error
}

type Type int

func (t Type) String() string {
	switch t {
	case TypeTarGz:
		return "tar.gz"
	default:
		return "unknown"
	}
}

type archiveAdder interface {
	Add(filename string, info os.FileInfo, targetPath string) error
	Close() error
}

type archiver interface {
	Extracter
	NewArchiveAdder(out io.WriteCloser) archiveAdder
}

type archivist struct {
	archiver
}

func (a *archivist) ArchiveDirectory(directory string, predicate func(string) bool, out io.WriteCloser) (err error) {
	archive := a.archiver.NewArchiveAdder(out)
	defer func() {
		closeErr := archive.Close()
		if err == nil {
			err = closeErr
		}
	}()

	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !predicate(path) {
			return nil
		}

		targetPath := strings.Trim(filepath.ToSlash(strings.TrimPrefix(path, directory)), "/")
		return archive.Add(path, info, targetPath)
	})

	return
}

type tarGzArchiver struct {
	out io.WriteCloser
	gw  *gzip.Writer
	tw  *tar.Writer
}

func (a *tarGzArchiver) Add(filename string, info os.FileInfo, targetPath string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	header, err := tar.FileInfoHeader(info, "") // TODO(bep) symlink handling?
	if err != nil {
		return err
	}
	header.Name = targetPath

	err = a.tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(a.tw, f)
	if err != nil {
		return err
	}

	return nil
}

func (a *tarGzArchiver) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	if err := a.gw.Close(); err != nil {
		return err
	}

	return a.out.Close()
}

type tarGzExtractor struct{}

func (e *tarGzExtractor) NewArchiveAdder(out io.WriteCloser) archiveAdder {
	a := &tarGzArchiver{
		out: out,
	}

	gw, _ := gzip.NewWriterLevel(out, gzip.BestCompression)
	tw := tar.NewWriter(gw)

	a.gw = gw
	a.tw = tw

	return struct {
		archiveAdder
		Extracter
	}{
		a,
		&tarGzExtractor{},
	}
}

func (a *tarGzExtractor) Extract(in io.ReadCloser, targetDir string) error {
	defer in.Close()

	gzr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil && !os.IsExist(err) {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		default:
			return fmt.Errorf("unable to untar type: %c in file %s", header.Typeflag, target)
		}
	}

	return nil
}
