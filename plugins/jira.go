package plugins

import (
	"fmt"

	"github.com/PossibleLlama/commit-check/model"
	"github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type Jira struct {
	client *jira.Client
}

func NewJira() Plugin {
	return &Jira{}
}

func (j *Jira) Init() error {
	url := viper.GetString("plugins.jira.url")
	username := viper.GetString("plugins.jira.username")
	password := viper.GetString("plugins.jira.apiKey")

	if username == "" || password == "" {
		return PluginErrorMissingCreds
	} else if url == "" {
		return PluginErrorMissingConfig
	}

	cl := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := jira.NewClient(cl.Client(), url)
	if err != nil {
		return err
	} else if !client.Authentication.Authenticated() {
		return PluginErrorInvalidCreds
	}

	j.client = client
	return nil
}

func (j *Jira) ListCards() tea.Msg {
	last := 0
	items := []model.ScopeItem{}

	jql := "assignee = currentUser()"
	proj := viper.GetString("plugins.jira.project")
	status := viper.GetString("plugins.jira.status")

	if proj != "" {
		jql = fmt.Sprintf("project = %s AND %s", proj, jql)
	}
	if status != "" {
		jql = fmt.Sprintf("status = %s AND %s", status, jql)
	}

	for {
		options := &jira.SearchOptions{
			StartAt:    last,
			MaxResults: 1000,
		}

		tasks, resp, err := j.client.Issue.Search(jql, options)
		if err != nil {
			return nil
		}

		for _, task := range tasks {
			items = append(items, model.ScopeItem{Heading: task.Key, Body: task.Fields.Summary})
		}

		last = resp.StartAt + len(tasks)
		if last >= resp.Total {
			break
		}
	}

	return items
}
