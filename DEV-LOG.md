# Development Log

This document tracks significant changes, decisions, and development progress for the nvim-plugin project.

## 2025-03-01: Fixed Template Structure & Plans for Future Commands

### Changes
- Fixed issue with the `.stylua.toml.tmpl` file by renaming it to `stylua.toml.tmpl` (without the leading dot)
- Updated the generator code to reference the new file path
- Updated README.md to reflect the correct directory structure of templates under `pkg/ui/templates/`
- Added DESIGN-DOC.md to outline future development plans
- Added DEV-LOG.md (this file) to track development history

### Technical Details
- The embed directive in Go had issues with files starting with a dot (hidden files)
- Templates are now correctly located in `pkg/ui/templates/` instead of a root-level templates directory
- All tests are now passing 

### Next Steps
- Refactor the CLI to support a command-based structure (new, list, go, etc.)
- Add configuration storage for tracking installed plugins
- Implement directory navigation to installed plugins

## 2025-02-NN: First Implementation

### Changes
- Created initial implementation of the nvim-plugin generator
- Set up basic project structure:
  - cmd/nvim-plugin: CLI entry point
  - pkg/ui: UI model and generator logic
  - templates: Template files for plugin generation
- Implemented interactive UI using Bubble Tea
- Added basic templates for Neovim plugin structure

### Technical Details
- Using Go's embed system to include templates in the binary
- Generating standard Neovim plugin structure:
  - lua/{plugin}/init.lua (main plugin code)
  - plugin/{plugin}.lua (plugin entry point)
  - doc/{plugin}.txt (documentation)
  - README.md
  - .stylua.toml (formatting configuration)
- Template variables provide customization of the plugin name and description

### Next Steps
- Add tests for the generator and UI components
- Add more comprehensive documentation
- Consider additional template options