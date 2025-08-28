package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// kuttlBinary is the path to the kuttl binary in the container
const kuttlBinary = "kubectl-kuttl"

// KuttlModule runs KUTTL for Kubernetes testing
type KuttlModule struct {
	client *dagger.Client
	name   string
}

// NewKuttlModule creates a new KUTTL module
func NewKuttlModule(client *dagger.Client) *KuttlModule {
	return &KuttlModule{
		client: client,
		name:   "kuttl",
	}
}

// RunTest runs KUTTL tests
func (m *KuttlModule) RunTest(ctx context.Context, testPath string, kubeconfig string, parallel int, skipDelete bool) (string, error) {
	container := m.client.Container().
		From("kudobuilder/kuttl:latest").
		WithDirectory("/tests", m.client.Host().Directory(testPath))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{kuttlBinary, "test", "/tests"}
	
	// Add parallel execution if specified
	if parallel > 0 {
		args = append(args, "--parallel", fmt.Sprintf("%d", parallel))
	}
	
	// Add skip delete if specified
	if skipDelete {
		args = append(args, "--skip-delete")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to run kuttl tests: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// RunTestWithKind runs KUTTL tests with kind cluster
func (m *KuttlModule) RunTestWithKind(ctx context.Context, testPath string, kindConfig string, parallel int) (string, error) {
	container := m.client.Container().
		From("kudobuilder/kuttl:latest").
		WithDirectory("/tests", m.client.Host().Directory(testPath))

	args := []string{kuttlBinary, "test", "/tests", "--start-kind"}
	
	// Add kind config if specified
	if kindConfig != "" {
		container = container.WithFile("/kind-config.yaml", m.client.Host().File(kindConfig))
		args = append(args, "--kind-config", "/kind-config.yaml")
	}
	
	// Add parallel execution if specified
	if parallel > 0 {
		args = append(args, "--parallel", fmt.Sprintf("%d", parallel))
	}

	// Mount Docker socket for kind
	container = container.
		WithUnixSocket("/var/run/docker.sock", m.client.Host().UnixSocket("/var/run/docker.sock")).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to run kuttl tests with kind: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// ValidateTest validates test configuration
func (m *KuttlModule) ValidateTest(ctx context.Context, testPath string) (string, error) {
	container := m.client.Container().
		From("kudobuilder/kuttl:latest").
		WithDirectory("/tests", m.client.Host().Directory(testPath)).
		WithExec([]string{
			kuttlBinary,
			"test",
			"--dry-run",
			"/tests",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to validate kuttl tests: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// GetVersion returns the version of KUTTL
func (m *KuttlModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("kudobuilder/kuttl:latest").
		WithExec([]string{kuttlBinary, "version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get kuttl version: %w", err)
	}

	return output, nil
}

// GetHelp returns the help information for KUTTL
func (m *KuttlModule) GetHelp(ctx context.Context, command string) (string, error) {
	container := m.client.Container().
		From("kudobuilder/kuttl:latest")

	args := []string{kuttlBinary, "help"}
	if command != "" {
		args = append(args, command)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to get kuttl help: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}