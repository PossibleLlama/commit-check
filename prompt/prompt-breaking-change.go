package prompt

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func SetupListOfBreakingChanges() list.Model {
	LItems := []list.Item{
		LItem("No"),
		LItem("Yes"),
	}
	breakingList := list.New(LItems, LItemDelegate{}, 0, 0)
	breakingList.Title = "Breaking change"
	breakingList.SetFilteringEnabled(false)
	return breakingList
}

func (p PromptCommit) UpdateBreakingChange(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			i, ok := p.breakingOptions.SelectedItem().(LItem)
			if ok {
				p.commit.IsBreakingChange = string(i) == "Yes"
			}
			return p, tea.Quit
		default:
			p.breakingOptions, cmd = p.breakingOptions.Update(msg)
		}
	}
	return p, cmd
}

func (p PromptCommit) ViewBreakingChange() string {
	return docStyle.Render(p.breakingOptions.View())
}
