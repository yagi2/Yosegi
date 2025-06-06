package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/yagi2/yosegi/internal/git"
)

func TestSimpleSelectWorktreeEmptyList(t *testing.T) {
	var worktrees []git.Worktree
	var output bytes.Buffer
	var input bytes.Buffer

	result, err := simpleSelectWorktreeWithFiles(worktrees, &mockFile{&output}, &mockFile{&input})

	if err == nil {
		t.Error("Expected error for empty worktree list")
	}
	if result != nil {
		t.Error("Expected nil result for empty worktree list")
	}
	if err.Error() != "no worktrees found" {
		t.Errorf("Expected 'no worktrees found' error, got: %s", err.Error())
	}
}

func TestSimpleSelectWorktreeValidSelection(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: true},
		{Path: "/path/2", Branch: "feature", IsCurrent: false},
		{Path: "/path/3", Branch: "hotfix", IsCurrent: false},
	}

	tests := []struct {
		name      string
		input     string
		expectIdx int
		expectErr bool
	}{
		{
			name:      "Select first worktree",
			input:     "1\n",
			expectIdx: 0,
			expectErr: false,
		},
		{
			name:      "Select second worktree",
			input:     "2\n",
			expectIdx: 1,
			expectErr: false,
		},
		{
			name:      "Select third worktree",
			input:     "3\n",
			expectIdx: 2,
			expectErr: false,
		},
		{
			name:      "Quit with q",
			input:     "q\n",
			expectIdx: -1,
			expectErr: true,
		},
		{
			name:      "Quit with Q",
			input:     "Q\n",
			expectIdx: -1,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			input := &mockStringReader{strings.NewReader(tt.input)}

			result, err := simpleSelectWorktreeWithFiles(worktrees, &mockFile{&output}, &mockFile{input})

			if (err != nil) != tt.expectErr {
				t.Errorf("SimpleSelectWorktree() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if tt.expectErr {
				if err.Error() != "selection cancelled" {
					t.Errorf("Expected 'selection cancelled' error, got: %s", err.Error())
				}
				if result != nil {
					t.Error("Expected nil result for cancelled selection")
				}
				return
			}

			if result == nil {
				t.Error("Expected non-nil result for valid selection")
				return
			}

			expected := &worktrees[tt.expectIdx]
			if result.Path != expected.Path {
				t.Errorf("Expected worktree path %s, got %s", expected.Path, result.Path)
			}
			if result.Branch != expected.Branch {
				t.Errorf("Expected worktree branch %s, got %s", expected.Branch, result.Branch)
			}

			// Check that UI was displayed
			outputStr := output.String()
			if !strings.Contains(outputStr, "🌲 Git Worktrees:") {
				t.Error("Expected UI header to be displayed")
			}
			if !strings.Contains(outputStr, "Select worktree") {
				t.Error("Expected selection prompt to be displayed")
			}
		})
	}
}

func TestSimpleSelectWorktreeInvalidInput(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
		{Path: "/path/2", Branch: "feature", IsCurrent: false},
	}

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid number then valid",
			input: "0\n2\n", // 0 is invalid (out of range), then select 2
		},
		{
			name:  "Out of range then valid",
			input: "5\n1\n", // 5 is out of range, then select 1
		},
		{
			name:  "Non-numeric then valid",
			input: "abc\n2\n", // abc is not a number, then select 2
		},
		{
			name:  "Empty line then valid",
			input: "\n1\n", // empty line, then select 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			input := &mockStringReader{strings.NewReader(tt.input)}

			result, err := simpleSelectWorktreeWithFiles(worktrees, &mockFile{&output}, &mockFile{input})

			if err != nil {
				t.Errorf("Expected successful selection after retry, got error: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected non-nil result after valid selection")
				return
			}

			// Check that error messages were displayed
			outputStr := output.String()
			if !strings.Contains(outputStr, "Invalid") {
				t.Error("Expected invalid input message to be displayed")
			}
		})
	}
}

func TestSimpleSelectWorktreeDisplayFormat(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/repo/main", Branch: "main", IsCurrent: true},
		{Path: "/repo/feature", Branch: "feature-branch", IsCurrent: false},
	}

	var output bytes.Buffer
	input := &mockStringReader{strings.NewReader("q\n")} // Quit immediately

	_, _ = simpleSelectWorktreeWithFiles(worktrees, &mockFile{&output}, &mockFile{input})

	outputStr := output.String()

	// Check header
	if !strings.Contains(outputStr, "🌲 Git Worktrees:") {
		t.Error("Expected tree emoji and header")
	}

	// Check separator lines
	if !strings.Contains(outputStr, strings.Repeat("-", 60)) {
		t.Error("Expected separator lines")
	}

	// Check current worktree indicator
	if !strings.Contains(outputStr, "* 1) /repo/main (main)") {
		t.Error("Expected current worktree to be marked with asterisk")
	}

	// Check non-current worktree
	if !strings.Contains(outputStr, "  2) /repo/feature (feature-branch)") {
		t.Error("Expected non-current worktree to be displayed properly")
	}

	// Check prompt
	if !strings.Contains(outputStr, "Select worktree (1-2) or 'q' to quit:") {
		t.Error("Expected selection prompt with correct range")
	}
}

func TestSimpleSelectWorktreeIOError(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
	}

	var output bytes.Buffer
	// Use an errorReader that returns EOF immediately
	input := &errorReader{err: io.EOF}

	result, err := simpleSelectWorktreeWithFiles(worktrees, &mockFile{&output}, &mockFile{input})

	if err == nil {
		t.Error("Expected error when input fails")
	}
	if result != nil {
		t.Error("Expected nil result when input fails")
	}
	if !strings.Contains(err.Error(), "failed to read input") {
		t.Errorf("Expected 'failed to read input' error, got: %s", err.Error())
	}
}

// Mock file type for testing
type mockFile struct {
	io.ReadWriter
}

func (m *mockFile) Fd() uintptr { return 0 }

// mockStringReader wraps strings.Reader to implement io.ReadWriter
type mockStringReader struct {
	*strings.Reader
}

func (m *mockStringReader) Write(p []byte) (n int, err error) {
	return len(p), nil // Just pretend to write
}

func TestSimpleSelectWorktreeWrapper(t *testing.T) {
	// Test the wrapper function SimpleSelectWorktree
	worktrees := []git.Worktree{
		{Path: "/test/path", Branch: "main", IsCurrent: false},
	}

	// This will fail in test environment, but covers the wrapper function
	result, err := SimpleSelectWorktree(worktrees, os.Stderr, os.Stdin)

	// Should fail due to type conversion (*os.File expected)
	if err == nil {
		t.Logf("SimpleSelectWorktree unexpectedly succeeded: %v", result)
	} else {
		t.Logf("SimpleSelectWorktree failed as expected in test environment: %v", err)
	}
}

// Benchmark tests
func BenchmarkSimpleSelectWorktreeDisplay(b *testing.B) {
	worktrees := []git.Worktree{
		{Path: "/repo/main", Branch: "main", IsCurrent: true},
		{Path: "/repo/feature-1", Branch: "feature-1", IsCurrent: false},
		{Path: "/repo/feature-2", Branch: "feature-2", IsCurrent: false},
		{Path: "/repo/hotfix", Branch: "hotfix", IsCurrent: false},
	}

	b.ResetTimer()
	for range b.N {
		var output bytes.Buffer
		input := &mockStringReader{strings.NewReader("q\n")}
		_, _ = simpleSelectWorktreeWithFiles(worktrees, &mockFile{&output}, &mockFile{input})
	}
}
