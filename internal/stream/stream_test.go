package stream

import (
	"reflect"
	"testing"
)

// Test command evaluation with correct text.
func TestEvalCommand(t *testing.T) {
	tab := []struct {
		in   string
		want Msg
	}{
		{"@LST", Msg{knd: Lst, cmd: true, ref: mngr, body: nil}},
		{"@RLD", Msg{knd: Rld, cmd: true, ref: srvr, body: nil}},
	}
	i := NewInterpreter()
	for _, param := range tab {
		has := i.Eval([]byte(param.in))
		if ok := reflect.DeepEqual(has, param.want); !ok {
			t.Errorf("has: %v; want %v", has, param.want)
		}
	}
}

// Test evaluation with some body to be kept.
func TestEvalWithBody(t *testing.T) {
	tab := []struct {
		in   string
		want Msg
	}{
		{
			"@FND gfunc",
			Msg{knd: Fnd, cmd: true, ref: mngr, body: []byte("gfunc")},
		},
		{
			`@INS startsnip func "sample function"
def main():
	print("Hello, world!")
endsnip`,
			Msg{knd: Ins, cmd: true, ref: mngr, body: []byte(`startsnip func "sample function"
def main():
	print("Hello, world!")
endsnip`)},
		},
	}
	i := NewInterpreter()
	for _, inst := range tab {
		has := i.Eval([]byte(inst.in))
		if ok := reflect.DeepEqual(has, inst.want); !ok {
			t.Errorf("has %v; want %v", has, inst.want)
		}
	}
}

// Test whether the token is able to tell its type in all formats.
func TestTokenTellsItsKind(t *testing.T) {
	type wnt struct {
		T  Kind
		Ts string
		Tb []byte
		bd []byte
	}
	tkns := []struct {
		inst Msg
		want wnt
	}{
		{
			Msg{knd: Ins, cmd: true, ref: mngr, body: []byte("some text")},
			wnt{Ins, "@INS", []byte("@INS"), []byte("some text")},
		},
		{
			Msg{knd: Fnd, cmd: true, ref: mngr, body: []byte("pyfunc")},
			wnt{Fnd, "@FND", []byte("@FND"), []byte("pyfunc")},
		},
		{
			Msg{knd: Rld, cmd: true, ref: mngr, body: []byte{}},
			wnt{Rld, "@RLD", []byte("@RLD"), []byte{}},
		},
		{
			Msg{knd: Del, cmd: true, ref: mngr, body: []byte("pyclass")},
			wnt{Del, "@DEL", []byte("@DEL"), []byte("pyclass")},
		},
		{
			Msg{knd: Lst, cmd: true, ref: mngr, body: []byte{}},
			wnt{Lst, "@LST", []byte("@LST"), []byte{}},
		},
	}
	for _, tk := range tkns {
		if has := tk.inst.T(); has != tk.want.T {
			t.Errorf("has %v; want %v", has, tk.want.T)
		}
		if has := tk.inst.TString(); has != tk.want.Ts {
			t.Errorf("has %v; want %v", has, tk.want.Ts)
		}
		has := tk.inst.TByte()
		if ok := reflect.DeepEqual(has, tk.want.Tb); !ok {
			t.Errorf("has %v; want %v", has, tk.want.Tb)
		}
		has = tk.inst.Contents()
		if ok := reflect.DeepEqual(has, tk.want.bd); !ok {
			t.Errorf("has %v; want %v", has, tk.want.bd)
		}
	}
}

// Test whether an empty string is rendered unbound.
func TestEmptyByteSliceUnbound(t *testing.T) {
	tab := []struct {
		inst string
		want bool
	}{
		{"", true},
		{"@FND goimports", false},
	}
	i := NewInterpreter()
	for _, inst := range tab {
		tkn := i.Eval([]byte(inst.inst))
		if has := tkn.IsUnbound(); has != inst.want {
			t.Errorf("has: %v; want: %v", has, inst.want)
		}
	}
}
