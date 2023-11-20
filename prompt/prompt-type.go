package prompt

import (
	"github.com/PossibleLlama/commit-check/model"

	tea "github.com/charmbracelet/bubbletea"
)

type PromptType struct {
	cursor          int
	typeNameOptions []model.CommitType

	width  int
	height int

	commit *model.Commit
}

func NewPromptType(typeNameOptions []model.CommitType, cmt *model.Commit) *PromptType {
	return &PromptType{
		cursor:          0,
		typeNameOptions: typeNameOptions,
		commit:          cmt,
	}
}

func (p PromptType) Init() tea.Cmd {
	return nil
}

func (p PromptType) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String(), tea.KeyCtrlC.String():
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
			p.commit.Type = p.typeNameOptions[p.cursor]
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p PromptType) View() string {
	if p.width == 0 || p.height == 0 {
		return "loading..."
	}
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
