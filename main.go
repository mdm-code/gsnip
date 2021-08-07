package main

import (
	"flag"
	"fmt"
)

func main() {
	var clsName string
	flag.StringVar(&clsName, "name", "A", "name the class")
	flag.Parse()
	fmt.Printf("class %s:\n    def __init__(self) -> None:\n        pass\n", clsName)
}
