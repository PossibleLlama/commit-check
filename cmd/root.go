package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var conventionType string
var VERSION string

var rootCmd = &cobra.Command{
	Use:   "commit-check",
	Short: "Verify your commits have a common format",
	Run: func(cmd *cobra.Command, args []string) {
		if conventionType != "angular" && conventionType != "conventional" {
			fmt.Println("convention type must be either angular or conventional but was", conventionType)
			os.Exit(1)
		}

		prefix := promptType()
		scope := promptScope()
		message := promptMessage()

		if scope != "" {
			prefix = prefix + "(" + scope + "): "
		} else {
			prefix = prefix + ": "
		}
		commitArgs := []string{"commit", "-m", prefix + message}
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
		"type",
		"t",
		"conventional",
		"accepts either 'conventional' or 'angular'")
}
