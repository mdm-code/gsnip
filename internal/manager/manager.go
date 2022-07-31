package manager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mdm-code/gsnip/internal/fs"
	"github.com/mdm-code/gsnip/internal/parsing"
	"github.com/mdm-code/gsnip/internal/snippets"
	"github.com/mdm-code/gsnip/internal/stream"
)

// Manager integrates operations on snippets stored in a file.
type Manager struct {
	fh      *fs.FileHandler
	c       snippets.Container
	p       *parsing.Parser
	actions map[stream.Opcode]interface{}
}

// NewManager creates a pointer to a Manager instance for a given file handle.
func NewManager(fh *fs.FileHandler) (*Manager, error) {
	parser := parsing.NewParser()
	snpts, err := parser.Parse(fh)
	actions := map[stream.Opcode]interface{}{
		stream.Find:   (*Manager).find,
		stream.Insert: (*Manager).insert,
		stream.Delete: (*Manager).delete,
		stream.Reload: (*Manager).reload,
		stream.List:   (*Manager).list,
	}
	if err != nil && !errors.Is(err, parsing.ErrEmptyFile) {
		return newManager(nil, nil, nil, actions), err
	}
	return newManager(fh, snpts, &parser, actions), nil
}

func newManager(fh *fs.FileHandler, snpts snippets.Container, p *parsing.Parser, actns map[stream.Opcode]interface{}) *Manager {
	return &Manager{fh, snpts, p, actns}
}

// Execute runs a server command against the snippet container.
//
// Allowed commands:
// 	* List out all stored snippets
// 	* Find a single snippet
// 	* Insert a snippet to the container
// 	* Delete a snippet from the container
//  * Reload the snippet container
func (m *Manager) Execute(request stream.Request, reply *stream.Reply) error {
	var body string
	var err error

	op, ok := m.actions[request.Operation]

	if !ok {
		err = fmt.Errorf("request %v not supported", request.Operation)
	}

	if ok {
		switch f := op.(type) {
		case func(*Manager) (string, error):
			body, err = f(m)
		case func(*Manager) error:
			err = f(m)
		case func(*Manager, string) (string, error):
			body, err = f(m, string(request.Body))
		}
	}

	if err != nil {
		reply.Result = stream.Failure
	} else {
		reply.Result = stream.Success
		reply.Body = []byte(body)
	}
	return err
}

func (m *Manager) list() (string, error) {
	result := ""
	listing, err := m.c.List()
	if err != nil {
		return "", fmt.Errorf("failed to list snippets")
	}
	for _, s := range listing {
		result = result + s + "\n"
	}
	return result, nil
}

func (m *Manager) find(s string) (string, error) {
	var searched snippets.Snippet
	var err error
	if searched, err = m.c.Find(s); err != nil {
		return "", fmt.Errorf("%s was not found", s)
	}
	return searched.Body, nil
}

func (m *Manager) insert(contents string) (string, error) {
	reader := strings.NewReader(contents)
	container, err := m.p.Parse(reader)
	if err != nil {
		return "ERROR", err
	}

	snips, err := container.ListObj()
	if err != nil {
		return "ERROR", err
	}

	for _, p := range snips {
		err = m.c.Insert(p)
		if err != nil {
			return "ERROR", err
		}
	}

	// NOTE: Rewrite file contents
	snips, err = m.c.ListObj()
	if err != nil {
		return "ERROR", err
	}
	err = m.fh.Truncate(0)
	for _, s := range snips {
		m.fh.Write([]byte(s.Repr()))
	}
	err = m.reload()
	if err != nil {
		return "ERROR", err
	}
	return "", nil
}

func (m *Manager) delete(s string) (string, error) {
	m.c.Delete(s)
	snips, err := m.c.ListObj()
	if err != nil {
		return "ERROR", err
	}
	err = m.fh.Truncate(0)
	for _, s := range snips {
		m.fh.Write([]byte(s.Repr()))
	}

	err = m.reload()
	if err != nil {
		return "ERROR", err
	}
	return "", nil
}

func (m *Manager) reload() error {
	err := m.fh.Reload()
	if err != nil {
		return err
	}
	parser := parsing.NewParser()
	snpts, err := parser.Parse(m.fh)
	if err != nil && !errors.Is(err, parsing.ErrEmptyFile) {
		return err
	}
	m.c = snpts
	return nil
}
