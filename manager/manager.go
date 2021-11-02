package manager

import (
	"fmt"
	"os"
	"strings"

	"github.com/mdm-code/gsnip/parsing"
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

func newManager(snpts *snippets.SnippetsMap) *Manager {
	return &Manager{c: snpts}
}

/* Execute a command on the snippet container.

At this moment, it is possible to perform two actions:

1. List out all snippets stored in a container
2. Retrieve the body of a searched snippet with optional replacements
*/
func (m *Manager) Execute(params ...string) (string, error) {
	var result string
	if len(params) == 0 {
		return "", fmt.Errorf("there is nothing to find")
	}
	cmd := strings.ToLower(params[0])
	if parsing.IsCommand(cmd) {
		if strings.ToLower(cmd) == "list" {
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
		return "", fmt.Errorf("unimplemented command")
	} else {
		if searched, err := m.c.Find(cmd); err != nil {
			return "", fmt.Errorf("%s was not found", cmd)
		} else {
			pat := `\${[0-9]+:\w*}`
			var repls []string
			if len(params) > 1 {
				repls = params[1:]
			}
			var ok bool
			result, ok = parsing.Replace(searched.Body, pat, repls...)
			if !ok {
				return "", fmt.Errorf("failed to compile regex pattern: %s", pat)
			}
			return result, nil
		}
	}
}
