package snippets

import (
	"fmt"
	"sort"
	"sync"
)

type Container interface {
	Insert(Snippet) error
	Find(string) (Snippet, error)
	List() ([]string, error)
	Delete(string) error
	ListObj() ([]Snippet, error)
}

type Snippet struct {
	Name string
	Desc string
	Body string
}

// In-file snippet text representation.
func (s Snippet) Repr() string {
	return fmt.Sprintf("startsnip %s \"%s\"\n%s\nendsnip\n\n", s.Name, s.Desc, s.Body)
}

type SnippetsMap struct {
	cntr map[string]Snippet
	sync.Mutex
}

// Create a fresh instance of snippets container.
//
// Allowed types (t): map
func NewSnippetsContainer(t string) (Container, error) {
	switch t {
	case "map":
		return NewSnippetsMap(), nil
	default:
		return nil, fmt.Errorf("container type (%s) is not implemented", t)
	}
}

func NewSnippetsMap() *SnippetsMap {
	return &SnippetsMap{
		cntr: make(map[string]Snippet),
	}
}

func (s *SnippetsMap) Insert(snip Snippet) (err error) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.cntr[snip.Name]; !exists {
		s.cntr[snip.Name] = snip
		err = nil
	} else {
		err = fmt.Errorf("snippet %s already exists", snip.Name)
	}
	return
}

func (s *SnippetsMap) Find(str string) (Snippet, error) {
	s.Lock()
	defer s.Unlock()
	var snip Snippet
	snip, ok := s.cntr[str]
	if !ok {
		return snip, fmt.Errorf("snippet was not found")
	}
	return snip, nil
}

func (s *SnippetsMap) List() ([]string, error) {
	s.Lock()
	defer s.Unlock()
	var result []string
	var str string
	for _, v := range s.cntr {
		str = fmt.Sprintf("%s\t%s", v.Name, v.Desc)
		result = append(result, str)
	}
	sort.Strings(result)
	return result, nil
}

func (s *SnippetsMap) Delete(key string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.cntr, key)
	return nil
}

func (s *SnippetsMap) ListObj() (result []Snippet, err error) {
	s.Lock()
	defer s.Unlock()
	for _, v := range s.cntr {
		result = append(result, v)
	}
	result = sorted(result)
	return
}

func sorted(s []Snippet) []Snippet {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Name < s[j].Name
	})
	return s
}
