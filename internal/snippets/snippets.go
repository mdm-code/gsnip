package snippets

import (
	"fmt"
	"sort"
	"sync"
)

// Container provides an interface for a type handling snippet storage.
type Container interface {
	Insert(Snippet) error
	Find(string) (Snippet, error)
	List() ([]string, error)
	Delete(string) error
	ListObj() ([]Snippet, error)
}

// Snippet carries information about a single code snippet.
type Snippet struct {
	Name string
	Desc string
	Body string
}

// Repr provides an in-file snippet text representation.
func (s Snippet) Repr() string {
	return fmt.Sprintf("startsnip %s \"%s\"\n%s\nendsnip\n\n", s.Name, s.Desc, s.Body)
}

// mapContainer is a map-based implementation of a snippet Container.
type mapContainer struct {
	cntr map[string]Snippet
	sync.RWMutex
}

// NewSnippetsContainer creates a fresh instance of snippets container.
//
// Allowed types (t): map
func NewSnippetsContainer(t string) (Container, error) {
	switch t {
	case "map":
		return newMap(), nil
	default:
		return nil, fmt.Errorf("container type (%s) is not implemented", t)
	}
}

// newMap creates an instance of the snippet container relying on a map composite
// type.
func newMap() *mapContainer {
	return &mapContainer{
		cntr: make(map[string]Snippet),
	}
}

// Insert inserts a snippet to the container.
func (s *mapContainer) Insert(snip Snippet) (err error) {
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

// Find searches for a snippet name in the container.
func (s *mapContainer) Find(str string) (Snippet, error) {
	s.RLock()
	defer s.RUnlock()
	var snip Snippet
	snip, ok := s.cntr[str]
	if !ok {
		return snip, fmt.Errorf("snippet was not found")
	}
	return snip, nil
}

// List lists out all stored snippet names.
func (s *mapContainer) List() ([]string, error) {
	s.RLock()
	defer s.RUnlock()
	var result []string
	var str string
	for _, v := range s.cntr {
		str = fmt.Sprintf("%s\t%s", v.Name, v.Desc)
		result = append(result, str)
	}
	sort.Strings(result)
	return result, nil
}

// Delete deletes a snippet from the container.
func (s *mapContainer) Delete(key string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.cntr, key)
	return nil
}

// ListObj lists out all snippets stored in the container.
func (s *mapContainer) ListObj() (result []Snippet, err error) {
	s.RLock()
	defer s.RUnlock()
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
