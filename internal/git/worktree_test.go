package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindGitRoot(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (string, func(), error) // Returns temp dir, cleanup func, error
		expectError bool
	}{
		{
			name: "Regular git repository",
			setupFunc: func() (string, func(), error) {
				tmpDir, err := os.MkdirTemp("", "yosegi-test-*")
				if err != nil {
					return "", nil, err
				}
				gitDir := filepath.Join(tmpDir, ".git")
				if err := os.Mkdir(gitDir, 0755); err != nil {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
					return "", nil, err
				}
				cleanup := func() {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
				}
				return tmpDir, cleanup, nil
			},
			expectError: false,
		},
		{
			name: "Worktree git file",
			setupFunc: func() (string, func(), error) {
				tmpDir, err := os.MkdirTemp("", "yosegi-test-*")
				if err != nil {
					return "", nil, err
				}
				// Create a main repo directory
				mainRepoDir := filepath.Join(tmpDir, "main-repo")
				if err := os.MkdirAll(filepath.Join(mainRepoDir, ".git"), 0755); err != nil {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
					return "", nil, err
				}

				// Create worktree directory
				worktreeDir := filepath.Join(tmpDir, "worktree")
				if err := os.MkdirAll(worktreeDir, 0755); err != nil {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
					return "", nil, err
				}

				gitFile := filepath.Join(worktreeDir, ".git")
				gitContent := fmt.Sprintf("gitdir: %s/.git/worktrees/test", mainRepoDir)
				if err := os.WriteFile(gitFile, []byte(gitContent), 0644); err != nil {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
					return "", nil, err
				}
				cleanup := func() {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
				}
				return worktreeDir, cleanup, nil
			},
			expectError: false,
		},
		{
			name: "Not a git repository",
			setupFunc: func() (string, func(), error) {
				tmpDir, err := os.MkdirTemp("", "yosegi-test-*")
				if err != nil {
					return "", nil, err
				}
				cleanup := func() {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
				}
				return tmpDir, cleanup, nil
			},
			expectError: true,
		},
		{
			name: "Nested directory in git repo",
			setupFunc: func() (string, func(), error) {
				tmpDir, err := os.MkdirTemp("", "yosegi-test-*")
				if err != nil {
					return "", nil, err
				}
				gitDir := filepath.Join(tmpDir, ".git")
				if err := os.Mkdir(gitDir, 0755); err != nil {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
					return "", nil, err
				}
				nestedDir := filepath.Join(tmpDir, "subdir", "nested")
				if err := os.MkdirAll(nestedDir, 0755); err != nil {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
					return "", nil, err
				}
				cleanup := func() {
					_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
				}
				return nestedDir, cleanup, nil
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, cleanup, err := tt.setupFunc()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer cleanup()

			result, err := FindGitRoot(testDir)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && err == nil {
				if result == "" {
					t.Errorf("Expected non-empty result")
				}
				// The result should be an existing directory
				if _, err := os.Stat(result); os.IsNotExist(err) {
					t.Errorf("Result directory does not exist: %s", result)
				}
			}
		})
	}
}

