package prompt

import (
	"github.com/PossibleLlama/commit-check/model"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PromptCommit struct {
	cursor int
	page   int

	typeNameOptions []model.CommitType
	scopeOptions    []string
	breakingOptions []string

	inputSingleLine textinput.Model
	inputMultiLine  textarea.Model

	width  int
	height int

	commit *model.Commit
}

func NewPromptCommit(typeNameOptions []model.CommitType, cmt *model.Commit) *PromptCommit {
	return &PromptCommit{
		cursor: 0,
		page:   0,

		typeNameOptions: typeNameOptions,
		scopeOptions:    []string{"None", "Other"},
		breakingOptions: []string{"No", "Yes"},

		inputSingleLine: textinput.New(),
		inputMultiLine:  textarea.New(),

		width:  0,
		height: 0,

		commit: cmt,
	}
}

func (p PromptCommit) Init() tea.Cmd {
	p.CheckJira()
	return nil
}

func (p PromptCommit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return p, tea.Quit
		default:
			switch p.page {
			case 0:
				return p.UpdateType(msg)
			case 1:
				return p.UpdateScope(msg)
			case 2:
				return p.UpdateDescription(msg)
			case 3:
				return p.UpdateBreakingChange(msg)
			}
		}
	}
	return p, nil
}

func (p PromptCommit) View() string {
	if p.width == 0 || p.height == 0 {
		return "loading..."
	}
	switch p.page {
	case 1:
		return p.ViewScope()
	case 2:
		return p.ViewDescription()
	case 3:
		return p.ViewBreakingChange()
	default:
		return p.ViewType()
	}
}
