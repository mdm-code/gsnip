package signals

import "testing"

// Test command evaluation with correct text.
func TestEvalCommand(t *testing.T) {
	tab := []struct {
		in   string
		want Token
	}{
		{"@LIST", cmds[0]},
		{"@RELOAD", cmds[1]},
	}
	i := NewInterpreter()
	for _, param := range tab {
		has := i.Eval(param.in)
		if has != param.want {
			t.Errorf("has: %v; want %v", has, param.want)
		}
	}
}

// Test evaluation of non-command signs.
func TestEvalToken(t *testing.T) {
	tab := []struct {
		in   string
		want Token
	}{
		{"pyclass", Token{sign: "pyclass", cmd: false, ref: mngr}},
		{"pyfunc ", Token{sign: "pyfunc", cmd: false, ref: mngr}},
		{"  gfunc  ", Token{sign: "gfunc", cmd: false, ref: mngr}},
	}
	i := NewInterpreter()
	for _, param := range tab {
		has := i.Eval(param.in)
		if has != param.want {
			t.Errorf("has %v; want %v", has, param.want)
		}
	}
}

func TestTokenTellsIfItsCommand(t *testing.T) {
	tab := []struct {
		clientSig string
		want      bool
	}{
		{
			"@RELOAD",
			true,
		},
		{
			"@LIST",
			true,
		},
		{
			"pprog",
			false,
		},
	}
	for _, c := range tab {
		interp := NewInterpreter()
		has := interp.Eval(c.clientSig)
		if out := has.IsCmd(); out != c.want {
			t.Errorf("has: %v; want: %v", out, c.want)
		}
	}
}

func TestTokenTellsIfItsReload(t *testing.T) {
	tab := []struct {
		tkn  Token
		want bool
	}{
		{Token{"@RELOAD", true, 0}, true},
		{Token{"@LIST", true, 1}, false},
		{Token{"pprog", false, 1}, false},
	}
	for _, cse := range tab {
		if has := cse.tkn.IsReload(); has != cse.want {
			t.Errorf("has: %v; want: %v", has, cse.want)
		}
	}
}
