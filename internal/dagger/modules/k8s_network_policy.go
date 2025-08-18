package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// K8sNetworkPolicyModule runs Kubernetes network policy tools
type K8sNetworkPolicyModule struct {
	client *dagger.Client
	name   string
}

const (
	k8sKubectlBinary = "/opt/bitnami/kubectl/bin/kubectl"
	netfetchBinary   = "/usr/local/bin/netfetch"
)

// NewK8sNetworkPolicyModule creates a new Kubernetes network policy module
func NewK8sNetworkPolicyModule(client *dagger.Client) *K8sNetworkPolicyModule {
	return &K8sNetworkPolicyModule{
		client: client,
		name:   "k8s-network-policy",
	}
}

// KubectlNetworkPolicy manages network policies using kubectl
func (m *K8sNetworkPolicyModule) KubectlNetworkPolicy(ctx context.Context, action string, resource string, namespace string, outputFormat string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("bitnami/kubectl:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	// Build kubectl command based on action
	var args []string
	switch action {
	case "get":
		args = []string{k8sKubectlBinary, "get", "networkpolicy"}
		if resource != "" {
			args = append(args, resource)
		}
	case "describe":
		args = []string{k8sKubectlBinary, "describe", "networkpolicy"}
		if resource != "" {
			args = append(args, resource)
		}
	case "delete":
		if resource == "" {
			return "", fmt.Errorf("resource is required for delete action")
		}
		args = []string{k8sKubectlBinary, "delete", "networkpolicy", resource}
	case "create":
		if resource == "" {
			return "", fmt.Errorf("resource file path is required for create action")
		}
		container = container.WithFile("/policy.yaml", m.client.Host().File(resource))
		args = []string{k8sKubectlBinary, "create", "-f", "/policy.yaml"}
	case "apply":
		if resource == "" {
			return "", fmt.Errorf("resource file path is required for apply action")
		}
		container = container.WithFile("/policy.yaml", m.client.Host().File(resource))
		args = []string{k8sKubectlBinary, "apply", "-f", "/policy.yaml"}
	default:
		return "", fmt.Errorf("unsupported action: %s", action)
	}

	// Add namespace if specified
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	// Add output format if specified
	if outputFormat != "" {
		args = append(args, "-o", outputFormat)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("kubectl command failed: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// NetfetchScan scans network policies using netfetch
func (m *K8sNetworkPolicyModule) NetfetchScan(ctx context.Context, namespace string, dryrun bool, cilium bool, target string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/deggja/netfetch/releases/latest/download/netfetch-linux-amd64 -o /usr/local/bin/netfetch && chmod +x /usr/local/bin/netfetch"})

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{netfetchBinary, "scan"}
	
	if namespace != "" {
		args = append(args, namespace)
	}
	if dryrun {
		args = append(args, "--dryrun")
	}
	if cilium {
		args = append(args, "--cilium")
	}
	if target != "" {
		args = append(args, "--target", target)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("netfetch scan failed: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// NetfetchDashboard launches netfetch dashboard for network policy visualization
func (m *K8sNetworkPolicyModule) NetfetchDashboard(ctx context.Context, port string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/deggja/netfetch/releases/latest/download/netfetch-linux-amd64 -o /usr/local/bin/netfetch && chmod +x /usr/local/bin/netfetch"})

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{netfetchBinary, "dash"}
	
	if port != "" {
		args = append(args, "--port", port)
	} else {
		args = append(args, "--port", "8080")
	}

	// Expose the port
	portNum := "8080"
	if port != "" {
		portNum = port
	}
	container = container.WithExposedPort(8080).WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("netfetch dashboard failed: %w\nStderr: %s", err, stderr)
	}

	return fmt.Sprintf("Dashboard started on port %s\n%s", portNum, output), nil
}

// NetpolEval evaluates network connectivity using netpol-analyzer
func (m *K8sNetworkPolicyModule) NetpolEval(ctx context.Context, dirpath string, source string, destination string, port string, verbose bool) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/np-guard/netpol-analyzer/releases/latest/download/netpol-linux-amd64 -o /usr/local/bin/netpol && chmod +x /usr/local/bin/netpol"})

	// Mount the directory containing Kubernetes resources
	if dirpath != "" {
		container = container.WithMountedDirectory("/workspace", m.client.Host().Directory(dirpath))
	}

	args := []string{"netpol", "eval"}
	
	if dirpath != "" {
		args = append(args, "--dirpath", "/workspace")
	}
	if source != "" {
		args = append(args, "-s", source)
	}
	if destination != "" {
		args = append(args, "-d", destination)
	}
	if port != "" {
		args = append(args, "-p", port)
	}
	if verbose {
		args = append(args, "-v")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("netpol eval failed: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// NetpolList lists all allowed connections using netpol-analyzer
func (m *K8sNetworkPolicyModule) NetpolList(ctx context.Context, dirpath string, verbose bool, quiet bool) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/np-guard/netpol-analyzer/releases/latest/download/netpol-linux-amd64 -o /usr/local/bin/netpol && chmod +x /usr/local/bin/netpol"})

	// Mount the directory containing Kubernetes resources
	if dirpath != "" {
		container = container.WithMountedDirectory("/workspace", m.client.Host().Directory(dirpath))
	}

	args := []string{"netpol", "list"}
	
	if dirpath != "" {
		args = append(args, "--dirpath", "/workspace")
	}
	if verbose {
		args = append(args, "-v")
	}
	if quiet {
		args = append(args, "-q")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("netpol list failed: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// NetpolDiff compares network policies between two directories
func (m *K8sNetworkPolicyModule) NetpolDiff(ctx context.Context, dir1 string, dir2 string, outputFormat string) (string, error) {
	container := m.client.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://github.com/np-guard/netpol-analyzer/releases/latest/download/netpol-linux-amd64 -o /usr/local/bin/netpol && chmod +x /usr/local/bin/netpol"})

	// Mount both directories
	if dir1 != "" {
		container = container.WithMountedDirectory("/workspace1", m.client.Host().Directory(dir1))
	}
	if dir2 != "" {
		container = container.WithMountedDirectory("/workspace2", m.client.Host().Directory(dir2))
	}

	args := []string{"netpol", "diff"}
	
	if dir1 != "" && dir2 != "" {
		args = append(args, "--dir1", "/workspace1", "--dir2", "/workspace2")
	}

	// Add output format if specified
	if outputFormat != "" {
		if strings.Contains(outputFormat, "md") || strings.Contains(outputFormat, "markdown") {
			args = append(args, "--output-format", "md")
		} else if strings.Contains(outputFormat, "csv") {
			args = append(args, "--output-format", "csv")
		} else if strings.Contains(outputFormat, "text") {
			args = append(args, "--output-format", "txt")
		}
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("netpol diff failed: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// AnalyzePolicies analyzes network policies in the cluster (legacy function for compatibility)
func (m *K8sNetworkPolicyModule) AnalyzePolicies(ctx context.Context, namespace string, kubeconfig string) (string, error) {
	return m.KubectlNetworkPolicy(ctx, "get", "", namespace, "json", kubeconfig)
}

// ValidatePolicy validates a network policy (legacy function for compatibility)
func (m *K8sNetworkPolicyModule) ValidatePolicy(ctx context.Context, policyPath string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("bitnami/kubectl:latest").
		WithFile("/policy.yaml", m.client.Host().File(policyPath))

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		k8sKubectlBinary, "apply",
		"--dry-run=client",
		"--validate=true",
		"-f", "/policy.yaml",
		"--output", "json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to validate network policy: %w", err)
	}

	return output, nil
}

// TestConnectivity tests network connectivity between pods (legacy function for compatibility)
func (m *K8sNetworkPolicyModule) TestConnectivity(ctx context.Context, sourceNamespace string, targetNamespace string, targetService string, kubeconfig string) (string, error) {
	container := m.client.Container().
		From("nicolaka/netshoot:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	container = container.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("nc -zv %s.%s.svc.cluster.local 80 || echo 'Connection failed'", targetService, targetNamespace),
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to test connectivity: %w", err)
	}

	return output, nil
}