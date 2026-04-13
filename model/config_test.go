package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestSetupConfigSetFieldValue(t *testing.T) {
	cfg := SetupConfig{}

	cfg.SetFieldValue("jira.url", " https://example.atlassian.net ")
	cfg.SetFieldValue("jira.username", " alice@example.com ")
	cfg.SetFieldValue("jira.apiKey", " test-key ")
	cfg.SetFieldValue("jira.projects", "ABC, DEF, , GHI")
	cfg.SetFieldValue("jira.status", "In Progress, Done")
	cfg.SetFieldValue("clickup.apiKey", " clickup-key ")
	cfg.SetFieldValue("clickup.listIds", "123, 456")
	cfg.SetFieldValue("clickup.assignee", " owner@example.com ")

	assert.Equal(t, "https://example.atlassian.net", cfg.JiraURL)
	assert.Equal(t, "alice@example.com", cfg.JiraUsername)
	assert.Equal(t, "test-key", cfg.JiraAPIKey)
	assert.Equal(t, []string{"ABC", "DEF", "GHI"}, cfg.JiraProjects)
	assert.Equal(t, []string{"In Progress", "Done"}, cfg.JiraStatus)
	assert.Equal(t, "clickup-key", cfg.ClickupAPIKey)
	assert.Equal(t, []string{"123", "456"}, cfg.ClickupListIDs)
	assert.Equal(t, "owner@example.com", cfg.ClickupAssignee)
}

func TestWriteConfigFile(t *testing.T) {
	cfg := SetupConfig{
		JiraURL:         "https://example.atlassian.net",
		JiraUsername:    "alice@example.com",
		JiraAPIKey:      "jira-key",
		JiraProjects:    []string{"ABC"},
		JiraStatus:      []string{"In Progress"},
		ClickupAPIKey:   "clickup-key",
		ClickupListIDs:  []string{"999"},
		ClickupAssignee: "alice@example.com",
	}

	output := filepath.Join(t.TempDir(), "config.yaml")
	err := WriteConfigFile(output, cfg)
	require.NoError(t, err)

	assert.FileExists(t, output)

	fileInfo, err := os.Stat(output)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), fileInfo.Mode().Perm())

	body, err := os.ReadFile(output)
	require.NoError(t, err)

	var parsed configFile
	err = yaml.Unmarshal(body, &parsed)
	require.NoError(t, err)

	assert.Equal(t, cfg.JiraURL, parsed.Plugins.Jira.URL)
	assert.Equal(t, cfg.JiraUsername, parsed.Plugins.Jira.Username)
	assert.Equal(t, cfg.JiraAPIKey, parsed.Plugins.Jira.APIKey)
	assert.Equal(t, cfg.JiraProjects, parsed.Plugins.Jira.Projects)
	assert.Equal(t, cfg.JiraStatus, parsed.Plugins.Jira.Status)
	assert.Equal(t, cfg.ClickupAPIKey, parsed.Plugins.Clickup.APIKey)
	assert.Equal(t, cfg.ClickupListIDs, parsed.Plugins.Clickup.ListIDs)
	assert.Equal(t, cfg.ClickupAssignee, parsed.Plugins.Clickup.Assignee)
}
