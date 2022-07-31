package fs

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

type mockFile struct {
	readWriteSeekCloseTruncator
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

func (m *mockOpener) open() (readWriteSeekCloseTruncator, error) {
	return nil, nil
}

func (m *mockOpener) name() string {
	return "mockName"
}

func TestFileHandlerInteraction(t *testing.T) {
	dst := []byte{}
	fh := &FileHandler{
		&mockFile{},
		&sync.RWMutex{},
		&mockOpener{},
	}
	fh.Read(dst)
	fh.Write(dst)
	fh.Seek(0, 0)
	fh.Close()
}

func TestFileHandlerLocking(t *testing.T) {
	fh := &FileHandler{
		&mockFile{},
		&sync.RWMutex{},
		&mockOpener{},
	}
	fh.mutex.Lock()
	fh.mutex.Unlock()
}

func TestFileHandlerTruncating(t *testing.T) {
	fh := &FileHandler{
		&mockFile{},
		&sync.RWMutex{},
		&mockOpener{},
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from: ", err)
		}
	}()
	err := fh.Open()
	if err != nil {
		t.Errorf("unexpected error was raised: %s", err)
	}
}

func TestFileHandlerOpen(t *testing.T) {
	fh := &FileHandler{
		&mockFile{},
		&sync.RWMutex{},
		&mockOpener{},
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from: ", err)
		}
	}()
	fh.Truncate(0)
}

func TestFileHandlerReload(t *testing.T) {
	fh := &FileHandler{
		&mockFile{},
		&sync.RWMutex{},
		&mockOpener{},
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from: ", err)
		}
	}()
	err := fh.Reload()
	if err != nil {
		t.Errorf("error was raised: %s", err)
	}
}

// Test if files are opened correctly.
func TestOpen(t *testing.T) {
	data := []struct {
		name   string
		opener opener
		closer func(string)
	}{
		{
			"perm",
			&openPerm{fname: "testfile"},
			func(fname string) { os.Remove(fname) },
		},
		{
			"temp",
			&openTemp{},
			func(fname string) { os.Remove(fname) },
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_, err := d.opener.open()
			_, err = d.opener.open()
			if err != nil {
				t.Errorf("failed to open %v", d.opener)
			}

			// NOTE: Could not defer it because the name is evaluated later
			d.closer(d.opener.name())
		})
	}
}

// Test if the file opener fails when the name contains illegal characters.
func TestOpenFail(t *testing.T) {
	wrongFname := "/failed"
	data := []struct {
		name   string
		opener opener
		closer func(string)
	}{
		{
			"perm",
			&openPerm{wrongFname},
			func(fname string) { os.Remove(fname) },
		},
		{
			"temp",
			&openTemp{wrongFname},
			func(fname string) { os.Remove(fname) },
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_, err := d.opener.open()
			_, err = d.opener.open()
			if err == nil {
				t.Errorf("failed to open %v", d.opener)
			}
			d.closer(d.opener.name())
		})
	}
}

// Make sure that file handlers are created properly.
func TestFileHandlerCreation(t *testing.T) {
	data := []struct {
		name  string
		fname string
		ft    FType
	}{
		{"perm", "testfile", Perm},
		{"temp", "", Temp},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			h, err := NewFileHandler(d.fname, d.ft)
			if err != nil {
				t.Errorf("want: %v; has: %v", nil, err)
			}
			h.Remove()
		})
	}
}

// Check if the file handler fails on illegal file name.
func TestTestFileHandlerWrongFName(t *testing.T) {
	_, err := NewFileHandler("/failed", Perm)
	if err == nil {
		t.Error("expected an error caused by wrong file name")
	}
}

// Verify if the wrong file type causes an error.
func TestFileHandlerWrongType(t *testing.T) {
	_, err := NewFileHandler("", FType(3))
	if err == nil {
		t.Error("expected an error caused by worng file type")
	}
}
