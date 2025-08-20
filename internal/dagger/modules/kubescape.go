package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)


type KubescapeModule struct {
	client *dagger.Client
}

func NewKubescapeModule(client *dagger.Client) *KubescapeModule {
	return &KubescapeModule{
		client: client,
	}
}

// ScanCluster scans a Kubernetes cluster for security issues
func (m *KubescapeModule) ScanCluster(ctx context.Context, kubeconfig string, framework string, format string, severityThreshold string) (string, error) {
	container := m.client.Container().
		From("quay.io/kubescape/kubescape-cli:v3.0.15").
		WithWorkdir("/workspace")

	// Mount kubeconfig if provided
	if kubeconfig != "" {
		container = container.WithMountedFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{"kubescape", "scan"}

	// Add framework (default to nsa if not specified)
	if framework == "" {
		framework = "nsa"
	}
	args = append(args, "framework", framework)

	// Add format
	if format != "" {
		args = append(args, "--format", format)
	}

	// Add severity threshold
	if severityThreshold != "" {
		args = append(args, "--severity-threshold", severityThreshold)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	// Kubescape might output to stderr
	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}
	
	return "", fmt.Errorf("failed to scan cluster: no output received")
}

// ScanManifests scans Kubernetes manifest files
func (m *KubescapeModule) ScanManifests(ctx context.Context, manifestsDir string, framework string, format string, severityThreshold string) (string, error) {
	container := m.client.Container().
		From("quay.io/kubescape/kubescape-cli:v3.0.15").
		WithWorkdir("/workspace")

	// Mount manifests directory
	if manifestsDir != "" {
		container = container.WithMountedDirectory("/workspace/manifests", m.client.Host().Directory(manifestsDir))
	}

	// Add framework (default to nsa if not specified)
	if framework == "" {
		framework = "nsa"
	}

	args := []string{"kubescape", "scan", "framework", framework}

	// Add format
	if format != "" {
		args = append(args, "--format", format)
	}

	// Add severity threshold
	if severityThreshold != "" {
		args = append(args, "--severity-threshold", severityThreshold)
	}

	// Add manifests directory
	args = append(args, "/workspace/manifests")

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	// Kubescape might output to stderr
	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}
	
	return "", fmt.Errorf("failed to scan manifests: no output received")
}

// ScanHelm scans Helm charts for security issues
func (m *KubescapeModule) ScanHelm(ctx context.Context, chartPath string, framework string, format string) (string, error) {
	container := m.client.Container().
		From("quay.io/kubescape/kubescape-cli:v3.0.15").
		WithWorkdir("/workspace")

	// Mount chart directory
	if chartPath != "" {
		container = container.WithMountedDirectory("/workspace/chart", m.client.Host().Directory(chartPath))
	}

	// Add framework (default to nsa if not specified)
	if framework == "" {
		framework = "nsa"
	}

	args := []string{"kubescape", "scan", "framework", framework}

	// Add format
	if format != "" {
		args = append(args, "--format", format)
	}

	// Add chart path with helm prefix
	args = append(args, "helm", "/workspace/chart")

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to scan Helm chart: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// ScanRepository scans a Git repository for Kubernetes manifests
func (m *KubescapeModule) ScanRepository(ctx context.Context, repoPath string, framework string, format string) (string, error) {
	container := m.client.Container().
		From("quay.io/kubescape/kubescape-cli:v3.0.15").
		WithWorkdir("/workspace")

	// Mount repository directory
	if repoPath != "" {
		container = container.WithMountedDirectory("/workspace/repo", m.client.Host().Directory(repoPath))
	}

	// Add framework (default to nsa if not specified)
	if framework == "" {
		framework = "nsa"
	}

	args := []string{"kubescape", "scan", "framework", framework}

	// Add format
	if format != "" {
		args = append(args, "--format", format)
	}

	// Add repository path
	args = append(args, "/workspace/repo")

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to scan repository: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// GetVersion returns the version of kubescape
func (m *KubescapeModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("quay.io/kubescape/kubescape-cli:v3.0.15").
		WithExec([]string{"kubescape", "version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	// Kubescape might output to stderr
	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}
	
	return "", fmt.Errorf("failed to get kubescape version: no output received")
}

// ListFrameworks lists all available security frameworks
func (m *KubescapeModule) ListFrameworks(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("quay.io/kubescape/kubescape-cli:v3.0.15").
		WithExec([]string{"kubescape", "list", "frameworks"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	// Kubescape might output to stderr
	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}
	
	return "", fmt.Errorf("failed to list frameworks: no output received")
}

// ListControls lists all available security controls
func (m *KubescapeModule) ListControls(ctx context.Context, framework string) (string, error) {
	container := m.client.Container().
		From("quay.io/kubescape/kubescape-cli:v3.0.15")

	args := []string{"kubescape", "list", "controls"}
	
	// Add framework filter if specified
	if framework != "" {
		args = append(args, "--framework", framework)
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to list controls: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// DownloadArtifacts downloads kubescape artifacts for offline use
func (m *KubescapeModule) DownloadArtifacts(ctx context.Context, outputDir string) (string, error) {
	container := m.client.Container().
		From("quay.io/kubescape/kubescape-cli:v3.0.15").
		WithWorkdir("/workspace")

	// Mount output directory
	if outputDir != "" {
		container = container.WithMountedDirectory("/workspace/output", m.client.Host().Directory(outputDir))
	}

	args := []string{"kubescape", "download", "artifacts", "--output", "/workspace/output"}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to download artifacts: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// KubescapeConfig holds configuration options - no longer needed for simplified functions
type KubescapeConfig struct {
	KubescapeVersion    string
	KubeconfigPath      string
	Framework           string
	Format              string
	Output              string
	SeverityThreshold   string
	ComplianceThreshold float64
	Namespace           string
	IncludeResources    []string
	IncludeKubeSystem   bool
	Verbose             bool
}

// KubescapeOption is a function that modifies the KubescapeConfig - no longer needed for simplified functions
type KubescapeOption func(*KubescapeConfig)