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

/* TODO: Token should have the BODY FIELD for insert and delete

Command sign (kind) type
	@LST
	@RLD
	@INS
	@FND
	@DEL
Body would keep the reminder of the received message
*/
type Token struct {
	sign string
	cmd  bool
	ref  dep
}

func (t Token) Contents() string {
	return t.sign
}

func (t Token) IsCmd() bool {
	return t.cmd
}

func (t Token) IsReload() bool {
	return t.sign == "@RELOAD" && t.IsCmd() && t.ref == srvr
}

func (t Token) IsList() bool {
	return t.sign == "@LIST" && t.IsCmd() && t.ref == mngr
}

func (t Token) IsUnbound() bool {
	if t.ref == unbound {
		return true
	}
	return false
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

func NewInterpreter() Interpreter {
	return Interpreter{
		cmds: cmds,
	}
}

// TODO: Eval Message reads first 4 bytes to interpret the kind
// The rest is passed to the body
func (i Interpreter) Eval(s string) Token {
	s = strings.TrimSpace(s)
	if s == "" {
		return Token{sign: "", cmd: false, ref: unbound}
	}
	for _, c := range i.cmds {
		if s == c.sign {
			return c
		}
	}
	return Token{sign: s, cmd: false, ref: mngr}
}
