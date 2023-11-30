package plugins

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

type Plugin interface {
	// Get any configuration from viper
	Init() error
	ListCards() tea.Msg
}

var (
	PluginErrorMissingCreds = errors.New("missing credentials")
	PluginErrorInvalidCreds = errors.New("invalid credentials")

	PluginErrorMissingConfig = errors.New("missing required configuration. check README for valid configuration")
)
