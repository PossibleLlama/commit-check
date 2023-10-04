package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"

	model "github.com/PossibleLlama/commit-check/model"
)

func promptType() model.CommitType {
	var prefixChoices []model.CommitType
	switch conventionType {
	case "angular":
		prefixChoices = model.TypeAngular
	case "conventionalcommit":
		prefixChoices = model.TypeConventionalCommit
	}

	prefixPrompt := promptui.Select{
		Label: "Select type of change",
		Items: prefixChoices,
	}
	selectedPosition, _, prefixErr := prefixPrompt.Run()
	if prefixErr != nil {
		fmt.Println("failed to select item from list", prefixErr)
		os.Exit(1)
	}
	return prefixChoices[selectedPosition]
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
