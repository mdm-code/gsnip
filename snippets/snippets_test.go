package snippets

import (
	"reflect"
	"testing"
)

func TestContainerFindMethod(t *testing.T) {
	var ss Container = make(SnippetsMap)
	ss.Insert(Snippet{
		Name: "func",
		Desc: "sample function",
		Body: "xxx",
	})
	_, ok := ss.Find("func")
	if !ok {
		t.Errorf("snippets fails to recover existing snippet")
	}
}

func TestSnippetsMapInsert(t *testing.T) {
	ss := make(SnippetsMap)
	ok := ss.Insert(Snippet{"name", "desc", "body"})
	if !ok {
		t.Error("Insert() fails to insert Snippet to map")
	}
}

func TestSnippetsMapFind(t *testing.T) {
	ss := make(SnippetsMap)
	ss["func"] = Snippet{"func", "Go function", "func ${1:name} () {}"}
	_, ok := ss.Find("func")
	if !ok {
		t.Error("existing snippet signature could not be retrieved")
	}
}

func TestSnippetsMapList(t *testing.T) {
	ss := SnippetsMap{
		"func":   {"func", "Go function", "func() {}"},
		"struct": {"struct", "Go struct", "type struct {}"},
		"map":    {"map", "Go map", "map[string]string"},
	}
	want := []string{"func\tGo function", "map\tGo map", "struct\tGo struct"}
	if has := ss.List(); !reflect.DeepEqual(has, want) {
		t.Errorf("want: %v; has: %v", want, has)
	}
}
