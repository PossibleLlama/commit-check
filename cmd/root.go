package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	model "github.com/PossibleLlama/commit-check/model"
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

		var commit = model.Commit{}

		commit.Type = promptType()
		commit.Scope = promptScope()
		commit.Description = promptMessage()

		commitArgs := []string{"commit", "-m", commit.String()}
		runOsCmd := exec.Command("git", commitArgs...)

		osCmdOutput, runErr := runOsCmd.CombinedOutput()
		if runErr != nil {
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
