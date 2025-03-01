package ui

import (
	// bubbletea is the main framework for building terminal user interfaces
	"github.com/charmbracelet/bubbletea"
	// lipgloss is a styling library for terminal applications
	"github.com/charmbracelet/lipgloss"
)

// status represents the different screens/states of the application
type status int

// Application states using iota for automatic incrementation
const (
	nameInput status = iota       // First screen: enter plugin name
	descriptionInput              // Second screen: enter plugin description
	confirmScreen                 // Third screen: confirm details
	done                          // Final screen: display result
)

// Model represents the application state
type Model struct {
	status      status  // Current screen of the application
	pluginName  string  // Stores the plugin name entered by the user
	description string  // Stores the plugin description entered by the user
	cursor      int     // Cursor position (not currently used but available for extension)
	err         error   // Stores any error that occurs during plugin generation
}

// NewModel creates a new Model with default values
func NewModel() Model {
	return Model{
		status: nameInput, // Start the application in the nameInput state
	}
}

// Init implements bubbletea.Model
// This is called when the program starts
// We don't need any initial commands, so we return nil
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements bubbletea.Model
// This function handles all user input and state transitions
// It takes a message (usually a keypress) and returns an updated model and command
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// First, handle global keypresses that work in any state
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Exit the application
			return m, tea.Quit
		}
	}

	// Then delegate to state-specific update functions based on current status
	switch m.status {
	case nameInput:
		return updateNameInput(msg, m)
	case descriptionInput:
		return updateDescriptionInput(msg, m)
	case confirmScreen:
		return updateConfirmScreen(msg, m)
	}

	return m, nil
}

// View implements bubbletea.Model
// This function renders the UI based on the current model state
func (m Model) View() string {
	// Create a styled title bar that appears at the top of every screen
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingLeft(2).
		PaddingRight(2).
		MarginBottom(1)

	title := titleStyle.Render("nvim-plugin: Neovim Plugin Generator")

	// Render different content based on the current screen/state
	var content string
	switch m.status {
	case nameInput:
		content = viewNameInput(m)
	case descriptionInput:
		content = viewDescriptionInput(m)
	case confirmScreen:
		content = viewConfirmScreen(m)
	case done:
		content = viewDone(m)
	}

	// Combine the title and content with a footer showing how to quit
	return title + "\n" + content + "\n\nPress q to quit\n"
}

// Input handlers for each screen/state

// updateNameInput handles user input on the plugin name screen
func updateNameInput(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Move to the next screen if there's a name
			if len(m.pluginName) > 0 {
				m.status = descriptionInput
			}
			return m, nil
		case "backspace":
			// Delete the last character from the plugin name
			if len(m.pluginName) > 0 {
				m.pluginName = m.pluginName[:len(m.pluginName)-1]
			}
			return m, nil
		default:
			// Add typed characters to the plugin name
			if msg.Type == tea.KeyRunes {
				m.pluginName += string(msg.Runes)
			}
			return m, nil
		}
	}
	return m, nil
}

// updateDescriptionInput handles user input on the description screen
func updateDescriptionInput(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Move to the confirmation screen
			m.status = confirmScreen
			return m, nil
		case "backspace":
			// Delete the last character from the description
			if len(m.description) > 0 {
				m.description = m.description[:len(m.description)-1]
			}
			return m, nil
		default:
			// Add typed characters to the description
			if msg.Type == tea.KeyRunes {
				m.description += string(msg.Runes)
			}
			return m, nil
		}
	}
	return m, nil
}

// updateConfirmScreen handles user input on the confirmation screen
func updateConfirmScreen(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			// If the user confirms, generate the plugin
			err := GeneratePlugin(m.pluginName, m.description)
			if err != nil {
				m.err = err
			}
			// Move to the done screen
			m.status = done
			return m, nil
		case "n", "N":
			// If the user declines, go back to the first screen
			m.status = nameInput
			return m, nil
		}
	}
	return m, nil
}

// View helpers - functions to render each screen

// viewNameInput renders the plugin name input screen
func viewNameInput(m Model) string {
	return lipgloss.NewStyle().MarginBottom(1).Render("Plugin Name:") + "\n" +
		m.pluginName + "█" + "\n\n" + // "█" represents the cursor
		"Enter the name of your Neovim plugin and press Enter"
}

// viewDescriptionInput renders the description input screen
func viewDescriptionInput(m Model) string {
	return lipgloss.NewStyle().MarginBottom(1).Render("Plugin Description:") + "\n" +
		m.description + "█" + "\n\n" + // "█" represents the cursor
		"Enter a short description and press Enter"
}

// viewConfirmScreen renders the confirmation screen
func viewConfirmScreen(m Model) string {
	summary := "Plugin Name: " + m.pluginName + "\n" +
		"Description: " + m.description + "\n\n" +
		"Is this correct? (y/n)"

	return lipgloss.NewStyle().MarginBottom(1).Render("Confirm Details:") + "\n" + summary
}

// viewDone renders the final screen showing success or error
func viewDone(m Model) string {
	// If there was an error, show it in red
	if m.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Render("✗ Error creating plugin:") + "\n\n" +
			m.err.Error()
	}

	// Otherwise show a success message in green
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true).
		Render("✓ Plugin created successfully!") + "\n\n" +
		"Your new plugin has been created at:\n" +
		"./" + m.pluginName
}