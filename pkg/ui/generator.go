package ui

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//go:embed templates
var templateFS embed.FS

// TemplateData holds all the variables used in templates
type TemplateData struct {
	Name           string // Plugin name
	Description    string // Plugin description
	Date           string // Current date
	VarName        string // Sanitized variable name (for Lua)
	CapitalizedCmd string // Capitalized first letter of name (for commands)
	HeaderTitle    string // Uppercase title for docs
	Underline      string // Underline for the header title
	DocHeader      string // Header for the docs file
}

// GeneratePlugin creates a new Neovim plugin with the given name and description
// It builds the directory structure and generates all necessary files
func GeneratePlugin(name, description string) error {
	// Create the main plugin directory
	pluginDir := "./" + name
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Create standard Neovim plugin directory structure:
	// - lua/{name}: contains the main plugin code
	// - plugin: contains the plugin entry point
	// - doc: contains plugin documentation
	dirs := []string{
		filepath.Join(pluginDir, "lua", name),
		filepath.Join(pluginDir, "plugin"),
		filepath.Join(pluginDir, "doc"),
	}

	// Create each directory with proper permissions
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Prepare template data
	data := TemplateData{
		Name:           name,
		Description:    description,
		Date:           time.Now().Format("2006-01-02"),
		VarName:        sanitizeVarName(name),
		CapitalizedCmd: capitalizeFirst(name),
		HeaderTitle:    strings.ToUpper(name),
		DocHeader:      strings.ToUpper(name) + ".TXT",
		Underline:      strings.Repeat("=", len(strings.ToUpper(name))),
	}

	// Define the files to generate
	filesToGenerate := []struct {
		outputPath string
		tmplPath   string
	}{
		{
			outputPath: filepath.Join(pluginDir, "lua", name, "init.lua"),
			tmplPath:   "templates/lua/plugin_name/init.lua.tmpl",
		},
		{
			outputPath: filepath.Join(pluginDir, "plugin", name+".lua"),
			tmplPath:   "templates/plugin/plugin.lua.tmpl",
		},
		{
			outputPath: filepath.Join(pluginDir, "README.md"),
			tmplPath:   "templates/README.md.tmpl",
		},
		{
			outputPath: filepath.Join(pluginDir, "doc", name+".txt"),
			tmplPath:   "templates/doc/plugin.txt.tmpl",
		},
	}

	// Generate each file from its template file
	for _, file := range filesToGenerate {
		content, err := renderTemplateFile(file.tmplPath, data)
		if err != nil {
			return fmt.Errorf("failed to render template for %s: %w", file.outputPath, err)
		}

		if err := writeFile(file.outputPath, content); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.outputPath, err)
		}
	}

	return nil
}

// renderTemplateFile loads a template from the embedded filesystem and renders it
func renderTemplateFile(tmplPath string, data TemplateData) (string, error) {
	// Read the template file from the embedded filesystem
	tmplContent, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", tmplPath, err)
	}

	// Parse the template
	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", tmplPath, err)
	}

	// Execute the template with the data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", tmplPath, err)
	}

	return buf.String(), nil
}

// writeFile is a helper function to write content to a file
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// Helper functions for string manipulation

// sanitizeVarName converts plugin name to a valid Lua variable name
// Replaces hyphens with underscores for use in Lua variables
func sanitizeVarName(name string) string {
	// Replace non-alphanumeric characters with underscore
	return strings.ReplaceAll(name, "-", "_")
}

// capitalizeFirst capitalizes the first letter of a string
// Used for command names and other user-facing identifiers
func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}