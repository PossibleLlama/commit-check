package prompt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PossibleLlama/commit-check/model"

	tea "github.com/charmbracelet/bubbletea"
)

type PromptScope struct {
	cursor       int
	scopeOptions []string

	jiraUrl    string
	clickupUrl string

	width  int
	height int

	commit *model.Commit
}

func NewPromptScope(jiraUrl, clickupUrl string, cmt *model.Commit) *PromptScope {
	return &PromptScope{
		cursor:       0,
		scopeOptions: []string{"None"},

		jiraUrl:    jiraUrl,
		clickupUrl: clickupUrl,

		commit: cmt,
	}
}

func (p PromptScope) Init() tea.Cmd {
	// Go to JIRA/Clickup
	return p.checkJira
}

func (p PromptScope) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return p, tea.Quit
		case "up":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down":
			if p.cursor < len(p.scopeOptions)-1 {
				p.cursor++
			}
		case "enter":
			if p.scopeOptions[p.cursor] != "None" {
				p.commit.Scope = p.scopeOptions[p.cursor]
			}
			return p, tea.Quit
		}
	case int:
		p.scopeOptions = append(p.scopeOptions, fmt.Sprint(msg))
	}
	return p, nil
}

func (p PromptScope) View() string {
	if p.width == 0 || p.height == 0 {
		return "loading..."
	}
	s := "Select the scope of the change:\n"
	for i, scope := range p.scopeOptions {
		if p.cursor == i {
			s += "> "
		} else {
			s += "  "
		}
		s += scope + "\n"
	}
	return s
}

func (p PromptScope) checkJira() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(p.jiraUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return res.StatusCode
}
