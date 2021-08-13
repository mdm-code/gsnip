package parsing

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

var funcSnip = Snippet{
	Name: "func",
	Desc: "Go function with no attributes and returns",
	Body: "func namedFunction() {\n\treturn\n}",
}

var structSnip = Snippet{
	Name: "struct",
	Desc: "Go struct template",
	Body: `type namedStruct struct {
	name string
	id int
}`,
}

var properReaders = []io.Reader{
	strings.NewReader(`
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
endsnip`),
}

var failingReaders = []io.Reader{
	strings.NewReader(`
startsnip funcr
func (r *receiver) func({$1:params}) ({$2:returns}) (
	{$3:body}
)
endsnip`), // missing comment
	strings.NewReader(`
startsnip "Where's the name?"
func someVeryImportantFunction () {
	return "This is very important"
}
endsnip`),
	strings.NewReader(`
startsnip
God knows what this is.
endsnip`),
}

var snips = map[string]Snippet{"func": funcSnip, "struct": structSnip}

func TestParsing(t *testing.T) {
	for _, r := range properReaders {
		parser := NewParser()
		has, _ := parser.Parse(r)
		if ok := reflect.DeepEqual(has, snips); !ok {
			t.Errorf("want: %s; has %s", snips, has)
		}
	}
}

func TestAssertReadersFail(t *testing.T) {
	parser := NewParser()
	for _, r := range failingReaders {
		_, err := parser.Parse(r)
		if err == nil {
			t.Errorf("Reader %v should fail", r)
		}
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

func TestSignatureSplitFails(t *testing.T) {
	inputs := []string{
		"startsnip struct",                // Missing comment
		"startsnip printf 'some comment'", // comment not in double quotes
		"",
	}
	for _, i := range inputs {
		_, ok := splitSignature(i)
		if ok {
			t.Errorf("Signature line : %s : should fail", i)
		}
	}
}

func TestSignatureSplitPasses(t *testing.T) {
	inputs := []string{
		"startsnip struct \"Go struct snippet\"",
		"startsnip func() \"\"", // Empty comment
	}
	for _, i := range inputs {
		_, ok := splitSignature(i)
		if !ok {
			t.Errorf("Signature line : %s : should not error out", i)
		}
	}
}
