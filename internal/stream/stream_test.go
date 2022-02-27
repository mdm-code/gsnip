package stream

import "testing"

// Test if instances of Reply are created correctly.
func TestReplyCreation(t *testing.T) {
	data := []struct {
		name string
		res  result
		body []byte
	}{
		{"success", Success, []byte("")},
		{"failure", Failure, []byte("")},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_ = Reply{d.res, d.body}
		})
	}
}

// Check if instances of Request are created correctly.
func TestRequestCreation(t *testing.T) {
	data := []struct {
		name string
		opcd Opcode
		body []byte
	}{
		{"success", Undefined, []byte("")},
		{"failure", Find, []byte("")},
		{"failure", List, []byte("")},
		{"failure", Insert, []byte("")},
		{"failure", Delete, []byte("")},
		{"failure", Reload, []byte("")},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_ = Request{d.opcd, d.body}
		})
	}
}
