package manager

import (
	"testing"

	"github.com/mdm-code/gsnip/snippets"
)

func TestProgramAcceptsFindCmd(t *testing.T) {
	c := snippets.NewSnippetsMap()

	s := snippets.Snippet{
		Name: "func",
		Desc: "desc",
		Body: "body",
	}
	c.Insert(s)
	m := newManager(c)
	params := []string{"func"}
	has, err := m.Execute(params...)
	want, _ := c.Find("func")
	if err != nil || has != want.Body {
		t.Error("executing find fails")
	}
}

func TestProgramAcceptsListCmd(t *testing.T) {
	c := snippets.NewSnippetsMap()
	c.Insert(snippets.Snippet{
		Name: "func",
		Desc: "simple function",
		Body: "body",
	})
	c.Insert(snippets.Snippet{
		Name: "method",
		Desc: "class method",
		Body: "body",
	})
	m := newManager(c)
	has, err := m.Execute("list")
	var want string
	listing, err := c.List()
	if err != nil {
		t.Error("failed to get a list of snippets")
	}
	for _, e := range listing {
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
	snips := snippets.NewSnippetsMap()
	m := newManager(snips)
	_, err := m.Execute("search")
	if err == nil {
		t.Error("unknown command or missing snippet does not raise an error")
	}
}
