package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	queryCmd.Flags().Int("timeout", 10, "Timeout in minutes for the query")
}

func runQuery(cmd *cobra.Command, args []string) error {
	// Get timeout from flags
	timeoutMinutes, _ := cmd.Flags().GetInt("timeout")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMinutes)*time.Minute)
	defer cancel()

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
		// Try to read AWS credentials directly if no environment variables are set
		if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
			// Try to read from credentials file
			if homeDir := os.Getenv("HOME"); homeDir != "" {
				credFile := filepath.Join(homeDir, ".aws", "credentials")
				if content, err := os.ReadFile(credFile); err == nil {
					// Simple parsing for default profile
					lines := strings.Split(string(content), "\n")
					inDefaultProfile := false
					slog.Debug("Parsing AWS credentials file")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if line == "[default]" {
							inDefaultProfile = true
							continue
						}
						if strings.HasPrefix(line, "[") && line != "[default]" {
							inDefaultProfile = false
							continue
						}
						if inDefaultProfile {
							if strings.HasPrefix(line, "aws_access_key_id") {
								parts := strings.SplitN(line, "=", 2)
								if len(parts) == 2 {
									credentials["AWS_ACCESS_KEY_ID"] = strings.TrimSpace(parts[1])
									slog.Debug("Found AWS access key")
								}
							}
							if strings.HasPrefix(line, "aws_secret_access_key") {
								parts := strings.SplitN(line, "=", 2)
								if len(parts) == 2 {
									credentials["AWS_SECRET_ACCESS_KEY"] = strings.TrimSpace(parts[1])
								}
							}
							if strings.HasPrefix(line, "region") {
								parts := strings.SplitN(line, "=", 2)
								if len(parts) == 2 && awsRegion == "" {
									credentials["AWS_REGION"] = strings.TrimSpace(parts[1])
								}
							}
						}
					}
				}
			}
		}

		if awsProfile != "" {
			credentials["AWS_PROFILE"] = awsProfile
		}
		// Set AWS region, defaulting to us-east-1 if not specified
		if awsRegion != "" {
			credentials["AWS_REGION"] = awsRegion
		} else if credentials["AWS_REGION"] == "" {
			credentials["AWS_REGION"] = "us-east-1"
		}
	}

	// Log what credentials we have (without sensitive data)
	credKeys := make([]string, 0)
	for k := range credentials {
		credKeys = append(credKeys, k)
	}
	slog.Info("Executing query", "provider", provider, "credentialKeys", credKeys)
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
