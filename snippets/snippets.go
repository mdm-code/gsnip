package snippets

import (
	"fmt"
	"sort"
)

type Container interface {
	Insert(Snippet) bool
	Find(string) (Snippet, bool)
	List() []string
}

type Snippet struct {
	Name string
	Desc string
	Body string
}

type SnippetsMap map[string]Snippet

func (s SnippetsMap) Insert(snip Snippet) (success bool) {
	s[snip.Name], success = snip, true
	return
}

func (s SnippetsMap) Find(str string) (snip Snippet, success bool) {
	snip, success = s[str]
	return
}

func (s SnippetsMap) List() (result []string) {
	var str string
	for _, v := range s {
		str = fmt.Sprintf("%s\t%s", v.Name, v.Desc)
		result = append(result, str)
	}
	sort.Strings(result)
	return
}
