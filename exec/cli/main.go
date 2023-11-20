package main

import (
	"fmt"
	"os"

	"github.com/PossibleLlama/commit-check/prompt"
	tea "github.com/charmbracelet/bubbletea"
)

var typeNameOptions = []string{
	"feat",
	"fix",
}

func main() {

	p := tea.NewProgram(prompt.NewPromptType(typeNameOptions))
	if _, err := p.Run(); err != nil {
		fmt.Println("An unexpected error:", err)
		os.Exit(1)
	}
}
