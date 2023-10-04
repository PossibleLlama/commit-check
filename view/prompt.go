package view

import (
	"fmt"

	model "github.com/PossibleLlama/commit-check/model"

	tea "github.com/charmbracelet/bubbletea"
)

const maxPages = 4

type commitModel struct {
	// Generic Bubble Tea fields
	cursor int // which item our cursor is pointing at
	page   int // which page we are on

	commit *model.Commit

	// 1st page
	cTypeOptions []model.CommitType // Available types of commit
	// 2nd page
	cScopeOptions []string // Available scopes of commit
	// 3rd page
	// This is a text field without options
	// 4th page
	cBreakingOptions []string // Available breaking changes
}

func InitCommitModel(commitTypes []model.CommitType, c *model.Commit) *commitModel {
	return &commitModel{
		cursor: 0,
		page:   1,

		commit: c,

		cTypeOptions: commitTypes,
		// In the future get scopes from Jira/Clickup
		cScopeOptions:    []string{"Other"},
		cBreakingOptions: []string{"No", "Yes"},
	}
}

func (c commitModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (c commitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return c, tea.Quit

		// These keys should move to the previous page
		case "backspace":
			if c.page > 1 {
				c.page--
				c.cursor = 0
			}
		}

		switch c.page {
		case 1:
			switch msg.String() {
			case "up", "k":
				if c.cursor > 0 {
					c.cursor--
				}
			case "down", "j":
				if c.cursor < len(c.cTypeOptions)-1 {
					c.cursor++
				}
			case "enter":
				c.commit.Type = c.cTypeOptions[c.cursor]
				c.page++
				c.cursor = 0
			}
		case 2:
			switch msg.String() {
			case "up", "k":
				if c.cursor > 0 {
					c.cursor--
				}
			case "down", "j":
				if c.cursor < len(c.cScopeOptions)-1 {
					c.cursor++
				}
			case "enter":
				c.commit.Scope = c.cScopeOptions[c.cursor]
				c.page++
				c.cursor = 0
			}
		case 3:
			switch msg.String() {
			case "enter":
				c.page++
				c.cursor = 0
			case "left":
				if len(c.commit.Description) > 0 {
					c.commit.Description = c.commit.Description[:len(c.commit.Description)-1]
				}
			default:
				c.commit.Description += msg.String()
			}
		case 4:
			switch msg.String() {
			case "up", "k":
				if c.cursor > 0 {
					c.cursor--
				}
			case "down", "j":
				if c.cursor < len(c.cBreakingOptions)-1 {
					c.cursor++
				}
			case "enter":
				c.commit.IsBreakingChange = c.cBreakingOptions[c.cursor] == "Yes"
				return c, tea.Quit
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return c, nil
}

func (c commitModel) View() string {
	s := fmt.Sprintf("Page: %d of %d\n\n", c.page, maxPages)

	switch c.page {
	case 1:
		s += c.RenderPage1()
	case 2:
		s += c.RenderPage2()
	case 3:
		s += c.RenderPage3()
	case 4:
		s += c.RenderPage4()
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (c commitModel) RenderPage1() string {
	s := fmt.Sprintln("Type of commit:")

	for i, t := range c.cTypeOptions {
		cursor := " "
		if c.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, t)
	}

	// Send the UI for rendering
	return s
}

func (c commitModel) RenderPage2() string {
	s := fmt.Sprintln("Scope of commit:")

	for i, t := range c.cScopeOptions {
		cursor := " "
		if c.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, t)
	}

	return s
}

func (c commitModel) RenderPage3() string {
	s := fmt.Sprintln("Description of commit:")

	s += fmt.Sprintf("%s\n", c.commit.Description)

	return s
}

func (c commitModel) RenderPage4() string {
	s := fmt.Sprintln("Is breaking change?")

	for i, t := range c.cBreakingOptions {
		cursor := " "
		if c.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, t)
	}

	return s
}
