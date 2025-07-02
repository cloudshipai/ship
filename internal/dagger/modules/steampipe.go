package modules

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

// SteampipeModule runs Steampipe for cloud infrastructure analysis
type SteampipeModule struct {
	client *dagger.Client
	name   string
}

// NewSteampipeModule creates a new Steampipe module
func NewSteampipeModule(client *dagger.Client) *SteampipeModule {
	return &SteampipeModule{
		client: client,
		name:   "steampipe",
	}
}

// RunQuery executes a Steampipe query with the specified plugin
func (m *SteampipeModule) RunQuery(ctx context.Context, plugin string, query string, credentials map[string]string, outputFormat ...string) (string, error) {
	output := "json"
	if len(outputFormat) > 0 {
		output = outputFormat[0]
	}

	// Start with base Steampipe container
	container := m.client.Container().
		From("turbot/steampipe:latest")

	// Install the required plugin if specified
	if plugin != "" {
		container = container.WithExec([]string{"steampipe", "plugin", "install", plugin})
	}

	// For AWS, we'll use environment variables only (no profile mounting)
	// This avoids the "failed to get shared config profile" error
	// The credentials are passed via the credentials map below

	// Set environment variables for credentials
	for key, value := range credentials {
		if value != "" {
			container = container.WithEnvVariable(key, value)
		}
	}

	// For AWS, check the Steampipe connection status
	if plugin == "aws" {
		// Check Steampipe AWS plugin connection
		checkContainer := container.WithExec([]string{
			"sh", "-c", "steampipe plugin list && echo '---' && steampipe connection list",
		})
		connectionStatus, _ := checkContainer.Stdout(ctx)
		fmt.Printf("Steampipe AWS connection status:\n%s\n", connectionStatus)
	}
	
	// Execute the query with explicit timeout and error capture
	queryContainer := container.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("steampipe query '%s' --output %s 2>&1 || echo 'EXIT_CODE:'$?", query, output),
	})

	stdout, err := queryContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run steampipe query: %w", err)
	}
	
	// Check if we got an error
	if strings.Contains(stdout, "EXIT_CODE:") || strings.Contains(stdout, "Error:") {
		return "", fmt.Errorf("steampipe query failed: %s", stdout)
	}

	return stdout, nil
}

// RunMultipleQueries executes multiple Steampipe queries
func (m *SteampipeModule) RunMultipleQueries(ctx context.Context, plugin string, queries []string, credentials map[string]string) (map[string]string, error) {
	results := make(map[string]string)

	// Start with base Steampipe container
	container := m.client.Container().
		From("turbot/steampipe:latest")

	// Install the required plugin
	container = container.WithExec([]string{"steampipe", "plugin", "install", plugin})

	// Mount AWS credentials from host if available and plugin is AWS
	if plugin == "aws" {
		if homeDir := os.Getenv("HOME"); homeDir != "" {
			awsCredsPath := filepath.Join(homeDir, ".aws")
			if _, err := os.Stat(awsCredsPath); err == nil {
				awsCreds := m.client.Host().Directory(awsCredsPath)
				// Mount to both /root/.aws and /home/steampipe/.aws since Steampipe runs as steampipe user
				container = container.WithDirectory("/home/steampipe/.aws", awsCreds)
				
				// Always set AWS SDK environment variables when credentials are mounted
				// This ensures the AWS SDK knows to look for the mounted credential files
				container = container.WithEnvVariable("AWS_SDK_LOAD_CONFIG", "1")
				container = container.WithEnvVariable("AWS_SHARED_CREDENTIALS_FILE", "/home/steampipe/.aws/credentials")
				container = container.WithEnvVariable("AWS_CONFIG_FILE", "/home/steampipe/.aws/config")
				
				// If no AWS_PROFILE is explicitly set, default to "default" profile
				if awsProfile, hasProfile := credentials["AWS_PROFILE"]; !hasProfile || awsProfile == "" {
					container = container.WithEnvVariable("AWS_PROFILE", "default")
				}
			}
		}
	}

	// Set environment variables for credentials
	for key, value := range credentials {
		if value != "" {
			container = container.WithEnvVariable(key, value)
		}
	}

	// Execute each query
	for i, query := range queries {
		queryContainer := container.WithExec([]string{
			"steampipe", "query", query, "--output", "json",
		})

		output, err := queryContainer.Stdout(ctx)
		if err != nil {
			results[fmt.Sprintf("query_%d_error", i)] = err.Error()
			continue
		}

		results[fmt.Sprintf("query_%d", i)] = output
	}

	return results, nil
}

// RunModCheck runs a Steampipe mod check (compliance framework)
func (m *SteampipeModule) RunModCheck(ctx context.Context, plugin string, modPath string, credentials map[string]string) (string, error) {
	// Start with base Steampipe container
	container := m.client.Container().
		From("turbot/steampipe:latest")

	// Install the required plugin
	container = container.WithExec([]string{"steampipe", "plugin", "install", plugin})

	// Set environment variables for credentials
	for key, value := range credentials {
		container = container.WithEnvVariable(key, value)
	}

	// Install the mod if it's a URL
	if strings.HasPrefix(modPath, "https://") || strings.HasPrefix(modPath, "github.com/") {
		container = container.WithExec([]string{
			"steampipe", "mod", "install", modPath,
		})
	}

	// Run the check
	container = container.WithExec([]string{
		"steampipe", "check", "all", "--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run steampipe mod check: %w", err)
	}

	return output, nil
}

// GetInstalledPlugins returns a list of installed Steampipe plugins
func (m *SteampipeModule) GetInstalledPlugins(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "plugin", "list", "--output", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list steampipe plugins: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Steampipe
func (m *SteampipeModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("turbot/steampipe:latest").
		WithExec([]string{"steampipe", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get steampipe version: %w", err)
	}

	return output, nil
}

// RunInteractiveQuery runs an interactive query session (for development)
func (m *SteampipeModule) RunInteractiveQuery(ctx context.Context, plugin string, credentials map[string]string) (*dagger.Container, error) {
	// Start with base Steampipe container
	container := m.client.Container().
		From("turbot/steampipe:latest")

	// Install the required plugin
	container = container.WithExec([]string{"steampipe", "plugin", "install", plugin})

	// Set environment variables for credentials
	for key, value := range credentials {
		container = container.WithEnvVariable(key, value)
	}

	// Start Steampipe service
	container = container.WithExec([]string{"steampipe", "service", "start"})

	return container, nil
}

// RunBenchmark runs a specific benchmark from a compliance mod
func (m *SteampipeModule) RunBenchmark(ctx context.Context, plugin string, benchmark string, credentials map[string]string) (string, error) {
	// Start with base Steampipe container
	container := m.client.Container().
		From("turbot/steampipe:latest")

	// Install the required plugin
	container = container.WithExec([]string{"steampipe", "plugin", "install", plugin})

	// Set environment variables for credentials
	for key, value := range credentials {
		container = container.WithEnvVariable(key, value)
	}

	// Run the specific benchmark
	container = container.WithExec([]string{
		"steampipe", "check", benchmark, "--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run steampipe benchmark: %w", err)
	}

	return output, nil
}
