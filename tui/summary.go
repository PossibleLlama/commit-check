package tui

import (
	"fmt"
	"strings"

	"github.com/PossibleLlama/commit-check/model"
	"github.com/PossibleLlama/commit-check/plugins"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

// Which section of the view we're on
type modelState uint

const (
	summaryState     modelState = iota
	typeState        modelState = iota
	scopeState       modelState = iota
	descriptionState modelState = iota
	helpState        modelState = iota

	cursor     = "> "
	headerText = ""
	footerText = "Press \"h\" for help"

	helpText = "Help!\n\n" +
		"Quit (ctrl + \"c\")\n" +
		"Help (\"h\")\n" +
		"Confirm commit (enter)\n" +
		"Dry run (\"D\")\n" +
		"Back (escape)\n" +
		"\n" +
		"Summary\n" +
		"Edit commit type (\"t\")\n" +
		"Edit commit scope (\"s\")\n" +
		"Edit commit description (\"d\")\n" +
		"Edit commit breaking change status (\"b\")\n"
)

type Summary struct {
	state modelState
	Quit  bool

	cmt         *model.Commit
	commitTypes []model.CommitType
	cTypeList   list.Model
	cScopeList  list.Model

	text textarea.Model
}

func NewCommitSummary(cmt *model.Commit, commitTypes []model.CommitType) *Summary {
	ta := textarea.New()
	ta.SetPromptFunc(1, func(lineIdx int) string { return cursor })
	ta.SetWidth(30)

	cType := []list.Item{}
	for _, c := range commitTypes {
		cType = append(cType, simpleListItem(c))
	}
	cTypeList := list.New(cType, simpleSelectableListItemFormatter{}, 60, len(cType))
	cTypeList.SetShowStatusBar(false)
	cTypeList.SetShowHelp(false)
	cTypeList.SetShowFilter(false)
	cTypeList.SetShowTitle(false)
	cTypeList.SetShowPagination(false)

	cScope := []list.Item{
		model.ScopeItem{ID: "None", Body: "No scope"},
		model.ScopeItem{ID: "Other", Body: "Manual input"},
	}
	cScopeList := list.New(cScope, scopeSelectableListItemFormatter{}, 60, len(cScope))
	cScopeList.SetShowStatusBar(false)
	cScopeList.SetShowHelp(false)
	cScopeList.SetShowFilter(false)
	cScopeList.SetShowTitle(false)
	cScopeList.SetShowPagination(false)

	return &Summary{
		state:       summaryState,
		cmt:         cmt,
		commitTypes: commitTypes,
		cTypeList:   cTypeList,

		cScopeList: cScopeList,

		text: ta,
	}
}

func (s *Summary) Init() tea.Cmd {
	pluginSources := []plugins.Plugin{}
	if viper.Sub("plugins.clickup") != nil {
		p := plugins.NewClickup()
		err := p.Init()
		if err == nil {
			pluginSources = append(pluginSources, p)
		}
		// TODO, log error
	}
	if viper.Sub("plugins.jira") != nil {
		p := plugins.NewJira()
		err := p.Init()
		if err == nil {
			pluginSources = append(pluginSources, p)
		}
		// TODO, log error
	}

	msgs := []tea.Cmd{
		tea.SetWindowTitle("commit-check"),
	}
	for _, p := range pluginSources {
		msgs = append(msgs, p.ListCards)
	}
	return tea.Batch(
		msgs...,
	)
}

func (s *Summary) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.text.SetWidth(msg.Width - focusedStyle.GetBorderLeftSize() - focusedStyle.GetBorderRightSize())
		s.text.SetHeight(msg.Height - focusedStyle.GetBorderTopSize() - focusedStyle.GetBorderBottomSize())
	case tea.KeyMsg:
		// Global keybindings
		switch msg.String() {
		case tea.KeyEsc.String():
			if s.state != summaryState {
				s.state = summaryState
				return s, nil
			}
			fallthrough
		case tea.KeyCtrlC.String():
			s.cmt.Quit(true)
			return s, tea.Quit
		}

		// Page keybindings
		switch s.state {
		// Summary page
		case summaryState:
			switch msg.String() {
			case "enter":
				return s, tea.Quit
			case "t":
				s.state = typeState
			case "s":
				s.state = scopeState
				s.text.Blur()
			case "d":
				s.state = descriptionState
				s.text.Placeholder = "Enter description"
				s.text.Focus()
				s.text.SetValue(s.cmt.Description)
			case "b":
				s.cmt.IsBreakingChange = !s.cmt.IsBreakingChange
			case "D":
				s.cmt.DryRun(!s.cmt.IsDryRun())
			case "h":
				s.state = helpState
			}
		// Type page
		case typeState:
			switch msg.String() {
			case "enter":
				s.cmt.Type = s.commitTypes[s.cTypeList.Cursor()]
				s.state = summaryState
			case "up", "down", "j", "k":
				s.cTypeList, cmd = s.cTypeList.Update(msg)
			}
		// Scope page
		case scopeState:
			if s.text.Focused() { // Manual input
				switch msg.String() {
				case "enter":
					s.cmt.Scope = s.text.Value()
					s.text.Reset()
					s.state = summaryState
				default:
					s.text, cmd = s.text.Update(msg)
					return s, cmd
				}
			} else { // List input
				switch msg.String() {
				case "enter":
					selectedScope := s.cScopeList.SelectedItem().(model.ScopeItem)
					switch selectedScope.ID {
					case "None":
						// No scope
						s.cmt.Scope = ""
						s.text.Reset()
						s.text.Blur()
						s.state = summaryState
					case "Other":
						// Manual input
						s.text.Placeholder = "Enter scope of change"
						s.text.Focus()
						s.text.SetValue(s.cmt.Scope)
					default:
						// From a plugin
						s.cmt.Scope = selectedScope.ID
						s.state = summaryState
					}
				case "up", "down", "j", "k":
					s.cScopeList, cmd = s.cScopeList.Update(msg)
				}
			}
		// Description page
		case descriptionState:
			switch msg.String() {
			case "enter":
				s.cmt.Description = s.text.Value()
				s.text.Reset()
				s.state = summaryState
			default:
				s.text, cmd = s.text.Update(msg)
				return s, cmd
			}
		}
	case []model.ScopeItem:
		for _, item := range msg {
			cmd = s.cScopeList.InsertItem(0, item)
		}
		s.cScopeList.SetHeight(len(msg) + s.cScopeList.Height())
	}
	return s, cmd
}

