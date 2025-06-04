#!/bin/zsh

# Yosegi shell integration for zsh
# Add this to your ~/.zshrc:
# source /path/to/yosegi/scripts/shell_integration.zsh

# Mark that shell integration is active
export YOSEGI_SHELL_INTEGRATION=1

yosegi() {
    local result
    result=$(command yosegi "$@")
    local exit_code=$?
    
    # Check if the output contains a directory change command
    if [[ $result == CD:* ]]; then
        local target_dir="${result#CD:}"
        if [[ -d "$target_dir" ]]; then
            cd "$target_dir" || echo "Failed to change directory to $target_dir"
            echo "Switched to worktree: $target_dir"
        else
            echo "Directory not found: $target_dir"
        fi
    else
        # Normal output, just print it
        echo "$result"
    fi
    
    return $exit_code
}

# Enable tab completion for yosegi
if command -v compdef &> /dev/null; then
    compdef '_values "yosegi commands" list new switch remove help version' yosegi
fi