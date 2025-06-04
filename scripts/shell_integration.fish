#!/usr/bin/env fish

# Yosegi shell integration for fish
# Add this to your ~/.config/fish/config.fish:
# source /path/to/yosegi/scripts/shell_integration.fish

# Mark that shell integration is active
set -gx YOSEGI_SHELL_INTEGRATION 1

function yosegi
    set result (command yosegi $argv)
    set exit_code $status
    
    # Check if the output contains a directory change command
    if string match -q "CD:*" $result
        set target_dir (string sub -s 4 $result)
        if test -d $target_dir
            cd $target_dir
            echo "Switched to worktree: $target_dir"
        else
            echo "Directory not found: $target_dir"
        end
    else
        # Normal output, just print it
        echo $result
    end
    
    return $exit_code
end

# Enable tab completion for yosegi
complete -c yosegi -f -a "list new switch remove help version"