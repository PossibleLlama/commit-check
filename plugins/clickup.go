package plugins

import (
	"context"
	"errors"
	"net/http"

	"github.com/PossibleLlama/commit-check/model"
	"github.com/raksul/go-clickup/clickup"
)

type Clickup struct {
	ctx    context.Context
	client *clickup.Client
}

func (c Clickup) Init() error {
	c.ctx = context.Background()
	httpClient := &http.Client{}
	// TODO, get api key from viper
	c.client = clickup.NewClient(httpClient, "api_key")
	if c.client == nil {
		return errors.New("Failed to create clickup client")
	}
	return nil
}

func (c Clickup) ListCards() ([]model.ScopeItem, error) {
	items := []model.ScopeItem{}

	taskOptions := clickup.GetTasksOptions{Statuses: []string{"to do", "in progress"}}
	// TODO, get list id from viper (how do we want this to be discovered?)
	tasks, _, err := c.client.Tasks.GetTasks(c.ctx, "list_id", &taskOptions)
	if err != nil {
		return items, err
	}

	for _, task := range tasks {
		items = append(items, model.ScopeItem{Heading: task.ID, Body: task.Name})
	}

	return items, nil
}
