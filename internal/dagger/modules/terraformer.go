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

// NewTerraformerModule creates a new Terraformer module
func NewTerraformerModule(client *dagger.Client) *TerraformerModule {
	return &TerraformerModule{
		client: client,
		name:   "terraformer",
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
