# Yosegi 🌲

Interactive git worktree management tool with a beautiful TUI interface.

## Overview

Yosegi is a CLI tool designed for the modern "Vibe Coding" era, providing intuitive and visual management of git worktrees. Like `tig` and `peco`, it offers an excellent visual interface for managing multiple git worktrees with ease.

## Features

- 🎯 **Interactive UI**: Beautiful terminal interface built with Bubble Tea and Lip Gloss
- 🌲 **Worktree Management**: Create, switch, and remove git worktrees seamlessly
- 🔄 **Shell Integration**: Automatic directory switching with bash/zsh/fish support
- 🎨 **Customizable Themes**: YAML-based configuration for colors and UI preferences
- ⚡ **Keyboard Navigation**: Vim-style navigation (j/k) and arrow keys
- 🛡️ **Safety Features**: Confirmation prompts and protection against accidental deletions

## Installation

### Build from Source

```bash
git clone https://github.com/yagi2/cli-vibe-go.git
cd cli-vibe-go
go build -o yosegi .
```

### Shell Integration Setup

To enable directory switching functionality, add the appropriate shell integration:

#### Bash
```bash
# Add to ~/.bashrc
source /path/to/yosegi/scripts/shell_integration.bash
```

#### Zsh
```bash
# Add to ~/.zshrc
source /path/to/yosegi/scripts/shell_integration.zsh
```

#### Fish
```bash
# Add to ~/.config/fish/config.fish
source /path/to/yosegi/scripts/shell_integration.fish
```

## Usage

### Basic Commands

#### List Worktrees
```bash
yosegi list     # or yosegi ls, yosegi l
```
Interactive list of all worktrees with current status indicators.

#### Create New Worktree
```bash
yosegi new [branch]              # Interactive creation
yosegi new feature-branch        # Create with specified branch
yosegi new -b new-feature        # Create new branch and worktree
yosegi new -p ../feature feature # Specify custom path
```

#### Switch Worktree
```bash
yosegi switch   # or yosegi sw, yosegi s
```
Interactive selection and automatic directory switching.

#### Remove Worktree
```bash
yosegi remove   # or yosegi rm, yosegi delete
```
Safe removal with confirmation prompts.

### Configuration

#### Initialize Configuration
```bash
yosegi config init
```
Creates a default configuration file at `~/.config/yosegi/config.yaml`.

#### View Current Configuration
```bash
yosegi config show
```

### Configuration File

Example `~/.config/yosegi/config.yaml`:

```yaml
default_worktree_path: "../"
theme:
  primary: "#7C3AED"
  secondary: "#06B6D4" 
  success: "#10B981"
  warning: "#F59E0B"
  error: "#EF4444"
  muted: "#6B7280"
  text: "#F9FAFB"
git:
  auto_create_branch: false
  default_remote: "origin"
  exclude_patterns: []
ui:
  show_icons: true
  confirm_delete: true
  max_path_length: 50
aliases:
  ls: "list"
  sw: "switch"
  rm: "remove"
```

## Keyboard Navigation

- `↑/k`: Move up
- `↓/j`: Move down  
- `Enter`: Select/Execute
- `d`: Delete (in remove mode)
- `q`: Quit
- `Tab/Shift+Tab`: Navigate input fields

## Examples

### Typical Workflow

```bash
# List current worktrees
yosegi list

# Create new worktree for feature development
yosegi new feature/user-auth

# Switch to the new worktree (automatically changes directory)
yosegi switch

# When done, remove the worktree
yosegi remove
```

### Advanced Usage

```bash
# Create worktree with custom path and new branch
yosegi new -b hotfix/urgent-fix -p ../hotfix

# Force remove worktree (skip confirmation)
yosegi remove --force
```

## Development

### Building
```bash
go build -o bin/yosegi .
```

### Testing
```bash
go test ./...
```

### Linting
```bash
go fmt ./...
go vet ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Requirements

- Go 1.21+
- Git with worktree support
- Terminal with color support

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by tools like `tig` and `peco` for their excellent visual interfaces
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Uses [Cobra](https://github.com/spf13/cobra) for CLI framework