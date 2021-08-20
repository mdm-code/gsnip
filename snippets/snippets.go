package snippets

type Snippet struct {
	Name string
	Desc string
	Body string
}

type Snippets map[string]Snippet

func (s Snippets) Insert(snip Snippet) (success bool) {
	s[snip.Name], success = snip, true
	return
}

func (s Snippets) Find(str string) (snip Snippet, success bool) {
	snip, success = s[str]
	return
}
