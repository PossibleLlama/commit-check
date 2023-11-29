package plugins

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Plugin interface {
	// Get any configuration from viper
	Init() error
	ListCards() tea.Msg
}
