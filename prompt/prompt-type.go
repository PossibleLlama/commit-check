package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PromptType struct {
	cursor          int
	typeNameOptions []string
	typeName        string
}

func NewPromptType(typeNameOptions []string) *PromptType {
	return &PromptType{
		cursor:          0,
		typeNameOptions: typeNameOptions,
		typeName:        "",
	}
}

func (p PromptType) Init() tea.Cmd {
	return nil
}

func (p PromptType) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		case "up":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down":
			if p.cursor < len(p.typeNameOptions)-1 {
				p.cursor++
			}
		case "enter":
			p.typeName = p.typeNameOptions[p.cursor]
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p PromptType) View() string {
	s := "Select the type of change you're committing:\n"
	for i, typeName := range p.typeNameOptions {
		if p.cursor == i {
			s += "> "
		} else {
			s += "  "
		}
		s += typeName + "\n"
	}
	return s
}
