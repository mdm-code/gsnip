package snippets

import (
	"reflect"
	"testing"
)

func TestMapLikeInterface(t *testing.T) {
	ss := make(Snippets)
	ss["func"] = Snippet{"func", "sample function", "xxx"}
	_, ok := ss["func"]
	if !ok {
		t.Error("Snippets cannot be used as a regular map")
	}
}

func TestSnippetsInsert(t *testing.T) {
	ss := make(Snippets)
	ok := ss.Insert(Snippet{"name", "desc", "body"})
	if !ok {
		t.Error("Insert() fails to insert Snippet to map")
	}
}

func TestSnippetsFind(t *testing.T) {
	ss := make(Snippets)
	ss["func"] = Snippet{"func", "Go function", "func ${1:name} () {}"}
	_, ok := ss.Find("func")
	if !ok {
		t.Error("existing snippet signature could not be retrieved")
	}
}

func TestSnippetsList(t *testing.T) {
	ss := Snippets{
		"func":   {"func", "Go function", "func() {}"},
		"struct": {"struct", "Go struct", "type struct {}"},
		"map":    {"map", "Go map", "map[string]string"},
	}
	want := []string{"func\tGo function", "map\tGo map", "struct\tGo struct"}
	if has := ss.List(); !reflect.DeepEqual(has, want) {
		t.Errorf("want: %v; has: %v", want, has)
	}
}
