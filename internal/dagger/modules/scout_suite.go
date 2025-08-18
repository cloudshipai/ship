package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// ScoutSuiteModule runs Scout Suite for multi-cloud security auditing
type ScoutSuiteModule struct {
	client *dagger.Client
	name   string
}

const scoutSuiteBinary = "/usr/local/bin/scout"

// NewScoutSuiteModule creates a new Scout Suite module
func NewScoutSuiteModule(client *dagger.Client) *ScoutSuiteModule {
	return &ScoutSuiteModule{
		client: client,
		name:   scoutSuiteBinary,
	}
}

// ScanAWS scans AWS environment
func (m *ScoutSuiteModule) ScanAWS(ctx context.Context, profile string) (string, error) {
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "scoutsuite"}).
		WithEnvVariable("AWS_PROFILE", profile).
		WithExec([]string{
			scoutSuiteBinary,
			"aws",
			"--report-dir", "/tmp/scout-report",
			"--force",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run scout suite AWS scan: %w", err)
	}

	return output, nil
}

// ScanAzure scans Azure environment
func (m *ScoutSuiteModule) ScanAzure(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "scoutsuite"}).
		WithExec([]string{
			scoutSuiteBinary,
			"azure",
			"--report-dir", "/tmp/scout-report",
			"--force",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run scout suite Azure scan: %w", err)
	}

	return output, nil
}

// ScanGCP scans Google Cloud Platform environment
func (m *ScoutSuiteModule) ScanGCP(ctx context.Context, projectID string) (string, error) {
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "scoutsuite"}).
		WithEnvVariable("GOOGLE_CLOUD_PROJECT", projectID).
		WithExec([]string{
			scoutSuiteBinary,
			"gcp",
			"--project-id", projectID,
			"--report-dir", "/tmp/scout-report",
			"--force",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run scout suite GCP scan: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Scout Suite
func (m *ScoutSuiteModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "scoutsuite"}).
		WithExec([]string{scoutSuiteBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get scout suite version: %w", err)
	}

	return output, nil
}

// ServeReport serves Scout Suite report via web server
func (m *ScoutSuiteModule) ServeReport(ctx context.Context, provider string, reportName string, host string, port string) (string, error) {
	args := []string{scoutSuiteBinary, provider, "--serve"}
	
	if reportName != "" {
		args = append(args, reportName)
	}
	if host != "" {
		args = append(args, "--host", host)
	}
	if port != "" {
		args = append(args, "--port", port)
	}

	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "scoutsuite"}).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to serve scout suite report: %w", err)
	}

	return output, nil
}

// Help returns Scout Suite help information
func (m *ScoutSuiteModule) Help(ctx context.Context, provider string) (string, error) {
	args := []string{scoutSuiteBinary, "--help"}
	
	if provider != "" {
		args = []string{scoutSuiteBinary, provider, "--help"}
	}

	container := m.client.Container().
		From("python:3.11-slim").
		WithExec([]string{"pip", "install", "scoutsuite"}).
		WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get scout suite help: %w", err)
	}

	return output, nil
}
