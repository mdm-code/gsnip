package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mdm-code/gsnip/parsing"
)

const Source = "/usr/local/share/gsnip/snippets"

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
	snippets, err := parser.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	// Check if there's anything to read on STDIN
	isPiped := func() bool {
		fi, _ := os.Stdin.Stat()
		return (fi.Mode() & os.ModeCharDevice) == 0
	}

	if isPiped() {
		s := bufio.NewScanner(os.Stdin)
		s.Split(bufio.ScanWords)
		var attrs []string
		for s.Scan() {
			attrs = append(attrs, s.Text())
		}

		var search string
		var repls []string
		if len(attrs) == 0 {
			log.Fatal("There's nothing to find")
		} else if len(attrs) == 1 {
			search = attrs[0]
		} else {
			search, repls = attrs[0], attrs[1:]
		}

		snip, ok := snippets[search]
		if !ok {
			os.Stderr.WriteString(search + " was not found")
		}
		pat := `\${[0-9]+:\w*}`
		out, ok := parsing.Replace(snip.Body, pat, repls...)
		if !ok {
			os.Stderr.WriteString("Failed to compile regex pattern: " + pat)
		}
		os.Stdout.WriteString(out)
	}
}
