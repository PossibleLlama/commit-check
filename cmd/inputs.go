package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

var (
	// Based on https://www.conventionalcommits.org/
	// TODO allow config to choose between these lists
	conventionTypes = []string{
		"fix",
		"feat",
		"BREAKING CHANGE",
	}
	angularConventionTypes = []string{
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
		"revert",
		"BREAKING CHANGE",
	}
)

func promptType() string {
	prefixPrompt := promptui.Select{
		Label: "Select type of change",
		Items: conventionTypes,
	}
	_, prefixAsStr, prefixErr := prefixPrompt.Run()
	if prefixErr != nil {
		fmt.Println("failed to select item from list", prefixErr)
		os.Exit(1)
	}
	return strings.TrimSpace(prefixAsStr)
}

func promptScope() string {
	scopePrompt := promptui.Prompt{
		Label: "(Optional) Scope of change - section of codebase, or a ticket reference",
	}

	scope, scopeErr := scopePrompt.Run()
	if scopeErr != nil {
		fmt.Println("failed to get scope", scopeErr)
		os.Exit(1)
	}

	return strings.TrimSpace(scope)
}

func promptMessage() string {
	messagePrompt := promptui.Prompt{
		Label: "Commit message",
		Validate: func(s string) error {
			if len(s) < 4 {
				return errors.New("message needs to be at least 4 characters")
			}
			return nil
		},
	}

	message, messageErr := messagePrompt.Run()
	if messageErr != nil {
		fmt.Println("failed to get message", messageErr)
		os.Exit(1)
	}

	return strings.TrimSpace(message)
}
