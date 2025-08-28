package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// litmusctlBinary is the path to the litmusctl binary in the container
const litmusctlBinary = "/usr/local/bin/litmusctl"

// LitmusModule runs Litmus for chaos engineering
type LitmusModule struct {
	client *dagger.Client
	name   string
}

// NewLitmusModule creates a new Litmus module
func NewLitmusModule(client *dagger.Client) *LitmusModule {
	return &LitmusModule{
		client: client,
		name:   "litmus",
	}
}

// CreateExperiment creates a chaos experiment
func (m *LitmusModule) CreateExperiment(ctx context.Context, experimentPath string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest").
		WithFile("/experiment.yaml", m.client.Host().File(experimentPath))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		litmusctlBinary,
		"create",
		"experiment",
		"-f", "/experiment.yaml",
		"--output", "json",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create litmus experiment: %w", err)
	}

	return output, nil
}

// GetExperiments lists chaos experiments
func (m *LitmusModule) GetExperiments(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		litmusctlBinary,
		"get",
		"experiments",
		"--output", "json",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get litmus experiments: %w", err)
	}

	return output, nil
}

// GetChaosResults gets chaos experiment results
func (m *LitmusModule) GetChaosResults(ctx context.Context, experimentName string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		litmusctlBinary,
		"get",
		"chaosresults",
		experimentName,
		"--output", "json",
	}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chaos results: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Litmus
func (m *LitmusModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest").
		WithExec([]string{litmusctlBinary, "version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to get litmus version: no output received")
}

// Install installs Litmus using Helm
func (m *LitmusModule) Install(ctx context.Context, namespace string, releaseName string, createNamespace bool) (string, error) {
	if namespace == "" {
		namespace = "litmus"
	}
	if releaseName == "" {
		releaseName = "chaos"
	}

	container := m.client.Container().
		From("alpine/helm:latest")

	// Add Litmus Helm repository
	container = container.WithExec([]string{"helm", "repo", "add", "litmuschaos", "https://litmuschaos.github.io/litmus-helm/"}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	}).
		WithExec([]string{"helm", "repo", "update"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	// Install Litmus
	args := []string{"helm", "install", releaseName, "litmuschaos/litmus", "--namespace", namespace}
	if createNamespace {
		args = append(args, "--create-namespace")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install litmus: %w", err)
	}

	return output, nil
}

// ConnectChaosInfra connects chaos infrastructure using litmusctl
func (m *LitmusModule) ConnectChaosInfra(ctx context.Context, projectID string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest")

	args := []string{litmusctlBinary, "connect", "chaos-infra"}
	if projectID != "" {
		args = append(args, "--project-id", projectID)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to connect chaos infrastructure: %w", err)
	}

	return output, nil
}

// CreateProject creates a new project using litmusctl
func (m *LitmusModule) CreateProject(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest").
		WithExec([]string{litmusctlBinary, "create", "project"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create litmus project: %w", err)
	}

	return output, nil
}

// CreateChaosExperiment creates chaos experiment using litmusctl
func (m *LitmusModule) CreateChaosExperiment(ctx context.Context, manifestFile string, projectID string, chaosInfraID string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest").
		WithFile("/manifest.yaml", m.client.Host().File(manifestFile))

	args := []string{litmusctlBinary, "create", "chaos-experiment", "-f", "/manifest.yaml"}
	if projectID != "" {
		args = append(args, "--project-id", projectID)
	}
	if chaosInfraID != "" {
		args = append(args, "--chaos-infra-id", chaosInfraID)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create chaos experiment: %w", err)
	}

	return output, nil
}

// RunChaosExperiment runs chaos experiment using litmusctl
func (m *LitmusModule) RunChaosExperiment(ctx context.Context, experimentID string, projectID string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest")

	args := []string{litmusctlBinary, "run", "chaos-experiment", experimentID}
	if projectID != "" {
		args = append(args, "--project-id", projectID)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run chaos experiment: %w", err)
	}

	return output, nil
}

// GetProjects lists projects using litmusctl
func (m *LitmusModule) GetProjects(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest").
		WithExec([]string{litmusctlBinary, "get", "projects"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to get litmus projects: no output received")
}

// GetChaosInfra lists chaos infrastructure using litmusctl
func (m *LitmusModule) GetChaosInfra(ctx context.Context, projectID string) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest")

	args := []string{litmusctlBinary, "get", "chaos-infra"}
	if projectID != "" {
		args = append(args, "--project-id", projectID)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chaos infrastructure: %w", err)
	}

	return output, nil
}

// ConfigSetAccount setup ChaosCenter account configuration using litmusctl
func (m *LitmusModule) ConfigSetAccount(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("litmuschaos/litmusctl:latest").
		WithExec([]string{litmusctlBinary, "config", "set-account"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to configure litmus account: %w", err)
	}

	return output, nil
}

// ApplyChaosExperiment applies chaos experiment manifest using kubectl
func (m *LitmusModule) ApplyChaosExperiment(ctx context.Context, manifestFile string, namespace string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithFile("/manifest.yaml", m.client.Host().File(manifestFile))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"kubectl", "apply", "-f", "/manifest.yaml"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to apply chaos experiment: %w", err)
	}

	return output, nil
}
