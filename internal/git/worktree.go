package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
}

type manager struct {
	repoRoot string
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
	// Check if branch exists
	if !createBranch {
		checkCmd := exec.Command("git", "rev-parse", "--verify", fmt.Sprintf("refs/heads/%s", branch))
		checkCmd.Dir = m.repoRoot
		if err := checkCmd.Run(); err != nil {
			// Branch doesn't exist, so create it
			createBranch = true
		}
	}

	args := []string{"worktree", "add"}
	if createBranch {
		args = append(args, "-b", branch)
	}
	args = append(args, path, branch)

	cmd := exec.Command("git", args...)
	cmd.Dir = m.repoRoot
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add worktree: %w", err)
	}
	return nil
}

// Remove removes a worktree
func (m *manager) Remove(path string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, path)

	cmd := exec.Command("git", args...)
	cmd.Dir = m.repoRoot
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
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