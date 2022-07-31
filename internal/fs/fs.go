package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

const (
	// Perm specifies a permanent file.
	Perm FType = iota
	// Temp specifies a temporary file.
	Temp
)

// FType represents a file type.
type FType uint8

type truncator interface {
	Truncate(size int64) error
}

type readWriteSeekCloseTruncator interface {
	io.Reader
	io.Writer
	io.Closer
	io.Seeker
	truncator
}

type opener interface {
	open() (readWriteSeekCloseTruncator, error)
	name() string
}

type openPerm struct {
	fname string
}

type openTemp struct {
	fname string
}

func (o *openPerm) open() (readWriteSeekCloseTruncator, error) {
	f, err := os.OpenFile(o.fname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (o *openPerm) name() string { return o.fname }

func (o *openTemp) open() (readWriteSeekCloseTruncator, error) {
	var err error
	var f *os.File

	if o.fname == "" {
		f, err = ioutil.TempFile("/tmp", "*.snip")
		if err != nil {
			return nil, err
		}
		o.fname = f.Name()
	} else {
		f, err = os.OpenFile(o.fname, os.O_APPEND|os.O_RDWR, 0644)
	}
	return f, err
}

func (o *openTemp) name() string { return o.fname }

// FileHandler defines a file handle structure responsible for controlling
// file access and operations.
type FileHandler struct {
	file  readWriteSeekCloseTruncator
	mutex *sync.RWMutex
	opener
}

// NewFileHandler returns a pointer to the file handler which opens either a
// temporary or a permanent file underneath.
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
	h := &FileHandler{f, &sync.RWMutex{}, o}
	return h, nil
}

// Open opens a file.
func (h *FileHandler) Open() (err error) {
	h.file, err = h.open()
	return
}

// Read reads a file.
func (h *FileHandler) Read(p []byte) (int, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.file.Read(p)
}

// Write writes bytes to a file.
func (h *FileHandler) Write(p []byte) (int, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.file.Write(p)
}

// Close closes down a file.
func (h *FileHandler) Close() error {
	return h.file.Close()
}

// Seek seeks to an offset in a file.
func (h *FileHandler) Seek(offset int64, whence int) (int64, error) {
	return h.file.Seek(offset, whence)
}

// Truncate truncates a file to size.
func (h *FileHandler) Truncate(size int64) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.file.Truncate(size)
	return
}

// Reload reloads a file.
func (h *FileHandler) Reload() error {
	err := h.Close()
	if err != nil {
		return err
	}
	err = h.Open()
	h.Seek(0, io.SeekStart)
	return err
}

// Remove attempts to remove a file from the file system.
func (h *FileHandler) Remove() error {
	err := os.Remove(h.Name())
	return err
}

// Name returns the name of the file.
func (h *FileHandler) Name() string {
	return h.name()
}
