package modules

import (
	"context"

	"dagger.io/dagger"
)

// TFLintModule runs TFLint for Terraform configuration linting
type TFLintModule struct {
	client *dagger.Client
	name   string
}

// TFLintOptions contains options for TFLint checking
type TFLintOptions struct {
	ConfigFile  string
	Format      string
	Recursive   bool
	EnableRule  string
	DisableRule string
	Only        string
	VarFile     string
	Var         string
	Fix         bool
}

// TFLintInitOptions contains options for TFLint initialization
type TFLintInitOptions struct {
	ConfigFile string
	Upgrade    bool
}

// NewTFLintModule creates a new TFLint module
func NewTFLintModule(client *dagger.Client) *TFLintModule {
	return &TFLintModule{
		client: client,
		name:   "tflint",
	}
}

// Check runs TFLint check on the provided Terraform configuration
func (m *TFLintModule) Check(ctx context.Context, sourcePath string, opts TFLintOptions) (string, error) {
	args := []string{"tflint"}

	// Add check options
	if opts.ConfigFile != "" {
		args = append(args, "--config", opts.ConfigFile)
	}
	if opts.Format != "" {
		args = append(args, "--format", opts.Format)
	}
	if opts.Recursive {
		args = append(args, "--recursive")
	}
	if opts.EnableRule != "" {
		args = append(args, "--enable-rule", opts.EnableRule)
	}
	if opts.DisableRule != "" {
		args = append(args, "--disable-rule", opts.DisableRule)
	}
	if opts.Only != "" {
		args = append(args, "--only", opts.Only)
	}
	if opts.VarFile != "" {
		args = append(args, "--var-file", opts.VarFile)
	}
	if opts.Var != "" {
		args = append(args, "--var", opts.Var)
	}
	if opts.Fix {
		args = append(args, "--fix")
	}

	container := m.client.Container().
		From(getImageTag("tflint", "ghcr.io/terraform-linters/tflint:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(sourcePath)).
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

	return "TFLint check completed successfully", nil
}

// Init initializes TFLint in a Terraform configuration directory
func (m *TFLintModule) Init(ctx context.Context, sourcePath string, opts TFLintInitOptions) (string, error) {
	args := []string{"tflint", "--init"}

	// Add init options
	if opts.ConfigFile != "" {
		args = append(args, "--config", opts.ConfigFile)
	}
	if opts.Upgrade {
		args = append(args, "--upgrade")
	}

	container := m.client.Container().
		From(getImageTag("tflint", "ghcr.io/terraform-linters/tflint:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(sourcePath)).
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

	return "TFLint initialization completed successfully", nil
}