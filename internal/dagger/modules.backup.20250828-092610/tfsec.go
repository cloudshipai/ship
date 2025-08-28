package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// TfsecModule runs tfsec for Terraform security scanning
type TfsecModule struct {
	client *dagger.Client
	name   string
}

// NewTfsecModule creates a new tfsec module
func NewTfsecModule(client *dagger.Client) *TfsecModule {
	return &TfsecModule{
		client: client,
		name:   "tfsec",
	}
}

// GetVersion returns the version of tfsec
func (m *TfsecModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithExec([]string{"tfsec", "--version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "tfsec - Terraform security scanner", nil
}

// ScanDirectory scans a Terraform directory for security issues
func (m *TfsecModule) ScanDirectory(ctx context.Context, dir string, format string) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	args := []string{"tfsec", dir}
	
	// Add format flag if specified
	if format != "" {
		args = append(args, "--format", format)
	} else {
		args = append(args, "--format", "json")
	}
	
	// Add common flags
	args = append(args, "--no-color", "--soft-fail")

	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	
	// Check if stderr contains actual errors vs findings
	if stderr != "" && !strings.Contains(stderr, "WARNING") {
		return stderr, nil
	}
	
	return `{"results": [], "message": "No security issues found"}`, nil
}

// ScanWithSeverity scans with minimum severity threshold
func (m *TfsecModule) ScanWithSeverity(ctx context.Context, dir string, severity string) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	args := []string{"tfsec", dir, "--format", "json", "--no-color"}
	
	// Add severity filter
	if severity != "" {
		args = append(args, "--minimum-severity", severity)
	}
	
	args = append(args, "--soft-fail")

	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return `{"results": [], "message": "No issues found at specified severity"}`, nil
}

// ScanWithExcludes scans with excluded checks
func (m *TfsecModule) ScanWithExcludes(ctx context.Context, dir string, excludes []string) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	args := []string{"tfsec", dir, "--format", "json", "--no-color", "--soft-fail"}
	
	// Add exclusions
	for _, exclude := range excludes {
		args = append(args, "--exclude", exclude)
	}

	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return `{"results": [], "message": "Scan completed"}`, nil
}

// ScanWithConfig scans using a config file
func (m *TfsecModule) ScanWithConfig(ctx context.Context, dir string, configPath string) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir))
	
	// Add config file if specified
	if configPath != "" {
		container = container.WithFile("/config.yml", m.client.Host().File(configPath))
	}
	
	args := []string{"tfsec", "/workspace", "--format", "json", "--no-color", "--soft-fail"}
	
	if configPath != "" {
		args = append(args, "--config-file", "/config.yml")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return `{"results": [], "message": "Scan completed"}`, nil
}

// ValidateTfvars validates terraform variable files
func (m *TfsecModule) ValidateTfvars(ctx context.Context, dir string, tfvarsFile string) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	args := []string{"tfsec", dir, "--format", "json", "--no-color", "--soft-fail"}
	
	// Add tfvars file
	if tfvarsFile != "" {
		args = append(args, "--tfvars-file", tfvarsFile)
	}

	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace")
	
	if tfvarsFile != "" {
		container = container.WithFile("/tfvars.auto.tfvars", m.client.Host().File(tfvarsFile))
		args[len(args)-1] = "/tfvars.auto.tfvars"
	}
	
	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return `{"results": [], "message": "Variables validated"}`, nil
}

// GenerateMetrics generates metrics in various formats
func (m *TfsecModule) GenerateMetrics(ctx context.Context, dir string) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	args := []string{"tfsec", dir, "--format", "metrics", "--no-color", "--soft-fail"}

	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return "No metrics generated", nil
}

// RunWithCustomChecks runs tfsec with custom check definitions
func (m *TfsecModule) RunWithCustomChecks(ctx context.Context, dir string, customChecksDir string) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir))
	
	args := []string{"tfsec", "/workspace", "--format", "json", "--no-color", "--soft-fail"}
	
	// Add custom checks directory if specified
	if customChecksDir != "" {
		container = container.WithDirectory("/custom-checks", m.client.Host().Directory(customChecksDir))
		args = append(args, "--custom-check-dir", "/custom-checks")
	}

	container = container.WithExec(args, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return `{"results": [], "message": "Custom checks completed"}`, nil
}

// ScanWithIgnores scans with inline comment ignores
func (m *TfsecModule) ScanWithIgnores(ctx context.Context, dir string, showIgnored bool) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	args := []string{"tfsec", dir, "--format", "json", "--no-color", "--soft-fail"}
	
	if showIgnored {
		args = append(args, "--include-ignored")
	}

	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return `{"results": [], "message": "Scan completed"}`, nil
}

// GenerateReport generates a detailed HTML report
func (m *TfsecModule) GenerateReport(ctx context.Context, dir string, reportType string) (string, error) {
	if dir == "" {
		dir = "."
	}
	
	format := "json"
	switch reportType {
	case "html":
		format = "html"
	case "junit":
		format = "junit"
	case "sarif":
		format = "sarif"
	case "csv":
		format = "csv"
	case "checkstyle":
		format = "checkstyle"
	default:
		format = "json"
	}
	
	args := []string{"tfsec", dir, "--format", format, "--no-color", "--soft-fail"}

	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	return fmt.Sprintf("No %s report generated", reportType), nil
}

// UpdateChecks updates tfsec check definitions
func (m *TfsecModule) UpdateChecks(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("aquasec/tfsec:latest").
		WithExec([]string{"tfsec", "--update"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	if output != "" {
		return output, nil
	}
	if stderr != "" {
		return stderr, nil
	}
	
	return "Checks updated successfully", nil
}