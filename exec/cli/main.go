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

	p := tea.NewProgram(prompt.NewPromptCommit(model.TypeAngular, commit), tea.WithAltScreen())
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
