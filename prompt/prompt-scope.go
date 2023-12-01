package prompt

import (
	"github.com/PossibleLlama/commit-check/model"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func SetupListOfScopes() list.Model {
	scopeItems := []list.Item{
		model.ScopeItem{ID: "None", Body: "No scope"},
		model.ScopeItem{ID: "Other", Body: "Manual input"},
	}
	scopeList := list.New(scopeItems, list.NewDefaultDelegate(), 0, 0)
	scopeList.Title = "Scope of change"
	return scopeList
}

func (p PromptCommit) UpdateScope(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedItem := p.scopeOptions.SelectedItem()
			selectedScope := selectedItem.(model.ScopeItem)
			switch selectedScope.ID {
			case "None":
				p.page++
				p.inputMultiLine.Focus()
				return p, nil
			case "Other":
				if !p.inputSingleLine.Focused() {
					// First time pressing enter, so we need to focus the input
					p.inputSingleLine.Focus()
					p.inputSingleLine.Placeholder = "Scope"
				} else {
					// Second time pressing enter goes to the next page
					p.inputSingleLine.Blur()
					p.commit.Scope = p.inputSingleLine.Value()
					p.page++
					p.inputMultiLine.Focus()
					return p, nil
				}
			default:
				p.commit.Scope = selectedScope.ID
				p.page++
				p.inputMultiLine.Focus()
				return p, nil
			}
		default:
			if p.inputSingleLine.Focused() {
				p.inputSingleLine, cmd = p.inputSingleLine.Update(msg)
			} else {
				p.scopeOptions, cmd = p.scopeOptions.Update(msg)
			}
		}
	}
	return p, cmd
}

func (p PromptCommit) ViewScope() string {
	s := ""
	if p.inputSingleLine.Focused() {
		s = "Scope of the change:\n"
		s += docStyle.Render(p.inputSingleLine.View())
	} else {
		s += docStyle.Render(p.scopeOptions.View())
	}
	return s
}
