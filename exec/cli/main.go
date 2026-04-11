package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/PossibleLlama/commit-check/model"
	"github.com/PossibleLlama/commit-check/tui"
	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Build time variables
var (
	VERSION string
)

// Cobra flags
var (
	conventionType string
	dryRun         bool

	configOutputPath    string
	configSkipUI        bool
	configJiraURL       string
	configJiraUsername  string
	configJiraAPIKey    string
	configJiraProjects  []string
	configJiraStatus    []string
	configClickupAPIKey string
	configClickupListID []string
	configClickupAssign string
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "commit-check",
	Short: "Verify your commits have a common format",
	Run: func(cmd *cobra.Command, args []string) {
		if conventionType != "angular" && conventionType != "conventionalcommit" {
			fmt.Println("convention type must be either 'angular' or 'conventionalcommit' but was", conventionType)
			os.Exit(1)
		}
		var cTypes []model.CommitType
		switch conventionType {
		case "angular":
			cTypes = model.TypeAngular
		case "conventionalcommit":
			cTypes = model.TypeConventionalCommit
		}

		commit := &model.Commit{}
		commit.DryRun(dryRun)

		p := tea.NewProgram(tui.NewCommitSummary(commit, cTypes))
		if _, err := p.Run(); err != nil {
			fmt.Println("An unexpected error:", err)
			os.Exit(1)
		}

		if !commit.IsCommittable() || commit.HasQuit() {
			escapedMessages := make([]string, len(commit.CommitStrings()))
			for i, msg := range commit.CommitStrings() {
				escapedMessages[i] = strconv.Quote(msg)
			}
			fmt.Printf("Did not commit changes. This would have been the command.\ngit commit %s\n", strings.Join(escapedMessages, " -m "))
		} else {
			if err := gitCommitGoGit(commit); err != nil {
				fmt.Println(err.Error())
			}
		}
	},
	Version: VERSION,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Setup and edit commit-check configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := model.NewSetupConfigFromViper(viper.GetString, viper.GetStringSlice)

		if cmd.Flags().Changed("jira-url") {
			cfg.JiraURL = configJiraURL
		}
		if cmd.Flags().Changed("jira-username") {
			cfg.JiraUsername = configJiraUsername
		}
		if cmd.Flags().Changed("jira-api-key") {
			cfg.JiraAPIKey = configJiraAPIKey
		}
		if cmd.Flags().Changed("jira-projects") {
			cfg.JiraProjects = configJiraProjects
		}
		if cmd.Flags().Changed("jira-status") {
			cfg.JiraStatus = configJiraStatus
		}
		if cmd.Flags().Changed("clickup-api-key") {
			cfg.ClickupAPIKey = configClickupAPIKey
		}
		if cmd.Flags().Changed("clickup-list-ids") {
			cfg.ClickupListIDs = configClickupListID
		}
		if cmd.Flags().Changed("clickup-assignee") {
			cfg.ClickupAssignee = configClickupAssign
		}

		if configSkipUI {
			if err := model.WriteConfigFile(configOutputPath, cfg); err != nil {
				fmt.Println("unable to write config:", err)
				os.Exit(1)
			}
			fmt.Printf("Configuration written to %s\n", configOutputPath)
			return
		}

		p := tea.NewProgram(tui.NewConfigEditor(&cfg))
		result, err := p.Run()
		if err != nil {
			fmt.Println("An unexpected error:", err)
			os.Exit(1)
		}

		editor, ok := result.(*tui.ConfigEditor)
		if !ok {
			fmt.Println("unexpected config editor state")
			os.Exit(1)
		}

		if !editor.Confirmed() {
			fmt.Println("Configuration was not saved")
			return
		}

		if err := model.WriteConfigFile(configOutputPath, cfg); err != nil {
			fmt.Println("unable to write config:", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration written to %s\n", configOutputPath)
	},
}

func init() {
	defaultPath, err := defaultConfigPath()
	if err != nil {
		defaultPath = "$HOME/.commit-check/config.yaml"
	}

	rootCmd.Flags().StringVarP(&conventionType,
		"type-list",
		"l",
		"angular",
		"accepts either 'angular' or 'conventionalcommit'")
	rootCmd.Flags().BoolVarP(&dryRun,
		"dry-run",
		"d",
		false,
		"run the program without committing")

	configCmd.Flags().StringVar(&configOutputPath,
		"output",
		defaultPath,
		"path to write the configuration file")
	configCmd.Flags().BoolVar(&configSkipUI,
		"skip-ui",
		false,
		"skip the interactive UI and write config from flags")
	configCmd.Flags().StringVar(&configJiraURL,
		"jira-url",
		"",
		"jira base URL")
	configCmd.Flags().StringVar(&configJiraUsername,
		"jira-username",
		"",
		"jira username/email")
	configCmd.Flags().StringVar(&configJiraAPIKey,
		"jira-api-key",
		"",
		"jira API key")
	configCmd.Flags().StringSliceVar(&configJiraProjects,
		"jira-projects",
		[]string{},
		"jira project keys to include")
	configCmd.Flags().StringSliceVar(&configJiraStatus,
		"jira-status",
		[]string{},
		"jira statuses to include")
	configCmd.Flags().StringVar(&configClickupAPIKey,
		"clickup-api-key",
		"",
		"clickup API key")
	configCmd.Flags().StringSliceVar(&configClickupListID,
		"clickup-list-ids",
		[]string{},
		"clickup list IDs to include")
	configCmd.Flags().StringVar(&configClickupAssign,
		"clickup-assignee",
		"",
		"clickup assignee email")

	rootCmd.AddCommand(configCmd)

	cobra.OnInitialize(initConfig)
}

func defaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".commit-check", "config.yaml"), nil
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/commit-check/")
	viper.AddConfigPath("$HOME/.commit-check")

	// Allows for variables such as CC_PLUGINS_CLICKUP_APIKEY
	viper.SetEnvPrefix("CC")

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("unable to use config file: '%s'. %s", viper.ConfigFileUsed(), err.Error())
	}
}

func gitCommitGoGit(commit *model.Commit) error {
	var err error
	var dir string

	// Get current directory of the running binary
	dir, err = os.Getwd()
	if err != nil {
		return err
	}

	// Open the git repository
	var repo *gogit.Repository
	repo, err = gogit.PlainOpenWithOptions(
		dir,
		&gogit.PlainOpenOptions{
			DetectDotGit:          true,
			EnableDotGitCommonDir: true,
		},
	)
	if err != nil {
		return err
	}

	// Create a new worktree
	var worktree *gogit.Worktree
	worktree, err = repo.Worktree()
	if err != nil {
		return err
	}

	// Commit changes
	_, err = worktree.Commit(strings.Join(commit.CommitStrings(), "\n"), &gogit.CommitOptions{})
	if err != nil {
		return err
	}

	// Commit via os so that signing and other git hooks can be used as gogit doesnt pick up on this configuration
	//#nosec G204 -- This is not editable by the user
	runOsCmd := exec.Command("git", "commit", "--amend", "--no-edit")
	var osCmdOutput []byte
	osCmdOutput, err = runOsCmd.CombinedOutput()
	if err != nil {
		return err
	} else {
		fmt.Println(string(osCmdOutput))
	}

	return nil
}
