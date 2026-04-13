package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"github.com/PossibleLlama/commit-check/model"
)

type configEditorState uint

const (
	configPluginState configEditorState = iota
	configFieldState
	configEditState
)

type configPlugin struct {
	Name   string
	Fields []model.ConfigField
}

type ConfigEditor struct {
	state configEditorState
	quit  bool

	cfg          *model.SetupConfig
	plugins      []configPlugin
	pluginCursor int
	fieldCursor  int

	input textarea.Model
}

func NewConfigEditor(cfg *model.SetupConfig) *ConfigEditor {
	input := textarea.New()
	input.SetWidth(70)
	input.SetHeight(1)
	input.ShowLineNumbers = false

	var jiraFields []model.ConfigField
	var clickupFields []model.ConfigField
	for _, field := range model.SetupConfigFields {
		switch {
		case strings.HasPrefix(field.Key, "jira."):
			jiraFields = append(jiraFields, field)
		case strings.HasPrefix(field.Key, "clickup."):
			clickupFields = append(clickupFields, field)
		}
	}

	return &ConfigEditor{
		state: configPluginState,
		cfg:   cfg,
		plugins: []configPlugin{
			{
				Name:   "Jira",
				Fields: jiraFields,
			},
			{
				Name:   "ClickUp",
				Fields: clickupFields,
			},
		},
		input: input,
	}
}

func (c *ConfigEditor) Init() tea.Cmd {
	return textarea.Blink
}

func (c *ConfigEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.input.SetWidth(msg.Width - focusedStyle.GetBorderLeftSize() - focusedStyle.GetBorderRightSize() - 10)
	case tea.KeyMsg:
		key := msg.Key()

		switch {
		case key.Code == 'c' && key.Mod == tea.ModCtrl:
			return c, tea.Quit
		}

		switch c.state {
		case configPluginState:
			switch msg.String() {
			case "up", "k":
				if c.pluginCursor > 0 {
					c.pluginCursor--
				}
			case "down", "j":
				if c.pluginCursor < len(c.plugins)+1 {
					c.pluginCursor++
				}
			case "enter", "e":
				switch c.pluginCursor {
				case len(c.plugins):
					c.quit = true
					return c, tea.Quit
				case len(c.plugins) + 1:
					return c, tea.Quit
				default:
					c.fieldCursor = 0
					c.state = configFieldState
				}
			}
		case configFieldState:
			switch msg.String() {
			case "esc":
				c.state = configPluginState
			case "up", "k":
				if c.fieldCursor > 0 {
					c.fieldCursor--
				}
			case "down", "j":
				selectedPlugin := c.plugins[c.pluginCursor]
				if c.fieldCursor < len(selectedPlugin.Fields) {
					c.fieldCursor++
				}
			case "enter", "e":
				selectedPlugin := c.plugins[c.pluginCursor]
				if c.fieldCursor == len(selectedPlugin.Fields) {
					c.state = configPluginState
				} else {
					field := selectedPlugin.Fields[c.fieldCursor]
					c.input.Placeholder = field.Label
					c.input.SetValue(c.cfg.FieldValue(field.Key))
					c.input.Focus()
					c.state = configEditState
				}
			}
		case configEditState:
			switch msg.String() {
			case "enter":
				field := c.plugins[c.pluginCursor].Fields[c.fieldCursor]
				c.cfg.SetFieldValue(field.Key, c.input.Value())
				c.input.Blur()
				c.input.Reset()
				c.state = configFieldState
			case "esc":
				c.input.Blur()
				c.input.Reset()
				c.state = configFieldState
			default:
				var cmd tea.Cmd
				c.input, cmd = c.input.Update(msg)
				return c, cmd
			}
		}
	}

	return c, nil
}

func (c *ConfigEditor) View() tea.View {
	var view string
	switch c.state {
	case configPluginState:
		lines := []string{"Configuration setup", "", "Plugins"}
		for idx, plugin := range c.plugins {
			line := "- " + plugin.Name
			if idx == c.pluginCursor {
				line = cursor + " " + plugin.Name
			}
			lines = append(lines, line)
		}

		confirmLine := "- Confirm and write config"
		if c.pluginCursor == len(c.plugins) {
			confirmLine = cursor + " Confirm and write config"
		}
		lines = append(lines, "", confirmLine)

		cancelLine := "- Cancel"
		if c.pluginCursor == len(c.plugins)+1 {
			cancelLine = cursor + " Cancel"
		}
		lines = append(lines, cancelLine)
		lines = append(lines, "", "Keys: up/down to move, enter or e to select, ctrl+c to quit")
		view = strings.Join(lines, "\n")
	case configFieldState:
		selectedPlugin := c.plugins[c.pluginCursor]
		lines := []string{"Configuration setup", "", selectedPlugin.Name + " fields"}

		for idx, field := range selectedPlugin.Fields {
			current := c.cfg.FieldValue(field.Key)
			if field.Secret && current != "" {
				current = strings.Repeat("*", 8)
			}
			if strings.TrimSpace(current) == "" {
				current = "-"
			}

			line := fmt.Sprintf("- %s: %s", field.Label, current)
			if idx == c.fieldCursor {
				line = fmt.Sprintf("%s %s: %s", cursor, field.Label, current)
			}
			lines = append(lines, line)
		}

		backLine := "- Back"
		if c.fieldCursor == len(selectedPlugin.Fields) {
			backLine = cursor + " Back"
		}
		lines = append(lines, "", backLine)
		lines = append(lines, "", "Keys: up/down to move, enter or e to edit/select, esc to return")
		view = strings.Join(lines, "\n")
	case configEditState:
		field := c.plugins[c.pluginCursor].Fields[c.fieldCursor]
		view = fmt.Sprintf("Edit: %s\n\n%s\n\nenter to save, esc to cancel", field.Label, c.input.View())
	}

	configuredView := tea.NewView(focusedStyle.Render(view) + "\n\n" + defaultStyle.Render(footerText))
	configuredView.WindowTitle = "commit-check config"
	configuredView.AltScreen = true
	return configuredView
}

func (c ConfigEditor) Confirmed() bool {
	return c.quit
}
