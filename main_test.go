package main

import (
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	f, err := os.Open("assets/sample.snippet")
	if err != nil {
		t.Error("test failed when opening snippet file")
	}
	defer f.Close()
	has, err := parse(f)
	want := snippet{
		name: "func",
		desc: "Go function with no attributes and returns",
		body: "func namedFunction() {\n\treturn\n}",
	}
	if has != want {
		t.Errorf("want: %s; has %s", want, has)
	}
}
