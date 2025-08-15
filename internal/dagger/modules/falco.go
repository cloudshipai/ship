package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// FalcoModule runs Falco for runtime security monitoring
type FalcoModule struct {
	client *dagger.Client
	name   string
}

// NewFalcoModule creates a new Falco module
func NewFalcoModule(client *dagger.Client) *FalcoModule {
	return &FalcoModule{
		client: client,
		name:   "falco",
	}
}

// RunWithDefaultRules runs Falco with default rules
func (m *FalcoModule) RunWithDefaultRules(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"falco",
		"-o", "json_output=true",
		"-o", "log_stderr=false",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run falco: %w", err)
	}

	return output, nil
}

// RunWithCustomRules runs Falco with custom rules
func (m *FalcoModule) RunWithCustomRules(ctx context.Context, rulesPath string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest").
		WithDirectory("/etc/falco/rules.d", m.client.Host().Directory(rulesPath))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"falco",
		"-o", "json_output=true",
		"-o", "log_stderr=false",
		"-r", "/etc/falco/rules.d",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run falco with custom rules: %w", err)
	}

	return output, nil
}

// ValidateRules validates Falco rules syntax
func (m *FalcoModule) ValidateRules(ctx context.Context, rulesPath string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest").
		WithDirectory("/rules", m.client.Host().Directory(rulesPath)).
		WithExec([]string{
			"falco",
			"--validate", "/rules",
			"-o", "json_output=true",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate falco rules: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Falco
func (m *FalcoModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest").
		WithExec([]string{"falco", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get falco version: %w", err)
	}

	return output, nil
}
