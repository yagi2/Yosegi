package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Worktree represents a git worktree
type Worktree struct {
	Path      string
	Branch    string
	Commit    string
	IsCurrent bool
}

// Manager handles git worktree operations
type Manager interface {
	List() ([]Worktree, error)
	Add(path, branch string, createBranch bool) error
	Remove(path string, force bool) error
	GetCurrentPath() (string, error)
	DeleteBranch(branch string, force bool) error
	HasUnpushedCommits(branch string) (bool, int, error)
}

type manager struct {
	repoRoot string
}

// Security validation functions
var (
	// validBranchName matches valid git branch names
	validBranchName = regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
	// dangerousShellChars contains characters that could be dangerous in shell commands
	dangerousShellChars = []string{";", "&", "|", "$", "`", "(", ")", "<", ">", "\"", "'", "\\", "\n", "\r"}
)

// validateBranchName ensures the branch name is safe for use in git commands
func validateBranchName(branch string) error {
	if branch == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	
	// Check for git-specific restrictions
	if strings.HasPrefix(branch, "-") {
		return fmt.Errorf("branch name cannot start with a dash")
	}
	
	if strings.HasPrefix(branch, ".") || strings.HasSuffix(branch, ".") {
		return fmt.Errorf("branch name cannot start or end with a dot")
	}
	
	if strings.Contains(branch, "..") {
		return fmt.Errorf("branch name cannot contain consecutive dots")
	}
	
	// Check for dangerous characters that could be interpreted as shell commands
	for _, char := range dangerousShellChars {
		if strings.Contains(branch, char) {
			return fmt.Errorf("branch name contains dangerous character: %s", char)
		}
	}
	
	// Ensure it matches the valid pattern
	if !validBranchName.MatchString(branch) {
		return fmt.Errorf("branch name contains invalid characters")
	}
	
	return nil
}

// validatePath ensures the path is safe for use in git commands
func validatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	
	// Check for dangerous characters
	for _, char := range dangerousShellChars {
		if strings.Contains(path, char) {
			return fmt.Errorf("path contains dangerous character: %s", char)
		}
	}
	
	// Prevent malicious directory traversal attempts
	// Check for patterns that could escape beyond the intended directory
	if strings.Contains(path, "../../../") {
		return fmt.Errorf("path contains directory traversal sequences")
	}
	
	// Convert to absolute path and check for suspicious patterns
	absPath, err := filepath.Abs(path)
	if err == nil {
		// Check if the absolute path tries to access system directories
		suspiciousPaths := []string{"/etc/", "/usr/bin/", "/bin/", "/sbin/", "/root/", "/home/root/"}
		for _, suspicious := range suspiciousPaths {
			if strings.HasPrefix(absPath, suspicious) {
				return fmt.Errorf("path attempts to access restricted directory: %s", suspicious)
			}
		}
	}
	
	return nil
}

// NewManager creates a new git worktree manager
func NewManager() (Manager, error) {
	repoRoot, err := FindGitRoot(".")
	if err != nil {
		return nil, err
	}
	return &manager{repoRoot: repoRoot}, nil
}

// FindGitRoot finds the git repository root directory
func FindGitRoot(startPath string) (string, error) {
	path, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	for {
		gitDir := filepath.Join(path, ".git")
		if info, err := os.Stat(gitDir); err == nil {
			if info.IsDir() {
				return path, nil // Regular git repository
			}
			// Handle worktree case - .git is a file pointing to the actual .git directory
			content, err := os.ReadFile(gitDir)
			if err != nil {
				return "", err
			}
			if strings.HasPrefix(string(content), "gitdir:") {
				// Extract the main repo path from the gitdir reference
				gitdirPath := strings.TrimSpace(strings.TrimPrefix(string(content), "gitdir:"))
				if filepath.IsAbs(gitdirPath) {
					// Find the main repo from the worktree gitdir path
					// e.g., /path/to/repo/.git/worktrees/branch -> /path/to/repo
					parts := strings.Split(gitdirPath, string(filepath.Separator))
					for i, part := range parts {
						if part == ".git" && i > 0 {
							return strings.Join(parts[:i], string(filepath.Separator)), nil
						}
					}
				}
			}
		}

		parent := filepath.Dir(path)
		if parent == path {
			break // Reached root directory
		}
		path = parent
	}
	return "", errors.New("not a git repository")
}

// List returns all worktrees in the repository
func (m *manager) List() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = m.repoRoot
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	return parseWorktreeList(string(output))
}

// Add creates a new worktree
func (m *manager) Add(path, branch string, createBranch bool) error {
	// Validate inputs for security
	if err := validatePath(path); err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	
	if err := validateBranchName(branch); err != nil {
		return fmt.Errorf("invalid branch name: %w", err)
	}
	
	// Check if branch exists
	checkCmd := exec.Command("git", "rev-parse", "--verify", fmt.Sprintf("refs/heads/%s", branch))
	checkCmd.Dir = m.repoRoot
	branchExists := checkCmd.Run() == nil

	args := []string{"worktree", "add"}
	if createBranch && !branchExists {
		// Create new branch
		args = append(args, "-b", branch, path)
	} else if branchExists && !createBranch {
		// Use existing branch
		args = append(args, path, branch)
	} else if !branchExists && !createBranch {
		// Branch doesn't exist and user doesn't want to create it
		return fmt.Errorf("branch '%s' does not exist. Use --create-branch flag to create it", branch)
	} else if branchExists && createBranch {
		// Branch exists but user wants to create it
		return fmt.Errorf("branch '%s' already exists. Remove --create-branch flag to use existing branch", branch)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = m.repoRoot

	// Get detailed error output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add worktree (command: git %v): %w\nOutput: %s", args, err, string(output))
	}
	return nil
}

