package signals

import (
	"bytes"
)

type dep int8

const (
	unbound dep = iota
	srvr
	mngr
)

type kind int8

const (
	undef kind = iota
	fnd
	lst
	ins
	del
	rld
)

var kindToStr = map[kind]string{
	fnd: "@FND",
	lst: "@LST",
	ins: "@INS",
	del: "@DEL",
	rld: "@RLD",
}

var strToKind = map[string]kind{
	"@FND": fnd,
	"@LST": lst,
	"@INS": ins,
	"@DEL": del,
	"@RLD": rld,
}

const hSize = 5

type Token struct {
	knd  kind
	cmd  bool
	ref  dep
	body []byte
}

func (t Token) T() kind {
	return t.knd
}

func (t Token) TString() string {
	return kindToStr[t.knd]
}

func (t Token) TByte() []byte {
	return []byte(t.TString())
}

func (t Token) Contents() []byte {
	return t.body
}

func (t Token) IsUnbound() bool {
	if t.ref == unbound {
		return true
	}
	return false
}

type Interpreter struct {
	kmap map[string]kind
}

func NewInterpreter() Interpreter {
	return Interpreter{
		kmap: strToKind,
	}
}

func (i Interpreter) Eval(b []byte) Token {
	b = bytes.TrimSpace(b)
	header := b[:hSize]
	header = header[:len(header)-1] // NOTE: Trailing 0 byte is removed
	var body []byte
	if len(b) > hSize {
		body = b[hSize:]
	}
	switch k := i.kmap[string(header)]; k {
	case lst, del, ins, fnd:
		return Token{knd: k, cmd: true, ref: mngr, body: body}
	case rld:
		return Token{knd: k, cmd: true, ref: srvr, body: body}
	default:
		return Token{knd: undef, cmd: false, ref: unbound}
	}
}
