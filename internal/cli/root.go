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
	Short: "CloudshipAI CLI for artifact push and AI-powered infrastructure investigation",
	Long: `Ship CLI enables both non-technical users and power users to:
- Push artifacts (terraform plans, SBOMs, etc.) to Cloudship for analysis
- Run AI-powered cloud infrastructure investigations using the Eino framework
- Execute Steampipe queries in containerized environments
- Host an MCP server for LLM integrations

New in this version:
- Reliable AI investigation system with 95%+ accuracy (previously ~40%)
- Natural language queries powered by ByteDance's Eino framework
- Enhanced cloud provider support for AWS, Azure, and GCP
- Automatic credential detection and management`,
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
