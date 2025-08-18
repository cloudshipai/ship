package modules

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

// InfraScanModule runs infrastructure scanning tools
// Supports both AWS infrastructure mapping (infrascan) and security scanning (Trivy)
type InfraScanModule struct {
	client *dagger.Client
	name   string
}

const (
	infrascanBinary     = "/usr/local/bin/infrascan"
	infrascanTrivyBinary = "/usr/local/bin/trivy"
)

// NewInfraScanModule creates a new InfraScan module
func NewInfraScanModule(client *dagger.Client) *InfraScanModule {
	return &InfraScanModule{
		client: client,
		name:   "infrascan",
	}
}

// ScanAWSInfrastructure scans AWS infrastructure and generates a system map
func (m *InfraScanModule) ScanAWSInfrastructure(ctx context.Context, regions []string, outputDir string) (string, error) {
	// Create container with infrascan CLI
	container := m.client.Container().
		From("node:18-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "aws-cli"}).
		WithExec([]string{"npm", "install", "-g", "@infrascan/cli"}).
		WithMountedDirectory("/workspace", m.client.Host().Directory(".")).
		WithWorkdir("/workspace")

	// Build command
	args := []string{infrascanBinary, "scan", "-o", outputDir}
	
	// Add regions
	for _, region := range regions {
		region = strings.TrimSpace(region)
		if region != "" {
			args = append(args, "--region", region)
		}
	}

	// Execute scan
	result := container.WithExec(args)

	output, err := result.Stdout(ctx)
	if err != nil {
		stderr, _ := result.Stderr(ctx)
		return "", fmt.Errorf("failed to scan AWS infrastructure: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// GenerateGraph generates a graph from scan results
func (m *InfraScanModule) GenerateGraph(ctx context.Context, inputDir string) (string, error) {
	// Create container with infrascan CLI
	container := m.client.Container().
		From("node:18-alpine").
		WithExec([]string{"npm", "install", "-g", "@infrascan/cli"}).
		WithMountedDirectory("/workspace", m.client.Host().Directory(".")).
		WithWorkdir("/workspace")

	// Execute graph generation
	result := container.WithExec([]string{infrascanBinary, "graph", "-i", inputDir})

	output, err := result.Stdout(ctx)
	if err != nil {
		stderr, _ := result.Stderr(ctx)
		return "", fmt.Errorf("failed to generate graph: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// RenderGraph renders an infrastructure graph
func (m *InfraScanModule) RenderGraph(ctx context.Context, inputFile string, openBrowser bool) (string, error) {
	// Create container with infrascan CLI
	container := m.client.Container().
		From("node:18-alpine").
		WithExec([]string{"npm", "install", "-g", "@infrascan/cli"}).
		WithMountedDirectory("/workspace", m.client.Host().Directory(".")).
		WithWorkdir("/workspace")

	// Build command
	args := []string{infrascanBinary, "render", "-i", inputFile}
	if openBrowser {
		args = append(args, "--browser")
	}

	// Execute render
	result := container.WithExec(args)

	output, err := result.Stdout(ctx)
	if err != nil {
		stderr, _ := result.Stderr(ctx)
		return "", fmt.Errorf("failed to render graph: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// ScanDirectory scans a directory for security issues (using Trivy)
func (m *InfraScanModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	// Mount the directory and run Trivy
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			infrascanTrivyBinary,
			"fs",
			".",
			"--scanners", "misconfig",
			"--format", "json",
			"--severity", "HIGH,CRITICAL,MEDIUM,LOW",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run trivy: %w", err)
	}

	return output, nil
}

// ScanFile scans a specific Terraform file (using Trivy)
func (m *InfraScanModule) ScanFile(ctx context.Context, filePath string) (string, error) {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)

	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			infrascanTrivyBinary,
			"fs",
			filename,
			"--scanners", "misconfig",
			"--format", "json",
			"--severity", "HIGH,CRITICAL,MEDIUM,LOW",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run trivy on file: %w", err)
	}

	return output, nil
}

// ScanWithRules scans using custom rule set (using Trivy)
func (m *InfraScanModule) ScanWithRules(ctx context.Context, dir string, rulesFile string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir))

	// If rules file is provided, mount it
	if rulesFile != "" {
		container = container.WithFile("/policy.rego", m.client.Host().File(rulesFile))
		container = container.WithExec([]string{
			infrascanTrivyBinary,
			"fs",
			"/workspace",
			"--scanners", "misconfig",
			"--format", "json",
			"--severity", "HIGH,CRITICAL,MEDIUM,LOW",
			"--config-policy", "/policy.rego",
		})
	} else {
		container = container.WithExec([]string{
			infrascanTrivyBinary,
			"fs",
			"/workspace",
			"--scanners", "misconfig",
			"--format", "json",
			"--severity", "HIGH,CRITICAL,MEDIUM,LOW",
		})
	}

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run trivy with rules: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of the scanner (infrascan)
func (m *InfraScanModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("node:18-alpine").
		WithExec([]string{"npm", "install", "-g", "@infrascan/cli"}).
		WithExec([]string{infrascanBinary, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		// Fallback to Trivy version if infrascan not available
		container = m.client.Container().
			From("aquasec/trivy:latest").
			WithExec([]string{infrascanTrivyBinary, "--version"})
		
		output, err = container.Stdout(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get version: %w", err)
		}
	}

	return output, nil
}