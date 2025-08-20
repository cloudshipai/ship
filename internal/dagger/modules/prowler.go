package modules

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// ProwlerModule runs Prowler for cloud security assessment
type ProwlerModule struct {
	client *dagger.Client
	name   string
}

const prowlerBinary = "prowler"

// NewProwlerModule creates a new Prowler module
func NewProwlerModule(client *dagger.Client) *ProwlerModule {
	return &ProwlerModule{
		client: client,
		name:   prowlerBinary,
	}
}

// ScanAWS scans AWS infrastructure for security issues
func (m *ProwlerModule) ScanAWS(ctx context.Context, provider string, region string) (string, error) {
	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_REGION", region).
		WithExec([]string{prowlerBinary, "aws", "--output-format", "json"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler AWS scan: %w", err)
	}

	return output, nil
}

// ScanAzure scans Azure infrastructure for security issues
func (m *ProwlerModule) ScanAzure(ctx context.Context) (string, error) {
	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithEnvVariable("AZURE_SUBSCRIPTION_ID", os.Getenv("AZURE_SUBSCRIPTION_ID")).
		WithEnvVariable("AZURE_TENANT_ID", os.Getenv("AZURE_TENANT_ID")).
		WithEnvVariable("AZURE_CLIENT_ID", os.Getenv("AZURE_CLIENT_ID")).
		WithEnvVariable("AZURE_CLIENT_SECRET", os.Getenv("AZURE_CLIENT_SECRET")).
		WithExec([]string{prowlerBinary, "azure", "--output-format", "json"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler Azure scan: %w", err)
	}

	return output, nil
}

// ScanGCP scans Google Cloud Platform for security issues
func (m *ProwlerModule) ScanGCP(ctx context.Context, projectId string) (string, error) {
	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithEnvVariable("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")).
		WithExec([]string{prowlerBinary, "gcp", "--project-id", projectId, "--output-format", "json"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler GCP scan: %w", err)
	}

	return output, nil
}

// ScanKubernetes scans Kubernetes cluster for security issues
func (m *ProwlerModule) ScanKubernetes(ctx context.Context, kubeconfigPath string) (string, error) {
	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithFile("/tmp/kubeconfig", m.client.Host().File(kubeconfigPath)).
		WithEnvVariable("KUBECONFIG", "/tmp/kubeconfig").
		WithExec([]string{prowlerBinary, "kubernetes", "--output-format", "json"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler Kubernetes scan: %w", err)
	}

	return output, nil
}

// ScanWithCompliance scans with specific compliance frameworks
func (m *ProwlerModule) ScanWithCompliance(ctx context.Context, provider string, compliance string, region string) (string, error) {
	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_REGION", region).
		WithExec([]string{prowlerBinary, provider, "--compliance", compliance, "--output-format", "json"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler compliance scan: %w", err)
	}

	return output, nil
}

// ScanSpecificServices scans specific cloud services
func (m *ProwlerModule) ScanSpecificServices(ctx context.Context, provider string, services string, region string) (string, error) {
	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_REGION", region).
		WithExec([]string{prowlerBinary, provider, "--services", services, "--output-format", "json"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler service scan: %w", err)
	}

	return output, nil
}

// ListChecks lists available Prowler checks
func (m *ProwlerModule) ListChecks(ctx context.Context, provider string) (string, error) {
	args := []string{prowlerBinary, "--list-checks"}
	if provider != "" {
		args = append(args, "--provider", provider)
	}

	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list checks: %w", err)
	}

	return output, nil
}

// ListServices lists available services for a provider
func (m *ProwlerModule) ListServices(ctx context.Context, provider string) (string, error) {
	args := []string{prowlerBinary, "--list-services"}
	if provider != "" {
		args = append(args, "--provider", provider)
	}

	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list services: %w", err)
	}

	return output, nil
}

// ListCompliance lists available compliance frameworks
func (m *ProwlerModule) ListCompliance(ctx context.Context, provider string) (string, error) {
	args := []string{prowlerBinary, "--list-compliance"}
	if provider != "" {
		args = append(args, "--provider", provider)
	}

	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list compliance frameworks: %w", err)
	}

	return output, nil
}

// GenerateDashboard generates HTML dashboard from scan results
func (m *ProwlerModule) GenerateDashboard(ctx context.Context, inputFile string, outputDir string) (string, error) {
	// Use Python base image and install Prowler since official image doesn't support ARM64
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "--no-cache-dir", "prowler"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		}).
		WithFile("/workspace/input.json", m.client.Host().File(inputFile)).
		WithWorkdir("/workspace").
		WithExec([]string{prowlerBinary, "dashboard", "--input", "input.json", "--output-dir", outputDir}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate dashboard: %w", err)
	}

	return output, nil
}