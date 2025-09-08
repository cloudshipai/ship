package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/cloudshipai/ship/internal/cli"
	"github.com/cloudshipai/ship/internal/telemetry"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// getVersion returns the version, trying build info first, then fallback to ldflags
func getVersion() string {
	if version != "dev" {
		return version
	}

	// Try to get version from build info (works with go install)
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "(devel)" && info.Main.Version != "" {
			return info.Main.Version
		}
	}

	return version
}

func main() {
	// Initialize telemetry
	if err := telemetry.Init(); err != nil {
		// Don't fail the CLI if telemetry fails
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize telemetry: %v\n", err)
	}
	defer telemetry.Close()

	// Set version for telemetry
	os.Setenv("SHIP_VERSION", getVersion())

	// Track app start
	command := "ship"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	telemetry.TrackAppStart(command)

	if err := cli.Execute(getVersion(), commit, date); err != nil {
		// Track error
		telemetry.TrackError("cli_execution", "main", err.Error())
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
