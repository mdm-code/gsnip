package editor

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mdm-code/gsnip/internal/fs"
)

// Editor defines the structure holding a handle to a file and the name of the
// program to open it in.
type Editor struct {
	handler *fs.FileHandler
	program string
}

// NewEditor is used to create a new text editor that is capable of editing the
// underlying text file.
//
// Fname, is a pointer to a string. A non-empty string would mean that a
// permanent file will be created. A nil pointer would mean that a temporary
// file will be created.
func NewEditor(fname *string) (*Editor, error) {
	var fh *fs.FileHandler
	var err error
	if fname == nil {
		fh, err = fs.NewFileHandler("", fs.Temp)
	} else {
		fh, err = fs.NewFileHandler(*fname, fs.Perm)
	}
	if err != nil {
		return nil, err
	}

	prog, ok := os.LookupEnv("EDITOR")
	if !ok {
		return nil, fmt.Errorf("$EDITOR is not set.")
	}

	e := Editor{fh, prog}
	return &e, nil
}

// Run opens the file in the text editor.
func (e *Editor) Run() ([]byte, error) {
	cmd := exec.Command(e.program, e.handler.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	err = e.handler.Reload()
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(e.handler)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Exit closes and removes the file.
func (e *Editor) Exit() error {
	err := e.handler.Close()
	if err != nil {
		return err
	}
	err = e.handler.Remove()
	return err
}
