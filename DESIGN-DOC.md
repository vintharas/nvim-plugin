# nvim-plugin Design Document

## Current Status

The `nvim-plugin` CLI tool currently provides basic functionality to generate Neovim plugin scaffolding with an interactive terminal UI. It creates the necessary files and directories for a Neovim plugin based on user input for the plugin name and description.

## Future Features and Enhancements

### Command Structure

Evolve the CLI to support a more sophisticated command structure:

```
nvim-plugin [command] [options]
```

#### Proposed Commands

1. **`new`** - Create a new plugin (current default behavior)
   ```
   nvim-plugin new [plugin-name] [--description "Description"] [--dir /path/to/dir]
   ```
   - Support both interactive mode (no arguments) and direct mode (with arguments)
   - Add template selection options (minimal, full, etc.)

2. **`list`** - List existing plugins 
   ```
   nvim-plugin list [--location /path/to/plugins] [--global]
   ```
   - Discover plugins in common locations (e.g., ~/.local/share/nvim/site/pack)
   - Show plugin name, location, and brief description

3. **`go`** - Navigate to a plugin directory
   ```
   nvim-plugin go <plugin-name>
   ```
   - Open the plugin directory or change the current directory to the plugin
   - Support tab completion for plugin names

4. **`update`** - Update plugin templates or add missing files
   ```
   nvim-plugin update <plugin-name> [--add-tests] [--add-ci]
   ```
   - Add CI configurations, tests, or other boilerplate to existing plugins

5. **`check`** - Validate a plugin's structure
   ```
   nvim-plugin check <plugin-name>
   ```
   - Check for recommended files and directory structure
   - Validate Lua syntax and common issues

### Template Enhancements

1. **Multiple Template Types**
   - Minimal templates for simple plugins
   - Full templates with tests, CI integration, etc.
   - Language-specific templates (Lua, VimL, mixed)

2. **Plugin Types**
   - UI enhancements (colorschemes, statuslines)
   - LSP extensions
   - Syntax/treesitter plugins
   - Integrations with external tools

3. **Additional Files**
   - GitHub Actions workflow templates
   - Testing scaffolding with Plenary.nvim
   - Makefile/build scripts
   - Type definitions for LSP integration

### Configuration and Customization

1. **User Configuration**
   - Allow users to set default author information
   - Support custom template directories
   - Remember frequently used options

2. **Template Customization**
   - Allow users to modify default templates
   - Support multiple template sets
   - Template versioning

### Integration Features

1. **Plugin Management**
   - Integration with plugin managers (packer, lazy.nvim, etc.)
   - Publishing to GitHub/GitLab with repository setup

2. **Development Tools**
   - Built-in testing commands
   - Linting and documentation tools
   - Release management assistance

## Technical Improvements

1. **Architecture**
   - Modular command structure for better extensibility
   - Improved error handling and logging
   - Configuration persistence

2. **Performance**
   - Optimize template rendering for larger templates
   - Improve file operations for larger plugin structures

3. **User Experience**
   - Improved CLI output with better formatting and colors
   - Interactive elements for all commands
   - Progress indicators for longer operations

## Implementation Priority

1. Core command structure (new, list, go)
2. Template enhancements (multiple template types)
3. User configuration system
4. Integration features
5. Advanced plugin management features

## Next Steps

1. Refactor main.go to support command substructure
2. Implement basic "list" and "go" commands
3. Add configuration file support
4. Expand template options