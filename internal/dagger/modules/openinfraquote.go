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
func (m *OpenInfraQuoteModule) AnalyzePlan(ctx context.Context, planFile string, region string) (string, error) {
	// Get the directory containing the plan file
	dir := filepath.Dir(planFile)
	filename := filepath.Base(planFile)
	
	// First, use a utility container to download the pricing sheet
	utilContainer := m.client.Container().
		From("alpine:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		// Download the pricing sheet
		WithExec([]string{"sh", "-c", "apk add --no-cache curl && curl -s https://oiq.terrateam.io/prices.csv.gz | gunzip > prices.csv"})
	
	// Get the workspace with pricing sheet
	_, err := utilContainer.Directory("/workspace").Export(ctx, dir+"_temp")
	if err != nil {
		return "", fmt.Errorf("failed to prepare workspace: %w", err)
	}
	defer func() {
		// Clean up temp directory
		m.client.Host().Directory(dir + "_temp")
	}()
	
	// First, get the match output
	matchContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir+"_temp")).
		WithWorkdir("/workspace").
		WithExec([]string{"oiq", "match", "--pricesheet", "prices.csv", filename})
	
	matchOutput, err := matchContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq match: %w", err)
	}
	
	// Now run oiq price with the match output as stdin using Dagger's stdin parameter
	priceContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithExec([]string{"oiq", "price", "--region", region, "--format", "json"}, dagger.ContainerWithExecOpts{
			Stdin: matchOutput,
		})
	
	priceOutput, err := priceContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq price: %w", err)
	}
	
	return priceOutput, nil
}

// AnalyzeDirectory analyzes all Terraform files in a directory
func (m *OpenInfraQuoteModule) AnalyzeDirectory(ctx context.Context, dir string, region string) (string, error) {
	// First, use a utility container to download the pricing sheet and generate Terraform plan
	terraformContainer := m.client.Container().
		From("hashicorp/terraform:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		// Download the pricing sheet
		WithExec([]string{"sh", "-c", "apk add --no-cache curl && curl -s https://oiq.terrateam.io/prices.csv.gz | gunzip > prices.csv"}).
		// Generate Terraform plan as JSON
		WithExec([]string{"sh", "-c", "terraform init && terraform plan -out=tf.plan && terraform show -json tf.plan > tfplan.json"})
	
	// Get the workspace with generated files
	_, err := terraformContainer.Directory("/workspace").Export(ctx, dir+"_temp")
	if err != nil {
		return "", fmt.Errorf("failed to prepare workspace: %w", err)
	}
	defer func() {
		// Clean up temp directory
		m.client.Host().Directory(dir + "_temp")
	}()
	
	// First, get the match output
	matchContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir+"_temp")).
		WithWorkdir("/workspace").
		WithExec([]string{"oiq", "match", "--pricesheet", "prices.csv", "tfplan.json"})
	
	matchOutput, err := matchContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq match: %w", err)
	}
	
	// Now run oiq price with the match output as stdin using Dagger's stdin parameter
	priceContainer := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithExec([]string{"oiq", "price", "--region", region, "--format", "json"}, dagger.ContainerWithExecOpts{
			Stdin: matchOutput,
		})
	
	priceOutput, err := priceContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run oiq price: %w", err)
	}
	
	return priceOutput, nil
}

// GetVersion returns the version of OpenInfraQuote
func (m *OpenInfraQuoteModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/terrateamio/openinfraquote:latest").
		WithExec([]string{"oiq", "--version"})
	
	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get openinfraquote version: %w", err)
	}
	
	return output, nil
}