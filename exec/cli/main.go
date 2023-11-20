package main

import (
	"fmt"
	"os"

	"github.com/PossibleLlama/commit-check/model"
	"github.com/PossibleLlama/commit-check/prompt"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	commit := &model.Commit{}

	p := tea.NewProgram(prompt.NewPromptType(model.TypeAngular, commit))
	if _, err := p.Run(); err != nil {
		fmt.Println("An unexpected error:", err)
		os.Exit(1)
	}
	fmt.Println("Finished commit as: ", commit.String())
}
