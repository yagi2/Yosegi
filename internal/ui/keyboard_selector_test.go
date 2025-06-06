package ui

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/yagi2/yosegi/internal/git"
)

func TestNewKeyboardSelector(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
	}
	
	var input, output bytes.Buffer
	selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{&input}, &mockFile{&output})
	
	if selector == nil {
		t.Error("Expected non-nil selector")
		return
	}
	
	if len(selector.worktrees) != 1 {
		t.Errorf("Expected 1 worktree, got %d", len(selector.worktrees))
	}
	
	if selector.cursor != 0 {
		t.Errorf("Expected cursor to start at 0, got %d", selector.cursor)
	}
	
	if selector.input == nil {
		t.Error("Expected non-nil input file")
	}
	
	if selector.output == nil {
		t.Error("Expected non-nil output file")
	}
}

func TestKeyboardSelectorEmptyWorktrees(t *testing.T) {
	var worktrees []git.Worktree
	var input, output bytes.Buffer
	
	selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{&input}, &mockFile{&output})
	
	// This test should focus on the structure since Run() requires terminal control
	if len(selector.worktrees) != 0 {
		t.Errorf("Expected 0 worktrees, got %d", len(selector.worktrees))
	}
}

func TestKeyboardSelectorReadKey(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
	}
	
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Enter key (CR)",
			input:    []byte{13},
			expected: "enter",
		},
		{
			name:     "Enter key (LF)",
			input:    []byte{10},
			expected: "enter",
		},
		{
			name:     "Ctrl+C",
			input:    []byte{3},
			expected: "ctrl+c",
		},
		{
			name:     "j key",
			input:    []byte{106},
			expected: "j",
		},
		{
			name:     "k key",
			input:    []byte{107},
			expected: "k",
		},
		{
			name:     "q key",
			input:    []byte{113},
			expected: "q",
		},
		{
			name:     "Up arrow",
			input:    []byte{27, 91, 65},
			expected: "up",
		},
		{
			name:     "Down arrow",
			input:    []byte{27, 91, 66},
			expected: "down",
		},
		{
			name:     "Unknown key",
			input:    []byte{120}, // 'x'
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			input := &mockKeyReader{data: tt.input}
			
			selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{input}, &mockFile{&output})
			
			key, err := selector.readKey()
			
			if err != nil {
				t.Errorf("readKey() error = %v", err)
				return
			}
			
			if key != tt.expected {
				t.Errorf("readKey() = %q, expected %q", key, tt.expected)
			}
		})
	}
}

func TestKeyboardSelectorRender(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/repo/main", Branch: "main", IsCurrent: true},
		{Path: "/repo/feature", Branch: "feature", IsCurrent: false},
	}
	
	var output bytes.Buffer
	var input bytes.Buffer
	
	selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{&input}, &mockFile{&output})
	selector.render()
	
	// Check for ANSI clear screen sequence
	if !bytes.Contains(output.Bytes(), []byte("\033[2J\033[H")) {
		t.Error("Expected clear screen sequence")
	}
	
	// Check for title
	if !bytes.Contains(output.Bytes(), []byte("🌲 Git Worktrees")) {
		t.Error("Expected title with tree emoji")
	}
	
	// Check for worktree entries
	if !bytes.Contains(output.Bytes(), []byte("* /repo/main (main)")) {
		t.Error("Expected current worktree to be marked with asterisk")
	}
	
	if !bytes.Contains(output.Bytes(), []byte("  /repo/feature (feature)")) {
		t.Error("Expected non-current worktree to be displayed")
	}
	
	// Check for help text
	if !bytes.Contains(output.Bytes(), []byte("↑/k up")) {
		t.Error("Expected help text for navigation")
	}
}

func TestKeyboardSelectorClearScreen(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
	}
	
	var output bytes.Buffer
	var input bytes.Buffer
	
	selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{&input}, &mockFile{&output})
	selector.clearScreen()
	
	// Check for ANSI clear screen sequence
	if !bytes.Contains(output.Bytes(), []byte("\033[2J\033[H")) {
		t.Error("Expected clear screen sequence")
	}
}

func TestKeyboardSelectorSetRawModeFailure(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
	}
	
	var output bytes.Buffer
	input := &errorReader{err: io.ErrUnexpectedEOF}
	
	selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{input}, &mockFile{&output})
	
	// setRawMode should fail with non-TTY input
	err := selector.setRawMode()
	if err == nil {
		t.Error("Expected setRawMode to fail with mock input")
	}
}

func TestKeyboardSelectorRestoreMode(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
	}
	
	var output bytes.Buffer
	var input bytes.Buffer
	
	selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{&input}, &mockFile{&output})
	
	// restoreMode should not panic even if setRawMode was never called
	selector.restoreMode()
	
	// This is mainly a smoke test to ensure no panic occurs
}

// Mock key reader for testing key input
type mockKeyReader struct {
	data []byte
	pos  int
}

func (m *mockKeyReader) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	
	n = copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockKeyReader) Write(p []byte) (n int, err error) {
	return len(p), nil // Just pretend to write
}

// errorReader for testing IO errors
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func (e *errorReader) Write(p []byte) (n int, err error) {
	return 0, e.err
}

func TestNewKeyboardSelectorWrapper(t *testing.T) {
	// Test the wrapper function NewKeyboardSelector
	worktrees := []git.Worktree{
		{Path: "/test/path", Branch: "main", IsCurrent: false},
	}
	
	// This should call the underlying implementation
	selector := NewKeyboardSelector(worktrees, os.Stdin, os.Stderr)
	
	if selector == nil {
		t.Error("NewKeyboardSelector returned nil")
		return
	}
	
	if len(selector.worktrees) != 1 {
		t.Errorf("Expected 1 worktree, got %d", len(selector.worktrees))
	}
	
	if selector.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", selector.cursor)
	}
}

// Test the terminal control commands (stty)
func TestKeyboardSelectorTerminalCommands(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
	}
	
	var output bytes.Buffer
	var input bytes.Buffer
	
	selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{&input}, &mockFile{&output})
	
	// Test that setRawMode and restoreMode don't panic
	// They will likely fail in test environment, but should not crash
	err := selector.setRawMode()
	// Don't assert on error since stty may not be available in test environment
	_ = err
	
	selector.restoreMode()
	// This should always complete without panic
}

// Benchmark tests
func BenchmarkKeyboardSelectorRender(b *testing.B) {
	worktrees := []git.Worktree{
		{Path: "/repo/main", Branch: "main", IsCurrent: true},
		{Path: "/repo/feature-1", Branch: "feature-1", IsCurrent: false},
		{Path: "/repo/feature-2", Branch: "feature-2", IsCurrent: false},
		{Path: "/repo/hotfix", Branch: "hotfix", IsCurrent: false},
	}
	
	var output bytes.Buffer
	var input bytes.Buffer
	selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{&input}, &mockFile{&output})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output.Reset()
		selector.render()
	}
}

func BenchmarkKeyboardSelectorReadKey(b *testing.B) {
	worktrees := []git.Worktree{
		{Path: "/path/1", Branch: "main", IsCurrent: false},
	}
	
	keyData := []byte{106} // 'j' key
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output bytes.Buffer
		input := &mockKeyReader{data: keyData}
		selector := newKeyboardSelectorWithFiles(worktrees, &mockFile{input}, &mockFile{&output})
		selector.readKey()
	}
}