package tui

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/PossibleLlama/commit-check/model"
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

type scopeSelectableListItemFormatter struct{}

func (d scopeSelectableListItemFormatter) Height() int                             { return 1 }
func (d scopeSelectableListItemFormatter) Spacing() int                            { return 0 }
func (d scopeSelectableListItemFormatter) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d scopeSelectableListItemFormatter) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(model.ScopeItem)
	if !ok {
		return
	}

	renderer := lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Center)
	if m.Index() == index {
		renderer = renderer.Bold(true)
		if i.ID == "Other" || i.ID == "None" {
			fmt.Fprint(w, renderer.Render(fmt.Sprintf("%s %s", cursor, i.Body)))
		} else {
			fmt.Fprint(w, renderer.Render(fmt.Sprintf("%s %s - %s", cursor, i.ID, i.Body)))
		}
		return
	}

	if i.ID == "Other" || i.ID == "None" {
		fmt.Fprint(w, renderer.Render(fmt.Sprintf("- %s", i.Body)))
	} else {
		fmt.Fprint(w, renderer.Render(fmt.Sprintf("- %s - %s", i.ID, i.Body)))
	}
}
