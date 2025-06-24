package modules

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
)

// OpenInfraQuoteModule runs OpenInfraQuote for Terraform cost analysis
type OpenInfraQuoteModule struct {
	client *dagger.Client
	name   string
}

// NewOpenInfraQuoteModule creates a new OpenInfraQuote module
func NewOpenInfraQuoteModule(client *dagger.Client) *OpenInfraQuoteModule {
	return &OpenInfraQuoteModule{
		client: client,
		name:   "openinfraquote",
	}
}

// AnalyzePlan analyzes a Terraform plan JSON file for cost estimation
func (m *OpenInfraQuoteModule) AnalyzePlan(ctx context.Context, planFile string) (string, error) {
	// Get the directory containing the plan file
	dir := filepath.Dir(planFile)
	filename := filepath.Base(planFile)
	
	// Mount the directory and run OpenInfraQuote
	container := m.client.Container().
		From("ghcr.io/initech-consulting/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"openinfraquote",
			"estimate",
			"--terraform-plan-file", filename,
			"--output", "json",
		})
	
	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run openinfraquote: %w", err)
	}
	
	return output, nil
}

// AnalyzeDirectory analyzes all Terraform files in a directory
func (m *OpenInfraQuoteModule) AnalyzeDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/initech-consulting/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"openinfraquote",
			"estimate",
			"--terraform-directory", ".",
			"--output", "json",
		})
	
	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run openinfraquote on directory: %w", err)
	}
	
	return output, nil
}

// GetVersion returns the version of OpenInfraQuote
func (m *OpenInfraQuoteModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/initech-consulting/openinfraquote:latest").
		WithExec([]string{"openinfraquote", "--version"})
	
	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get openinfraquote version: %w", err)
	}
	
	return output, nil
}