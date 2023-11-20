package prompt

import (
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (p PromptCommit) CheckJira() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get("https://google.com")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return res.StatusCode
}

func (p PromptCommit) UpdateScope(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down":
			if p.cursor < len(p.scopeOptions)-1 {
				p.cursor++
			}
		case "enter":
			if p.scopeOptions[p.cursor] == "Other" {
				if p.inputSingleLine.Focused() {
					p.inputSingleLine.Blur()
					p.commit.Scope = p.inputSingleLine.Value()
					p.page++
					p.cursor = 0
					p.inputMultiLine.Focus()
					return p, nil
				} else {
					p.inputSingleLine.Focus()
					p.inputSingleLine.Placeholder = "Scope"
				}
			} else if p.scopeOptions[p.cursor] != "None" {
				p.commit.Scope = p.scopeOptions[p.cursor]
				p.page++
				p.cursor = 0
				p.inputMultiLine.Focus()
				return p, nil
			} else {
				p.page++
				p.cursor = 0
				p.inputMultiLine.Focus()
				return p, nil
			}
		default:
			if p.inputSingleLine.Focused() {
				var cmd tea.Cmd
				p.inputSingleLine, cmd = p.inputSingleLine.Update(msg)
				return p, cmd
			}
		}
	case int:
		p.scopeOptions = append(p.scopeOptions, fmt.Sprint(msg))
	}
	return p, nil
}

func (p PromptCommit) ViewScope() string {
	s := ""
	if p.inputSingleLine.Focused() {
		s = "Scope of the change:\n"
		s += p.inputSingleLine.View()
	} else {
		s = "Select the scope of the change:\n"
		for i, scope := range p.scopeOptions {
			if p.cursor == i {
				s += "> "
			} else {
				s += "  "
			}
			s += scope + "\n"
		}
	}
	return s
}
