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

const packerBinary = "packer"

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

	args := []string{packerBinary, "build"}
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
			packerBinary, "validate",
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
			packerBinary, "fmt",
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
		WithExec([]string{packerBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get packer version: %w", err)
	}

	return output, nil
}

// InspectTemplate inspects and analyzes Packer template configuration
func (m *PackerModule) InspectTemplate(ctx context.Context, templatePath string, machineReadable bool) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithFile("/template.pkr.hcl", m.client.Host().File(templatePath))

	args := []string{packerBinary, "inspect"}
	if machineReadable {
		args = append(args, "-machine-readable")
	}
	args = append(args, "/template.pkr.hcl")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to inspect template: %w", err)
	}

	return output, nil
}

// FixTemplate fixes and upgrades Packer template to current version
func (m *PackerModule) FixTemplate(ctx context.Context, templatePath string, validate bool) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithFile("/template.pkr.hcl", m.client.Host().File(templatePath))

	args := []string{packerBinary, "fix"}
	if validate {
		args = append(args, "-validate")
	}
	args = append(args, "/template.pkr.hcl")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fix template: %w", err)
	}

	return output, nil
}

// InitConfiguration initializes Packer configuration and installs required plugins
func (m *PackerModule) InitConfiguration(ctx context.Context, configFile string, upgrade bool) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithFile("/config.pkr.hcl", m.client.Host().File(configFile))

	args := []string{packerBinary, "init"}
	if upgrade {
		args = append(args, "-upgrade")
	}
	args = append(args, "/config.pkr.hcl")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to initialize configuration: %w", err)
	}

	return output, nil
}

// ManagePlugins manages Packer plugins (install, remove, required)
func (m *PackerModule) ManagePlugins(ctx context.Context, subcommand string, pluginName string, version string, configFile string) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest")

	if configFile != "" {
		container = container.WithFile("/config.pkr.hcl", m.client.Host().File(configFile))
	}

	args := []string{packerBinary, "plugins", subcommand}
	switch subcommand {
	case "install", "remove":
		if pluginName != "" {
			args = append(args, pluginName)
		}
		if version != "" {
			args = append(args, version)
		}
	case "required":
		if configFile != "" {
			args = append(args, "/config.pkr.hcl")
		}
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to manage plugins: %w", err)
	}

	return output, nil
}

// HCL2Upgrade upgrades JSON Packer template to HCL2
func (m *PackerModule) HCL2Upgrade(ctx context.Context, templateFile string, outputFile string, withAnnotations bool) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithFile("/template.json", m.client.Host().File(templateFile))

	args := []string{packerBinary, "hcl2_upgrade"}
	if outputFile != "" {
		args = append(args, "-output-file", "/output.pkr.hcl")
	}
	if withAnnotations {
		args = append(args, "-with-annotations")
	}
	args = append(args, "/template.json")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to upgrade to HCL2: %w", err)
	}

	return output, nil
}

// Console opens Packer console for template debugging
func (m *PackerModule) Console(ctx context.Context, templateFile string, vars string, varFile string) (string, error) {
	container := m.client.Container().
		From("hashicorp/packer:latest").
		WithFile("/template.pkr.hcl", m.client.Host().File(templateFile))

	if varFile != "" {
		container = container.WithFile("/vars.pkrvars.hcl", m.client.Host().File(varFile))
	}

	args := []string{packerBinary, "console"}
	if vars != "" {
		args = append(args, "-var", vars)
	}
	if varFile != "" {
		args = append(args, "-var-file", "/vars.pkrvars.hcl")
	}
	args = append(args, "/template.pkr.hcl")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to open console: %w", err)
	}

	return output, nil
}
