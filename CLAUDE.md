# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go CLI application template repository designed for GitHub Codespaces with Claude Code integration. The project uses standard Go tooling and follows conventional Go project structure. It includes Claude Code for AI-assisted development directly in the Codespace terminal.

## Development Commands

### Building and Running
- `go build -o bin/yosegi .` - Build the CLI application
- `go run main.go` - Run the application directly
- `./bin/yosegi` - Run the built binary

### Testing and Quality
- `go test ./...` - Run all tests
- `go test -v ./...` - Run tests with verbose output
- `go test -race ./...` - Run tests with race detection
- `go mod tidy` - Clean up go.mod and go.sum
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run Go vet for static analysis

### Module Management
- `go mod init` - Initialize go module (if needed)
- `go mod download` - Download dependencies
- `go get <package>` - Add new dependency

## Performance Optimization for Claude Code

### Quick Commands
- `go build -ldflags="-s -w" -o bin/yosegi .` - Build with optimization flags (reduces binary size)
- `go test -short ./...` - Run only short tests for quick feedback
- `go mod why <package>` - Check why a dependency exists
- `go list -m all` - List all dependencies quickly

### Code Generation Helpers
- Use `//go:generate` directives for repetitive code
- Leverage `stringer` for enum string methods: `//go:generate stringer -type=MyEnum`
- Use `mockgen` for test mocks: `//go:generate mockgen -source=interface.go -destination=mock.go`

### Preferred Patterns
- Use table-driven tests for better readability and maintenance
- Implement interfaces for testability and loose coupling
- Use context.Context for cancellation and timeout control
- Prefer error wrapping with `fmt.Errorf("%w", err)` for better error traces
- Use `sync.Pool` for frequently allocated objects
- Leverage `go:embed` for embedding static files

### Quick Development Workflow
- Use `task` commands for common tasks (see Taskfile.yml)
- Run `task dev` for development mode
- Run `task test-short` for fast test feedback
- Run `task lint` for quick code quality checks
- Run `task ci` for all CI checks
- Use `go work` for multi-module development

## Architecture

This is a template project for Go CLI applications. The main entry point is in `main.go`. When developing CLI applications from this template:

- Use standard Go project layout conventions
- CLI functionality should typically use libraries like `cobra` or `flag` for argument parsing
- Consider using `viper` for configuration management
- Follow Go naming conventions and package organization

## Development Best Practices
- Run `task lint` before you create commit.
- Use `task --list-all` to see all available tasks.

## Available Tasks
Run `task --list-all` to see all available tasks. Key tasks include:
- `task dev` - Run in development mode
- `task test` - Run all tests
- `task test-short` - Run short tests with coverage
- `task lint` - Run all linting checks
- `task build` - Build the application
- `task build-release` - Build optimized release binary
- `task clean` - Clean build artifacts
- `task ci` - Run all CI checks