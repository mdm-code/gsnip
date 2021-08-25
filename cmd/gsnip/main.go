package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdm-code/gsnip/parsing"
	"github.com/mdm-code/gsnip/snippets"
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
	var snippets snippets.Container
	snippets, err = parser.Parse(f)
	if err != nil {
		log.Fatal(err)
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

		cmd := strings.ToLower(search)
		if parsing.IsCommand(cmd) && cmd == "list" {
			out := snippets.List()
			for _, s := range out {
				os.Stdout.WriteString(s + "\n")
			}
		} else {
			snip, ok := snippets.Find(search)
			if !ok {
				log.Fatal(search + " was not found")
			}
			pat := `\${[0-9]+:\w*}`
			out, ok := parsing.Replace(snip.Body, pat, repls...)
			if !ok {
				log.Fatal("Failed to compile regex pattern: " + pat)
			}
			os.Stdout.WriteString(out)
		}
	}
}

// Check if there's anything to read on STDIN
func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}
