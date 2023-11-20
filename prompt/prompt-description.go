package prompt

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (p PromptCommit) UpdateDescription(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		p.inputMultiLine.Focus()
		switch msg.String() {
		case "enter":
			if strings.HasSuffix(p.inputMultiLine.Value(), "\n") {
				p.inputMultiLine.Blur()
				p.page++
				p.cursor = 0
			}
			p.inputMultiLine.SetValue(p.inputMultiLine.Value() + "\n")
			p.commit.Description = p.inputMultiLine.Value()

		default:
			p.inputMultiLine, cmd = p.inputMultiLine.Update(msg)
			p.commit.Description = p.inputMultiLine.Value()
		}
	}
	return p, cmd
}

func (p PromptCommit) ViewDescription() string {
	s := "Context for why you are making the change:\nEnter twice to finish\n"
	s += p.inputMultiLine.View()
	if p.inputMultiLine.Focused() {
		s += "is focused"
	} else {
		s += "is not focused"
	}
	return s
}
