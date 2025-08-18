package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// CustodianModule runs Cloud Custodian for cloud governance
type CustodianModule struct {
	client *dagger.Client
	name   string
}

const custodianBinary = "/src/.venv/bin/custodian"

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
			custodianBinary, "run",
			"-s", "/output",
			"/policy.yml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run custodian policy: %w", err)
	}

	return output, nil
}

// ValidatePolicy validates a custodian policy
func (m *CustodianModule) ValidatePolicy(ctx context.Context, policyPath string) (string, error) {
	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithFile("/policy.yml", m.client.Host().File(policyPath)).
		WithExec([]string{
			custodianBinary, "validate",
			"/policy.yml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate custodian policy: %w", err)
	}

	return output, nil
}

// DryRun performs a dry run of a policy
func (m *CustodianModule) DryRun(ctx context.Context, policyPath string) (string, error) {
	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithFile("/policy.yml", m.client.Host().File(policyPath)).
		WithExec([]string{
			custodianBinary, "run",
			"--dryrun",
			"-s", "/output",
			"/policy.yml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run custodian dry run: %w", err)
	}

	return output, nil
}

// Schema shows schema for a particular resource type
func (m *CustodianModule) Schema(ctx context.Context, resourceType string) (string, error) {
	args := []string{custodianBinary, "schema"}
	if resourceType != "" {
		args = append(args, resourceType)
	}

	container := m.client.Container().
		From("cloudcustodian/c7n:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get custodian schema: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Cloud Custodian
func (m *CustodianModule) GetVersion(ctx context.Context) (string, error) {
	// The container has custodian as the entrypoint, so we need to use WithEntrypoint
	container := m.client.Container().
		From("cloudcustodian/c7n:latest")
	
	// Execute the version command
	container = container.WithExec([]string{custodianBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get custodian version: %w", err)
	}

	return output, nil
}
