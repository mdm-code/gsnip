package manager

import (
	"fmt"

	"github.com/mdm-code/gsnip/access"
	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/signals"
	"github.com/mdm-code/gsnip/snippets"
)

type Manager struct {
	fh *access.FileHandler
	c  snippets.Container
}

func NewManager(fh *access.FileHandler) (*Manager, error) {
	parser := parsing.NewParser()
	snpts, err := parser.Parse(fh)
	if err != nil {
		return newManager(nil, nil), err
	}
	return newManager(fh, snpts), nil
}

func newManager(fh *access.FileHandler, snpts snippets.Container) *Manager {
	return &Manager{fh, snpts}
}

/* Execute a command on the snippet container.

At this moment, it is possible to perform two actions:

1. List out all snippets stored in a container
2. Retrieve the body of a searched snippet with optional replacements
*/
func (m *Manager) Execute(token signals.Token) (string, error) {
	if token.IsUnbound() {
		return "", fmt.Errorf("empty strings are unbound")
	}
	switch token.IsList() {
	case true:
		result := ""
		listing, err := m.c.List()
		if err != nil {
			return "", fmt.Errorf("failed to list snippets")
		}
		for _, s := range listing {
			result = result + s + "\n"
		}
		return result, nil
	default:
		if searched, err := m.c.Find(token.Contents()); err != nil {
			return "", fmt.Errorf("%s was not found", token.Contents())
		} else {
			return searched.Body, nil
		}
	}
}

// Reload all snippets from the source file.
func (m *Manager) Reload() error {
	_, err := m.fh.Reload()
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
