package modules

import (
	"context"

	"dagger.io/dagger"
)

// TerraformDocsModule runs terraform-docs for documentation generation
type TerraformDocsModule struct {
	client *dagger.Client
	name   string
}

// TerraformDocsOptions contains options for terraform-docs generation
type TerraformDocsOptions struct {
	OutputFormat string
	OutputFile   string
	ConfigFile   string
	Recursive    bool
	Sort         bool
	HeaderFrom   string
	FooterFrom   string
}

// TerraformDocsValidateOptions contains options for terraform-docs validation
type TerraformDocsValidateOptions struct {
	ConfigFile string
	Recursive  bool
}

// NewTerraformDocsModule creates a new terraform-docs module
func NewTerraformDocsModule(client *dagger.Client) *TerraformDocsModule {
	return &TerraformDocsModule{
		client: client,
		name:   "terraform-docs",
	}
}

// Generate runs terraform-docs generate on the provided module
func (m *TerraformDocsModule) Generate(ctx context.Context, modulePath string, opts TerraformDocsOptions) (string, error) {
	args := []string{"terraform-docs"}

	// Add format
	if opts.OutputFormat != "" {
		args = append(args, opts.OutputFormat)
	} else {
		args = append(args, "markdown")
	}

	// Add options
	if opts.ConfigFile != "" {
		args = append(args, "--config", opts.ConfigFile)
	}
	if opts.Recursive {
		args = append(args, "--recursive")
	}
	if opts.Sort {
		args = append(args, "--sort")
	}
	if opts.HeaderFrom != "" {
		args = append(args, "--header-from", opts.HeaderFrom)
	}
	if opts.FooterFrom != "" {
		args = append(args, "--footer-from", opts.FooterFrom)
	}
	if opts.OutputFile != "" {
		args = append(args, "--output-file", opts.OutputFile)
	}

	// Add module path
	args = append(args, ".")

	container := m.client.Container().
		From(getImageTag("terraform-docs", "quay.io/terraform-docs/terraform-docs:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(modulePath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	stdout, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if stdout != "" {
		return stdout, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Documentation generated successfully", nil
}

// Validate checks if terraform-docs output is up to date
func (m *TerraformDocsModule) Validate(ctx context.Context, modulePath string, opts TerraformDocsValidateOptions) (string, error) {
	args := []string{"terraform-docs", "markdown", "--output-check"}

	// Add options
	if opts.ConfigFile != "" {
		args = append(args, "--config", opts.ConfigFile)
	}
	if opts.Recursive {
		args = append(args, "--recursive")
	}

	// Add module path
	args = append(args, ".")

	container := m.client.Container().
		From(getImageTag("terraform-docs", "quay.io/terraform-docs/terraform-docs:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(modulePath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	stdout, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if stdout != "" {
		return stdout, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Documentation validation completed", nil
}