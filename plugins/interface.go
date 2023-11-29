package plugins

import "github.com/PossibleLlama/commit-check/model"

type Plugin interface {
	// Get any configuration from viper
	Init() error
	ListCards() ([]model.ScopeItem, error)
}
