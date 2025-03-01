package main

import (
	"fmt"
	"os"

	// Bubble Tea is a framework for building terminal user interfaces based on The Elm Architecture
	tea "github.com/charmbracelet/bubbletea"
	// Import our UI package that contains the model and generator
	"github.com/vintharas/nvim-plugin/pkg/ui"
)

func main() {
	// Initialize a new Bubble Tea program with our model
	// Bubble Tea follows the Model-View-Update (MVU) architecture pattern
	p := tea.NewProgram(ui.NewModel())
	
	// Run the program and handle any errors
	// The program will run until a tea.Quit command is received
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}