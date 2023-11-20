package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PossibleLlama/commit-check/model"
	"github.com/PossibleLlama/commit-check/prompt"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f := setupLoggingToFile()
	defer f.Close()

	commit := &model.Commit{}
	var p *tea.Program

	p = tea.NewProgram(prompt.NewPromptType(model.TypeAngular, commit), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("An unexpected error:", err)
		os.Exit(1)
	}
	p = tea.NewProgram(prompt.NewPromptScope("https://google.com", "", commit), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("An unexpected error:", err)
		os.Exit(1)
	}
	p = tea.NewProgram(prompt.NewPromptDescription(commit), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("An unexpected error:", err)
		os.Exit(1)
	}
	log.Output(1, commit.String())
}

func setupLoggingToFile() *os.File {
	f, err := tea.LogToFile("debug.log", "commit-check")
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	return f
}
