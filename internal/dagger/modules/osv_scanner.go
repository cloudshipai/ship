package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// OSVScannerModule runs OSV Scanner for vulnerability detection
type OSVScannerModule struct {
	client *dagger.Client
	name   string
}

// NewOSVScannerModule creates a new OSV Scanner module
func NewOSVScannerModule(client *dagger.Client) *OSVScannerModule {
	return &OSVScannerModule{
		client: client,
		name:   "osv-scanner",
	}
}

// ScanDirectory scans a directory for vulnerabilities
func (m *OSVScannerModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"osv-scanner",
			"--format", "json",
			".",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run OSV scanner: %w", err)
	}

	return output, nil
}

// ScanLockfile scans a specific lockfile
func (m *OSVScannerModule) ScanLockfile(ctx context.Context, lockfilePath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithFile("/lockfile", m.client.Host().File(lockfilePath)).
		WithExec([]string{
			"osv-scanner",
			"--format", "json",
			"--lockfile", "/lockfile",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan lockfile: %w", err)
	}

	return output, nil
}

// ScanSBOM scans an SBOM file
func (m *OSVScannerModule) ScanSBOM(ctx context.Context, sbomPath string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithFile("/sbom.json", m.client.Host().File(sbomPath)).
		WithExec([]string{
			"osv-scanner",
			"--format", "json",
			"--sbom", "/sbom.json",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan SBOM: %w", err)
	}

	return output, nil
}

// GetVersion returns the version of OSV Scanner
func (m *OSVScannerModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithExec([]string{"osv-scanner", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get OSV scanner version: %w", err)
	}

	return output, nil
}

// ScanImage scans a container image for vulnerabilities
func (m *OSVScannerModule) ScanImage(ctx context.Context, image string, output string, format string, config string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest")

	args := []string{"osv-scanner", "scan", "image", image}
	if output != "" {
		args = append(args, "--output", output)
	}
	if format != "" {
		args = append(args, "--format", format)
	}
	if config != "" {
		container = container.WithFile("/config.yaml", m.client.Host().File(config))
		args = append(args, "--config", "/config.yaml")
	}

	container = container.WithExec(args)

	output_result, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan image: %w", err)
	}

	return output_result, nil
}

// ScanManifest scans a package manifest file
func (m *OSVScannerModule) ScanManifest(ctx context.Context, manifestPath string, output string, format string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithFile("/manifest", m.client.Host().File(manifestPath))

	args := []string{"osv-scanner", "-M", "/manifest"}
	if output != "" {
		args = append(args, "--output", output)
	}
	if format != "" {
		args = append(args, "--format", format)
	}

	container = container.WithExec(args)

	output_result, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan manifest: %w", err)
	}

	return output_result, nil
}

// LicenseScan scans for license compliance
func (m *OSVScannerModule) LicenseScan(ctx context.Context, path string, allowedLicenses string, output string) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithDirectory("/workspace", m.client.Host().Directory(path))

	args := []string{"osv-scanner", "--licenses"}
	if allowedLicenses != "" {
		args = []string{"osv-scanner", "--licenses=" + allowedLicenses}
	}
	if output != "" {
		args = append(args, "--output", output)
	}
	args = append(args, "/workspace")

	container = container.WithExec(args)

	output_result, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to scan licenses: %w", err)
	}

	return output_result, nil
}

// OfflineScan scans using offline vulnerability databases
func (m *OSVScannerModule) OfflineScan(ctx context.Context, path string, offlineDbPath string, downloadDatabases bool, recursive bool) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithDirectory("/workspace", m.client.Host().Directory(path))

	args := []string{"osv-scanner", "--offline"}
	if downloadDatabases {
		args = append(args, "--download-offline-databases")
	}
	if offlineDbPath != "" {
		container = container.WithDirectory("/db", m.client.Host().Directory(offlineDbPath))
		args = append(args, "--offline-vulnerabilities", "/db")
	}
	if recursive {
		args = append(args, "-r")
	}
	args = append(args, "/workspace")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to perform offline scan: %w", err)
	}

	return output, nil
}

// Fix applies guided remediation for vulnerabilities
func (m *OSVScannerModule) Fix(ctx context.Context, manifestPath string, lockfilePath string, strategy string, maxDepth string, minSeverity string, ignoreDev bool) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest")

	if manifestPath != "" {
		container = container.WithFile("/manifest", m.client.Host().File(manifestPath))
	}
	if lockfilePath != "" {
		container = container.WithFile("/lockfile", m.client.Host().File(lockfilePath))
	}

	args := []string{"osv-scanner", "fix"}
	if manifestPath != "" {
		args = append(args, "-M", "/manifest")
	}
	if lockfilePath != "" {
		args = append(args, "-L", "/lockfile")
	}
	if strategy != "" {
		args = append(args, "--strategy", strategy)
	}
	if maxDepth != "" {
		args = append(args, "--max-depth", maxDepth)
	}
	if minSeverity != "" {
		args = append(args, "--min-severity", minSeverity)
	}
	if ignoreDev {
		args = append(args, "--ignore-dev")
	}

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fix vulnerabilities: %w", err)
	}

	return output, nil
}

// ServeReport generates and serves HTML vulnerability report locally
func (m *OSVScannerModule) ServeReport(ctx context.Context, path string, port string, recursive bool) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithDirectory("/workspace", m.client.Host().Directory(path))

	args := []string{"osv-scanner", "--serve"}
	if port != "" {
		args = append(args, "--port", port)
	}
	if recursive {
		args = append(args, "-r")
	}
	args = append(args, "/workspace")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to serve report: %w", err)
	}

	return output, nil
}

// VerboseScan runs OSV Scanner with verbose logging
func (m *OSVScannerModule) VerboseScan(ctx context.Context, path string, verbosity string, recursive bool) (string, error) {
	container := m.client.Container().
		From("ghcr.io/google/osv-scanner:latest").
		WithDirectory("/workspace", m.client.Host().Directory(path))

	args := []string{"osv-scanner"}
	if verbosity != "" {
		args = append(args, "--verbosity", verbosity)
	}
	if recursive {
		args = append(args, "-r")
	}
	args = append(args, "/workspace")

	container = container.WithExec(args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to perform verbose scan: %w", err)
	}

	return output, nil
}
