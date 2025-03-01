package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Helper function to simulate key presses
func pressKeys(model tea.Model, keys ...string) tea.Model {
	m := model
	for _, key := range keys {
		var msg tea.Msg
		if key == "enter" {
			msg = tea.KeyMsg{Type: tea.KeyEnter}
		} else if key == "backspace" {
			msg = tea.KeyMsg{Type: tea.KeyBackspace}
		} else {
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
		}
		var cmd tea.Cmd
		m, cmd = m.Update(msg)
		if cmd != nil {
			// We're ignoring commands in tests for simplicity
		}
	}
	return m
}

func TestNewModel(t *testing.T) {
	model := NewModel()

	// Check initial state
	if model.status != nameInput {
		t.Errorf("New model should start in nameInput state, got %v", model.status)
	}

	if model.pluginName != "" {
		t.Errorf("New model should have empty plugin name, got %q", model.pluginName)
	}

	if model.description != "" {
		t.Errorf("New model should have empty description, got %q", model.description)
	}
}

func TestModelInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	if cmd != nil {
		t.Errorf("Init should return nil command, got non-nil")
	}
}

func TestModelUpdateNameInput(t *testing.T) {
	model := NewModel()

	// Test typing a plugin name
	m := pressKeys(model, "t", "e", "s", "t")
	updatedModel := m.(Model)

	if updatedModel.pluginName != "test" {
		t.Errorf("Expected plugin name to be 'test', got %q", updatedModel.pluginName)
	}

	// Test backspace
	m = pressKeys(updatedModel, "backspace")
	updatedModel = m.(Model)

	if updatedModel.pluginName != "tes" {
		t.Errorf("After backspace, expected plugin name to be 'tes', got %q", updatedModel.pluginName)
	}

	// Test Enter to move to next state
	m = pressKeys(updatedModel, "enter")
	updatedModel = m.(Model)

	if updatedModel.status != descriptionInput {
		t.Errorf("After Enter, expected to move to descriptionInput state, got %v", updatedModel.status)
	}
}

func TestModelUpdateDescriptionInput(t *testing.T) {
	// Start with a model in the descriptionInput state
	model := Model{
		status:     descriptionInput,
		pluginName: "test-plugin",
	}

	// Test typing a description
	m := pressKeys(model, "a", " ", "t", "e", "s", "t")
	updatedModel := m.(Model)

	if updatedModel.description != "a test" {
		t.Errorf("Expected description to be 'a test', got %q", updatedModel.description)
	}

	// Test Enter to move to confirm screen
	m = pressKeys(updatedModel, "enter")
	updatedModel = m.(Model)

	if updatedModel.status != confirmScreen {
		t.Errorf("After Enter, expected to move to confirmScreen state, got %v", updatedModel.status)
	}
}

func TestModelUpdateConfirmScreen(t *testing.T) {
	// Start with a model in the confirmScreen state
	model := Model{
		status:      confirmScreen,
		pluginName:  "test-plugin",
		description: "A test plugin",
	}

	// In a real test, we would need to mock or replace the GeneratePlugin function
	// to avoid actually generating files during tests.
	// For simplicity, we'll just check state transitions.

	// Test 'n' to go back to name input
	m := pressKeys(model, "n")
	updatedModel := m.(Model)

	if updatedModel.status != nameInput {
		t.Errorf("After 'n', expected to move back to nameInput state, got %v", updatedModel.status)
	}

	// Reset to confirm screen
	updatedModel.status = confirmScreen

	// For the 'y' case, we need to mock GeneratePlugin or the test will try to create real files
	// This is simplified for the example
}

func TestModelView(t *testing.T) {
	// Test nameInput view
	nameModel := Model{
		status:     nameInput,
		pluginName: "test",
	}

	nameView := nameModel.View()
	if !strings.Contains(nameView, "Plugin Name:") {
		t.Errorf("nameInput view should contain 'Plugin Name:' header")
	}
	if !strings.Contains(nameView, "test") {
		t.Errorf("nameInput view should contain the plugin name")
	}

	// Test descriptionInput view
	descModel := Model{
		status:      descriptionInput,
		pluginName:  "test",
		description: "description",
	}

	descView := descModel.View()
	if !strings.Contains(descView, "Plugin Description:") {
		t.Errorf("descriptionInput view should contain 'Plugin Description:' header")
	}
	if !strings.Contains(descView, "description") {
		t.Errorf("descriptionInput view should contain the description")
	}

	// Test confirmScreen view
	confirmModel := Model{
		status:      confirmScreen,
		pluginName:  "test",
		description: "description",
	}

	confirmView := confirmModel.View()
	if !strings.Contains(confirmView, "Confirm Details:") {
		t.Errorf("confirmScreen view should contain 'Confirm Details:' header")
	}
	if !strings.Contains(confirmView, "Plugin Name: test") {
		t.Errorf("confirmScreen view should contain the plugin name")
	}
	if !strings.Contains(confirmView, "Description: description") {
		t.Errorf("confirmScreen view should contain the description")
	}

	// Test done view with success
	doneSuccessModel := Model{
		status:      done,
		pluginName:  "test",
		description: "description",
		err:         nil,
	}

	doneSuccessView := doneSuccessModel.View()
	if !strings.Contains(doneSuccessView, "Plugin created successfully") {
		t.Errorf("done view with success should indicate success")
	}

	// Test done view with error
	doneErrorModel := Model{
		status:      done,
		pluginName:  "test",
		description: "description",
		err:         &mockError{message: "test error"},
	}

	doneErrorView := doneErrorModel.View()
	if !strings.Contains(doneErrorView, "Error creating plugin") {
		t.Errorf("done view with error should indicate error")
	}
	if !strings.Contains(doneErrorView, "test error") {
		t.Errorf("done view with error should contain the error message")
	}
}

// Mock error type for testing
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

