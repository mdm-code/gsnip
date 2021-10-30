package manager

import (
	"fmt"
	"strings"

	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/snippets"
)

type Manager struct {
	c snippets.Container
}

// Create a fresh instance of a program manager.
func NewManager(c snippets.Container) (Manager, bool) {
	return Manager{c: c}, true
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
		if searched, err := m.c.Find(cmd); err == nil {
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
		} else {
			return "", fmt.Errorf("%s was not found", cmd)
		}
	}
}
