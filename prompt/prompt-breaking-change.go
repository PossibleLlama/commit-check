package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (p PromptCommit) UpdateBreakingChange(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down":
			if p.cursor < len(p.breakingOptions)-1 {
				p.cursor++
			}
		case "enter":
			p.commit.IsBreakingChange = p.breakingOptions[p.cursor] == "Yes"
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p PromptCommit) ViewBreakingChange() string {
	s := "Is the commit a breaking change:\n"
	for i, option := range p.breakingOptions {
		if p.cursor == i {
			s += "> "
		} else {
			s += "  "
		}
		s += option + "\n"
	}
	return s
}
