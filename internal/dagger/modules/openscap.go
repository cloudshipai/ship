package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// OpenSCAPModule runs OpenSCAP for security compliance scanning
type OpenSCAPModule struct {
	client *dagger.Client
	name   string
}

// NewOpenSCAPModule creates a new OpenSCAP module
func NewOpenSCAPModule(client *dagger.Client) *OpenSCAPModule {
	return &OpenSCAPModule{
		client: client,
		name:   "openscap",
	}
}

// EvaluateProfile evaluates a system against SCAP content
func (m *OpenSCAPModule) EvaluateProfile(ctx context.Context, contentPath string, profile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/content.xml", m.client.Host().File(contentPath)).
		WithExec([]string{
			"oscap",
			"xccdf", "eval",
			"--profile", profile,
			"--results", "/results.xml",
			"--report", "/report.html",
			"/content.xml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate profile: %w", err)
	}

	return output, nil
}

// ScanImage scans a container image for compliance
func (m *OpenSCAPModule) ScanImage(ctx context.Context, imageName string, profile string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithExec([]string{
			"oscap-podman",
			imageName,
			"xccdf", "eval",
			"--profile", profile,
			"--report", "/report.html",
			"--results", "/results.xml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan image: %w", err)
	}

	return output, nil
}

// GenerateReport generates compliance report
func (m *OpenSCAPModule) GenerateReport(ctx context.Context, resultsPath string) (string, error) {
	container := m.client.Container().
		From("cloudshipai/openscap:latest").
		WithFile("/results.xml", m.client.Host().File(resultsPath)).
		WithExec([]string{
			"oscap",
			"xccdf", "generate", "report",
			"/results.xml",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	return output, nil
}
