package xdg

import (
	"os"
	"path"
)

var HOME string
var HOME_IS_SET bool
var XDG_DATA_HOME string
var XDG_DATA_DIRS []string = []string{"/usr/local/share/gsnip/", "/usr/share/gsnip/"}

func init() {
	HOME, HOME_IS_SET = os.LookupEnv("HOME")

	if HOME_IS_SET {
		XDG_DATA_HOME = path.Join(HOME, ".local/share/gsnip/")
	}
}

type queue struct {
	head   *node
	tail   *node
	length int
}

type node struct {
	item string
	next *node
}

func NewQueue() queue {
	return queue{}
}

func (q *queue) Enqueue(item string) *node {
	n := newNode(item)

	if q.head == nil {
		q.head = n
		q.tail = q.head
		q.length++
		return n
	}

	q.tail.SetNext(n)
	q.tail = n
	q.length++
	return n
}

func (q *queue) Dequeue() (*node, bool) {
	if q.head == nil {
		return nil, false
	}

	res := q.head
	q.head = res.Next()
	return res, true
}

func (q *queue) Peek() *node {
	return q.head
}

func newNode(item string) *node {
	return &node{item, nil}
}

func (n *node) Item() string {
	return n.item
}

func (n *node) Next() *node {
	return n.next
}

func (n *node) SetItem(item string) string {
	n.item = item
	return item
}

func (n *node) SetNext(next *node) *node {
	n.next = next
	return next
}

func Exists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func Arrange() queue {
	q := NewQueue()

	if Exists(XDG_DATA_HOME) {
		q.Enqueue(XDG_DATA_HOME)
	}
	for _, dir := range XDG_DATA_DIRS {
		if Exists(dir) {
			q.Enqueue(dir)
		}
	}
	return q

}

func Discover(q queue) (result *node, ok bool) {
	ok = true
	for ok {
		if result, ok = q.Dequeue(); ok {
			return
		}
	}
	return
}
