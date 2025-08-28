package modules

import (
	"context"

	"dagger.io/dagger"
)

// GitleaksModule runs Gitleaks for fast secret scanning
type GitleaksModule struct {
	client *dagger.Client
	name   string
}

// GitleaksDetectOptions contains options for gitleaks detection
type GitleaksDetectOptions struct {
	ConfigPath   string
	ReportFormat string
	ReportPath   string
	Verbose      bool
	NoGit        bool
}

// GitleaksProtectOptions contains options for gitleaks protection
type GitleaksProtectOptions struct {
	ConfigPath string
	Staged     bool
	Verbose    bool
}

// NewGitleaksModule creates a new Gitleaks module
func NewGitleaksModule(client *dagger.Client) *GitleaksModule {
	return &GitleaksModule{
		client: client,
		name:   "gitleaks",
	}
}

// Detect runs gitleaks detect on the provided directory
func (m *GitleaksModule) Detect(ctx context.Context, sourcePath string, opts GitleaksDetectOptions) (string, error) {
	args := []string{"gitleaks", "detect"}

	// Add detection options
	if opts.ConfigPath != "" {
		args = append(args, "--config", opts.ConfigPath)
	}
	if opts.ReportFormat != "" {
		args = append(args, "--report-format", opts.ReportFormat)
	}
	if opts.ReportPath != "" {
		args = append(args, "--report-path", opts.ReportPath)
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}
	if opts.NoGit {
		args = append(args, "--no-git")
	}

	// Add source path
	args = append(args, "--source", ".")

	container := m.client.Container().
		From(getImageTag("gitleaks", "zricethezav/gitleaks:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(sourcePath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	stdout, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if stdout != "" {
		return stdout, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "No secrets detected", nil
}

// Protect runs gitleaks protect for pre-commit scanning
func (m *GitleaksModule) Protect(ctx context.Context, sourcePath string, opts GitleaksProtectOptions) (string, error) {
	args := []string{"gitleaks", "protect"}

	// Add protection options
	if opts.ConfigPath != "" {
		args = append(args, "--config", opts.ConfigPath)
	}
	if opts.Staged {
		args = append(args, "--staged")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}

	container := m.client.Container().
		From(getImageTag("gitleaks", "zricethezav/gitleaks:latest")).
		WithDirectory("/workspace", m.client.Host().Directory(sourcePath)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	stdout, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)

	if stdout != "" {
		return stdout, nil
	}
	if stderr != "" {
		return stderr, nil
	}

	return "Protection scan completed successfully", nil
}