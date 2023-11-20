package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	model "github.com/PossibleLlama/commit-check/model"
	"github.com/PossibleLlama/commit-check/view"
)

var conventionType string
var VERSION string

var rootCmd = &cobra.Command{
	Use:   "commit-check",
	Short: "Verify your commits have a common format",
	Run: func(cmd *cobra.Command, args []string) {
		if conventionType != "angular" && conventionType != "conventionalcommit" {
			fmt.Println("convention type must be either 'angular' or 'conventionalcommit' but was", conventionType)
			os.Exit(1)
		}
		var prefixChoices []model.CommitType
		switch conventionType {
		case "angular":
			prefixChoices = model.TypeAngular
		case "conventionalcommit":
			prefixChoices = model.TypeConventionalCommit
		}

		var commit = model.Commit{}
		viewModel := view.InitCommitModel(prefixChoices, &commit)
		tea.NewProgram(viewModel).Run()

		if commit.String() == "" {
			os.Exit(1)
		}

		cString := strings.Split(commit.String(), "\n")
		commitArgs := []string{"commit"}
		for _, line := range cString {
			commitArgs = append(commitArgs, "-m", "\""+line+"\"")
		}

		runOsCmd := exec.Command("git", commitArgs...)

		osCmdOutput, runErr := runOsCmd.CombinedOutput()
		if runErr != nil {
			fmt.Println(runOsCmd.String())
			fmt.Println("failed to commit with error:", string(osCmdOutput))
			os.Exit(1)
		}
		fmt.Println(string(osCmdOutput))
	},
	Version: VERSION,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&conventionType,
		"type-list",
		"l",
		"conventionalcommit",
		"accepts either 'angular' or 'conventionalcommit'")
}
