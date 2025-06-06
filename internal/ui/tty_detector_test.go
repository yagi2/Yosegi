package ui

import (
	"os"
	"runtime"
	"testing"
)

func TestTTYCapabilityConstants(t *testing.T) {
	// Test that constants have expected values
	if NoTTYControl != 0 {
		t.Errorf("Expected NoTTYControl to be 0, got %d", NoTTYControl)
	}
	if BasicTTYControl != 1 {
		t.Errorf("Expected BasicTTYControl to be 1, got %d", BasicTTYControl)
	}
	if FullTTYControl != 2 {
		t.Errorf("Expected FullTTYControl to be 2, got %d", FullTTYControl)
	}
}

func TestDetectTTYCapability(t *testing.T) {
	// This test verifies the function runs without error
	// Actual result depends on the test environment
	capability := DetectTTYCapability()

	// Should return one of the valid values
	switch capability {
	case NoTTYControl, BasicTTYControl, FullTTYControl:
		// Valid capability detected
	default:
		t.Errorf("DetectTTYCapability returned invalid capability: %d", capability)
	}
}

func TestGetTTYFiles(t *testing.T) {
	tests := []struct {
		name       string
		capability TTYCapability
		expectErr  bool
	}{
		{
			name:       "FullTTYControl",
			capability: FullTTYControl,
			expectErr:  false,
		},
		{
			name:       "BasicTTYControl",
			capability: BasicTTYControl,
			expectErr:  runtime.GOOS != "windows", // On Windows, falls back to stdin/stderr
		},
		{
			name:       "NoTTYControl",
			capability: NoTTYControl,
			expectErr:  false,
		},
		{
			name:       "InvalidCapability",
			capability: TTYCapability(999),
			expectErr:  false, // Should fall through to default case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, output, cleanup, err := GetTTYFiles(tt.capability)

			if (err != nil) != tt.expectErr {
				t.Errorf("GetTTYFiles() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if err == nil {
				// Verify we got valid file handles
				if input == nil {
					t.Error("GetTTYFiles() returned nil input file")
				}
				if output == nil {
					t.Error("GetTTYFiles() returned nil output file")
				}
				if cleanup == nil {
					t.Error("GetTTYFiles() returned nil cleanup function")
				}

				// Call cleanup function
				cleanup()
			}
		})
	}
}

func TestGetTTYFilesFullTTYControl(t *testing.T) {
	input, output, cleanup, err := GetTTYFiles(FullTTYControl)

	if err != nil {
		t.Fatalf("GetTTYFiles(FullTTYControl) failed: %v", err)
	}

	// Should return stdin and stderr
	if input != os.Stdin {
		t.Error("Expected input to be os.Stdin for FullTTYControl")
	}
	if output != os.Stderr {
		t.Error("Expected output to be os.Stderr for FullTTYControl")
	}

	// Cleanup should be a no-op
	cleanup()
}

func TestGetTTYFilesBasicTTYControlWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	input, output, cleanup, err := GetTTYFiles(BasicTTYControl)

	if err != nil {
		t.Fatalf("GetTTYFiles(BasicTTYControl) failed on Windows: %v", err)
	}

	// On Windows, should return stdin and stderr
	if input != os.Stdin {
		t.Error("Expected input to be os.Stdin for BasicTTYControl on Windows")
	}
	if output != os.Stderr {
		t.Error("Expected output to be os.Stderr for BasicTTYControl on Windows")
	}

	cleanup()
}

func TestGetTTYFilesNoTTYControl(t *testing.T) {
	input, output, cleanup, err := GetTTYFiles(NoTTYControl)

	if err != nil {
		t.Fatalf("GetTTYFiles(NoTTYControl) failed: %v", err)
	}

	// Should return stdin and stderr as fallback
	if input != os.Stdin {
		t.Error("Expected input to be os.Stdin for NoTTYControl")
	}
	if output != os.Stderr {
		t.Error("Expected output to be os.Stderr for NoTTYControl")
	}

	cleanup()
}

// Benchmark tests
func BenchmarkDetectTTYCapability(b *testing.B) {
	for b.Loop() {
		DetectTTYCapability()
	}
}

func BenchmarkGetTTYFiles(b *testing.B) {
	capabilities := []TTYCapability{
		NoTTYControl,
		BasicTTYControl,
		FullTTYControl,
	}

	for _, cap := range capabilities {
		b.Run(cap.String(), func(b *testing.B) {
			for b.Loop() {
				_, _, cleanup, err := GetTTYFiles(cap)
				if err == nil {
					cleanup()
				}
			}
		})
	}
}

// Helper method for TTYCapability to string conversion (for benchmarks)
func (c TTYCapability) String() string {
	switch c {
	case NoTTYControl:
		return "NoTTYControl"
	case BasicTTYControl:
		return "BasicTTYControl"
	case FullTTYControl:
		return "FullTTYControl"
	default:
		return "Unknown"
	}
}
