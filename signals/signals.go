package signals

import (
	"strings"
)

type dep int8

const (
	srvr dep = iota
	mngr
	unbound
)

type Token struct {
	Sign string
	cmd  bool
	ref  dep
}

func (t Token) IsCmd() bool {
	return t.cmd
}

func (t Token) IsReload() bool {
	return t.Sign == "@RELOAD" && t.IsCmd() && t.ref == srvr
}

func (t Token) IsList() bool {
	return t.Sign == "@LIST" && t.IsCmd() && t.ref == mngr
}

func (t Token) IsUnbound() bool {
	if t.ref == unbound {
		return true
	}
	return false
}

var cmds = []Token{
	{
		Sign: "@LIST",
		cmd:  true,
		ref:  mngr,
	},
	{
		Sign: "@RELOAD",
		cmd:  true,
		ref:  srvr,
	},
}

// Handles signal evaluation.
type Interpreter struct {
	cmds []Token
}

func NewInterpreter() Interpreter {
	return Interpreter{
		cmds: cmds,
	}
}

func (i Interpreter) Eval(s string) Token {
	s = strings.TrimSpace(s)
	if s == "" {
		return Token{Sign: "", cmd: false, ref: unbound}
	}
	for _, c := range i.cmds {
		if s == c.Sign {
			return c
		}
	}
	return Token{Sign: s, cmd: false, ref: mngr}
}
