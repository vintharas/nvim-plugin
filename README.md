# nvim-plugin

A CLI tool for generating Neovim plugin scaffolding with ease.

## Overview

nvim-plugin is a terminal-based application that helps you quickly bootstrap a new Neovim plugin with the correct structure and boilerplate code. It provides an interactive interface to specify your plugin's name and description, then generates the necessary files and directories.

## Architecture

The project follows a simple architecture:

```
nvim-plugin/
├── cmd/                     # Command-line application entry points
│   └── nvim-plugin/         # Main CLI application
│       └── main.go          # Application entry point
├── pkg/                     # Reusable packages
│   └── ui/                  # UI components and logic
│       ├── generator.go     # Plugin generation functionality
│       └── model.go         # Application state and UI model
└── templates/               # Templates for generated files
    ├── README.md.tmpl       # Template for plugin README
    ├── doc/                 # Templates for documentation
    ├── lua/                 # Templates for Lua modules
    └── plugin/              # Templates for plugin entry points
```

### Core Components

1. **Terminal UI (using Bubble Tea)**: The application uses the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework for creating a terminal user interface based on The Elm Architecture.

2. **UI Model**: The UI state is managed in `pkg/ui/model.go`, which implements the Bubble Tea model interface with the following methods:

   - `Init()`: Sets up the initial model state
   - `Update()`: Handles user input and state transitions
   - `View()`: Renders the UI based on the current state

3. **Plugin Generator**: The `pkg/ui/generator.go` file contains the logic for creating the plugin directory structure and generating all required files based on the user's input.

### Template System

The plugin generator uses Go's `text/template` package for creating all plugin files, with templates stored in separate files that mirror the structure of a generated plugin.

1. **Template Organization**: Templates are organized in a directory structure that matches the plugin structure:
   ```
   templates/
   ├── README.md.tmpl             # Template for the plugin README
   ├── doc/
   │   └── plugin.txt.tmpl        # Template for Neovim help docs
   ├── lua/
   │   └── plugin_name/
   │       └── init.lua.tmpl      # Template for the main Lua module 
   └── plugin/
       └── plugin.lua.tmpl        # Template for the plugin entry point
   ```

2. **Template Data Structure**: A `TemplateData` struct holds all variables needed for the templates:
   ```go
   type TemplateData struct {
       Name           string    // Plugin name
       Description    string    // Plugin description
       Date           string    // Current date
       VarName        string    // Sanitized variable name (for Lua)
       CapitalizedCmd string    // Capitalized first letter of name (for commands)
       HeaderTitle    string    // Uppercase title for docs
       Underline      string    // Underline for the header title
       DocHeader      string    // Header for the docs file
   }
   ```

3. **Template Embedding**: Templates are embedded in the binary using Go's `embed` package:
   ```go
   //go:embed ../../templates/*
   var templateFS embed.FS
   ```

4. **Template Rendering**: The `renderTemplateFile` function loads and processes templates from the embedded filesystem:
   ```go
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
   ```

This approach makes it easy to maintain and modify the generated files while keeping a consistent structure. Storing templates as separate files that mirror the plugin structure provides better clarity and makes it easier to understand the output that will be generated.

## Getting Started

### Installation

```bash
go install github.com/vintharas/nvim-plugin/cmd/nvim-plugin@latest
```

### Usage

Simply run:

```bash
nvim-plugin
```

Follow the interactive prompts to:

1. Enter your plugin name
2. Provide a short description
3. Confirm the details
4. Generate your plugin

## Generated Plugin Structure

The tool generates a complete Neovim plugin structure including:

- Entry point for the plugin
- Lua module structure
- Proper documentation
- README with installation instructions
- Necessary boilerplate code

## Development

### Local Development

When developing nvim-plugin, you can run the application locally to test changes:

```bash
# Run the application directly from the source code
go run ./cmd/nvim-plugin/main.go

# Build and run the binary
go build -o nvim-plugin ./cmd/nvim-plugin
./nvim-plugin
```

### Testing Changes Manually

To test the plugin generation functionality:

1. Run the application in development mode:

   ```bash
   go run ./cmd/nvim-plugin/main.go
   ```

2. Follow the prompts to create a test plugin:

   - Enter a test plugin name (e.g., "test-plugin")
   - Enter a description
   - Confirm the details with "y"

3. Examine the generated plugin structure:

   ```bash
   ls -la ./test-plugin
   cat ./test-plugin/lua/test-plugin/init.lua
   cat ./test-plugin/README.md
   ```

4. Clean up after testing:

   ```bash
   rm -rf ./test-plugin
   ```

### Running Tests

The project has a comprehensive test suite that covers both the UI components and the generator functionality:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./pkg/ui

# Run a specific test
go test -run TestGeneratePlugin ./pkg/ui
```

### Test Structure

The tests are organized by component:

1. **Generator Tests**: Tests in `pkg/ui/generator_test.go` verify the plugin generation logic:

   - File and directory creation
   - Template rendering
   - Helper functions

2. **UI Model Tests**: Tests in `pkg/ui/model_test.go` verify the UI component behavior:

   - State transitions
   - User input handling
   - View rendering

3. **Integration Tests**: A basic test in `cmd/nvim-plugin/main_test.go` ensures the application compiles and starts correctly.

### Adding New Tests

When adding new features, follow these guidelines for testing:

1. For generator features, add tests that verify file content and structure
2. For UI features, add tests that simulate user interactions and verify state changes
3. Use temporary directories for tests that create files
4. Mock external dependencies when appropriate

### Development Workflow

A typical development workflow might look like:

1. Make changes to the code
2. Run tests to ensure everything still works:

   ```bash
   go test ./...
   ```

3. Manually test the application:

   ```bash
   go run ./cmd/nvim-plugin/main.go
   ```

4. Examine the generated plugin to ensure it meets expectations
5. Repeat

## License

MIT

