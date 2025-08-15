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
			"custodian", "run",
			"--output-dir", "/output",
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
			"custodian", "validate",
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
			"custodian", "run",
			"--dryrun",
			"/policy.yml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run custodian dry run: %w", err)
	}

	return output, nil
}
