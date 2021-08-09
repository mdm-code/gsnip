package main

import "os"

func main() {
}

type snippet struct {
	name, desc, body string
}

func parse(f *os.File) (snippet, error) {
	return snippet{
		name: "func",
		desc: "Go function with no attributes and returns",
		body: "func namedFunction() {\n\treturn\n}",
	}, nil
}
