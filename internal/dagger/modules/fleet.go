package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// FleetModule runs Fleet for GitOps with Kubernetes
type FleetModule struct {
	client *dagger.Client
	name   string
}

const kubectlBinary = "kubectl"

// NewFleetModule creates a new Fleet module
func NewFleetModule(client *dagger.Client) *FleetModule {
	return &FleetModule{
		client: client,
		name:   "fleet",
	}
}

// GetClusters lists Fleet clusters
func (m *FleetModule) GetClusters(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("rancher/fleet:v0.10.4")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		kubectlBinary, "get", "clusters",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get fleet clusters: %w", err)
	}

	return output, nil
}

// GetGitRepos lists Git repositories managed by Fleet
func (m *FleetModule) GetGitRepos(ctx context.Context, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("rancher/fleet:v0.10.4")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		kubectlBinary, "get", "gitrepos",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get git repos: %w", err)
	}

	return output, nil
}

// CreateGitRepo creates a new Git repository resource
func (m *FleetModule) CreateGitRepo(ctx context.Context, name string, repoURL string, branch string, path string, kubeconfig string) (string, error) {
	gitRepoYAML := fmt.Sprintf(`
apiVersion: fleet.cattle.io/v1alpha1
kind: GitRepo
metadata:
  name: %s
spec:
  repo: %s
  branch: %s
  paths:
  - %s
`, name, repoURL, branch, path)

	container := m.client.Container().
		From("rancher/fleet:v0.10.4").
		WithNewFile("/gitrepo.yaml", gitRepoYAML)

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		kubectlBinary, "apply", "-f", "/gitrepo.yaml",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create git repo: %w", err)
	}

	return output, nil
}

// GetBundles lists Fleet bundles
func (m *FleetModule) GetBundles(ctx context.Context, kubeconfig string, namespace string) (string, error) {
	args := []string{"kubectl", "get", "bundles"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "--all-namespaces")
	}
	args = append(args, "--output", "json")

	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get bundles: %w", err)
	}

	return output, nil
}

// GetBundleDeployments lists Fleet bundle deployments
func (m *FleetModule) GetBundleDeployments(ctx context.Context, kubeconfig string, namespace string) (string, error) {
	args := []string{"kubectl", "get", "bundledeployments"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "--all-namespaces")
	}
	args = append(args, "--output", "json")

	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get bundle deployments: %w", err)
	}

	return output, nil
}

// DescribeGitRepo describes a Fleet GitRepo
func (m *FleetModule) DescribeGitRepo(ctx context.Context, name string, namespace string, kubeconfig string) (string, error) {
	args := []string{"kubectl", "describe", "gitrepo", name}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to describe gitrepo: %w", err)
	}

	return output, nil
}

// Install installs Fleet using Helm
func (m *FleetModule) Install(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	if namespace == "" {
		namespace = "cattle-fleet-system"
	}

	container := m.client.Container().
		From("alpine/helm:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"helm", "repo", "add", "fleet", "https://rancher.github.io/fleet-helm-charts/",
	}).WithExec([]string{
		"helm", "install", "fleet", "fleet/fleet",
		"--namespace", namespace,
		"--create-namespace",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install fleet: %w", err)
	}

	return output, nil
}

// ApplyGitRepo applies Fleet GitRepo configuration using kubectl (MCP compatible)
func (m *FleetModule) ApplyGitRepo(ctx context.Context, gitrepoFile string, namespace string) (string, error) {
	gitrepoFileObj := m.client.Host().File(gitrepoFile)
	
	args := []string{"kubectl", "apply", "-f", "/gitrepo.yaml"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithFile("/gitrepo.yaml", gitrepoFileObj).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to apply gitrepo: %w", err)
	}

	return output, nil
}

// GetGitReposSimple lists Git repositories managed by Fleet (MCP compatible)
func (m *FleetModule) GetGitReposSimple(ctx context.Context, namespace string) (string, error) {
	if namespace == "" {
		namespace = "fleet-local"
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithExec([]string{"kubectl", "get", "gitrepo", "-n", namespace})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gitrepos: %w", err)
	}

	return output, nil
}

// GetBundlesSimple lists Fleet bundles (MCP compatible)
func (m *FleetModule) GetBundlesSimple(ctx context.Context, namespace string) (string, error) {
	args := []string{"kubectl", "get", "bundles"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "--all-namespaces")
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get bundles: %w", err)
	}

	return output, nil
}

// GetBundleDeploymentsSimple lists Fleet bundle deployments (MCP compatible)
func (m *FleetModule) GetBundleDeploymentsSimple(ctx context.Context, namespace string) (string, error) {
	args := []string{"kubectl", "get", "bundledeployments"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "--all-namespaces")
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get bundle deployments: %w", err)
	}

	return output, nil
}

// DescribeGitRepoSimple describes a Fleet GitRepo (MCP compatible)
func (m *FleetModule) DescribeGitRepoSimple(ctx context.Context, gitrepoName string, namespace string) (string, error) {
	if namespace == "" {
		namespace = "fleet-local"
	}

	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithExec([]string{"kubectl", "describe", "gitrepo", gitrepoName, "-n", namespace})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to describe gitrepo: %w", err)
	}

	return output, nil
}

// InstallWithVersion installs Fleet using Helm with specific version (MCP compatible)
func (m *FleetModule) InstallWithVersion(ctx context.Context, version string) (string, error) {
	if version == "" {
		version = "v0.13.0"
	}

	// Install Fleet CRD first
	container := m.client.Container().
		From("alpine/helm:latest").
		WithExec([]string{
			"helm", "-n", "cattle-fleet-system", "install", "--create-namespace", "--wait",
			"fleet-crd", "https://github.com/rancher/fleet/releases/download/" + version + "/fleet-crd-" + version[1:] + ".tgz",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to install fleet with version: %w", err)
	}

	return output, nil
}
