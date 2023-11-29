package model

type ScopeItem struct {
	Heading string
	Body    string
}

func (s ScopeItem) Title() string       { return s.Heading }
func (s ScopeItem) Description() string { return s.Body }
func (s ScopeItem) FilterValue() string { return s.Heading }
