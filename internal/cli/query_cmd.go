package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/cloudshipai/ship/internal/dagger"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query [SQL_QUERY]",
	Short: "Run a Steampipe query directly",
	Long:  `Executes a single Steampipe SQL query string in a containerized environment.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runQuery,
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().String("provider", "aws", "Cloud provider (aws, azure, gcp)")
	queryCmd.Flags().String("output", "json", "Output format (json, csv, table)")
	queryCmd.Flags().String("aws-profile", "", "AWS profile to use (from ~/.aws/config)")
	queryCmd.Flags().String("aws-region", "", "AWS region to use")
	queryCmd.Flags().String("log-level", "info", "Log level (debug, info, warn, error)")
}

func runQuery(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	query := args[0]

	if err := checkDockerRunning(); err != nil {
		return err
	}

	provider, _ := cmd.Flags().GetString("provider")
	output, _ := cmd.Flags().GetString("output")
	awsProfile, _ := cmd.Flags().GetString("aws-profile")
	awsRegion, _ := cmd.Flags().GetString("aws-region")
	logLevel, _ := cmd.Flags().GetString("log-level")

	slog.Info("Initializing Dagger engine...")
	engine, err := dagger.NewEngine(ctx, logLevel)
	if err != nil {
		return fmt.Errorf("failed to initialize dagger: %w", err)
	}
	defer engine.Close()

	steampipeModule := engine.NewSteampipeModule()

	credentials := getProviderCredentials(provider)
	if provider == "aws" {
		if awsProfile != "" {
			credentials["AWS_PROFILE"] = awsProfile
		}
		if awsRegion != "" {
			credentials["AWS_REGION"] = awsRegion
		}
	}

	slog.Info("Executing query", "provider", provider)
	result, err := steampipeModule.RunQuery(ctx, provider, query, credentials, output)
	if err != nil {
		return fmt.Errorf("failed to run query: %w", err)
	}

	fmt.Println(result)

	return nil
}

func checkDockerRunning() error {
	if _, err := os.Stat("/var/run/docker.sock"); os.IsNotExist(err) {
		return fmt.Errorf("docker daemon is not running. Please start Docker and try again")
	}
	return nil
}