// Remove removes a worktree
func (m *manager) Remove(path string, force bool) error {
	// Validate input for security
	if err := validatePath(path); err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	
	// First, try to get the absolute path
	absPath, err := filepath.Abs(path)
	if err == nil {
		path = absPath
	}

	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, path)

	cmd := exec.Command("git", args...)
	cmd.Dir = m.repoRoot

	// Get detailed error output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMsg := string(output)

		// Common error patterns and solutions
		if strings.Contains(errorMsg, "is dirty") {
			return fmt.Errorf("worktree contains uncommitted changes. Use --force flag to remove anyway: %s", errorMsg)
		}

		if strings.Contains(errorMsg, "does not exist") {
			// Worktree might be already removed but git doesn't know
			// Try to prune first
			pruneCmd := exec.Command("git", "worktree", "prune")
			pruneCmd.Dir = m.repoRoot
			pruneOutput, pruneErr := pruneCmd.CombinedOutput()
			if pruneErr == nil {
				// Check if the worktree still exists after pruning
				listCmd := exec.Command("git", "worktree", "list")
				listCmd.Dir = m.repoRoot
				listOutput, _ := listCmd.Output()
				if !strings.Contains(string(listOutput), path) {
					// Worktree was successfully pruned
					return nil
				}
				// Try remove again after pruning
				cmd := exec.Command("git", args...)
				cmd.Dir = m.repoRoot
				output, err = cmd.CombinedOutput()
				if err == nil {
					return nil
				}
			} else {
				return fmt.Errorf("failed to prune worktrees: %s", string(pruneOutput))
			}
		}

		if strings.Contains(errorMsg, "is locked") {
			// Worktree is locked, need force flag
			if !force {
				return fmt.Errorf("worktree is locked. Use --force flag to remove anyway: %s", errorMsg)
			}
		}

		return fmt.Errorf("failed to remove worktree: %s", errorMsg)
	}
	return nil
}

// GetCurrentPath returns the current working directory
func (m *manager) GetCurrentPath() (string, error) {
	return os.Getwd()
}

// parseWorktreeList parses the output of 'git worktree list --porcelain'
func parseWorktreeList(output string) ([]Worktree, error) {
	var worktrees []Worktree
	var current Worktree

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return worktrees, nil
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			// End of current worktree entry
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = Worktree{}
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "HEAD ") {
			current.Commit = strings.TrimPrefix(line, "HEAD ")
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch refs/heads/")
		} else if line == "bare" {
			current.Branch = "(bare)"
		} else if line == "detached" {
			current.Branch = "(detached)"
		}
	}

	// Add the last worktree if exists
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	// Mark current worktree
	currentPath, _ := os.Getwd()
	for i := range worktrees {
		if worktrees[i].Path == currentPath {
			worktrees[i].IsCurrent = true
			break
		}
	}

	return worktrees, nil
}

// DeleteBranch deletes a local branch
func (m *manager) DeleteBranch(branch string, force bool) error {
	// Validate input for security
	if err := validateBranchName(branch); err != nil {
		return fmt.Errorf("invalid branch name: %w", err)
	}
	
	args := []string{"branch", "-d"}
	if force {
		args[1] = "-D"
	}
	args = append(args, branch)

	cmd := exec.Command("git", args...)
	cmd.Dir = m.repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMsg := string(output)
		if strings.Contains(errorMsg, "not found") {
			return fmt.Errorf("branch '%s' not found", branch)
		}
		if strings.Contains(errorMsg, "not fully merged") {
			return fmt.Errorf("branch '%s' is not fully merged. Use force flag to delete anyway", branch)
		}
		return fmt.Errorf("failed to delete branch: %s", errorMsg)
	}
	return nil
}

// HasUnpushedCommits checks if a branch has unpushed commits
func (m *manager) HasUnpushedCommits(branch string) (bool, int, error) {
	// Validate input for security
	if err := validateBranchName(branch); err != nil {
		return false, 0, fmt.Errorf("invalid branch name: %w", err)
	}
	
	// First check if the branch has an upstream
	upstreamCmd := exec.Command("git", "rev-parse", "--abbrev-ref", fmt.Sprintf("%s@{upstream}", branch))
	upstreamCmd.Dir = m.repoRoot
	upstreamOutput, err := upstreamCmd.Output()
	if err != nil {
		// No upstream configured, consider all commits as unpushed
		countCmd := exec.Command("git", "rev-list", "--count", branch)
		countCmd.Dir = m.repoRoot
		countOutput, countErr := countCmd.Output()
		if countErr != nil {
			return false, 0, fmt.Errorf("failed to count commits: %w", countErr)
		}
		count := 0
		fmt.Sscanf(strings.TrimSpace(string(countOutput)), "%d", &count)
		return count > 0, count, nil
	}

	upstream := strings.TrimSpace(string(upstreamOutput))

	// Count commits ahead of upstream
	cmd := exec.Command("git", "rev-list", "--count", fmt.Sprintf("%s..%s", upstream, branch))
	cmd.Dir = m.repoRoot

	output, err := cmd.Output()
	if err != nil {
		return false, 0, fmt.Errorf("failed to check unpushed commits: %w", err)
	}

	count := 0
	fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &count)

	return count > 0, count, nil
}
