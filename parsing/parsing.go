package parsing

import (
	"bufio"
	"fmt"
	"io"
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

// Parses input files with snippets.
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

// Create new parser.
func NewParser() Parser {
	return Parser{
		sm: newStateMachine(),
	}
}

func (s Snippet) String() string {
	return fmt.Sprintf("%s\n%s\n\n%s", s.Name, s.Desc, s.Body)
}

// Parse file with snippets.
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

func splitSignature(sig string) ([]string, bool) {
	tokens := strings.Split(sig, " ")
	var result []string
	var comment []string
out:
	for i, tkn := range tokens {
		if strings.HasPrefix(tkn, "\"") {
			for _, ctkn := range tokens[i:] {
				comment = append(comment, ctkn)
			}
			break out
		}
		if tkn == "startsnip" {
			continue
		}
		result = append(result, tkn)
	}
	stripped := strings.Trim(strings.Join(comment, " "), `"`)
	if len(comment) != 0 { // Append even if it's an empty double quote
		result = append(result, stripped)
	}
	if len(result) != 2 {
		return []string{}, false
	}
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
