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

const falcoBinary = "falco"

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
		falcoBinary,
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
		falcoBinary,
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
			falcoBinary,
			"--validate", "/rules",
			"-o", "json_output=true",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate falco rules: %w", err)
	}

	return output, nil
}

// DryRun performs dry run without monitoring
func (m *FalcoModule) DryRun(ctx context.Context, rulesPath string, configPath string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest")

	args := []string{"falco", "--dry-run"}
	if configPath != "" {
		container = container.WithFile("/etc/falco/falco.yaml", m.client.Host().File(configPath))
		args = append(args, "--config", "/etc/falco/falco.yaml")
	}
	if rulesPath != "" {
		container = container.WithFile("/etc/falco/custom_rules.yaml", m.client.Host().File(rulesPath))
		args = append(args, "--rules", "/etc/falco/custom_rules.yaml")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run falco dry run: %w", err)
	}

	return output, nil
}

// ListFields lists available fields for Falco rules
func (m *FalcoModule) ListFields(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest").
		WithExec([]string{"falco", "--list"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list falco fields: %w", err)
	}

	return output, nil
}

// ListRules lists all loaded Falco rules
func (m *FalcoModule) ListRules(ctx context.Context, rulesPath string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest")

	args := []string{"falco", "--list-rules"}
	if rulesPath != "" {
		container = container.WithFile("/etc/falco/custom_rules.yaml", m.client.Host().File(rulesPath))
		args = append(args, "--rules", "/etc/falco/custom_rules.yaml")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "No rules found or error listing rules", nil
}

// DescribeRule describes a specific Falco rule
func (m *FalcoModule) DescribeRule(ctx context.Context, ruleName string, rulesPath string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest")

	args := []string{"falco", "--describe-rule", ruleName}
	if rulesPath != "" {
		container = container.WithFile("/etc/falco/custom_rules.yaml", m.client.Host().File(rulesPath))
		args = append(args, "--rules", "/etc/falco/custom_rules.yaml")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to describe falco rule: %w", err)
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

// StartMonitoring starts Falco runtime security monitoring (MCP compatible)
func (m *FalcoModule) StartMonitoring(ctx context.Context, configPath string, rulesPath string, outputFormat string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest")

	args := []string{"falco"}
	
	if configPath != "" {
		container = container.WithFile("/etc/falco/custom_falco.yaml", m.client.Host().File(configPath))
		args = append(args, "-c", "/etc/falco/custom_falco.yaml")
	}
	if rulesPath != "" {
		container = container.WithFile("/etc/falco/custom_rules.yaml", m.client.Host().File(rulesPath))
		args = append(args, "-r", "/etc/falco/custom_rules.yaml")
	}
	if outputFormat == "json" {
		args = append(args, "-o", "json_output=true")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start falco monitoring: %w", err)
	}

	return output, nil
}

// ValidateRulesSimple validates Falco rules syntax (MCP compatible)
func (m *FalcoModule) ValidateRulesSimple(ctx context.Context, rulesPath string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest").
		WithFile("/etc/falco/rules_to_validate.yaml", m.client.Host().File(rulesPath)).
		WithExec([]string{"falco", "-V", "/etc/falco/rules_to_validate.yaml"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate falco rules: %w", err)
	}

	return output, nil
}

// DryRunSimple performs dry run without monitoring (MCP compatible)
func (m *FalcoModule) DryRunSimple(ctx context.Context, configPath string, rulesPath string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest")

	args := []string{"falco", "--dry-run"}
	
	if configPath != "" {
		container = container.WithFile("/etc/falco/custom_falco.yaml", m.client.Host().File(configPath))
		args = append(args, "-c", "/etc/falco/custom_falco.yaml")
	}
	if rulesPath != "" {
		container = container.WithFile("/etc/falco/custom_rules.yaml", m.client.Host().File(rulesPath))
		args = append(args, "-r", "/etc/falco/custom_rules.yaml")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run falco dry run: %w", err)
	}

	return output, nil
}

// ListFieldsWithSource lists supported fields for Falco rules (MCP compatible)
func (m *FalcoModule) ListFieldsWithSource(ctx context.Context, source string) (string, error) {
	args := []string{"falco", "--list"}
	
	if source != "" {
		args = append(args, source)
	}

	container := m.client.Container().
		From("falcosecurity/falco:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list falco fields: %w", err)
	}

	return output, nil
}

// GetVersionSimple returns the version of Falco (MCP compatible)
func (m *FalcoModule) GetVersionSimple(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest").
		WithExec([]string{"falco", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get falco version: %w", err)
	}

	return output, nil
}

// ListRulesSimple lists all loaded Falco rules (MCP compatible)
func (m *FalcoModule) ListRulesSimple(ctx context.Context, rulesPath string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest")

	args := []string{"falco", "-L"}
	if rulesPath != "" {
		container = container.WithFile("/etc/falco/custom_rules.yaml", m.client.Host().File(rulesPath))
		args = append(args, "-r", "/etc/falco/custom_rules.yaml")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list falco rules: %w", err)
	}

	return output, nil
}

// DescribeRuleSimple describes a specific Falco rule (MCP compatible)
func (m *FalcoModule) DescribeRuleSimple(ctx context.Context, ruleName string, rulesPath string) (string, error) {
	container := m.client.Container().
		From("falcosecurity/falco:latest")

	args := []string{"falco", "-l", ruleName}
	if rulesPath != "" {
		container = container.WithFile("/etc/falco/custom_rules.yaml", m.client.Host().File(rulesPath))
		args = append(args, "-r", "/etc/falco/custom_rules.yaml")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to describe falco rule: %w", err)
	}

	return output, nil
}
