package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
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

func mockReader() io.Reader {
	reader := strings.NewReader(`
startsnip func "Go function with no attributes and returns"
func namedFunction() {
	return
}
endsnip

startsnip struct "Go struct template"
type namedStruct struct {
	name string
	id int
}
endsnip`)
	return reader
}

var snips = map[string]Snip{"func": funcSnip, "struct": structSnip}

func TestParsing(t *testing.T) {
	reader := mockReader()
	has, _ := parse(reader)
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

// TODO: Make errors assertions -- replace FILE with fake buffer
