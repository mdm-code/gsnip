package parsing

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/mdm-code/gsnip/snippets"
)

var funcSnip = snippets.Snippet{
	Name: "func",
	Desc: "Go function with no attributes and returns",
	Body: "func namedFunction() {\n\treturn\n}",
}

var structSnip = snippets.Snippet{
	Name: "struct",
	Desc: "Go struct template",
	Body: "type namedStruct struct {\n\tname string\n\tid int\n}",
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

func TestParsing(t *testing.T) {
	snips, err := snippets.NewSnippetsContainer("map")
	if err != nil {
		t.Error("failed to create snippet container")
	}
	snips.Insert(funcSnip)
	snips.Insert(structSnip)
	for _, r := range properReaders {
		parser := NewParser()
		has, _ := parser.Parse(r)
		if ok := reflect.DeepEqual(has, snips); !ok {
			t.Errorf("want: %v; has %v", snips, has)
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
		"	  startsnip struct \"sample comment\"   ", // whitespace on both sides
		"startsnip func() \"\"", // Empty comment
	}
	for _, i := range inputs {
		_, ok := splitSignature(i)
		if !ok {
			t.Errorf("Signature line : %s : should not error out", i)
		}
	}
}

func TestTakeBetweenPasses(t *testing.T) {
	inputs := []struct {
		text, want string
		delim      rune
	}{
		{"\"This text works just fine\"", "This text works just fine", '"'},
		{"`What about using ticks?`", "What about using ticks?", '`'},
		{"'Three single ' quotes return the longest'", "Three single ' quotes return the longest", '\''},
	}
	for _, i := range inputs {
		has, ok := takeBetween(i.text, i.delim)
		if !ok {
			t.Errorf("String :: %s :: is malformed", i.text)
		}
		if has != i.want {
			t.Errorf("Want: %s; has %s", i.want, has)
		}
	}
}

func TestTakeBetweenFails(t *testing.T) {
	inputs := []struct {
		text  string
		delim rune
	}{
		{"This has no delimiters", '"'},
		{"", '`'},
		{"' Has only one delimiter", '\''},
	}
	for _, i := range inputs {
		_, ok := takeBetween(i.text, i.delim)
		if ok {
			t.Errorf("Input :: %s :: should error out", i.text)
		}
	}
}
