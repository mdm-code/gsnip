package xdg

import (
	"testing"
)

func TestQueueOperations(t *testing.T) {
	tab := []struct {
		value string
		ok    bool
	}{
		{"1", true},
		{"2", true},
		{"3", true},
		{"4", true},
		{"5", true},
	}
	q := NewQueue()

	for _, elem := range tab {
		q.Enqueue(elem.value)
	}

	for _, elem := range tab {
		if val, ok := q.Dequeue(); val.Item() != elem.value || ok != elem.ok {
			t.Errorf("has: %s; want %s", val.Item(), elem.value)
		}
	}
}

func TestDequeueEmpty(t *testing.T) {
	q := NewQueue()
	if _, ok := q.Dequeue(); ok {
		t.Error("empty queue does not return false")
	}
}

func TestIfExists(t *testing.T) {
	dirs := []string{"/usr"}
	for _, d := range dirs {
		if ok := Exists(d); !ok {
			t.Errorf("want: %t; has: %t for %s", true, ok, d)
		}
	}
}

func TestDiscoverNotOk(t *testing.T) {
	q := NewQueue()

	if _, ok := Discover(q); ok {
		t.Errorf("has: %t; want %t", ok, false)
	}
}

func TestDiscoverIsOk(t *testing.T) {
	q := NewQueue()
	q.Enqueue("/usr/local/share/gsnip/")

	if _, ok := Discover(q); !ok {
		t.Errorf("has: %t; want %t", ok, true)
	}
}
