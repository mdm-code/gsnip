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

func NewEditor(prog, fname string) (*editor, error) {
	fh, err := fs.NewFileHandler(fname, fs.Temp)
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