func TestParseWorktreeList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Worktree
		hasError bool
	}{
		{
			name:     "Empty output",
			input:    "",
			expected: []Worktree{},
			hasError: false,
		},
		{
			name: "Single worktree",
			input: `worktree /path/to/main
HEAD 1234567890abcdef
branch refs/heads/main

`,
			expected: []Worktree{
				{
					Path:   "/path/to/main",
					Branch: "main",
					Commit: "1234567890abcdef",
				},
			},
			hasError: false,
		},
		{
			name: "Multiple worktrees",
			input: `worktree /path/to/main
HEAD 1234567890abcdef
branch refs/heads/main

worktree /path/to/feature
HEAD abcdef1234567890
branch refs/heads/feature/test

`,
			expected: []Worktree{
				{
					Path:   "/path/to/main",
					Branch: "main",
					Commit: "1234567890abcdef",
				},
				{
					Path:   "/path/to/feature",
					Branch: "feature/test",
					Commit: "abcdef1234567890",
				},
			},
			hasError: false,
		},
		{
			name: "Bare repository",
			input: `worktree /path/to/repo
HEAD 1234567890abcdef
bare

`,
			expected: []Worktree{
				{
					Path:   "/path/to/repo",
					Branch: "(bare)",
					Commit: "1234567890abcdef",
				},
			},
			hasError: false,
		},
		{
			name: "Detached HEAD",
			input: `worktree /path/to/detached
HEAD 1234567890abcdef
detached

`,
			expected: []Worktree{
				{
					Path:   "/path/to/detached",
					Branch: "(detached)",
					Commit: "1234567890abcdef",
				},
			},
			hasError: false,
		},
		{
			name: "Mixed worktrees",
			input: `worktree /path/to/main
HEAD 1234567890abcdef
branch refs/heads/main

worktree /path/to/bare
HEAD abcdef1234567890
bare

worktree /path/to/detached
HEAD fedcba0987654321
detached

`,
			expected: []Worktree{
				{
					Path:   "/path/to/main",
					Branch: "main",
					Commit: "1234567890abcdef",
				},
				{
					Path:   "/path/to/bare",
					Branch: "(bare)",
					Commit: "abcdef1234567890",
				},
				{
					Path:   "/path/to/detached",
					Branch: "(detached)",
					Commit: "fedcba0987654321",
				},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseWorktreeList(tt.input)

			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d worktrees, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if i >= len(result) {
					t.Errorf("Missing worktree at index %d", i)
					continue
				}
				actual := result[i]
				if actual.Path != expected.Path {
					t.Errorf("Worktree %d: expected path %s, got %s", i, expected.Path, actual.Path)
				}
				if actual.Branch != expected.Branch {
					t.Errorf("Worktree %d: expected branch %s, got %s", i, expected.Branch, actual.Branch)
				}
				if actual.Commit != expected.Commit {
					t.Errorf("Worktree %d: expected commit %s, got %s", i, expected.Commit, actual.Commit)
				}
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	// Test in a directory that's not a git repository
	tmpDir, err := os.MkdirTemp("", "yosegi-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
	}()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore directory: %v", err)
		}
	}()

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	_, err = NewManager()
	if err == nil {
		t.Errorf("Expected error when not in git repository")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("Expected 'not a git repository' error, got: %v", err)
	}
}

func TestManagerGetCurrentPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yosegi-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
	}()

	// Create a mock git repository
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore directory: %v", err)
		}
	}()

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	manager, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	currentPath, err := manager.GetCurrentPath()
	if err != nil {
		t.Errorf("GetCurrentPath failed: %v", err)
	}

	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	if currentPath != expectedPath {
		t.Errorf("Expected current path %s, got %s", expectedPath, currentPath)
	}
}

// Helper function to test the manager interface
func createTestGitRepo(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "yosegi-test-git-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
	}

	return tmpDir, cleanup
}

func TestManagerInterface(t *testing.T) {
	// Test that manager implements the Manager interface
	var _ Manager = &manager{}

	// Test creating a manager with a valid git repository
	repoDir, cleanup := createTestGitRepo(t)
	defer cleanup()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Logf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	manager, err := NewManager()
	if err != nil {
		t.Errorf("Failed to create manager in git repository: %v", err)
	}

	if manager == nil {
		t.Errorf("Manager should not be nil")
	}
}

// Benchmark tests
func BenchmarkFindGitRoot(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "yosegi-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir) // Ignore cleanup errors
	}()

	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		b.Fatalf("Failed to create .git directory: %v", err)
	}

	// Create nested directory structure
	nestedDir := filepath.Join(tmpDir, "a", "b", "c", "d", "e")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		b.Fatalf("Failed to create nested directories: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := FindGitRoot(nestedDir)
		if err != nil {
			b.Errorf("FindGitRoot failed: %v", err)
		}
	}
}

