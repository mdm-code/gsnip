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

// Kind labels the allowed operation.
type Kind int8

const (
	// Undef represents an undefined operation.
	Undef Kind = iota
	// Fnd represents the find operation.
	Fnd
	// Lst represents the list operation.
	Lst
	// Ins represents the insert operation.
	Ins
	// Del represents the delete operation.
	Del
	// Rld represents the reload operation.
	Rld
)

var kindToStr = map[Kind]string{
	Fnd: "@FND",
	Lst: "@LST",
	Ins: "@INS",
	Del: "@DEL",
	Rld: "@RLD",
}

var strToKind = map[string]Kind{
	"@FND": Fnd,
	"@LST": Lst,
	"@INS": Ins,
	"@DEL": Del,
	"@RLD": Rld,
}

const hSize = 4

// Msg represents the evaluated client message.
type Msg struct {
	knd  Kind
	cmd  bool
	ref  dep
	body []byte
}

// T tells the type of the message.
func (m Msg) T() Kind {
	return m.knd
}

// TString return the type of the operation as a string.
func (m Msg) TString() string {
	return kindToStr[m.knd]
}

// TByte returns the type of the operation as a slice of bytes.
func (m Msg) TByte() []byte {
	return []byte(m.TString())
}

// Contents returns the body of the message.
func (m Msg) Contents() []byte {
	return m.body
}

// IsUnbound tells whether the message is unbound.
func (m Msg) IsUnbound() bool {
	if m.ref == unbound {
		return true
	}
	return false
}

// Interpreter interprets the message received from the client.
type Interpreter struct {
	kmap map[string]Kind
}

// NewInterpreter creates a new instance of the Interpreter type.
func NewInterpreter() Interpreter {
	return Interpreter{
		kmap: strToKind,
	}
}

// Eval evalues the message sent into an appropriate Msg type.
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
		return Msg{knd: Undef, cmd: false, ref: unbound}
	}
}
