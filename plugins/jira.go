package plugins

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/PossibleLlama/commit-check/model"
	"github.com/andygrunwald/go-jira"
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
		return ErrorPluginMissingCreds
	} else if url == "" {
		return ErrorPluginMissingConfig
	}

	cl := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := jira.NewClient(cl.Client(), url)
	if err != nil {
		return err
	}

	j.client = client
	return nil
}

func (j *Jira) ListCards() tea.Msg {
	last := 0
	items := []model.ScopeItem{}

	jql := "assignee = currentUser()"

	projects := []string{}
	for _, p := range viper.GetStringSlice("plugins.jira.projects") {
		if strings.TrimSpace(p) != "" {
			projects = append(projects, fmt.Sprintf("\"%s\"", p))
		}
	}
	status := []string{}
	for _, s := range viper.GetStringSlice("plugins.jira.status") {
		if strings.TrimSpace(s) != "" {
			status = append(status, fmt.Sprintf("\"%s\"", s))
		}
	}

	if len(projects) > 0 {
		jql = fmt.Sprintf("project IN (%s) AND %s", strings.Join(projects, ", "), jql)
	}
	if len(status) > 0 {
		jql = fmt.Sprintf("status IN (%s) AND %s", strings.Join(status, ", "), jql)
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
			items = append(items, model.ScopeItem{ID: task.Key, Body: task.Fields.Summary})
		}

		last = resp.StartAt + len(tasks)
		if last >= resp.Total {
			break
		}
	}

	return items
}
