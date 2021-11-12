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
		{"pyclass", Token{Sign: "pyclass", cmd: false, ref: mngr}},
		{"pyfunc ", Token{Sign: "pyfunc", cmd: false, ref: mngr}},
		{"  gfunc  ", Token{Sign: "gfunc", cmd: false, ref: mngr}},
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

func TestTokenTellsIfItsList(t *testing.T) {
	tab := []struct {
		tkn  Token
		want bool
	}{
		{Token{"@LIST", true, 1}, true},
		{Token{"@RELOAD", true, 0}, false},
		{Token{"pprog", false, 1}, false},
	}
	for _, cse := range tab {
		if has := cse.tkn.IsList(); has != cse.want {
			t.Errorf("has: %v; want: %v", has, cse.want)
		}
	}
}

func TestEmptyStringIsUnbound(t *testing.T) {
	input := ""
	want := true
	interp := NewInterpreter()
	tkn := interp.Eval(input)
	if has := tkn.IsUnbound(); has != want {
		t.Errorf("has: %v; want: %v", has, want)
	}
}
