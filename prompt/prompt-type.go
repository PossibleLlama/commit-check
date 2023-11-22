package prompt

import (
	"fmt"

	"github.com/PossibleLlama/commit-check/model"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func SetupListOfTypes(typeNameOptions []model.CommitType) list.Model {
	typeNameItems := []list.Item{}
	for _, t := range typeNameOptions {
		typeNameItems = append(typeNameItems, LItem(t.String()))
	}
	typeNameList := list.New(typeNameItems, LItemDelegate{}, 0, 0)
	typeNameList.Title = "Type of change"
	return typeNameList
}

func (p PromptCommit) UpdateType(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			p.commit.Type = model.CommitType(fmt.Sprint(p.typeNameOptions.SelectedItem()))
			p.page++
		default:
			p.typeNameOptions, cmd = p.typeNameOptions.Update(msg)
		}
	}
	return p, cmd
}

func (p PromptCommit) ViewType() string {
	return docStyle.Render(p.typeNameOptions.View())
}
