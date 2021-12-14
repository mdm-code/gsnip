package snippets

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSnippetRepr(t *testing.T) {
	s := Snippet{"func", "a function", "def func(): return None"}
	want := fmt.Sprintf("startsnip %s \"%s\"\n%s\nendsnip\n\n", s.Name, s.Desc, s.Body)
	if want != s.Repr() {
		t.Errorf("want: %s; has %s", want, s.Repr())
	}
}

func TestContainerFindMethod(t *testing.T) {
	ss := NewSnippetsMap()
	ss.Insert(Snippet{
		Name: "anonfunc",
		Desc: "anonymous function in the Go programming language",
		Body: "func () {${1:body}}()",
	})
	_, err := ss.Find("anonfunc")
	if err != nil {
		t.Errorf("snippets fails to recover existing snippet")
	}
}

func TestSnippetsMapInsert(t *testing.T) {
	ss := NewSnippetsMap()
	err := ss.Insert(Snippet{"name", "desc", "body"})
	if err != nil {
		t.Error("Insert() fails to insert Snippet to map")
	}
}

func TestSnippetsMapFind(t *testing.T) {
	ss := NewSnippetsMap()
	ss.cntr["func"] = Snippet{"func", "Go function", "func ${1:name} () {}"}
	_, err := ss.Find("func")
	if err != nil {
		t.Error("existing snippet signature could not be retrieved")
	}
}

func TestSnippetsMapList(t *testing.T) {
	ss := NewSnippetsMap()
	ss.cntr = map[string]Snippet{
		"func":   {"func", "Go function", "func() {}"},
		"struct": {"struct", "Go struct", "type struct {}"},
		"map":    {"map", "Go map", "map[string]string"},
	}
	want := []string{"func\tGo function", "map\tGo map", "struct\tGo struct"}
	if has, err := ss.List(); !reflect.DeepEqual(has, want) || err != nil {
		t.Errorf("want: %v; has: %v", want, has)
	}
}

func TestSnippetsMapDelete(t *testing.T) {
	sm := NewSnippetsMap()
	sm.cntr = map[string]Snippet{
		"func":   {"func", "Go function", "func() {}"},
		"struct": {"struct", "Go struct", "type struct {}"},
		"map":    {"map", "Go map", "map[string]string"},
	}
	toDel := "map"
	sm.Delete(toDel)
	if _, err := sm.Find(toDel); err == nil {
		t.Errorf("snippet `%s` is still in map", toDel)
	}
}
