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
