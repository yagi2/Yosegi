package cmd

import "os"

// isatty checks if stdout is connected to a terminal
func isatty() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
