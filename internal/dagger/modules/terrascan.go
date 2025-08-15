package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// TerrascanModule runs Terrascan for IaC security scanning
type TerrascanModule struct {
	client *dagger.Client
	name   string
}

// NewTerrascanModule creates a new Terrascan module
func NewTerrascanModule(client *dagger.Client) *TerrascanModule {
	return &TerrascanModule{
		client: client,
		name:   "terrascan",
	}
}

// ScanDirectory scans a directory for IaC security issues using Terrascan
func (m *TerrascanModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Terrascan returns non-zero exit code when violations are found
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan: %w", err)
	}

	return output, nil
}

// ScanTerraform scans Terraform files specifically
func (m *TerrascanModule) ScanTerraform(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-t", "terraform", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan on terraform: %w", err)
	}

	return output, nil
}

// ScanKubernetes scans Kubernetes manifests
func (m *TerrascanModule) ScanKubernetes(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-t", "k8s", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan on kubernetes: %w", err)
	}

	return output, nil
}

// ScanCloudFormation scans CloudFormation templates
func (m *TerrascanModule) ScanCloudFormation(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-t", "cloudformation", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan on cloudformation: %w", err)
	}

	return output, nil
}

// ScanDockerfiles scans Dockerfile for security issues
func (m *TerrascanModule) ScanDockerfiles(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{"terrascan", "scan", "-t", "docker", "-d", ".", "-o", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan on docker: %w", err)
	}

	return output, nil
}

// ScanWithSeverity scans with a specific severity threshold
func (m *TerrascanModule) ScanWithSeverity(ctx context.Context, dir string, severity string, iacType string) (string, error) {
	args := []string{"terrascan", "scan", "-d", ".", "-o", "json"}
	
	if iacType != "" {
		args = append(args, "-t", iacType)
	}
	
	if severity != "" {
		args = append(args, "--severity", severity)
	}

	container := m.client.Container().
		From("tenable/terrascan:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to run terrascan with severity %s: %w", severity, err)
	}

	return output, nil
}