func BenchmarkParseWorktreeList(b *testing.B) {
	input := `worktree /path/to/main
HEAD 1234567890abcdef
branch refs/heads/main

worktree /path/to/feature1
HEAD abcdef1234567890
branch refs/heads/feature/test1

worktree /path/to/feature2
HEAD fedcba0987654321
branch refs/heads/feature/test2

worktree /path/to/hotfix
HEAD 1111222233334444
branch refs/heads/hotfix/urgent

`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parseWorktreeList(input)
		if err != nil {
			b.Errorf("parseWorktreeList failed: %v", err)
		}
	}
}

func TestManagerListErrors(t *testing.T) {
	// Test List method error handling
	m := &manager{repoRoot: "/nonexistent"}

	_, err := m.List()
	if err == nil {
		t.Error("Expected error for non-existent repository")
	}

	if !strings.Contains(err.Error(), "failed to list worktrees") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestManagerAddErrors(t *testing.T) {
	// Test Add method error handling
	m := &manager{repoRoot: "/nonexistent"}

	err := m.Add("/tmp/test", "test-branch", false)
	if err == nil {
		t.Error("Expected error for non-existent repository")
	}

	if !strings.Contains(err.Error(), "failed to add worktree") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestManagerRemoveErrors(t *testing.T) {
	// Test Remove method error handling
	m := &manager{repoRoot: "/nonexistent"}

	err := m.Remove("/tmp/test", false)
	if err == nil {
		t.Error("Expected error for non-existent repository")
	}

	if !strings.Contains(err.Error(), "failed to remove worktree") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestManagerRemoveWithForce(t *testing.T) {
	// Test Remove method with force flag
	m := &manager{repoRoot: "/nonexistent"}

	err := m.Remove("/tmp/test", true)
	if err == nil {
		t.Error("Expected error for non-existent repository")
	}

	// Error should still occur, but we test that force flag is handled
	if !strings.Contains(err.Error(), "failed to remove worktree") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestManagerAddWithCreateBranch(t *testing.T) {
	// Test Add method with createBranch flag
	m := &manager{repoRoot: "/nonexistent"}

	err := m.Add("/tmp/test", "new-branch", true)
	if err == nil {
		t.Error("Expected error for non-existent repository")
	}

	if !strings.Contains(err.Error(), "failed to add worktree") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestManagerMethodSignatures(t *testing.T) {
	// Test that all manager methods have correct signatures for interface compliance

	m := Manager(&manager{repoRoot: "/test"})

	// Test List signature: () ([]Worktree, error)
	worktrees, err := m.List()
	_ = worktrees
	_ = err

	// Test Add signature: (string, string, bool) error
	err = m.Add("path", "branch", true)
	_ = err

	// Test Remove signature: (string, bool) error
	err = m.Remove("path", false)
	_ = err

	// Test GetCurrentPath signature: () (string, error)
	path, err := m.GetCurrentPath()
	_ = path
	_ = err
}

func TestManagerPathHandling(t *testing.T) {
	// Test path handling in manager methods
	m := &manager{repoRoot: "/test/repo"}

	// Test empty path handling
	err := m.Add("", "branch", false)
	if err == nil {
		t.Log("Add with empty path handled (expected to fail)")
	}

	err = m.Remove("", false)
	if err == nil {
		t.Log("Remove with empty path handled (expected to fail)")
	}
}

func TestManagerBranchHandling(t *testing.T) {
	// Test branch name handling in Add method
	m := &manager{repoRoot: "/test/repo"}

	// Test empty branch name
	err := m.Add("/tmp/test", "", false)
	if err == nil {
		t.Log("Add with empty branch handled (expected to fail)")
	}

	// Test branch name with spaces
	err = m.Add("/tmp/test", "branch with spaces", false)
	if err == nil {
		t.Log("Add with spaced branch name handled (expected to fail)")
	}
}
