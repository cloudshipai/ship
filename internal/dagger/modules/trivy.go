package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// TrivyModule runs Trivy for comprehensive vulnerability scanning
type TrivyModule struct {
	client *dagger.Client
	name   string
}

const trivyBinaryMain = "/usr/local/bin/trivy"

// NewTrivyModule creates a new Trivy module
func NewTrivyModule(client *dagger.Client) *TrivyModule {
	return &TrivyModule{
		client: client,
		name:   trivyBinaryMain,
	}
}

// ScanImage scans a container image for vulnerabilities
func (m *TrivyModule) ScanImage(ctx context.Context, imageName string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			trivyBinaryMain,
			"image",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			imageName,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan image: %w", err)
	}

	return output, nil
}

// ScanFilesystem scans a filesystem for vulnerabilities
func (m *TrivyModule) ScanFilesystem(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			trivyBinaryMain,
			"fs",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			".",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan filesystem: %w", err)
	}

	return output, nil
}

// ScanRepository scans a git repository
func (m *TrivyModule) ScanRepository(ctx context.Context, repoURL string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{
			trivyBinaryMain,
			"repo",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			repoURL,
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan repository: %w", err)
	}

	return output, nil
}

// ScanConfig scans configuration files for misconfigurations
func (m *TrivyModule) ScanConfig(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			trivyBinaryMain,
			"config",
			"--format", "json",
			"--severity", "HIGH,CRITICAL",
			".",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan config: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of Trivy
func (m *TrivyModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithExec([]string{trivyBinaryMain, "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get trivy version: %w", err)
	}

	return output, nil
}

// ScanSBOM scans SBOM file for vulnerabilities
func (m *TrivyModule) ScanSBOM(ctx context.Context, sbomPath string, severity string, outputFormat string, outputFile string, ignoreUnfixed bool) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithFile("/sbom.json", m.client.Host().File(sbomPath))

	args := []string{trivyBinaryMain, "sbom", "/sbom.json"}
	if severity != "" {
		args = append(args, "--severity", severity)
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	if ignoreUnfixed {
		args = append(args, "--ignore-unfixed")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan SBOM: %w", err)
	}

	return output, nil
}

// ScanKubernetes scans Kubernetes cluster for vulnerabilities
func (m *TrivyModule) ScanKubernetes(ctx context.Context, target string, clusterContext string, namespace string, severity string, outputFormat string, scanners string, includeImages bool) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	if target == "" {
		target = "cluster"
	}

	args := []string{trivyBinaryMain, "k8s", target}
	if clusterContext != "" {
		args = append(args, "--context", clusterContext)
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	if severity != "" {
		args = append(args, "--severity", severity)
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}
	if scanners != "" {
		args = append(args, "--scanners", scanners)
	}
	if includeImages {
		args = append(args, "--include-images")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan Kubernetes: %w", err)
	}

	return output, nil
}

// GenerateSBOM generates Software Bill of Materials
func (m *TrivyModule) GenerateSBOM(ctx context.Context, target string, targetType string, sbomFormat string, outputFile string, includeDevDeps bool) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	if targetType == "fs" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		target = "."
	}

	args := []string{trivyBinaryMain, targetType, "--format", sbomFormat, target}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}
	if includeDevDeps {
		args = append(args, "--include-dev-deps")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate SBOM: %w", err)
	}

	return output, nil
}

// ScanWithFilters scans with advanced filtering options
func (m *TrivyModule) ScanWithFilters(ctx context.Context, target string, targetType string, severity string, vulnType string, ignoreFile string, ignoreUnfixed bool, exitCode string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	if targetType == "fs" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		target = "."
	}
	if ignoreFile != "" {
		container = container.WithFile("/.trivyignore", m.client.Host().File(ignoreFile))
	}

	args := []string{trivyBinaryMain, targetType, target}
	if severity != "" {
		args = append(args, "--severity", severity)
	}
	if vulnType != "" {
		args = append(args, "--vuln-type", vulnType)
	}
	if ignoreFile != "" {
		args = append(args, "--ignorefile", "/.trivyignore")
	}
	if ignoreUnfixed {
		args = append(args, "--ignore-unfixed")
	}
	if exitCode != "" {
		args = append(args, "--exit-code", exitCode)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan with filters: %w", err)
	}

	return output, nil
}

// DatabaseUpdate performs database operations
func (m *TrivyModule) DatabaseUpdate(ctx context.Context, operation string, skipUpdate bool, cacheDir string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	var args []string
	switch operation {
	case "download", "update":
		args = []string{trivyBinaryMain, "image", "--download-db-only", "alpine"}
		if cacheDir != "" {
			args = append(args, "--cache-dir", cacheDir)
		}
	case "reset", "clean":
		args = []string{trivyBinaryMain, "clean", "--all"}
		if cacheDir != "" {
			args = append(args, "--cache-dir", cacheDir)
		}
	default:
		args = []string{trivyBinaryMain, "--help"}
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to perform database operation: %w", err)
	}

	return output, nil
}

// ServerMode runs Trivy in server mode
func (m *TrivyModule) ServerMode(ctx context.Context, listenPort string, listenAddress string, debug bool, token string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	args := []string{trivyBinaryMain, "server"}
	if listenPort != "" {
		args = append(args, "--listen", "0.0.0.0:"+listenPort)
	}
	if listenAddress != "" {
		args = append(args, "--listen", listenAddress)
	}
	if debug {
		args = append(args, "--debug")
	}
	if token != "" {
		args = append(args, "--token", token)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start server mode: %w", err)
	}

	return output, nil
}

// ClientScan scans using client mode
func (m *TrivyModule) ClientScan(ctx context.Context, target string, targetType string, serverURL string, token string, outputFormat string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	if targetType == "fs" {
		container = container.WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
		target = "."
	}

	args := []string{trivyBinaryMain, targetType, "--server", serverURL, target}
	if token != "" {
		args = append(args, "--token", token)
	}
	if outputFormat != "" {
		args = append(args, "--format", outputFormat)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan using client mode: %w", err)
	}

	return output, nil
}

// PluginManagement manages Trivy plugins
func (m *TrivyModule) PluginManagement(ctx context.Context, action string, pluginName string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest")

	args := []string{trivyBinaryMain, "plugin", action}
	if pluginName != "" && action != "list" {
		args = append(args, pluginName)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to manage plugins: %w", err)
	}

	return output, nil
}

// ConvertSBOM converts SBOM between different formats
func (m *TrivyModule) ConvertSBOM(ctx context.Context, inputSBOM string, outputFormat string, outputFile string) (string, error) {
	container := m.client.Container().
		From("aquasec/trivy:latest").
		WithFile("/input.sbom", m.client.Host().File(inputSBOM))

	args := []string{trivyBinaryMain, "convert", "--format", outputFormat, "/input.sbom"}
	if outputFile != "" {
		args = append(args, "--output", outputFile)
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to convert SBOM: %w", err)
	}

	return output, nil
}
