package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// PackerModule runs Packer for image building
type PackerModule struct {
	client *dagger.Client
	name   string
}

// NewPackerModule creates a new Packer module
func NewPackerModule(client *dagger.Client) *PackerModule {
	return &PackerModule{
		client: client,
		name:   "packer",
	}
}

// BuildImage builds an image using Packer
func (m *PackerModule) BuildImage(ctx context.Context, templatePath string, varsFile string) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithFile("/template.pkr.hcl", m.client.Host().File(templatePath))

	if varsFile != "" {
		container = container.WithFile("/vars.pkrvars.hcl", m.client.Host().File(varsFile))
	}

	args := []string{"packer", "build"}
	if varsFile != "" {
		args = append(args, "-var-file=/vars.pkrvars.hcl")
	}
	args = append(args, "/template.pkr.hcl")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}

	return output, nil
}

// ValidateTemplate validates a Packer template
func (m *PackerModule) ValidateTemplate(ctx context.Context, templatePath string) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithFile("/template.pkr.hcl", m.client.Host().File(templatePath)).
		WithExec([]string{
			"packer", "validate",
			"/template.pkr.hcl",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate template: %w", err)
	}

	return output, nil
}

// FormatTemplate formats a Packer template
func (m *PackerModule) FormatTemplate(ctx context.Context, templatePath string) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithFile("/template.pkr.hcl", m.client.Host().File(templatePath)).
		WithExec([]string{
			"packer", "fmt",
			"/template.pkr.hcl",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to format template: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Packer
func (m *PackerModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithExec([]string{"packer", "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get packer version: %w", err)
	}

	return output, nil
}
