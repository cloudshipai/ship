package modules

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger.io/dagger"
)

// TFLintModule runs TFLint for Terraform linting
type TFLintModule struct {
	client *dagger.Client
	name   string
}

const tflintBinary = "/usr/local/bin/tflint"

// NewTFLintModule creates a new TFLint module
func NewTFLintModule(client *dagger.Client) *TFLintModule {
	return &TFLintModule{
		client: client,
		name:   tflintBinary,
	}
}

// LintDirectory lints all Terraform files in a directory
func (m *TFLintModule) LintDirectory(ctx context.Context, dir string) (string, error) {
	// Initialize TFLint first
	initContainer := m.client.Container().
		From("ghcr.io/terraform-linters/tflint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{tflintBinary, "--init"})

	// Sync to ensure init completes
	_, err := initContainer.Sync(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to initialize tflint: %w", err)
	}

	// Run TFLint - use a bash wrapper to capture output regardless of exit code
	lintContainer := initContainer.WithExec([]string{
		"sh", "-c", "tflint --format json || true",
	})

	output, err := lintContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run tflint: %w", err)
	}

	return output, nil
}

// LintFile lints a specific Terraform file
func (m *TFLintModule) LintFile(ctx context.Context, filePath string) (string, error) {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)

	container := m.client.Container().
		From("ghcr.io/terraform-linters/tflint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			tflintBinary,
			"--format", "json",
			filename,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run tflint on file: %w", err)
	}

	return output, nil
}

// LintWithConfig lints using a custom configuration file
func (m *TFLintModule) LintWithConfig(ctx context.Context, dir string, configFile string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/terraform-linters/tflint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir))

	// If config file is provided, mount it
	if configFile != "" {
		container = container.WithFile("/.tflint.hcl", m.client.Host().File(configFile))
		container = container.WithExec([]string{
			tflintBinary,
			"--format", "json",
			"--config", "/.tflint.hcl",
			"/workspace",
		})
	} else {
		container = container.WithExec([]string{
			tflintBinary,
			"--format", "json",
			"/workspace",
		})
	}

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run tflint with config: %w", err)
	}

	return output, nil
}

// LintWithRules runs TFLint with specific rule sets enabled
func (m *TFLintModule) LintWithRules(ctx context.Context, dir string, enableRules []string, disableRules []string) (string, error) {
	args := []string{tflintBinary, "--format", "json"}

	// Add enabled rules
	for _, rule := range enableRules {
		args = append(args, "--enable-rule", rule)
	}

	// Add disabled rules
	for _, rule := range disableRules {
		args = append(args, "--disable-rule", rule)
	}

	container := m.client.Container().
		From("ghcr.io/terraform-linters/tflint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run tflint with rules: %w", err)
	}

	return output, nil
}

// InitPlugins initializes TFLint plugins
func (m *TFLintModule) InitPlugins(ctx context.Context, dir string) error {
	container := m.client.Container().
		From("ghcr.io/terraform-linters/tflint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			tflintBinary,
			"--init",
		})

	_, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize tflint plugins: %w", err)
	}

	return nil
}

// GetVersion returns the version of TFLint
func (m *TFLintModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/terraform-linters/tflint:latest").
		WithExec([]string{tflintBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get tflint version: %w", err)
	}

	return output, nil
}
