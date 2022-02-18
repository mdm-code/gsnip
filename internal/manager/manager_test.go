package manager

import (
	"fmt"
	"testing"

	"github.com/mdm-code/gsnip/internal/fs"
	"github.com/mdm-code/gsnip/internal/parsing"
	"github.com/mdm-code/gsnip/internal/snippets"
	"github.com/mdm-code/gsnip/internal/stream"
)

var c snippets.Container
var p parsing.Parser

func init() {
	c, _ = snippets.NewSnippetsContainer("map")
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
	p = parsing.NewParser()
}

func TestProgramAcceptsFindCmd(t *testing.T) {
	s := snippets.Snippet{
		Name: "func",
		Desc: "desc",
		Body: "body",
	}
	c.Insert(s)
	m := newManager(&fs.FileHandler{}, c, &p)
	input := "@FND func"
	interp := stream.NewInterpreter()
	tkn := interp.Eval([]byte(input))
	has, err := m.Execute(tkn)
	want, _ := c.Find("func")
	if err != nil || has != want.Body {
		t.Error("executing find fails")
	}
}

func TestProgramAcceptsListCmd(t *testing.T) {
	m := newManager(&fs.FileHandler{}, c, &p)
	interp := stream.NewInterpreter()
	msg := "@LST"
	tkn := interp.Eval([]byte(msg))
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
	m := newManager(&fs.FileHandler{}, c, &p)
	interp := stream.NewInterpreter()
	msg := "search"
	tkn := interp.Eval([]byte(msg))
	_, err := m.Execute(tkn)
	if err == nil {
		t.Error("unknown command or missing snippet does not raise an error")
	}
}

func TestExecuteList(t *testing.T) {
	m := newManager(&fs.FileHandler{}, c, &p)
	result, err := m.list()
	if err != nil {
		t.Errorf("got %v", result)
	}
}

func TestExecuteFind(t *testing.T) {
	m := newManager(&fs.FileHandler{}, c, &p)
	result, err := m.find("func")
	if err != nil {
		t.Errorf("got: %v", result)
	}
}

func TestExecuteFindFails(t *testing.T) {
	m := newManager(&fs.FileHandler{}, c, &p)
	result, err := m.find("non-existent")
	if err == nil {
		t.Errorf("got: %v", result)
	}
}

func TestExecuteInsert(t *testing.T) {
	m := newManager(&fs.FileHandler{}, c, &p)

	// NOTE: Recover from nil pointer FileHandler.file.Write panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from: ", err)
		}
	}()

	_, err := m.insert("startsnip test \"\"\ntesting\nendsnip")

	if err != nil {
		t.Errorf("failed to insert snippet to the manager")
	}

	_, err = m.insert("bugsnig gf \"gonna fail\"\nfailing\nendbug")

	if err == nil {
		t.Errorf("managed to insert faulty-formatted snippet")
	}
}

func TestExecuteDelete(t *testing.T) {
	m := newManager(&fs.FileHandler{}, c, &p)

	// NOTE: Recover from nil pointer FileHandler.file.Write panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from:", err)
		}
	}()

	_, err := m.delete("func")

	if err != nil {
		t.Error("failed to delete a snippet: ", err)
	}
}