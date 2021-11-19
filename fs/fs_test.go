package fs

import (
	"os"
	"sync"
	"testing"
)

type mockFile struct {
	os.File
}

func (m mockFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (m mockFile) Read(p []byte) (int, error) {
	return 0, nil
}

func (m mockFile) Write(p []byte) (int, error) {
	return 0, nil
}

func (m mockFile) Close() error { return nil }

type mockOpener struct {
	fname string
}

func (m *mockOpener) open() (*os.File, error) {
	return nil, nil
}

func (m *mockOpener) name() string {
	return "mockName"
}

func TestFileHandlerInteraction(t *testing.T) {
	dst := []byte{}
	fh := &FileHandler{
		&mockFile{},
		&sync.Mutex{},
		&mockOpener{},
	}
	fh.Reload()
	fh.Read(dst)
	fh.Write(dst)
	fh.Seek(0, 0)
	fh.Close()
}
