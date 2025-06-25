package cli

import (
	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

func Execute(v, c, d string) error {
	version = v
	commit = c
	date = d
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "ship",
	Short: "CloudshipAI CLI for artifact push and infrastructure investigation",
	Long: `Ship CLI enables both non-technical users and power users to:
- Push artifacts (terraform plans, SBOMs, etc.) to Cloudship for analysis
- Run automated cloud infrastructure investigations using Steampipe
- Host an MCP server for LLM integrations`,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = version
}