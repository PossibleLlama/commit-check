package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/PossibleLlama/commit-check/model"
	"github.com/PossibleLlama/commit-check/prompt"
	tea "github.com/charmbracelet/bubbletea"
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
			break
		case "conventionalcommit":
			cTypes = model.TypeConventionalCommit
			break
		}

		commit := &model.Commit{}

		p := tea.NewProgram(prompt.NewPromptCommit(cTypes, commit), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("An unexpected error:", err)
			os.Exit(1)
		}

		if err := gitCommitGoGit(commit); err != nil {
			fmt.Println(err.Error())
		}
	},
	Version: VERSION,
}

func init() {
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

	cobra.OnInitialize(initConfig)
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

func gitCommitOs(commit *model.Commit) error {
	if dryRun || !commit.IsValid() {
		fmt.Printf("Did not commit changes. This would have been the command.\ngit commit -m \"%s\"\n", commit.String())
	} else {
		//#nosec G204 -- The point of this app is to run git commands
		runOsCmd := exec.Command("git", "commit", "-m", commit.String())

		osCmdOutput, runErr := runOsCmd.CombinedOutput()
		if runErr != nil {
			return fmt.Errorf("Did not commit changes due to git error.\n%s", osCmdOutput)
		} else {
			fmt.Println(string(osCmdOutput))
		}
	}
	return nil
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
	_, err = worktree.Commit(commit.String(), &gogit.CommitOptions{})
	if err != nil {
		return err
	}

	return nil
}
