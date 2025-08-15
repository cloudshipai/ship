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

// NewProwlerModule creates a new Prowler module
func NewProwlerModule(client *dagger.Client) *ProwlerModule {
	return &ProwlerModule{
		client: client,
		name:   "prowler",
	}
}

// ScanAWS scans AWS infrastructure for security issues
func (m *ProwlerModule) ScanAWS(ctx context.Context, provider string, region string) (string, error) {
	container := m.client.Container().
		From("toniblyx/prowler:latest").
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_REGION", region).
		WithExec([]string{"prowler", "aws", "--output-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler AWS scan: %w", err)
	}

	return output, nil
}

// ScanAzure scans Azure infrastructure for security issues
func (m *ProwlerModule) ScanAzure(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("toniblyx/prowler:latest").
		WithEnvVariable("AZURE_SUBSCRIPTION_ID", os.Getenv("AZURE_SUBSCRIPTION_ID")).
		WithEnvVariable("AZURE_TENANT_ID", os.Getenv("AZURE_TENANT_ID")).
		WithEnvVariable("AZURE_CLIENT_ID", os.Getenv("AZURE_CLIENT_ID")).
		WithEnvVariable("AZURE_CLIENT_SECRET", os.Getenv("AZURE_CLIENT_SECRET")).
		WithExec([]string{"prowler", "azure", "--output-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler Azure scan: %w", err)
	}

	return output, nil
}

// ScanGCP scans Google Cloud Platform for security issues
func (m *ProwlerModule) ScanGCP(ctx context.Context, projectId string) (string, error) {
	container := m.client.Container().
		From("toniblyx/prowler:latest").
		WithEnvVariable("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")).
		WithExec([]string{"prowler", "gcp", "--project-id", projectId, "--output-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler GCP scan: %w", err)
	}

	return output, nil
}

// ScanKubernetes scans Kubernetes cluster for security issues
func (m *ProwlerModule) ScanKubernetes(ctx context.Context, kubeconfigPath string) (string, error) {
	container := m.client.Container().
		From("toniblyx/prowler:latest").
		WithFile("/tmp/kubeconfig", m.client.Host().File(kubeconfigPath)).
		WithEnvVariable("KUBECONFIG", "/tmp/kubeconfig").
		WithExec([]string{"prowler", "kubernetes", "--output-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler Kubernetes scan: %w", err)
	}

	return output, nil
}

// ScanWithCompliance scans with specific compliance frameworks
func (m *ProwlerModule) ScanWithCompliance(ctx context.Context, provider string, compliance string, region string) (string, error) {
	container := m.client.Container().
		From("toniblyx/prowler:latest").
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_REGION", region).
		WithExec([]string{"prowler", provider, "--compliance", compliance, "--output-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler compliance scan: %w", err)
	}

	return output, nil
}

// ScanSpecificServices scans specific cloud services
func (m *ProwlerModule) ScanSpecificServices(ctx context.Context, provider string, services string, region string) (string, error) {
	container := m.client.Container().
		From("toniblyx/prowler:latest").
		WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
		WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
		WithEnvVariable("AWS_REGION", region).
		WithExec([]string{"prowler", provider, "--services", services, "--output-format", "json"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run prowler service scan: %w", err)
	}

	return output, nil
}