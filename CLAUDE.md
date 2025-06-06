# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Yosegi is an interactive git worktree management tool with a beautiful TUI (Terminal User Interface). It provides an intuitive visual interface similar to tools like `tig` or `peco` for managing multiple git worktrees efficiently. Built with Go using the Cobra CLI framework and Bubble Tea for the TUI, Yosegi makes git worktree operations simple and enjoyable.

## CLI Commands

### Yosegi Commands
- `yosegi` - Run list command by default (interactive worktree selector)
- `yosegi list` (aliases: `ls`, `l`) - List all worktrees interactively
  - `-p, --print` - Print selected worktree path without interactive UI (for use in scripts)
- `yosegi new [branch]` (aliases: `add`, `create`, `n`) - Create a new worktree
  - `-b, --create-branch` - Create a new branch
  - `-p, --path` - Specify worktree path
- `yosegi remove` (aliases: `rm`, `delete`, `del`, `r`) - Remove a worktree
  - `-f, --force` - Force removal
- `yosegi config init` - Create default config file
- `yosegi config show` - Display current configuration

### Non-Interactive Mode
Yosegi supports non-interactive mode for shell scripting and command substitution:
- Use `--print` flag or `-p` to show interactive TUI on stderr and output selected path to stdout
- Without `--print`, automatically detects non-TTY environments and returns first non-current worktree
- Example: `cd $(yosegi list --print)` - shows TUI selector, then changes to selected directory
- Example: `cd $(yosegi list)` - automatically selects first non-current worktree

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

### Key Dependencies
- **Cobra** (github.com/spf13/cobra) - CLI framework for command structure
- **Bubble Tea** (github.com/charmbracelet/bubbletea) - TUI framework
- **Bubbles** (github.com/charmbracelet/bubbles) - Pre-built TUI components
- **Lip Gloss** (github.com/charmbracelet/lipgloss) - Terminal styling
- **yaml.v3** (gopkg.in/yaml.v3) - Configuration file handling

### Preferred Patterns
- Use table-driven tests for better readability and maintenance
- Implement interfaces for testability and loose coupling
- Use context.Context for cancellation and timeout control
- Prefer error wrapping with `fmt.Errorf("%w", err)` for better error traces
- Mock external dependencies (git commands) in tests
- Follow Bubble Tea patterns for TUI components

### Quick Development Workflow
- Use `task` commands for common tasks (see Taskfile.yml)
- Run `task dev` for development mode
- Run `task test-short` for fast test feedback
- Run `task lint` for quick code quality checks
- Run `task ci` for all CI checks
- Use `go work` for multi-module development

## Architecture

Yosegi follows a clean architecture pattern with clear separation of concerns:

- **main.go**: Application entry point
- **cmd/**: Cobra command definitions (list, new, remove, config)
- **internal/config/**: Configuration management with YAML support
- **internal/git/**: Git worktree operations wrapper
- **internal/ui/**: Bubble Tea TUI components (input, selector, styles)

Key design decisions:
- Uses Cobra for robust CLI command handling
- Implements interactive TUI with Bubble Tea framework
- Provides Vim-style keybindings (j/k navigation)
- YAML-based configuration for customization
- Comprehensive test coverage for all components

## Development Best Practices
- Run `task lint` before you create commit.
- Use `task --list-all` to see all available tasks.
- Always create commit message in English

## Available Tasks
Run `task --list-all` to see all available tasks. Key tasks include:
- `task dev` - Run in development mode
- `task test` - Run all tests with race detection
- `task test-short` - Run short tests with coverage
- `task lint` - Run all linting checks (fmt, vet, golangci-lint)
- `task build` - Build the application
- `task build-release` - Build optimized release binary
- `task install` - Install to GOPATH/bin
- `task clean` - Clean build artifacts
- `task ci` - Run all CI checks

## Yosegi-Specific Development Guidelines

### TUI Development
- Use Bubble Tea patterns for TUI components (Model, Update, View)
- Follow existing key binding conventions (j/k for navigation, Enter for selection)
- Maintain consistent styling using Lip Gloss
- Test TUI components with mock git operations

### Git Worktree Operations
- All git operations should go through `internal/git/worktree.go`
- Handle edge cases (current worktree, invalid paths, etc.)
- Provide clear error messages for git operation failures

### Configuration
- Configuration is stored in `~/.config/yosegi/config.yaml`
- Support both default and custom configurations
- Validate configuration values appropriately

### Testing Strategy
- Write table-driven tests for command logic
- Mock git operations for unit tests
- Test both success and error paths
- Ensure TUI components handle all user inputs gracefully

## Interaction Guidelines
- Always response in Japanese