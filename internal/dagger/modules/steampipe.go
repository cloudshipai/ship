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
func (m *SteampipeModule) RunQuery(ctx context.Context, plugin string, query string, credentials map[string]string) (string, error) {
	// Start with base Steampipe container
	container := m.client.Container().
		From("turbot/steampipe:latest")

	// Install the required plugin if specified
	if plugin != "" {
		container = container.WithExec([]string{"steampipe", "plugin", "install", plugin})
	}

	// Mount AWS credentials from host if available and plugin is AWS
	if plugin == "aws" {
		if homeDir := os.Getenv("HOME"); homeDir != "" {
			awsCredsPath := filepath.Join(homeDir, ".aws")
			if _, err := os.Stat(awsCredsPath); err == nil {
				awsCreds := m.client.Host().Directory(awsCredsPath)
				container = container.WithDirectory("/root/.aws", awsCreds)
			}
		}
	}

	// Set environment variables for credentials
	for key, value := range credentials {
		if value != "" {
			container = container.WithEnvVariable(key, value)
		}
	}

	// If AWS_PROFILE is set, we need to configure Steampipe to use it
	if profile, ok := credentials["AWS_PROFILE"]; ok && profile != "" {
		container = container.WithEnvVariable("AWS_SDK_LOAD_CONFIG", "1")
		// Also set the shared config file location
		container = container.WithEnvVariable("AWS_SHARED_CREDENTIALS_FILE", "/root/.aws/credentials")
		container = container.WithEnvVariable("AWS_CONFIG_FILE", "/root/.aws/config")
	}

	// Execute the query
	container = container.WithExec([]string{
		"steampipe", "query", query, "--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run steampipe query: %w", err)
	}

	return output, nil
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
				container = container.WithDirectory("/root/.aws", awsCreds)
			}
		}
	}

	// Set environment variables for credentials
	for key, value := range credentials {
		if value != "" {
			container = container.WithEnvVariable(key, value)
		}
	}

	// If AWS_PROFILE is set, we need to configure Steampipe to use it
	if profile, ok := credentials["AWS_PROFILE"]; ok && profile != "" {
		container = container.WithEnvVariable("AWS_SDK_LOAD_CONFIG", "1")
		// Also set the shared config file location
		container = container.WithEnvVariable("AWS_SHARED_CREDENTIALS_FILE", "/root/.aws/credentials")
		container = container.WithEnvVariable("AWS_CONFIG_FILE", "/root/.aws/config")
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
