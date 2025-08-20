package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// ActionlintModule runs actionlint for GitHub Actions workflow validation
type ActionlintModule struct {
	client *dagger.Client
	name   string
}

const actionlintBinary = "actionlint"

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
		From("wolff2023/actionlint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec([]string{
			actionlintBinary,
			"-format", "{{json .}}",
			"-color",
		}, dagger.ContainerWithExecOpts{
			// actionlint returns non-zero exit code when it finds issues
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	// If there's stderr output, there might be an error
	if stderr != "" {
		return "", fmt.Errorf("actionlint stderr: %s", stderr)
	}
	
	// If there's no output, it means no issues were found (success case)
	if output == "" {
		return "No workflow issues found", nil
	}

	return output, nil
}

// ScanFile scans a specific workflow file
func (m *ActionlintModule) ScanFile(ctx context.Context, filePath string) (string, error) {
	container := m.client.Container().
		From("wolff2023/actionlint:latest").
		WithFile("/workspace/workflow.yml", m.client.Host().File(filePath)).
		WithWorkdir("/workspace").
		WithExec([]string{
			actionlintBinary,
			"-format", "{{json .}}",
			"workflow.yml",
		}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	// If there's stderr output, there might be an error
	if stderr != "" {
		return "", fmt.Errorf("actionlint stderr: %s", stderr)
	}
	
	// If there's no output, it means no issues were found (success case)
	if output == "" {
		return "No workflow issues found", nil
	}

	return output, nil
}

// ScanDirectoryWithOptions scans a directory with advanced options (format template, ignore patterns, color)
func (m *ActionlintModule) ScanDirectoryWithOptions(ctx context.Context, dir string, formatTemplate, ignorePatterns string, color bool) (string, error) {
	args := []string{actionlintBinary}

	if formatTemplate != "" {
		args = append(args, "-format", formatTemplate)
	} else {
		args = append(args, "-format", "{{json .}}")
	}

	if ignorePatterns != "" {
		// Split comma-separated patterns and add each with -ignore flag
		patterns := strings.Split(ignorePatterns, ",")
		for _, pattern := range patterns {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				args = append(args, "-ignore", pattern)
			}
		}
	}

	if color {
		args = append(args, "-color")
	}

	container := m.client.Container().
		From("wolff2023/actionlint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	// If there's stderr output, there might be an error
	if stderr != "" {
		return "", fmt.Errorf("actionlint stderr: %s", stderr)
	}
	
	// If there's no output, it means no issues were found (success case)
	if output == "" {
		return "No workflow issues found", nil
	}

	return output, nil
}

// ScanWithExternalTools scans workflows with shellcheck and pyflakes integration
func (m *ActionlintModule) ScanWithExternalTools(ctx context.Context, dir, shellcheckPath, pyflakesPath string, color bool) (string, error) {
	args := []string{actionlintBinary, "-format", "{{json .}}"}

	if shellcheckPath != "" {
		args = append(args, "-shellcheck", shellcheckPath)
	}
	if pyflakesPath != "" {
		args = append(args, "-pyflakes", pyflakesPath)
	}
	if color {
		args = append(args, "-color")
	}

	container := m.client.Container().
		From("wolff2023/actionlint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	// If there's stderr output, there might be an error
	if stderr != "" {
		return "", fmt.Errorf("actionlint stderr: %s", stderr)
	}
	
	// If there's no output, it means no issues were found (success case)
	if output == "" {
		return "No workflow issues found", nil
	}

	return output, nil
}

// ScanSpecificFiles scans specific workflow files with options
func (m *ActionlintModule) ScanSpecificFiles(ctx context.Context, dir string, workflowFiles []string, formatTemplate, ignorePatterns string, color bool) (string, error) {
	args := []string{actionlintBinary}

	if formatTemplate != "" {
		args = append(args, "-format", formatTemplate)
	} else {
		args = append(args, "-format", "{{json .}}")
	}

	if ignorePatterns != "" {
		// Split comma-separated patterns and add each with -ignore flag
		patterns := strings.Split(ignorePatterns, ",")
		for _, pattern := range patterns {
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				args = append(args, "-ignore", pattern)
			}
		}
	}

	if color {
		args = append(args, "-color")
	}

	// Add specific workflow files
	for _, file := range workflowFiles {
		file = strings.TrimSpace(file)
		if file != "" {
			args = append(args, file)
		}
	}

	container := m.client.Container().
		From("wolff2023/actionlint:latest").
		WithDirectory("/workspace", m.client.Host().Directory(dir)).
		WithWorkdir("/workspace").
		WithExec(args, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	stderr, _ := container.Stderr(ctx)
	
	// If there's stderr output, there might be an error
	if stderr != "" {
		return "", fmt.Errorf("actionlint stderr: %s", stderr)
	}
	
	// If there's no output, it means no issues were found (success case)
	if output == "" {
		return "No workflow issues found", nil
	}

	return output, nil
}

// GetVersion returns the version of actionlint
func (m *ActionlintModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("wolff2023/actionlint:latest").
		WithExec([]string{actionlintBinary, "-version"})

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get actionlint version: %w", err)
	}

	return output, nil
}
