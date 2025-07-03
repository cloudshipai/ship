package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"dagger.io/dagger"
)

// SteampipeModule runs Steampipe for cloud infrastructure analysis
type SteampipeModule struct {
	client *dagger.Client
	name   string
}

// TableColumn represents a column in a Steampipe table
type TableColumn struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
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
	
	// Configure Steampipe directories
	container = container.WithExec([]string{
		"sh", "-c", 
		"mkdir -p /home/steampipe/.steampipe/config",
	})

	// For AWS, we'll use environment variables only (no profile mounting)
	// This avoids the "failed to get shared config profile" error
	// The credentials are passed via the credentials map below

	// Set environment variables for credentials
	for key, value := range credentials {
		if value != "" {
			container = container.WithEnvVariable(key, value)
		}
	}

	// For AWS, configure region-specific connection to avoid scanning all regions
	if plugin == "aws" {
		// Get the AWS region from credentials or default to us-east-1
		awsRegion := credentials["AWS_REGION"]
		if awsRegion == "" {
			awsRegion = "us-east-1"
		}

		// Create AWS connection config with specific region
		awsConfig := fmt.Sprintf(`connection "aws" {
  plugin = "aws"
  regions = ["%s"]
}`, awsRegion)

		// Write the config file
		container = container.WithExec([]string{
			"sh", "-c", 
			fmt.Sprintf("echo '%s' > /home/steampipe/.steampipe/config/aws.spc", awsConfig),
		})

		// Check Steampipe AWS plugin connection
		checkContainer := container.WithExec([]string{
			"sh", "-c", "steampipe plugin list && echo '---' && steampipe connection list",
		})
		connectionStatus, _ := checkContainer.Stdout(ctx)
		fmt.Printf("Steampipe AWS connection status (region: %s):\n%s\n", awsRegion, connectionStatus)
	}
	
	// Validate query doesn't contain multiple statements
	if strings.Contains(query, ";") && strings.Count(query, ";") > 1 {
		// Multiple statements detected - split and execute only the first
		statements := strings.Split(query, ";")
		query = strings.TrimSpace(statements[0])
		slog.Warn("Multiple SQL statements detected, executing only the first", "original_count", len(statements))
	}
	
	// Execute the query with explicit timeout and error capture
	// Escape single quotes in the query to prevent shell injection issues
	escapedQuery := strings.ReplaceAll(query, "'", "'\"'\"'")
	queryContainer := container.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("steampipe query '%s' --output %s 2>&1 || echo 'EXIT_CODE:'$?", escapedQuery, output),
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

	// Configure Steampipe directories
	container = container.WithExec([]string{
		"sh", "-c", 
		"mkdir -p /home/steampipe/.steampipe/config",
	})

	// Set environment variables for credentials
	for key, value := range credentials {
		if value != "" {
			container = container.WithEnvVariable(key, value)
		}
	}

	// For AWS, configure region-specific connection to avoid scanning all regions
	if plugin == "aws" {
		// Get the AWS region from credentials or default to us-east-1
		awsRegion := credentials["AWS_REGION"]
		if awsRegion == "" {
			awsRegion = "us-east-1"
		}

		// Create AWS connection config with specific region
		awsConfig := fmt.Sprintf(`connection "aws" {
  plugin = "aws"
  regions = ["%s"]
}`, awsRegion)

		// Write the config file
		container = container.WithExec([]string{
			"sh", "-c", 
			fmt.Sprintf("echo '%s' > /home/steampipe/.steampipe/config/aws.spc", awsConfig),
		})
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

// GetTableSchema queries the schema of a specific table
func (m *SteampipeModule) GetTableSchema(ctx context.Context, plugin string, tableName string, credentials map[string]string) ([]TableColumn, error) {
	// Query to get column information from information_schema
	schemaQuery := fmt.Sprintf(`
		SELECT 
			column_name,
			data_type
		FROM information_schema.columns 
		WHERE table_name = '%s'
		ORDER BY ordinal_position
	`, tableName)

	_, err := m.RunQuery(ctx, plugin, schemaQuery, credentials, "json")
	if err != nil {
		return nil, fmt.Errorf("failed to query table schema: %w", err)
	}

	// Parse JSON result to extract columns
	var columns []TableColumn
	// TODO: Parse the JSON result here
	return columns, nil
}

// GetTableColumns queries the columns of a specific table dynamically
func (m *SteampipeModule) GetTableColumns(ctx context.Context, plugin string, tableName string, credentials map[string]string) ([]string, error) {
	// Query the information_schema to get actual column names
	schemaQuery := fmt.Sprintf(`
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = '%s' 
		ORDER BY ordinal_position
	`, tableName)
	
	result, err := m.RunQuery(ctx, plugin, schemaQuery, credentials, "json")
	if err != nil {
		// Fallback to empty if query fails
		slog.Debug("Failed to query table schema", "table", tableName, "error", err)
		return []string{}, nil
	}
	
	// Parse the JSON result to extract column names
	var schemaResult []map[string]interface{}
	if err := json.Unmarshal([]byte(result), &schemaResult); err != nil {
		slog.Debug("Failed to parse schema result", "error", err)
		return []string{}, nil
	}
	
	columns := make([]string, 0, len(schemaResult))
	for _, row := range schemaResult {
		if colName, ok := row["column_name"].(string); ok {
			columns = append(columns, colName)
		}
	}
	
	slog.Debug("Retrieved table columns", "table", tableName, "columns", columns)
	return columns, nil
}

// GetAvailableTables queries all available tables for a provider
func (m *SteampipeModule) GetAvailableTables(ctx context.Context, plugin string, credentials map[string]string) ([]string, error) {
	// Query to list all tables for the provider
	tablesQuery := fmt.Sprintf(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name LIKE '%s_%%'
		ORDER BY table_name
	`, plugin)

	_, err := m.RunQuery(ctx, plugin, tablesQuery, credentials, "json")
	if err != nil {
		return nil, fmt.Errorf("failed to query available tables: %w", err)
	}

	// Parse result to extract table names
	// TODO: Parse the JSON result here
	return []string{}, nil
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
