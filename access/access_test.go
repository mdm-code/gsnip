package access

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

func TestFileHandlerInteraction(t *testing.T) {
	dst := []byte{}
	fh := &FileHandler{&mockFile{}, &sync.Mutex{}}
	fh.Reload()
	fh.Read(dst)
	fh.Write(dst)
	fh.Seek(0, 0)
	fh.Close()
}
