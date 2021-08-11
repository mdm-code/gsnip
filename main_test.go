package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var funcSnip = Snip{
	name: "func",
	desc: "Go function with no attributes and returns",
	body: "func namedFunction() {\n\treturn\n}",
}

var structSnip = Snip{
	name: "struct",
	desc: "Go struct template",
	body: `type namedStruct struct {
	name string
	id int
}`,
}

var snips = map[string]Snip{"func": funcSnip, "struct": structSnip}

func TestParsing(t *testing.T) {
	f, err := os.Open("assets/sample.snip")
	if err != nil {
		t.Error("test failed when opening file with snippets")
	}
	defer f.Close()
	has, err := parse(f)
	if ok := reflect.DeepEqual(has, snips); !ok {
		t.Errorf("want: %s; has %s", snips, has)
	}
}

func TestSnipStructDisplay(t *testing.T) {
	want := `func
Go function with no attributes and returns

func namedFunction() {
	return
}`
	has := fmt.Sprintf("%s", funcSnip)
	if has != want {
		t.Errorf("has: %s; want: %s", has, want)
	}
}
