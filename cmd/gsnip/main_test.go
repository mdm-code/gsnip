package main

import "testing"

func TestCheckIfIsCommand(t *testing.T) {
	inputs := []struct {
		cmd      string
		expected bool
	}{
		{"list", true},
		{"prune", false},
		{"", false},
	}
	for _, i := range inputs {
		ok := isCommand(i.cmd)
		if ok != i.expected {
			t.Errorf("command string was misidentified: %s", i.cmd)
		}
	}
}
