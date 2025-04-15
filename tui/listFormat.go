package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type simpleListItem string

func (e simpleListItem) FilterValue() string { return "" }

type simpleSelectableListItemFormatter struct{}

func (d simpleSelectableListItemFormatter) Height() int                             { return 1 }
func (d simpleSelectableListItemFormatter) Spacing() int                            { return 0 }
func (d simpleSelectableListItemFormatter) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d simpleSelectableListItemFormatter) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(simpleListItem)
	if !ok {
		return
	}

	renderer := lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Center)
	if m.Index() == index {
		renderer = renderer.Bold(true)
		fmt.Fprint(w, renderer.Render(fmt.Sprintf("%s %s", cursor, i)))
		return
	}

	fmt.Fprint(w, renderer.Render(fmt.Sprintf("- %s", i)))
}

type simpleListItemFormatter struct{}

func (d simpleListItemFormatter) Height() int                             { return 1 }
func (d simpleListItemFormatter) Spacing() int                            { return 0 }
func (d simpleListItemFormatter) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d simpleListItemFormatter) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(simpleListItem)
	if !ok {
		return
	}

	renderer := lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Center)

	fmt.Fprint(w, renderer.Render(fmt.Sprintf("- %s", i)))
}
