package testutil

import (
	"bytes"
	"io"
	"io/fs"
	"os"
)

type FakeFileSystem struct {
	Files map[string][]byte
}

type fakeFile struct {
	fs.File
	content io.Reader
}

func (f FakeFileSystem) Open(name string) (fs.File, error) {
	content, found := f.Files[name]
	if !found {
		return nil, os.ErrNotExist
	}

	return &fakeFile{
		content: bytes.NewBuffer(content),
	}, nil
}

func (f *fakeFile) Read(p []byte) (n int, err error) {
	return f.content.Read(p)
}

func (f FakeFileSystem) ReadFile(name string) ([]byte, error) {
	content, found := f.Files[name]
	if !found {
		return nil, os.ErrNotExist
	}

	return content, nil
}

func (f *fakeFile) Close() error {
	return nil
}
