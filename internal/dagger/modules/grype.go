package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GrypeModule runs Grype for vulnerability scanning
type GrypeModule struct {
	client *dagger.Client
	name   string
}

// NewGrypeModule creates a new Grype module
func NewGrypeModule(client *dagger.Client) *GrypeModule {
	return &GrypeModule{
		client: client,
		name:   "grype",
	}
}

// ScanDirectory scans a directory for vulnerabilities using Grype
func (m *GrypeModule) ScanDirectory(ctx context.Context, dir string, format ...string) (string, error) {
	outputFormat := "json"
	if len(format) > 0 && format[0] != "" {
		outputFormat = format[0]
	}
	
	container := m.client.Container().
		From("anchore/grype:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"/grype", ".", "-o", outputFormat,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run grype scan: no output received")
}

// ScanImage scans a container image for vulnerabilities
func (m *GrypeModule) ScanImage(ctx context.Context, imageName string, format ...string) (string, error) {
	outputFormat := "json"
	if len(format) > 0 && format[0] != "" {
		outputFormat = format[0]
	}
	
	container := m.client.Container().
		From("anchore/grype:latest").
		WithExec([]string{
			"/grype", imageName, "-o", outputFormat,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run grype image scan: no output received")
}

// ScanSBOM scans an SBOM file for vulnerabilities
func (m *GrypeModule) ScanSBOM(ctx context.Context, sbomPath string, format ...string) (string, error) {
	outputFormat := "json"
	if len(format) > 0 && format[0] != "" {
		outputFormat = format[0]
	}
	dir := "/workspace"
	container := m.client.Container().
		From("anchore/grype:latest").
		WithFile(dir+"/sbom.json", m.client.Host().File(sbomPath)).
		WithWorkdir(dir).
		WithExec([]string{
			"/grype", "sbom:sbom.json", "-o", outputFormat,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run grype SBOM scan: no output received")
}

// ScanWithSeverity scans with a specific severity threshold
func (m *GrypeModule) ScanWithSeverity(ctx context.Context, target string, severity string) (string, error) {
	var args []string
	if target[:4] == "dir:" || target[:6] == "image:" {
		args = []string{"/grype", target, "--fail-on", severity, "-o", "json"}
	} else {
		// Assume it's a directory path
		args = []string{"/grype", ".", "--fail-on", severity, "-o", "json"}
	}

	container := m.client.Container().From("anchore/grype:latest")
	
	if target[:4] != "dir:" && target[:6] != "image:" {
		// It's a directory path
		container = container.
			WithDirectory("/workspace", m.client.Host().Directory(target)).
			WithWorkdir("/workspace")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run grype scan with severity: no output received")
}

// DBStatus checks database status
func (m *GrypeModule) DBStatus(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("anchore/grype:latest").
		WithExec([]string{
			"/grype", "db", "status",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to check grype db status: no output received")
}

// DBCheck checks if database update is available
func (m *GrypeModule) DBCheck(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("anchore/grype:latest").
		WithExec([]string{
			"/grype", "db", "check",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to check grype db updates: no output received")
}

// DBUpdate updates vulnerability database
func (m *GrypeModule) DBUpdate(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("anchore/grype:latest").
		WithExec([]string{
			"/grype", "db", "update",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to update grype db: no output received")
}

// DBList lists available databases
func (m *GrypeModule) DBList(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("anchore/grype:latest").
		WithExec([]string{
			"/grype", "db", "list",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list grype dbs: %w", err)
	}

	return output, nil
}

// Explain explains vulnerability findings
func (m *GrypeModule) Explain(ctx context.Context, id string) (string, error) {
	container := m.client.Container().
		From("anchore/grype:latest").
		WithExec([]string{
			"/grype", "explain", id,
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to explain vulnerability: %w", err)
	}

	return output, nil
}

// GetVersion returns Grype version
func (m *GrypeModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("anchore/grype:latest").
		WithExec([]string{
			"/grype", "--version",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to get grype version: no output received")
}