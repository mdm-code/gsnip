package editor

import (
	"io"
	"os"
	"os/exec"

	"github.com/mdm-code/gsnip/fs"
)

type editor struct {
	handler *fs.FileHandler
	program string
}

// NewEditor is used to create a new text editor that is capable of editing the
// underlying text file.
//
// The function takes the prog argument, which can be `vim` or `nano`, for instance.
// The second argument, fname, is a pointer to a string. A non-empty string would
// mean that a permanent file will be created. A nil pointer would mean that a
// temporary file will be created.
func NewEditor(prog string, fname *string) (*editor, error) {
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
	e := editor{fh, prog}
	return &e, nil
}

// Open the file in the text editor.
func (e *editor) Run() ([]byte, error) {
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

// Close and remove the file.
func (e *editor) Exit() error {
	err := e.handler.Close()
	if err != nil {
		return err
	}
	err = e.handler.Remove()
	return err
}