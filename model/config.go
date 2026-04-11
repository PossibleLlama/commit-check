package model

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type SetupConfig struct {
	JiraURL         string
	JiraUsername    string
	JiraAPIKey      string
	JiraProjects    []string
	JiraStatus      []string
	ClickupAPIKey   string
	ClickupListIDs  []string
	ClickupAssignee string
}

type ConfigField struct {
	Key    string
	Label  string
	Secret bool
}

var SetupConfigFields = []ConfigField{
	{Key: "jira.url", Label: "Jira URL"},
	{Key: "jira.username", Label: "Jira Username"},
	{Key: "jira.apiKey", Label: "Jira API Key", Secret: true},
	{Key: "jira.projects", Label: "Jira Projects (comma separated)"},
	{Key: "jira.status", Label: "Jira Statuses (comma separated)"},
	{Key: "clickup.apiKey", Label: "ClickUp API Key", Secret: true},
	{Key: "clickup.listIds", Label: "ClickUp List IDs (comma separated)"},
	{Key: "clickup.assignee", Label: "ClickUp Assignee Email"},
}

func NewSetupConfigFromViper(getString func(string) string, getStringSlice func(string) []string) SetupConfig {
	return SetupConfig{
		JiraURL:         getString("plugins.jira.url"),
		JiraUsername:    getString("plugins.jira.username"),
		JiraAPIKey:      getString("plugins.jira.apiKey"),
		JiraProjects:    getStringSlice("plugins.jira.projects"),
		JiraStatus:      getStringSlice("plugins.jira.status"),
		ClickupAPIKey:   getString("plugins.clickup.apiKey"),
		ClickupListIDs:  getStringSlice("plugins.clickup.listIds"),
		ClickupAssignee: getString("plugins.clickup.assignee"),
	}
}

func (c SetupConfig) FieldValue(key string) string {
	switch key {
	case "jira.url":
		return c.JiraURL
	case "jira.username":
		return c.JiraUsername
	case "jira.apiKey":
		return c.JiraAPIKey
	case "jira.projects":
		return strings.Join(c.JiraProjects, ",")
	case "jira.status":
		return strings.Join(c.JiraStatus, ",")
	case "clickup.apiKey":
		return c.ClickupAPIKey
	case "clickup.listIds":
		return strings.Join(c.ClickupListIDs, ",")
	case "clickup.assignee":
		return c.ClickupAssignee
	default:
		return ""
	}
}

func (c *SetupConfig) SetFieldValue(key string, value string) {
	switch key {
	case "jira.url":
		c.JiraURL = strings.TrimSpace(value)
	case "jira.username":
		c.JiraUsername = strings.TrimSpace(value)
	case "jira.apiKey":
		c.JiraAPIKey = strings.TrimSpace(value)
	case "jira.projects":
		c.JiraProjects = csvToSlice(value)
	case "jira.status":
		c.JiraStatus = csvToSlice(value)
	case "clickup.apiKey":
		c.ClickupAPIKey = strings.TrimSpace(value)
	case "clickup.listIds":
		c.ClickupListIDs = csvToSlice(value)
	case "clickup.assignee":
		c.ClickupAssignee = strings.TrimSpace(value)
	}
}

type configFile struct {
	Plugins configPlugins `yaml:"plugins"`
}

type configPlugins struct {
	Jira    configJira    `yaml:"jira"`
	Clickup configClickup `yaml:"clickup"`
}

type configJira struct {
	URL      string   `yaml:"url"`
	Username string   `yaml:"username"`
	APIKey   string   `yaml:"apiKey"`
	Projects []string `yaml:"projects"`
	Status   []string `yaml:"status"`
}

type configClickup struct {
	APIKey   string   `yaml:"apiKey"`
	ListIDs  []string `yaml:"listIds"`
	Assignee string   `yaml:"assignee"`
}

func (c SetupConfig) ToConfigFile() configFile {
	return configFile{
		Plugins: configPlugins{
			Jira: configJira{
				URL:      c.JiraURL,
				Username: c.JiraUsername,
				APIKey:   c.JiraAPIKey,
				Projects: c.JiraProjects,
				Status:   c.JiraStatus,
			},
			Clickup: configClickup{
				APIKey:   c.ClickupAPIKey,
				ListIDs:  c.ClickupListIDs,
				Assignee: c.ClickupAssignee,
			},
		},
	}
}

func WriteConfigFile(path string, cfg SetupConfig) error {
	body, err := yaml.Marshal(cfg.ToConfigFile())
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, body, 0o600)
}

func csvToSlice(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}

	split := strings.Split(value, ",")
	cleaned := []string{}
	for _, s := range split {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}
