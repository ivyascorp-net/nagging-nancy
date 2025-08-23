package main

import (
	"fmt"
	"os"

	"github.com/ivyascorp-net/nagging-nancy/internal/app"
	"github.com/ivyascorp-net/nagging-nancy/internal/cli"
)

var (
	// These will be set during build time
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Set version information in the app package
	app.Version = version
	app.BuildTime = buildTime
	app.GitCommit = gitCommit

	// Initialize application directories on first run
	if err := app.InitApp(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Execute CLI commands
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
