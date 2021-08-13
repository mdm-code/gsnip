package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mdm-code/gsnip/parsing"
)

func main() {
	var file string
	flag.StringVar(&file, "snippets", "", "Flat file with snippets")
	flag.Parse()
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	parser := parsing.NewParser()
	result, err := parser.Parse(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
