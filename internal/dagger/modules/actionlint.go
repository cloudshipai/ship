package modules

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// ActionlintModule runs actionlint for GitHub Actions workflow validation
type ActionlintModule struct {
	client *dagger.Client
	name   string
}

// NewActionlintModule creates a new actionlint module
func NewActionlintModule(client *dagger.Client) *ActionlintModule {
	return &ActionlintModule{
		client: client,
		name:   "actionlint",
	}
}

// ScanDirectory scans a directory for GitHub Actions workflow issues
func (m *ActionlintModule) ScanDirectory(ctx context.Context, dir string) (string, error) {
	container := m.client.Container().
		From("rhysd/actionlint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"actionlint",
			"-format", "{{json .}}",
			"-color",
		}, dagger.ContainerWithExecOpts{
			// actionlint returns non-zero exit code when it finds issues
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run actionlint: no output received")
}

// ScanFile scans a specific workflow file
func (m *ActionlintModule) ScanFile(ctx context.Context, filePath string) (string, error) {
	container := m.client.Container().
		From("rhysd/actionlint:latest").
		WithFile("/workspace/workflow.yml", m.client.Host().File(filePath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			"actionlint",
			"-format", "{{json .}}",
			"workflow.yml",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}

	return "", fmt.Errorf("failed to run actionlint on file: no output received")
}

// GetVersion returns the version of actionlint
func (m *ActionlintModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("rhysd/actionlint:latest").
		WithExec([]string{"actionlint", "-version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get actionlint version: %w", err)
	}

	return output, nil
}
