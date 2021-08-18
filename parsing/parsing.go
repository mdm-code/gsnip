package parsing

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	SCANNING State = iota
	SIGNATURE
	SCANBODY
	ERROR
)

type State uint8

type Snippet struct {
	Name string
	Desc string
	Body string
}

type StateMachine struct {
	transitions map[State]func(*StateMachine, string) (State, string)
	parsed      []Snippet
	body        []string
	state       State
}

// Parser parses input files with snippets.
type Parser struct {
	sm *StateMachine
}

func newStateMachine() *StateMachine {
	return &StateMachine{
		transitions: map[State]func(*StateMachine, string) (State, string){
			SCANNING:  (*StateMachine).scanLine,
			SIGNATURE: (*StateMachine).readSignature,
			SCANBODY:  (*StateMachine).scanBody,
		},
		state: SCANNING,
	}
}

// NewParser creates a fresh new parser.
func NewParser() Parser {
	return Parser{
		sm: newStateMachine(),
	}
}

// Parse file with snippets. The result is a map
// of of snippets with name as key and body as value.
func (p *Parser) Parse(i io.Reader) (map[string]Snippet, error) {
	result := make(map[string]Snippet)
	parsed, err := p.sm.run(i)
	if err != nil {
		return result, err
	}
	for _, s := range parsed {
		result[s.Name] = s
	}
	return result, nil
}

func (sm *StateMachine) scanLine(line string) (State, string) {
	if strings.HasPrefix(line, "startsnip") {
		return SIGNATURE, line
	}
	return SCANNING, ""
}

func (sm *StateMachine) readSignature(line string) (State, string) {
	elems, ok := splitSignature(line)
	if !ok {
		return ERROR, line
	}
	snip := Snippet{Name: elems[0], Desc: elems[1]}
	sm.parsed = append(sm.parsed, snip)
	return SCANBODY, ""
}

func splitSignature(s string) ([]string, bool) {
	var startToken, name, comment string
	splits := strings.SplitN(s, " ", 3)
	unpack(splits, &startToken, &name, &comment)
	comment, ok := takeBetween(comment, '"')
	if !ok {
		return []string{}, false
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

func (sm *StateMachine) scanBody(line string) (State, string) {
	if strings.HasPrefix(line, "endsnip") {
		sm.parsed[len(sm.parsed)-1].Body = strings.Join(sm.body, "\n")
		sm.body = sm.body[:0]
		return SCANNING, ""
	}
	sm.body = append(sm.body, line)
	return SCANBODY, ""
}

func (sm *StateMachine) run(f io.Reader) ([]Snippet, error) {
	s := bufio.NewScanner(f)
	var line string
	for {
		if sm.state == ERROR {
			return []Snippet{}, fmt.Errorf("Error on line: %s", line)
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
	return sm.parsed, nil
}

// Replace attempts to substitute placeholders with strings.
func Replace(str string, pat string, repls ...string) (string, bool) {
	re, err := regexp.Compile(pat)
	if err != nil {
		return str, false
	}
	for _, r := range repls {
		sms := re.FindStringSubmatch(str)
		if len(sms) == 0 {
			break
		}
		sm := sms[0]
		str = strings.Replace(str, sm, r, 1)
	}
	return str, true
}
