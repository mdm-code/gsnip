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

var snips = snippets.SnippetsMap{"func": funcSnip, "struct": structSnip}

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

func TestTestReplaceOutput(t *testing.T) {
	inputs := []struct {
		text, want string
		repls      []string
	}{
		{
			"${1:foo} ${2:bar} ${3:baz}",
			"foo bar baz",
			[]string{"foo", "bar", "baz"},
		},
	}
	pat := `\${[0-9]+:\w*}`
	for _, i := range inputs {
		has, ok := Replace(i.text, pat, i.repls...)
		if !ok {
			t.Errorf("Failed to compile regex pattern: " + pat)
		}
		if has != i.want {
			t.Errorf("String '%s' should look like '%s'", has, i.want)
		}
	}
}

func TestCheckIfIsCommand(t *testing.T) {
	inputs := []struct {
		cmd      string
		expected bool
	}{
		{"list", true},
		{"prune", false},
		{"", false},
	}
	for _, i := range inputs {
		ok := IsCommand(i.cmd)
		if ok != i.expected {
			t.Errorf("command string was misidentified: %s", i.cmd)
		}
	}
}

func TestParsingFailsOnCmd(t *testing.T) {
	line := "startsnip list \"Signature of a command to fail\""
	sm := newStateMachine()
	state, line := sm.readSignature(line)
	if state != ERROR {
		t.Errorf("error was not raised on line: %s", line)
	}
}
