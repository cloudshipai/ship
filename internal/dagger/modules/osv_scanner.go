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
