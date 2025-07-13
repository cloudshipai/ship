package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/cloudshipai/ship/internal/cli"
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
	if err := cli.Execute(getVersion(), commit, date); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
