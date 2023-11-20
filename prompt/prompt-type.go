package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (p PromptCommit) UpdateType(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down":
			if p.cursor < len(p.typeNameOptions)-1 {
				p.cursor++
			}
		case "enter":
			p.commit.Type = p.typeNameOptions[p.cursor]
			p.page++
			p.cursor = 0
			return p, nil
		}
	}
	return p, nil
}

func (p PromptCommit) ViewType() string {
	s := "Select the type of change you're committing:\n"
	for i, typeName := range p.typeNameOptions {
		if p.cursor == i {
			s += "> "
		} else {
			s += "  "
		}
		s += string(typeName) + "\n"
	}
	return s
}
