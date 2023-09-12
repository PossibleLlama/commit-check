package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	// Based on https://www.conventionalcommits.org/
	// TODO allow config to choose between these lists
	conventionPrefix = []string{
		"fix",
		"feat",
		"BREAKING CHANGE",
	}
	angularConventionPrefix = []string{
		"fix",
		"feat",
		"build",
		"chore",
		"ci",
		"docs",
		"style",
		"refactor",
		"perf",
		"test",
		"BREAKING CHANGE",
	}
)

var rootCmd = &cobra.Command{
	Use:   "commit-check",
	Short: "Verify your commits have a common format",
	Run: func(cmd *cobra.Command, args []string) {
		prefixPrompt := promptui.Select{
			Label: "Select type of change",
			Items: conventionPrefix,
		}
		_, prefixAsStr, prefixErr := prefixPrompt.Run()
		if prefixErr != nil {
			fmt.Println("failed to select item from list", prefixErr)
			os.Exit(1)
		}

		// TODO get input of messages

		commitArgs := []string{"commit", "-m", "\"" + prefixAsStr + ": " + strings.Join(args, " ") + "\""}
		runOsCmd := exec.Command("git", commitArgs...)

		osCmdOutput, runErr := runOsCmd.CombinedOutput()
		if runErr != nil {
			fmt.Println("failed to commit with error:", string(osCmdOutput))
			os.Exit(1)
		}
		fmt.Println(string(osCmdOutput))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
