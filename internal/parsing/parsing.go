package parsing

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mdm-code/gsnip/internal/snippets"
)

const (
	scanning state = iota
	signature
	scanBody
	errored
)

var (
	// ErrEmptyFile is raised when file is not parsed successfully.
	ErrEmptyFile = errors.New("nothing to parse")
	// ErrLine is raised when there is an error on a line.
	ErrLine = errors.New("line contains an error")
)

type state uint8

type stateMachine struct {
	transitions map[state]func(*stateMachine, string) (state, string)
	parsed      []snippets.Snippet
	body        []string
	state       state
}

// Parser parses input files with snippets.
type Parser struct {
	sm *stateMachine
}

func newStateMachine() *stateMachine {
	return &stateMachine{
		transitions: map[state]func(*stateMachine, string) (state, string){
			scanning:  (*stateMachine).scanLine,
			signature: (*stateMachine).readSignature,
			scanBody:  (*stateMachine).scanBody,
		},
		state: scanning,
	}
}

// NewParser creates a new parser.
func NewParser() Parser {
	return Parser{
		sm: newStateMachine(),
	}
}

// Parse parses file with snippets. The result is a map
// of of snippets with name as key and body as value.
func (p *Parser) Parse(i io.Reader) (snippets.Container, error) {
	smap, err := snippets.NewSnippetsContainer("map")
	if err != nil {
		return nil, err
	}
	parsed, err := p.run(i)
	if err != nil {
		return smap, err
	}
	for _, s := range parsed {
		smap.Insert(s)
	}
	return smap, nil
}

// run runs the parser against input text.
func (p *Parser) run(i io.Reader) ([]snippets.Snippet, error) {
	result, err := p.sm.run(i)
	return result, err
}

func (sm *stateMachine) scanLine(line string) (state, string) {
	if l := strings.TrimSpace(line); strings.HasPrefix(l, "startsnip") {
		return signature, line
	}
	return scanning, ""
}

func (sm *stateMachine) readSignature(line string) (state, string) {
	elems, ok := splitSignature(line)
	if !ok {
		return errored, line
	}
	snip := snippets.Snippet{Name: elems[0], Desc: elems[1]}
	sm.parsed = append(sm.parsed, snip)
	return scanBody, ""
}

func splitSignature(s string) ([]string, bool) {
	var startToken, name, comment string
	splits := strings.SplitN(s, " ", 3)
	unpack(splits, &startToken, &name, &comment)
	comment, ok := takeBetween(comment, '"')
	if !ok {
		return nil, false
	}
	return []string{name, comment}, true
}

func unpack(s []string, vars ...*string) {
	for i, str := range s {
		*vars[i] = str
	}
}

// Grab text between two delimiters.
func takeBetween(s string, delim rune) (string, bool) {
	var idxs []int
	var result string
	for i, c := range s {
		if c == delim {
			idxs = append(idxs, i)
		}
	}
	if len(idxs) < 2 {
		return result, false
	}
	result = s[idxs[0]+1 : idxs[len(idxs)-1]]
	return result, true
}

func (sm *stateMachine) scanBody(line string) (state, string) {
	if l := strings.TrimSpace(line); strings.HasPrefix(l, "endsnip") {
		sm.parsed[len(sm.parsed)-1].Body = strings.Join(sm.body, "\n")
		sm.body = sm.body[:0]
		return scanning, ""
	}
	sm.body = append(sm.body, line)
	return scanBody, ""
}

func (sm *stateMachine) run(f io.Reader) ([]snippets.Snippet, error) {
	s := bufio.NewScanner(f)
	sm.reset()

	var line string
	for {
		if sm.state == errored {
			return nil, fmt.Errorf("%w: %s", ErrLine, line)
		}
		if line == "" {
			if ok := s.Scan(); !ok {
				break
			}
			line = s.Text()
		}
		callable := sm.transitions[sm.state]
		sm.state, line = callable(sm, line)
	}
	if sm.parsed == nil {
		return sm.parsed, fmt.Errorf("%w", ErrEmptyFile)
	}
	return sm.parsed, nil
}

// Reset the state of the object.
func (sm *stateMachine) reset() {
	sm.parsed = nil
	sm.body = nil
	sm.state = scanning
}
