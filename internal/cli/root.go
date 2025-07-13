package cli

import (
	"github.com/spf13/cobra"

	"github.com/cloudshipai/ship/internal/logger"
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
	Short: "CloudshipAI CLI for Terraform analysis and infrastructure tools",
	Long: `Ship CLI enables both non-technical users and power users to:
- Push artifacts (terraform plans, SBOMs, etc.) to Cloudship for analysis
- Run comprehensive Terraform analysis tools in containerized environments
- Generate infrastructure documentation and diagrams
- Host an MCP server for AI assistant integrations

Key capabilities:
- Terraform linting, security scanning, and cost analysis
- Infrastructure diagram generation from HCL files or state
- Documentation generation for Terraform modules
- Containerized tool execution with Dagger for consistency
- MCP server for Claude Code, Cursor, and other AI assistants`,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = version

	// Set up logging
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-file", "", "Log file path")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Read log level and log file from flags
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFile, _ := cmd.Flags().GetString("log-file")

		// Configure logger
		return logger.Init(logLevel, logFile)
	}
}
