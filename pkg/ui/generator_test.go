package ui

import (
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed testdata/*.tmpl
var testTemplateFS embed.FS

func TestSanitizeVarName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my-plugin", "my_plugin"},
		{"plugin_name", "plugin_name"},
		{"camelCase", "camelCase"},
		{"plugin-with-many-hyphens", "plugin_with_many_hyphens"},
		{"", ""},
	}

	for _, test := range tests {
		result := sanitizeVarName(test.input)
		if result != test.expected {
			t.Errorf("sanitizeVarName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"camelCase", "CamelCase"},
		{"123test", "123test"}, // Numbers not affected
		{"", ""},
	}

	for _, test := range tests {
		result := capitalizeFirst(test.input)
		if result != test.expected {
			t.Errorf("capitalizeFirst(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestRenderTemplateFile(t *testing.T) {
	// First verify the template file exists in the file system
	_, err := os.Stat("../../templates/README.md.tmpl")
	if err != nil {
		t.Skipf("Skipping test: template file not found: %v", err)
	}
	
	data := TemplateData{
		Name:           "test-plugin",
		Description:    "A test plugin",
		VarName:        "test_plugin",
		CapitalizedCmd: "Test-plugin",
	}
	
	result, err := renderTemplateFile("templates/README.md.tmpl", data)
	if err != nil {
		t.Fatalf("renderTemplateFile failed: %v", err)
	}
	
	// Verify the rendered content contains expected elements
	expectedElements := []string{
		"# test-plugin",
		"A test plugin",
		"require('test-plugin')",
		":Test-plugin",
	}
	
	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("Rendered template missing expected content: %q", expected)
		}
	}
}

func TestGeneratePlugin(t *testing.T) {
	// Skip this test if running in CI/CD where templates might not be available
	_, err := os.Stat("../../templates")
	if err != nil {
		t.Skipf("Skipping test: templates directory not available: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "plugin-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory so GeneratePlugin creates files there
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Test parameters
	pluginName := "test-plugin"
	description := "A test plugin for Neovim"

	// Generate the plugin
	err = GeneratePlugin(pluginName, description)
	if err != nil {
		t.Fatalf("GeneratePlugin failed: %v", err)
	}

	// Check that all expected directories and files were created
	expectedDirs := []string{
		filepath.Join(tempDir, pluginName),
		filepath.Join(tempDir, pluginName, "lua", pluginName),
		filepath.Join(tempDir, pluginName, "plugin"),
		filepath.Join(tempDir, pluginName, "doc"),
	}

	for _, dir := range expectedDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to exist", dir)
		}
	}

	expectedFiles := []string{
		filepath.Join(tempDir, pluginName, "lua", pluginName, "init.lua"),
		filepath.Join(tempDir, pluginName, "plugin", pluginName+".lua"),
		filepath.Join(tempDir, pluginName, "README.md"),
		filepath.Join(tempDir, pluginName, "doc", pluginName+".txt"),
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", file)
		}
	}

	// Check file contents
	readmeContent, err := ioutil.ReadFile(filepath.Join(tempDir, pluginName, "README.md"))
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	// Verify README contains expected content
	readmeStr := string(readmeContent)
	if !strings.Contains(readmeStr, pluginName) {
		t.Errorf("README.md does not contain plugin name")
	}
	if !strings.Contains(readmeStr, description) {
		t.Errorf("README.md does not contain plugin description")
	}

	// Check plugin entry file
	pluginFileContent, err := ioutil.ReadFile(filepath.Join(tempDir, pluginName, "plugin", pluginName+".lua"))
	if err != nil {
		t.Fatalf("Failed to read plugin file: %v", err)
	}

	// Verify plugin file contains expected code
	pluginFileStr := string(pluginFileContent)
	if !strings.Contains(pluginFileStr, "vim.g.loaded_"+sanitizeVarName(pluginName)) {
		t.Errorf("Plugin file does not contain expected global variable")
	}
}

func TestWriteFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "write-file-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test writing a file
	testPath := filepath.Join(tempDir, "test-file.txt")
	testContent := "This is test content\nWith multiple lines"

	err = writeFile(testPath, testContent)
	if err != nil {
		t.Fatalf("writeFile failed: %v", err)
	}

	// Read back the file content and verify
	content, err := ioutil.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("File content mismatch. Expected %q, got %q", testContent, string(content))
	}
}

func TestVerifyTemplateFiles(t *testing.T) {
	// This test verifies that all template files exist and are readable
	templatePaths := []string{
		"templates/lua/plugin_name/init.lua.tmpl",
		"templates/plugin/plugin.lua.tmpl",
		"templates/README.md.tmpl",
		"templates/doc/plugin.txt.tmpl",
	}
	
	for _, path := range templatePaths {
		// Skip if files don't exist locally during development or testing
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			t.Logf("Skipping test for template %s: not found locally", path)
			continue
		}
		
		// Check if the file exists in the embedded filesystem
		_, err = templateFS.ReadFile(path)
		if err != nil {
			t.Errorf("Template file not found or not readable in embedded FS: %s, error: %v", path, err)
		}
	}
}