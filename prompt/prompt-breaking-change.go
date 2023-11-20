package prompt

import (
	"github.com/PossibleLlama/commit-check/model"

	tea "github.com/charmbracelet/bubbletea"
)

type PromptBreakingChange struct {
	cursor                int
	breakingChangeOptions []string

	width  int
	height int

	commit *model.Commit
}

func NewPromptBreakingChange(cmt *model.Commit) *PromptBreakingChange {
	return &PromptBreakingChange{
		cursor:                0,
		breakingChangeOptions: []string{"No", "Yes"},
		commit:                cmt,
	}
}

func (p PromptBreakingChange) Init() tea.Cmd {
	return nil
}

func (p PromptBreakingChange) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if p.cursor < len(p.breakingChangeOptions)-1 {
				p.cursor++
			}
		case "enter":
			p.commit.IsBreakingChange = p.breakingChangeOptions[p.cursor] == "Yes"
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p PromptBreakingChange) View() string {
	if p.width == 0 || p.height == 0 {
		return "loading..."
	}
	s := "Is the commit a breaking change:\n"
	for i, option := range p.breakingChangeOptions {
		if p.cursor == i {
			s += "> "
		} else {
			s += "  "
		}
		s += option + "\n"
	}
	return s
}
