package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/PossibleLlama/commit-check/model"
	"github.com/PossibleLlama/commit-check/prompt"
	tea "github.com/charmbracelet/bubbletea"
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

		commitArgs := []string{"commit", "-m", commit.String()}
		if dryRun || !commit.IsValid() {
			fmt.Println("git", strings.Join(commitArgs, " "))
		} else {
			//#nosec G204 -- The point of this app is to run git commands
			runOsCmd := exec.Command("git", commitArgs...)

			osCmdOutput, runErr := runOsCmd.CombinedOutput()
			if runErr != nil {
				fmt.Println("failed to commit with error:", string(osCmdOutput))
				os.Exit(1)
			}
			fmt.Println(string(osCmdOutput))
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