func (s *Summary) View() string {
	var v string
	switch s.state {
	case summaryState:
		columns := []table.Column{
			{Title: "(t)ype", Width: 10},
			{Title: "(s)cope", Width: 10},
			{Title: "(d)escription", Width: 40},
			{Title: "(b)reaking change", Width: 18},
			{Title: "(D)ry run", Width: 10},
		}
		rows := []table.Row{
			{
				StringToStringOrDash(s.cmt.Type.String()),
				StringToStringOrDash(s.cmt.Scope),
				StringToStringOrDash(s.cmt.Description),
				BoolToYesNo(s.cmt.IsBreakingChange),
				BoolToYesNo(s.cmt.IsDryRun()),
			},
		}
		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithHeight(2),
		)

		ts := table.DefaultStyles()
		ts.Header = ts.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(false).
			Bold(false)
		ts.Selected = ts.Selected.
			Foreground(lipgloss.Color("240")).
			Bold(false)
		t.SetStyles(ts)

		cmtAdditionalInformation := ""
		if len(s.cmt.IsCommittableReason()) > 0 {
			reasons := []list.Item{}
			for _, reason := range s.cmt.IsCommittableReason() {
				reasons = append(reasons, simpleListItem(reason))
			}
			l := list.New(reasons, simpleListItemFormatter{}, 60, len(reasons))
			l.SetShowStatusBar(false)
			l.SetShowHelp(false)
			l.SetShowFilter(false)
			l.SetShowTitle(false)
			l.SetShowPagination(false)

			cmtAdditionalInformation += fmt.Sprintf("\n\nCommit is not valid because:\n%s", l.View())
		} else {
			cmtAdditionalInformation += "\n\n(enter) to commit"
		}
		v = fmt.Sprintf("%s\n\nGit log: '%s'%s", t.View(), StringToStringOrDash(s.cmt.String()), cmtAdditionalInformation)
	case typeState:
		v = fmt.Sprintf("Type:\n%s", s.cTypeList.View())
	case scopeState:
		if s.text.Focused() {
			v = fmt.Sprintf("Scope:\n%s", s.text.Value())
		} else {
			v = fmt.Sprintf("Scope:\n%s", s.cScopeList.View())
		}
	case descriptionState:
		v = fmt.Sprintf("Description:\n%s", s.text.Value())
	case helpState:
		v = helpText
	default:
		v = "loading"
	}
	return focusedStyle.Render(v) + "\n\n" + defaultStyle.Render(footerText)
}

func BoolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func StringToStringOrDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return strings.TrimSpace(s)
}
