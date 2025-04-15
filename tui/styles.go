package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	defaultStyle = lipgloss.NewStyle().
			Align(lipgloss.Left, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder()).
			Padding(0, 1)
	focusedStyle = defaultStyle.
			BorderStyle(lipgloss.ASCIIBorder()).
			BorderForeground(lipgloss.Color("20"))
)
