package fs

import (
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
		&sync.Mutex{},
		&mockOpener{},
	}
	fh.Read(dst)
	fh.Write(dst)
	fh.Seek(0, 0)
	fh.Close()
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
