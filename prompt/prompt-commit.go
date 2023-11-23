package prompt

import (
	"github.com/PossibleLlama/commit-check/model"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(0, 2)

type PromptCommit struct {
	page int

	typeNameOptions list.Model
	scopeOptions    list.Model
	breakingOptions list.Model

	inputSingleLine textinput.Model
	inputMultiLine  textarea.Model

	width  int
	height int

	commit *model.Commit
}

func NewPromptCommit(typeNameOptions []model.CommitType, cmt *model.Commit) *PromptCommit {
	return &PromptCommit{
		page: 0,

		typeNameOptions: SetupListOfTypes(typeNameOptions),
		scopeOptions:    SetupListOfScopes(),
		breakingOptions: SetupListOfBreakingChanges(),

		inputSingleLine: textinput.New(),
		inputMultiLine:  textarea.New(),

		width:  0,
		height: 0,

		commit: cmt,
	}
}

func (p PromptCommit) Init() tea.Cmd {
	return tea.Batch(
	// p.CheckJira,
	)
}

func (p PromptCommit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w, h := docStyle.GetFrameSize()
		p.width = msg.Width - w
		p.height = msg.Height - h

		p.typeNameOptions.SetWidth(p.width)
		p.typeNameOptions.SetHeight(p.height - 2)

		p.scopeOptions.SetWidth(p.width)
		p.scopeOptions.SetHeight(p.height - 2)

		p.breakingOptions.SetWidth(p.width)
		p.breakingOptions.SetHeight(p.height - 2)

		p.inputMultiLine.SetWidth(p.width)
		p.inputMultiLine.SetHeight(p.height - 2)

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
	case ScopeItem:
		cmd := p.scopeOptions.InsertItem(0, msg)
		return p, cmd
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
