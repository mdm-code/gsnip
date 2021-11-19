package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type FType uint8

const (
	Perm FType = iota
	Temp
)

type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Closer
	io.Seeker
}

type opener interface {
	open() (*os.File, error)
	name() string
}

type openPerm struct {
	fname string
}

type openTemp struct {
	fname string
}

func (o *openPerm) open() (*os.File, error) {
	f, err := os.OpenFile(o.fname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (o *openPerm) name() string { return o.fname }

func (o *openTemp) open() (f *os.File, err error) {
	if o.fname == "" {
		f, err = ioutil.TempFile("/tmp", "*.snip")
		if err != nil {
			return nil, err
		}
		o.fname = f.Name()
	} else {
		f, err = os.OpenFile(o.fname, os.O_APPEND|os.O_RDWR, 0644)
	}
	return f, nil
}

func (o *openTemp) name() string { return o.fname }

type FileHandler struct {
	file  ReadWriteSeekCloser
	mutex *sync.Mutex
	opener
}

func NewFileHandler(fname string, ft FType) (*FileHandler, error) {
	var o opener
	switch ft {
	case Perm:
		o = &openPerm{fname}
	case Temp:
		o = &openTemp{} // NOTE: fname is set upon calling open()
	default:
		return nil, fmt.Errorf("unsupported file type: %T", ft)
	}
	f, err := o.open()
	if err != nil {
		return nil, err
	}
	h := &FileHandler{f, &sync.Mutex{}, o}
	return h, nil
}

func (h *FileHandler) Open() (err error) {
	h.file, err = h.open()
	return
}

func (h *FileHandler) Read(p []byte) (int, error) {
	return h.file.Read(p)
}

func (h *FileHandler) Write(p []byte) (int, error) {
	return h.file.Write(p)
}

func (h *FileHandler) Close() error {
	return h.file.Close()
}

func (h *FileHandler) Seek(offset int64, whence int) (int64, error) {
	return h.file.Seek(offset, whence)
}

func (h *FileHandler) Lock() {
	h.mutex.Lock()
}

func (h *FileHandler) Unlock() {
	h.mutex.Unlock()
}

func (h *FileHandler) Reload() error {
	err := h.Close()
	if err != nil {
		return err
	}
	err = h.Open()
	h.Seek(0, io.SeekStart)
	return err
}

func (h *FileHandler) Remove() error {
	err := os.Remove(h.Name())
	return err
}

func (h *FileHandler) Name() string {
	return h.name()
}
