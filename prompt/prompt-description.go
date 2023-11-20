package prompt

import (
	"github.com/PossibleLlama/commit-check/model"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type PromptDescription struct {
	input textarea.Model

	width  int
	height int

	commit *model.Commit
}

func NewPromptDescription(cmt *model.Commit) *PromptDescription {
	ta := textarea.New()
	ta.Placeholder = "Reason"
	ta.Focus()

	return &PromptDescription{
		input:  ta,
		commit: cmt,
	}
}

func (p PromptDescription) Init() tea.Cmd {
	return nil
}

func (p PromptDescription) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return p, tea.Quit
		default:
			p.input, cmd = p.input.Update(msg)
			p.commit.Description = p.input.Value()
		}
	}
	return p, cmd
}

func (p PromptDescription) View() string {
	if p.width == 0 || p.height == 0 {
		return "loading..."
	}
	s := "Context for why you are making the change:\n"
	s += p.input.View()
	return s
}
