package modules

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// CloudsplainingModule runs Cloudsplaining for AWS IAM security assessment
type CloudsplainingModule struct {
	client *dagger.Client
	name   string
}

// NewCloudsplainingModule creates a new Cloudsplaining module
func NewCloudsplainingModule(client *dagger.Client) *CloudsplainingModule {
	return &CloudsplainingModule{
		client: client,
		name:   "cloudsplaining",
	}
}

// ScanAccountAuthorization scans account authorization details
func (m *CloudsplainingModule) ScanAccountAuthorization(ctx context.Context, profile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	container = container.WithExec([]string{
		"cloudsplaining", "scan-account-authorization-details",
		"--output", "/workspace/results.json",
	})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan account authorization: %w", err)
	}

	return output, nil
}

// ScanPolicyFile scans a specific IAM policy file
func (m *CloudsplainingModule) ScanPolicyFile(ctx context.Context, policyPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithFile("/workspace/policy.json", m.client.Host().File(policyPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"cloudsplaining", "scan-policy-file",
			"--input-file", "policy.json",
			"--output", "results.json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan policy file: %w", err)
	}

	return output, nil
}

// CreateReportFromResults creates an HTML report from scan results
func (m *CloudsplainingModule) CreateReportFromResults(ctx context.Context, resultsPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithFile("/workspace/results.json", m.client.Host().File(resultsPath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"cloudsplaining", "create-report",
			"--input-file", "results.json",
			"--output", "report.html",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to create report: %w", err)
	}

	return output, nil
}

// ScanWithMinimization scans with policy minimization recommendations
func (m *CloudsplainingModule) ScanWithMinimization(ctx context.Context, profile string, minimizeStatementId string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/cloudsplaining:latest").
		WithEnvVariable("AWS_PROFILE", profile)

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		container = container.
			WithEnvVariable("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID")).
			WithEnvVariable("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY")).
			WithEnvVariable("AWS_REGION", os.Getenv("AWS_REGION"))
	}

	args := []string{
		"cloudsplaining", "scan-account-authorization-details",
		"--output", "/workspace/results.json",
	}

	if minimizeStatementId != "" {
		args = append(args, "--minimize-statement-id", minimizeStatementId)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		stderr, _ := container.Stderr(ctx)
		if stderr != "" {
			return stderr, nil
		}
		return "", fmt.Errorf("failed to scan with minimization: %w", err)
	}

	return output, nil
}