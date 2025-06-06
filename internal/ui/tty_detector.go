package ui

import (
	"os"
	"runtime"

	"github.com/mattn/go-isatty"
)

// TTYCapability represents the level of TTY control available
type TTYCapability int

const (
	NoTTYControl TTYCapability = iota
	BasicTTYControl
	FullTTYControl
)

// DetectTTYCapability determines what level of TTY control is available
func DetectTTYCapability() TTYCapability {
	// Check if we're in a command substitution or pipe
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		// We're in command substitution, check if we can access TTY directly
		if runtime.GOOS != "windows" {
			// Try to open /dev/tty
			if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
				tty.Close()
				// We can access TTY directly - basic control available
				return BasicTTYControl
			}
		}
		
		// Check if stderr is a terminal (fallback method)
		if isatty.IsTerminal(os.Stderr.Fd()) {
			return BasicTTYControl
		}
		
		// For development/testing: try to force BasicTTYControl when appropriate
		// If stdin looks like it could be a terminal, try basic control
		if isatty.IsTerminal(os.Stdin.Fd()) {
			return BasicTTYControl
		}
		
		// No TTY access at all
		return NoTTYControl
	}
	
	// Full terminal environment available
	return FullTTYControl
}

// GetTTYFiles returns appropriate input/output files based on capability
func GetTTYFiles(capability TTYCapability) (*os.File, *os.File, func(), error) {
	switch capability {
	case FullTTYControl:
		// Normal terminal environment
		return os.Stdin, os.Stderr, func() {}, nil
		
	case BasicTTYControl:
		if runtime.GOOS != "windows" {
			// Unix: try to open /dev/tty
			tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
			if err != nil {
				return nil, nil, nil, err
			}
			cleanup := func() { tty.Close() }
			return tty, tty, cleanup, nil
		} else {
			// Windows: use stdin/stderr
			return os.Stdin, os.Stderr, func() {}, nil
		}
		
	case NoTTYControl:
		// Fallback - return stdin/stderr even though they might not work
		return os.Stdin, os.Stderr, func() {}, nil
		
	default:
		return os.Stdin, os.Stderr, func() {}, nil
	}
}