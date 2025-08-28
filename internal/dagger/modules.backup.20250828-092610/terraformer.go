package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// TerraformerModule runs Terraformer for infrastructure discovery
type TerraformerModule struct {
	client *dagger.Client
	name   string
}

const terraformerBinary = "/usr/local/bin/terraformer"

// NewTerraformerModule creates a new Terraformer module
func NewTerraformerModule(client *dagger.Client) *TerraformerModule {
	return &TerraformerModule{
		client: client,
		name:   terraformerBinary,
	}
}

// ImportAWS imports AWS resources
func (m *TerraformerModule) ImportAWS(ctx context.Context, region string, services []string) (string, error) {
	args := []string{
		"terraformer",
		"import", "aws",
		"--regions", region,
		"--output", "/output",
	}

	for _, service := range services {
		args = append(args, "--resources", service)
	}

	container := m.client.Container().
		From("quay.io/GoogleCloudPlatform/terraformer:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to import AWS resources: %w", err)
	}

	return output, nil
}

// ImportGCP imports GCP resources
func (m *TerraformerModule) ImportGCP(ctx context.Context, project string, services []string) (string, error) {
	args := []string{
		"terraformer",
		"import", "google",
		"--projects", project,
		"--output", "/output",
	}

	for _, service := range services {
		args = append(args, "--resources", service)
	}

	container := m.client.Container().
		From("quay.io/GoogleCloudPlatform/terraformer:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to import GCP resources: %w", err)
	}

	return output, nil
}

// ImportAzure imports Azure resources
func (m *TerraformerModule) ImportAzure(ctx context.Context, subscription string, services []string) (string, error) {
	args := []string{
		"terraformer",
		"import", "azure",
		"--subscription", subscription,
		"--output", "/output",
	}

	for _, service := range services {
		args = append(args, "--resources", service)
	}

	container := m.client.Container().
		From("quay.io/GoogleCloudPlatform/terraformer:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to import Azure resources: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of terraformer
func (m *TerraformerModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("quay.io/GoogleCloudPlatform/terraformer:latest").
		WithExec([]string{terraformerBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get terraformer version: %w", err)
	}

	return output, nil
}

// ListResources lists available resources for a provider
func (m *TerraformerModule) ListResources(ctx context.Context, provider string) (string, error) {
	container := m.client.Container().
		From("quay.io/GoogleCloudPlatform/terraformer:latest").
		WithExec([]string{
			terraformerBinary,
			"import", provider,
			"--list",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list resources for provider %s: %w", provider, err)
	}

	return output, nil
}

// Plan shows what resources would be imported without actually importing
func (m *TerraformerModule) Plan(ctx context.Context, provider string, region string, resources []string) (string, error) {
	args := []string{
		terraformerBinary,
		"plan", provider,
		"--regions", region,
		"--output", "/output",
	}

	for _, resource := range resources {
		args = append(args, "--resources", resource)
	}

	container := m.client.Container().
		From("quay.io/GoogleCloudPlatform/terraformer:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to plan import for provider %s: %w", provider, err)
	}

	return output, nil
}

// Import imports resources from any supported provider
func (m *TerraformerModule) Import(ctx context.Context, provider string, region string, resources []string, extraArgs map[string]string) (string, error) {
	args := []string{
		terraformerBinary,
		"import", provider,
		"--regions", region,
		"--output", "/output",
	}

	for _, resource := range resources {
		args = append(args, "--resources", resource)
	}

	// Add extra provider-specific arguments
	for key, value := range extraArgs {
		args = append(args, "--"+key, value)
	}

	container := m.client.Container().
		From("quay.io/GoogleCloudPlatform/terraformer:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to import resources from provider %s: %w", provider, err)
	}

	return output, nil
}
