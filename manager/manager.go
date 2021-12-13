package manager

import (
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
	if err != nil {
		return newManager(nil, nil, nil), err
	}
	return newManager(fh, snpts, &parser), nil
}

func newManager(fh *fs.FileHandler, snpts snippets.Container, p *parsing.Parser) *Manager {
	return &Manager{fh, snpts, p}
}

// TODO: Replace each case scope with a method.
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
		result := ""
		listing, err := m.c.List()
		if err != nil {
			return "", fmt.Errorf("failed to list snippets")
		}
		for _, s := range listing {
			result = result + s + "\n"
		}
		return result, nil
	case stream.Fnd:
		if searched, err := m.c.Find(string(msg.Contents())); err != nil {

			return "", fmt.Errorf("%s was not found", string(msg.Contents()))
		} else {
			return searched.Body, nil
		}
	case stream.Ins:
		reader := strings.NewReader(string(msg.Contents()))
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
	case stream.Del:
		m.c.Delete(string(msg.Contents()))
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
	default:
		return "ERROR", fmt.Errorf("message kind %s is not supported", msg.TString())
	}
}

// Reload all snippets from the source file.
func (m *Manager) Reload() error {
	err := m.fh.Reload()
	if err != nil {
		return err
	}
	parser := parsing.NewParser()
	snpts, err := parser.Parse(m.fh)
	if err != nil {
		return err
	}
	m.c = snpts
	return nil
}
