package signals

import "strings"

type dep int8

const (
	srvr dep = iota
	mngr
)

type Token struct {
	sign string
	cmd  bool
	ref  dep
}

var cmds = []Token{
	{
		sign: "@LIST",
		cmd:  true,
		ref:  mngr,
	},
	{
		sign: "@RELOAD",
		cmd:  true,
		ref:  srvr,
	},
}

// Handles signal evaluation.
type Interpreter struct {
	cmds []Token
}

func (i Interpreter) Eval(s string) (Token, error) {
	s = strings.TrimSpace(s)
	for _, c := range i.cmds {
		if s == c.sign {
			return c, nil
		}
	}
	return Token{sign: s, cmd: false, ref: mngr}, nil
}

func NewInterpreter() Interpreter {
	return Interpreter{
		cmds: cmds,
	}
}
