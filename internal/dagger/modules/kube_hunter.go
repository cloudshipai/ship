package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// KubeHunterModule runs kube-hunter for Kubernetes penetration testing
type KubeHunterModule struct {
	client *dagger.Client
	name   string
}

const kubeHunterBinary = "/usr/local/bin/kube-hunter"

// NewKubeHunterModule creates a new kube-hunter module
func NewKubeHunterModule(client *dagger.Client) *KubeHunterModule {
	return &KubeHunterModule{
		client: client,
		name:   kubeHunterBinary,
	}
}

// ScanRemote scans remote Kubernetes cluster
func (m *KubeHunterModule) ScanRemote(ctx context.Context, remote string, active bool, reportFormat string) (string, error) {
	args := []string{kubeHunterBinary, "--remote", remote}
	
	if active {
		args = append(args, "--active")
	}
	
	if reportFormat == "" {
		reportFormat = "json"
	}
	args = append(args, "--report", reportFormat)

	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to run kube-hunter remote scan: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// ScanCIDR scans CIDR range for Kubernetes clusters
func (m *KubeHunterModule) ScanCIDR(ctx context.Context, cidr string, active bool, reportFormat string) (string, error) {
	args := []string{kubeHunterBinary, "--cidr", cidr}
	
	if active {
		args = append(args, "--active")
	}
	
	if reportFormat == "" {
		reportFormat = "json"
	}
	args = append(args, "--report", reportFormat)

	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to run kube-hunter CIDR scan: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// ScanInterface scans network interface
func (m *KubeHunterModule) ScanInterface(ctx context.Context, networkInterface string, active bool, reportFormat string) (string, error) {
	args := []string{kubeHunterBinary, "--interface", networkInterface}
	
	if active {
		args = append(args, "--active")
	}
	
	if reportFormat == "" {
		reportFormat = "json"
	}
	args = append(args, "--report", reportFormat)

	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to run kube-hunter interface scan: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// ScanPod runs kube-hunter as pod in cluster
func (m *KubeHunterModule) ScanPod(ctx context.Context, kubeconfig string, active bool, reportFormat string) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-hunter:latest")

	if kubeconfig != "" {
		container = container.WithFile("/root/.kube/config", m.client.Host().File(kubeconfig))
	}

	args := []string{kubeHunterBinary, "--pod"}
	
	if active {
		args = append(args, "--active")
	}
	
	if reportFormat == "" {
		reportFormat = "json"
	}
	args = append(args, "--report", reportFormat)

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to run kube-hunter pod scan: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// ListTests lists all available tests
func (m *KubeHunterModule) ListTests(ctx context.Context, showActive bool) (string, error) {
	args := []string{kubeHunterBinary, "--list"}
	
	if showActive {
		args = append(args, "--active")
	}

	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to list kube-hunter tests: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// RunCustomHunters runs kube-hunter with specific hunters enabled/disabled
func (m *KubeHunterModule) RunCustomHunters(ctx context.Context, target string, includeHunters []string, excludeHunters []string, active bool) (string, error) {
	args := []string{kubeHunterBinary}
	
	// Add target (could be remote, cidr, or interface)
	if strings.Contains(target, "/") {
		// Looks like CIDR
		args = append(args, "--cidr", target)
	} else if strings.Contains(target, ".") || strings.Contains(target, ":") {
		// Looks like IP or hostname
		args = append(args, "--remote", target)
	} else {
		// Assume it's an interface name
		args = append(args, "--interface", target)
	}
	
	// Add include hunters
	for _, hunter := range includeHunters {
		args = append(args, "--include-hunter-type", hunter)
	}
	
	// Add exclude hunters
	for _, hunter := range excludeHunters {
		args = append(args, "--exclude-hunter-type", hunter)
	}
	
	if active {
		args = append(args, "--active")
	}
	
	args = append(args, "--report", "json")

	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		return "", fmt.Errorf("failed to run custom hunters: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// GetVersion returns the version of kube-hunter
func (m *KubeHunterModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/kube-hunter:latest").
		WithExec([]string{kubeHunterBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get kube-hunter version: %w", err)
	}

	return output, nil
}