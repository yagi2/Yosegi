# Shell Integration Guide

## Overview

Yosegi provides shell integration that enables automatic directory switching when using the `switch` command. However, this comes with some trade-offs regarding interactive UI behavior.

## Behavior Differences

### Without Shell Integration
- **UI**: Full interactive interface with visual selection
- **Directory Change**: Manual (copy/paste the provided `cd` command)
- **Commands**: All commands show their full interactive UI

### With Shell Integration
- **UI**: Automatic plain text mode (to prevent output buffering issues)
- **Directory Change**: Automatic via shell function
- **Commands**: `list` and `switch` use plain text output for better shell integration

## The Buffering Issue

When shell integration is active, commands are executed via:
```bash
result=$(command yosegi "$@")
```

This causes:
1. **Output Buffering**: Interactive UI output gets buffered
2. **Delayed Display**: UI appears only after program termination
3. **Poor User Experience**: Nothing visible during selection process

## Solutions

### 1. Use Plain Mode (Default with Shell Integration)
```bash
# Automatically enabled when YOSEGI_SHELL_INTEGRATION=1
yosegi list        # Shows plain text list
yosegi switch      # Shows plain text list with usage instructions
```

### 2. Force Interactive Mode
```bash
# Override shell integration behavior
yosegi list --interactive        # Force interactive UI
yosegi switch --interactive      # Force interactive UI (may have buffering issues)
```

### 3. Use Direct Arguments
```bash
# Bypass interactive selection entirely
yosegi switch <branch-name>      # Direct switch to specific branch
yosegi switch <worktree-path>    # Direct switch to specific path
```

## Recommended Workflow

### For Daily Use (with Shell Integration)
```bash
# Quick list
yosegi list

# Direct switch (fastest)
yosegi switch main
yosegi switch feature-branch

# Browse and select (plain mode)
yosegi switch  # Shows list, then use: yosegi switch <target>
```

### For Visual Experience (without Shell Integration)
```bash
# Full interactive experience
unset YOSEGI_SHELL_INTEGRATION
yosegi list      # Beautiful interactive UI
yosegi switch    # Beautiful interactive selection

# Manual directory change
cd $(yosegi switch target-branch | grep "CD:" | cut -d: -f2)
```

## Setup Instructions

### Bash
```bash
# Add to ~/.bashrc
source /path/to/yosegi/scripts/shell_integration.bash
```

### Zsh
```bash
# Add to ~/.zshrc
source /path/to/yosegi/scripts/shell_integration.zsh
```

### Fish
```bash
# Add to ~/.config/fish/config.fish
source /path/to/yosegi/scripts/shell_integration.fish
```

## Troubleshooting

### Problem: No UI shown, only output after quit
**Cause**: Shell integration buffering issue
**Solution**: This is expected behavior. Use plain mode or direct arguments instead.

### Problem: Want visual UI with automatic directory change
**Solution**: Currently not possible due to shell/process limitations. Choose either:
- Visual UI without auto-cd (no shell integration)
- Auto-cd without visual UI (with shell integration)

### Problem: Interactive mode doesn't work in terminal
**Cause**: No TTY available (common in CI/scripts)
**Solution**: Use `--plain` flag or direct arguments

## Technical Limitations

The interactive UI + shell integration limitation is due to:

1. **Process Isolation**: Child processes cannot change parent shell's working directory
2. **Output Capture**: Shell functions capture stdout, causing buffering
3. **TTY Requirements**: Interactive UI requires proper terminal allocation

This is a fundamental Unix/shell limitation, not a bug in Yosegi.