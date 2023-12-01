package plugins

import (
	"context"
	"sync"

	"github.com/PossibleLlama/commit-check/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/raksul/go-clickup/clickup"
	"github.com/spf13/viper"
)

type Clickup struct {
	ctx    context.Context
	client *clickup.Client
}

func NewClickup() Plugin {
	return &Clickup{}
}

func (c *Clickup) Init() error {
	c.ctx = context.Background()

	key := viper.GetString("plugins.clickup.apiKey")
	if key == "" {
		return ErrorPluginMissingCreds
	} else if len(viper.GetStringSlice("plugins.clickup.listIds")) == 0 {
		return ErrorPluginMissingConfig
	}

	c.client = clickup.NewClient(nil, key)
	if c.client == nil {
		return ErrorPluginInvalidCreds
	}

	return nil
}

func (c *Clickup) ListCards() tea.Msg {
	items := []model.ScopeItem{}

	taskOptions := clickup.GetTasksOptions{Statuses: []string{"to do", "in progress"}}

	listIds := viper.GetStringSlice("plugins.clickup.listIds")
	email := viper.GetString("plugins.clickup.assignee")

	var wg sync.WaitGroup

	for _, listId := range listIds {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			tasks, _, err := c.client.Tasks.GetTasks(c.ctx, id, &taskOptions)
			if err != nil {
				return
			}

			for _, task := range tasks {
				if email != "" {
					for _, assignee := range task.Assignees {
						if assignee.Email == email {
							items = append(items, model.ScopeItem{Heading: task.ID, Body: task.Name})
							break
						}
					}
				} else {
					items = append(items, model.ScopeItem{Heading: task.ID, Body: task.Name})
				}
			}
		}(listId)
	}

	wg.Wait()
	return items
}
