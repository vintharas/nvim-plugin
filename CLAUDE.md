# nvim-plugin Development Guide

## Build/Test Commands
- Run app: `go run ./cmd/nvim-plugin/main.go`
- Build binary: `go build -o nvim-plugin ./cmd/nvim-plugin`
- Run all tests: `go test ./...`
- Run verbose tests: `go test -v ./...`
- Test specific package: `go test ./pkg/ui`
- Test specific function: `go test -run TestGeneratePlugin ./pkg/ui`
- Special tests: `RUN_ALL_TESTS=1 go test ./pkg/ui`

## Code Style Guidelines
- **Imports**: Standard library first, third-party next, with comment blocks
- **Error Handling**: Use `fmt.Errorf()` with `%w` for error wrapping
- **Naming**: PascalCase for exported, camelCase for unexported, package names short
- **Types**: Use iota for enums, add comments for structs and fields
- **Documentation**: Document "why" not just "what", all functions should have comments
- **Testing**: Use table-driven tests, descriptive names, proper cleanup with defer
- **Structure**: Follow `cmd/` for entrypoints, `pkg/` for packages convention