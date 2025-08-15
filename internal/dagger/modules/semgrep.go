package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// SemgrepModule runs Semgrep for static analysis
type SemgrepModule struct {
	client *dagger.Client
	name   string
}

// NewSemgrepModule creates a new Semgrep module
func NewSemgrepModule(client *dagger.Client) *SemgrepModule {
	return &SemgrepModule{
		client: client,
		name:   "semgrep",
	}
}

// ScanDirectory scans a directory with Semgrep rules
func (m *SemgrepModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("returntocorp/semgrep:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"semgrep",
			"--config=auto",
			"--json",
			"--severity=ERROR",
			".",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run semgrep: no output received")
}

// ScanWithRuleset scans with specific ruleset
func (m *SemgrepModule) ScanWithRuleset(ctx context.Context, dir string, ruleset string) (string, error) {
	container := m.client.Container().
		From("returntocorp/semgrep:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"semgrep",
			"--config", ruleset,
			"--json",
			".",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run semgrep with ruleset: no output received")
}

// ScanFile scans a specific file
func (m *SemgrepModule) ScanFile(ctx context.Context, filePath string) (string, error) {
	container := m.client.Container().
		From("returntocorp/semgrep:latest").
		WithFile("/workspace/target.file", m.client.Host().File(filePath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"semgrep",
			"--config=auto",
			"--json",
			"target.file",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run semgrep on file: no output received")
}

// GetVersion returns the version of Semgrep
func (m *SemgrepModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("returntocorp/semgrep:latest").
		WithExec([]string{"semgrep", "--version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get semgrep version: %w", err)
	}

	return output, nil
}
