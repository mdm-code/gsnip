package fhandle

import (
	"io/ioutil"
	"os"
	"os/exec"
)

type FSBuffer struct {
	buff *os.File
}

// Create new temporary file buffer.
func (f *FSBuffer) Open() (*os.File, error) {
	tmpf, err := ioutil.TempFile("/tmp", "*.snip")
	if err != nil {
		return nil, err
	}
	f.buff = tmpf
	return f.buff, nil
}

// Close the temporary file and remove it from the file system.
func (f *FSBuffer) Close() error {
	err := f.buff.Close()
	if err != nil {
		return err
	}
	os.Remove(f.buff.Name())
	if err != nil {
		return err
	}
	return nil
}

// Vim serves the role of the user's input interface.
type Vim struct {
	Prog string
	File string
}

// Open the file in the Vim editor.
func (v *Vim) Exec() error {
	cmd := exec.Command(v.Prog, v.File)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
