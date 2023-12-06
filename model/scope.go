package model

type ScopeItem struct {
	ID   string
	Body string
}

func (s ScopeItem) Title() string       { return s.ID }
func (s ScopeItem) Description() string { return s.Body }
func (s ScopeItem) FilterValue() string { return s.ID }
