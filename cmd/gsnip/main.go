package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mdm-code/gsnip/manager"
	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/snippets"
)

const Source = "/usr/local/share/gsnip/snippets"

// Check if there's anything to read on STDIN
func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func main() {
	fn := flag.String("f", Source, "Snippets source file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Snippet manager written in Go.\n\nUsage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var f *os.File
	f, err := os.Open(*fn)
	if err != nil {
		log.Fatal("Could not open " + *fn)
		f, err = os.Open(Source)
		if err != nil {
			log.Fatal("Could not find " + Source)
		}
	}
	defer f.Close()

	parser := parsing.NewParser()
	var snippets snippets.Container
	snippets, err = parser.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	mgr, ok := manager.NewManager(snippets)
	if !ok {
		log.Fatal("failed to initialized snippet manager")
	}

	if isPiped() {
		s := bufio.NewScanner(os.Stdin)
		s.Split(bufio.ScanWords)
		var params []string
		for s.Scan() {
			params = append(params, s.Text())
		}

		output, err := mgr.Execute(params...)
		if err != nil {
			log.Fatal("failed to execute command: ", err)
		}
		os.Stdout.WriteString(output)
	}
}
