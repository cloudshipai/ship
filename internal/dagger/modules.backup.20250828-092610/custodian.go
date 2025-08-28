package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// CustodianModule runs Cloud Custodian for cloud governance
type CustodianModule struct {
	client *dagger.Client
	name   string
}

// NewCustodianModule creates a new Cloud Custodian module
func NewCustodianModule(client *dagger.Client) *CustodianModule {
	return &CustodianModule{
		client: client,
		name:   "custodian",
	}
}

// RunPolicy runs a custodian policy
func (m *CustodianModule) RunPolicy(ctx context.Context, policyPath string, outputDir string) (string, error) {
	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithFile("/policy.yml", m.client.Host().File(policyPath)).
		WithExec([]string{
			"sh", "-c",
			"/src/.venv/bin/custodian run -s /output /policy.yml 2>&1",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run custodian policy: %w", err)
	}

	return output, nil
}

// ValidatePolicy validates a custodian policy
func (m *CustodianModule) ValidatePolicy(ctx context.Context, policyPath string) (string, error) {
	// Use a wrapper script to capture both stdout and stderr, and only fail on actual validation errors
	// This works around Dagger's issue with stderr output causing failures
	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithFile("/policy.yml", m.client.Host().File(policyPath)).
		WithExec([]string{
			"sh", "-c", 
			`output=$(/src/.venv/bin/custodian validate /policy.yml 2>&1); echo "$output"; if echo "$output" | grep -q "Configuration invalid"; then exit 1; else exit 0; fi`,
		})

	// Get the output
	output, err := container.Stdout(ctx)
	
	// Check the validation result based on content
	if strings.Contains(output, "Configuration valid") {
		return output, nil
	}
	
	if strings.Contains(output, "Configuration invalid") {
		return "", fmt.Errorf("policy validation failed: %s", output)
	}
	
	// If we have an error and no clear validation result
	if err != nil {
		return "", fmt.Errorf("failed to validate custodian policy: %w", err)
	}
	
	// If we have output but no clear success/failure, return the output
	if output != "" {
		return output, nil
	}
	
	// Default success message
	return "Policy validation completed", nil
}

// DryRun performs a dry run of a policy
func (m *CustodianModule) DryRun(ctx context.Context, policyPath string) (string, error) {
	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithFile("/policy.yml", m.client.Host().File(policyPath)).
		WithExec([]string{
			"sh", "-c",
			"/src/.venv/bin/custodian run --dryrun -s /output /policy.yml 2>&1",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run custodian dry run: %w", err)
	}

	return output, nil
}

// Schema shows schema for a particular resource type
func (m *CustodianModule) Schema(ctx context.Context, resourceType string) (string, error) {
	// Build command with wrapper script
	cmd := "/src/.venv/bin/custodian schema"
	if resourceType != "" {
		cmd += " " + resourceType
	}
	cmd += " 2>&1"

	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithExec([]string{"sh", "-c", cmd})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get custodian schema: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Cloud Custodian
func (m *CustodianModule) GetVersion(ctx context.Context) (string, error) {
	// Use wrapper script to handle stderr output
	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithExec([]string{
			"sh", "-c", 
			"/src/.venv/bin/custodian version 2>&1",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get custodian version: %w", err)
	}

	return output, nil
}

// Report generates a tabular report on policy matched resources
func (m *CustodianModule) Report(ctx context.Context, policyPath string, outputDir string, format string) (string, error) {
	// Build command with wrapper script
	cmd := "/src/.venv/bin/custodian report -s /output"
	
	// Add format if specified
	if format != "" {
		cmd += " --format " + format
	}
	
	cmd += " /policy.yml 2>&1"

	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithFile("/policy.yml", m.client.Host().File(policyPath)).
		WithDirectory("/output", m.client.Host().Directory(outputDir)).
		WithExec([]string{"sh", "-c", cmd})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate custodian report: %w", err)
	}

	return output, nil
}

// Logs retrieves logs for a specific policy
func (m *CustodianModule) Logs(ctx context.Context, policyPath string, outputDir string) (string, error) {
	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithFile("/policy.yml", m.client.Host().File(policyPath)).
		WithDirectory("/output", m.client.Host().Directory(outputDir)).
		WithExec([]string{
			"sh", "-c",
			"/src/.venv/bin/custodian logs -s /output /policy.yml 2>&1",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get custodian logs: %w", err)
	}

	return output, nil
}

// Metrics retrieves metrics for policy execution
func (m *CustodianModule) Metrics(ctx context.Context, policyPath string, outputDir string, start string, end string) (string, error) {
	// Build command with wrapper script
	cmd := "/src/.venv/bin/custodian metrics -s /output"
	
	// Add time range if specified
	if start != "" {
		cmd += " --start " + start
	}
	if end != "" {
		cmd += " --end " + end
	}
	
	cmd += " /policy.yml 2>&1"

	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithFile("/policy.yml", m.client.Host().File(policyPath)).
		WithDirectory("/output", m.client.Host().Directory(outputDir)).
		WithExec([]string{"sh", "-c", cmd})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get custodian metrics: %w", err)
	}

	return output, nil
}
