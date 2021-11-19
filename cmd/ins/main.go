package main

import (
	"fmt"

	"github.com/mdm-code/gsnip/editor"
)

func main() {
	e, _ := editor.NewEditor("vim", "text.txt")
	defer e.Exit()
	data, _ := e.Run()
	fmt.Println(string(data))
}
