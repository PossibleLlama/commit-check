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

		commit := &model.Commit{}

		p := tea.NewProgram(prompt.NewPromptCommit(model.TypeAngular, commit), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("An unexpected error:", err)
			os.Exit(1)
		}

		commitArgs := []string{"commit", "-m", commit.String()}
		if dryRun || !commit.IsValid() {
			fmt.Println("git", strings.Join(commitArgs, " "))
		} else {
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
		"conventionalcommit",
		"accepts either 'angular' or 'conventionalcommit'")
	rootCmd.Flags().BoolVarP(&dryRun,
		"dry-run",
		"d",
		false,
		"run the program without committing")
}
