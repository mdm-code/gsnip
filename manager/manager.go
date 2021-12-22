package manager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mdm-code/gsnip/fs"
	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/snippets"
	"github.com/mdm-code/gsnip/stream"
)

type Manager struct {
	fh *fs.FileHandler
	c  snippets.Container
	p  *parsing.Parser
}

func NewManager(fh *fs.FileHandler) (*Manager, error) {
	parser := parsing.NewParser()
	snpts, err := parser.Parse(fh)
	if err != nil && !errors.Is(err, parsing.ErrEmptyFile) {
		return newManager(nil, nil, nil), err
	}
	return newManager(fh, snpts, &parser), nil
}

func newManager(fh *fs.FileHandler, snpts snippets.Container, p *parsing.Parser) *Manager {
	return &Manager{fh, snpts, p}
}

/* Run a server command against the snippet container.

Allowed commands:
	* @LST: list out all stored snippets
	* @FND: retrieve a snippet
	* @INS: insert a snippet to container
	* @DEL: delete a snippet
*/
func (m *Manager) Execute(msg stream.Msg) (string, error) {
	if msg.IsUnbound() {
		return "", fmt.Errorf("empty strings are unbound")
	}
	switch msg.T() {
	case stream.Lst:
		return m.list()
	case stream.Fnd:
		return m.find(string(msg.Contents()))
	case stream.Ins:
		return m.insert(string(msg.Contents()))
	case stream.Del:
		return m.delete(string(msg.Contents()))
	default:
		return "ERROR", fmt.Errorf("message kind %s is not supported", msg.TString())
	}
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
	if searched, err := m.c.Find(s); err != nil {
		return "", fmt.Errorf("%s was not found", s)
	} else {
		return searched.Body, nil
	}
}

func (m *Manager) insert(c string) (string, error) {
	reader := strings.NewReader(c)
	parsed, err := m.p.Run(reader)
	if err != nil {
		return "ERROR", err
	}
	for _, p := range parsed {
		err = m.c.Insert(p)
		if err != nil {
			return "ERROR", err
		}
	}
	snips, err := m.c.ListObj()
	if err != nil {
		return "ERROR", err
	}
	err = m.fh.Truncate(0)
	for _, s := range snips {
		m.fh.Write([]byte(s.Repr()))
	}
	err = m.Reload()
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

	err = m.Reload()
	if err != nil {
		return "ERROR", err
	}
	return "", nil
}

// Reload all snippets from the source file.
func (m *Manager) Reload() error {
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
