package access

import (
	"io"
	"os"
	"sync"
)

type ReadWriteSeekCloser interface {
	io.Writer
	io.Reader
	io.Seeker
	io.Closer
}

type FileHandler struct {
	ReadWriteSeekCloser
	*sync.Mutex
}

// TODO: The type of file handler is specified with a numerical constant
func NewFileHandler(fname string) (*FileHandler, error) {
	// TODO: Editor should be given a write-only file -- make it a proper factory
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &FileHandler{f, &sync.Mutex{}}, nil
}

func (f *FileHandler) Reload() (ret int64, err error) {
	ret, err = f.Seek(0, io.SeekStart)
	return
}
