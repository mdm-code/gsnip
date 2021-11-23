package stream

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
	Fnd
	Lst
	Ins
	Del
	Rld
)

var kindToStr = map[kind]string{
	Fnd: "@FND",
	Lst: "@LST",
	Ins: "@INS",
	Del: "@DEL",
	Rld: "@RLD",
}

var strToKind = map[string]kind{
	"@FND": Fnd,
	"@LST": Lst,
	"@INS": Ins,
	"@DEL": Del,
	"@RLD": Rld,
}

const hSize = 4

type Msg struct {
	knd  kind
	cmd  bool
	ref  dep
	body []byte
}

func (m Msg) T() kind {
	return m.knd
}

func (m Msg) TString() string {
	return kindToStr[m.knd]
}

func (m Msg) TByte() []byte {
	return []byte(m.TString())
}

func (m Msg) Contents() []byte {
	return m.body
}

func (m Msg) IsUnbound() bool {
	if m.ref == unbound {
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

func (i Interpreter) Eval(b []byte) Msg {
	b = bytes.TrimSpace(b)

	var header []byte
	if len(b) >= hSize {
		header = b[:hSize]
	}
	var body []byte
	if len(b) > hSize {
		body = b[hSize+1:] // NOTE: skip space between kind and body
	}
	switch k := i.kmap[string(header)]; k {
	case Lst, Del, Ins, Fnd:
		return Msg{knd: k, cmd: true, ref: mngr, body: body}
	case Rld:
		return Msg{knd: k, cmd: true, ref: srvr, body: body}
	default:
		return Msg{knd: undef, cmd: false, ref: unbound}
	}
}
