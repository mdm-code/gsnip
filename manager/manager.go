package manager

import (
	"fmt"
	"os"

	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/signals"
	"github.com/mdm-code/gsnip/snippets"
)

type Manager struct {
	c snippets.Container
}

func NewManager(fname string) (*Manager, error) {
	f, err := os.Open(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gsnipd ERROR: %s", err)
		os.Exit(1)
	}
	defer f.Close()

	parser := parsing.NewParser()
	snpts, err := parser.Parse(f)
	if err != nil {
		return newManager(nil), err
	}
	return newManager(snpts), nil
}

func newManager(snpts snippets.Container) *Manager {
	return &Manager{c: snpts}
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
		if searched, err := m.c.Find(token.Sign); err != nil {
			return "", fmt.Errorf("%s was not found", token.Sign)
		} else {
			return searched.Body, nil
		}
	}
}
