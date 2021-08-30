package manager

import (
	"testing"

	"github.com/mdm-code/gsnip/snippets"
)

func TestNewProgramCreated(t *testing.T) {
	_, ok := NewManager(make(snippets.SnippetsMap))
	if !ok {
		t.Error("Program object cannot be instantiated")
	}
}

func TestProgramAcceptsFindCmd(t *testing.T) {
	c := snippets.SnippetsMap{
		"func": snippets.Snippet{
			Name: "func",
			Desc: "desc",
			Body: "body",
		},
	}
	m, _ := NewManager(c)
	params := []string{"func"}
	has, err := m.Execute(params...)
	want, _ := c.Find("func")
	if err != nil || has != want.Body {
		t.Error("executing find fails")
	}
}

func TestProgramAcceptsListCmd(t *testing.T) {
	c := snippets.SnippetsMap{
		"func": snippets.Snippet{
			Name: "func",
			Desc: "simple function",
			Body: "body",
		},
		"method": snippets.Snippet{
			Name: "method",
			Desc: "class method",
			Body: "body",
		},
	}
	m, _ := NewManager(c)
	has, err := m.Execute("list")
	var want string
	for _, e := range c.List() {
		want = want + e + "\n"
	}
	if err != nil {
		t.Error("failed to execute the list command")
	}
	if has != want {
		t.Errorf("has: %s; want %s", has, want)
	}
}

func TestUnrecognizedInputFails(t *testing.T) {
	m, _ := NewManager(snippets.SnippetsMap{})
	_, err := m.Execute("search")
	if err == nil {
		t.Error("unknown command or missing snippet does not raise an error")
	}
}
