package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// ConftestModule runs Conftest for OPA policy testing
type ConftestModule struct {
	client *dagger.Client
	name   string
}

// NewConftestModule creates a new Conftest module
func NewConftestModule(client *dagger.Client) *ConftestModule {
	return &ConftestModule{
		client: client,
		name:   "conftest",
	}
}

// TestWithPolicy tests files against OPA policies
func (m *ConftestModule) TestWithPolicy(ctx context.Context, dir string, policyPath string) (string, error) {
	container := m.client.Container().
		From("openpolicyagent/conftest:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithDirectory("/policies", m.client.Host().Directory(policyPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"conftest",
			"test",
			"--policy", "/policies",
			"--output", "json",
			".",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run conftest: no output received")
}

// TestFile tests a specific file against policies
func (m *ConftestModule) TestFile(ctx context.Context, filePath string, policyPath string) (string, error) {
	container := m.client.Container().
		From("openpolicyagent/conftest:latest").
		WithFile("/workspace/target.yaml", m.client.Host().File(filePath)).
		WithDirectory("/policies", m.client.Host().Directory(policyPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"conftest",
			"test",
			"--policy", "/policies",
			"--output", "json",
			"target.yaml",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run conftest on file: no output received")
}

// VerifyPolicies runs policy unit tests
func (m *ConftestModule) VerifyPolicies(ctx context.Context, policyPath string) (string, error) {
	container := m.client.Container().
		From("openpolicyagent/conftest:latest").
		WithDirectory("/policies", m.client.Host().Directory(policyPath)).
		WithExec([]string{
			"conftest",
			"verify",
			"--policy", "/policies",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run conftest verify: no output received")
}

// ParseFile parses and prints structured data from input files
func (m *ConftestModule) ParseFile(ctx context.Context, filePath string, parser string) (string, error) {
	args := []string{"conftest", "parse", "/workspace/target.yaml"}
	if parser != "" {
		args = append(args, "--parser", parser)
	}

	container := m.client.Container().
		From("openpolicyagent/conftest:latest").
		WithFile("/workspace/target.yaml", m.client.Host().File(filePath)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	return output, nil
}

// PushPolicies pushes OPA policy bundles to OCI registry
func (m *ConftestModule) PushPolicies(ctx context.Context, registryURL string, policyPath string) (string, error) {
	args := []string{"conftest", "push", registryURL}
	if policyPath != "" {
		args = append(args, "/policies")
	}

	container := m.client.Container().
		From("openpolicyagent/conftest:latest")

	if policyPath != "" {
		container = container.WithDirectory("/policies", m.client.Host().Directory(policyPath))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to push policies: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Conftest
func (m *ConftestModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("openpolicyagent/conftest:latest").
		WithExec([]string{"conftest", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get conftest version: %w", err)
	}

	return output, nil
}
