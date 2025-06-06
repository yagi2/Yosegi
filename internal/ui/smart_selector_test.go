package ui

import (
	"testing"

	"github.com/yagi2/yosegi/internal/git"
)

func TestSmartSelectWorktreeEmptyList(t *testing.T) {
	// Test with empty worktree list
	var worktrees []git.Worktree
	
	result, err := SmartSelectWorktree(worktrees)
	
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

func TestFallbackSelector(t *testing.T) {
	tests := []struct {
		name      string
		worktrees []git.Worktree
		expectErr bool
		expectIdx int
	}{
		{
			name:      "Empty worktrees",
			worktrees: []git.Worktree{},
			expectErr: true,
			expectIdx: -1,
		},
		{
			name: "Single non-current worktree",
			worktrees: []git.Worktree{
				{Path: "/path/1", Branch: "main", IsCurrent: false},
			},
			expectErr: false,
			expectIdx: 0,
		},
		{
			name: "Single current worktree",
			worktrees: []git.Worktree{
				{Path: "/path/1", Branch: "main", IsCurrent: true},
			},
			expectErr: false,
			expectIdx: 0,
		},
		{
			name: "Mixed worktrees - first non-current",
			worktrees: []git.Worktree{
				{Path: "/path/1", Branch: "main", IsCurrent: false},
				{Path: "/path/2", Branch: "feature", IsCurrent: true},
			},
			expectErr: false,
			expectIdx: 0,
		},
		{
			name: "Mixed worktrees - current first",
			worktrees: []git.Worktree{
				{Path: "/path/1", Branch: "main", IsCurrent: true},
				{Path: "/path/2", Branch: "feature", IsCurrent: false},
			},
			expectErr: false,
			expectIdx: 1,
		},
		{
			name: "All current worktrees",
			worktrees: []git.Worktree{
				{Path: "/path/1", Branch: "main", IsCurrent: true},
				{Path: "/path/2", Branch: "feature", IsCurrent: true},
			},
			expectErr: false,
			expectIdx: 0,
		},
		{
			name: "Multiple non-current worktrees",
			worktrees: []git.Worktree{
				{Path: "/path/1", Branch: "main", IsCurrent: false},
				{Path: "/path/2", Branch: "feature", IsCurrent: false},
				{Path: "/path/3", Branch: "hotfix", IsCurrent: false},
			},
			expectErr: false,
			expectIdx: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fallbackSelector(tt.worktrees)
			
			if (err != nil) != tt.expectErr {
				t.Errorf("fallbackSelector() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			
			if tt.expectErr {
				if result != nil {
					t.Error("Expected nil result for error case")
				}
				return
			}
			
			if result == nil {
				t.Error("Expected non-nil result for success case")
				return
			}
			
			if tt.expectIdx >= 0 && tt.expectIdx < len(tt.worktrees) {
				expected := &tt.worktrees[tt.expectIdx]
				if result.Path != expected.Path {
					t.Errorf("Expected worktree path %s, got %s", expected.Path, result.Path)
				}
				if result.Branch != expected.Branch {
					t.Errorf("Expected worktree branch %s, got %s", expected.Branch, result.Branch)
				}
			}
		})
	}
}

func TestFallbackSelectorLogic(t *testing.T) {
	// Test specific logic cases
	t.Run("PreferNonCurrent", func(t *testing.T) {
		worktrees := []git.Worktree{
			{Path: "/current", Branch: "main", IsCurrent: true},
			{Path: "/non-current-1", Branch: "feature-1", IsCurrent: false},
			{Path: "/non-current-2", Branch: "feature-2", IsCurrent: false},
		}
		
		result, err := fallbackSelector(worktrees)
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if result.Path != "/non-current-1" {
			t.Errorf("Expected first non-current worktree, got %s", result.Path)
		}
		
		if result.IsCurrent {
			t.Error("Selected worktree should not be current")
		}
	})
	
	t.Run("AllCurrentFallsBackToFirst", func(t *testing.T) {
		worktrees := []git.Worktree{
			{Path: "/first", Branch: "main", IsCurrent: true},
			{Path: "/second", Branch: "feature", IsCurrent: true},
		}
		
		result, err := fallbackSelector(worktrees)
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if result.Path != "/first" {
			t.Errorf("Expected first worktree when all current, got %s", result.Path)
		}
	})
}

// Test helper function to create sample worktrees
func createSampleWorktrees() []git.Worktree {
	return []git.Worktree{
		{Path: "/repo/main", Branch: "main", IsCurrent: true},
		{Path: "/repo/feature-1", Branch: "feature-1", IsCurrent: false},
		{Path: "/repo/feature-2", Branch: "feature-2", IsCurrent: false},
	}
}

func TestSmartSelectWorktreeIntegration(t *testing.T) {
	// Integration test - this will use actual TTY detection
	// The specific selector used depends on the test environment
	worktrees := createSampleWorktrees()
	
	// This test verifies the function doesn't panic and returns a result
	result, err := SmartSelectWorktree(worktrees)
	
	// In CI/testing environment, this will likely use fallbackSelector
	// which should return the first non-current worktree
	if err != nil {
		// Some selectors might fail in CI environment, that's okay
		t.Logf("SmartSelectWorktree failed (expected in CI): %v", err)
		return
	}
	
	if result == nil {
		t.Error("Expected non-nil result when no error")
		return
	}
	
	// Verify we got a valid worktree
	found := false
	for _, wt := range worktrees {
		if wt.Path == result.Path && wt.Branch == result.Branch {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Result worktree not found in original list")
	}
}

// Benchmark tests
func BenchmarkFallbackSelector(b *testing.B) {
	worktrees := createSampleWorktrees()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fallbackSelector(worktrees)
	}
}

func BenchmarkSmartSelectWorktree(b *testing.B) {
	worktrees := createSampleWorktrees()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SmartSelectWorktree(worktrees)
	}
}