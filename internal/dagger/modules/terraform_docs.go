package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// TerraformDocsModule runs terraform-docs for documentation generation
type TerraformDocsModule struct {
	client *dagger.Client
	name   string
}

const terraformDocsBinary = "/usr/local/bin/terraform-docs"

// NewTerraformDocsModule creates a new terraform-docs module
func NewTerraformDocsModule(client *dagger.Client) *TerraformDocsModule {
	return &TerraformDocsModule{
		client: client,
		name:   terraformDocsBinary,
	}
}

// GenerateMarkdown generates markdown documentation for Terraform modules
func (m *TerraformDocsModule) GenerateMarkdown(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("quay.io/terraform-docs/terraform-docs:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			terraformDocsBinary,
			"markdown",
			".",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate terraform docs: %w", err)
	}

	return output, nil
}

// GenerateJSON generates JSON documentation for Terraform modules
func (m *TerraformDocsModule) GenerateJSON(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("quay.io/terraform-docs/terraform-docs:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			terraformDocsBinary,
			"json",
			".",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate terraform docs json: %w", err)
	}

	return output, nil
}

// GenerateWithConfig generates documentation using a config file
func (m *TerraformDocsModule) GenerateWithConfig(ctx context.Context, dir string, configFile string) (string, error) {
	container := m.client.Container().
		From("quay.io/terraform-docs/terraform-docs:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir))

	// If config file is provided, mount it
	if configFile != "" {
		container = container.WithFile("/.terraform-docs.yml", m.client.Host().File(configFile))
		container = container.WithExec([]string{
			terraformDocsBinary,
			"--config", "/.terraform-docs.yml",
			"markdown",
			"/workspace",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})
	} else {
		container = container.WithExec([]string{
			terraformDocsBinary,
			"markdown",
			"/workspace",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})
	}

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate terraform docs with config: %w", err)
	}

	return output, nil
}

// GenerateTable generates a markdown table of inputs and outputs
func (m *TerraformDocsModule) GenerateTable(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("quay.io/terraform-docs/terraform-docs:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			terraformDocsBinary,
			"markdown", "table",
			".",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate terraform docs table: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of terraform-docs
func (m *TerraformDocsModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("quay.io/terraform-docs/terraform-docs:latest").
		WithExec([]string{terraformDocsBinary, "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get terraform-docs version: %w", err)
	}

	return output, nil
}
