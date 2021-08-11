package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("assets/sample.snip")
	defer f.Close()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(parse(f))
}

type Snip struct {
	name, desc, body string
}

func (s Snip) String() string {
	return fmt.Sprintf("%s\n%s\n\n%s", s.name, s.desc, s.body)
}

type State uint8

const (
	SCANLINE State = iota
	SIGNATURE
	SCANBODY
)

func parse(f *os.File) (map[string]Snip, error) {
	// Finite state automaton for reading file contents.
	in := bufio.NewScanner(f)
	snips := make(map[string]Snip)
	state := SCANLINE

	var line string
	var name, desc string
	var body []string

	for in.Scan() {
		if state == SCANLINE {
			line = in.Text()
			if strings.HasPrefix(line, "startsnip") {
				state = SIGNATURE
			}
		}
		if state == SCANBODY {
			line = in.Text()
			if strings.HasPrefix(line, "endsnip") {
				snips[name] = Snip{
					name: name,
					desc: desc,
					body: strings.Join(body, "\n"),
				}
				state = SCANLINE
				body = body[:0]
				continue
			}
			body = append(body, line)
		}
		if state == SIGNATURE {
			elems := strings.SplitN(line, " ", 3)
			name = elems[1]
			desc = strings.Trim(elems[2], "\"")
			state = SCANBODY
			continue
		}
	}
	return snips, nil
}
