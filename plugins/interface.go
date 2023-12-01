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
	ErrorPluginMissingCreds = errors.New("missing credentials")
	ErrorPluginInvalidCreds = errors.New("invalid credentials")

	ErrorPluginMissingConfig = errors.New("missing required configuration. check README for valid configuration")
)
