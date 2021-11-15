package main

import (
	"fmt"

	"github.com/mdm-code/gsnip/fhandle"
	"github.com/mdm-code/gsnip/parsing"
)

func main() {
	fb := fhandle.FSBuffer{}
	f, _ := fb.Open()
	defer f.Close()
	v := fhandle.Vim{Prog: "vim", File: f.Name()}
	v.Exec()
	parser := parsing.NewParser()
	snpts, _ := parser.Parse(f)
	fmt.Println(snpts)
}
