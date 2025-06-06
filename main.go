package main

import (
	"github.com/yagi2/yosegi/cmd"
)

// Build information (set by goreleaser)
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Set version information for the CLI
	cmd.SetVersionInfo(version, commit, date, builtBy)
	cmd.Execute()
}
