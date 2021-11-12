package manager

import (
	"testing"

	"github.com/mdm-code/gsnip/signals"
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
	input := "func"
	interp := signals.NewInterpreter()
	tkn := interp.Eval(input)
	has, err := m.Execute(tkn)
	want, _ := c.Find("func")
	if err != nil || has != want.Body {
		t.Error("executing find fails")
	}
}

func TestProgramAcceptsListCmd(t *testing.T) {
	c, err := snippets.NewSnippetsContainer("map")
	if err != nil {
		t.Error("failed to create snippet container")
	}
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
	interp := signals.NewInterpreter()
	tkn := interp.Eval("@LIST")
	has, err := m.Execute(tkn)
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
	snips, err := snippets.NewSnippetsContainer("map")
	if err != nil {
		t.Error("failed to create snippet container")
	}
	m := newManager(snips)
	interp := signals.NewInterpreter()
	tkn := interp.Eval("search")
	_, err = m.Execute(tkn)
	if err == nil {
		t.Error("unknown command or missing snippet does not raise an error")
	}
}
