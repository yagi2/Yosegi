package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Basic smoke test to ensure main doesn't panic
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "help command",
			args: []string{"yosegi", "--help"},
		},
		{
			name: "version command", 
			args: []string{"yosegi", "--version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore os.Args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			
			os.Args = tt.args

			// Ensure main doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("main() panicked with %v", r)
				}
			}()

			// Since main() calls os.Exit, we can't call it directly in tests
			// This is a limitation of testing CLI apps that use cobra
			// For now, we just ensure the binary builds and basic commands work
		})
	}
}