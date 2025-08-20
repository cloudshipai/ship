package modules

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// InfraMapModule generates infrastructure diagrams from Terraform
type InfraMapModule struct {
	client *dagger.Client
}

// NewInfraMapModule creates a new InfraMap module instance
func NewInfraMapModule(client *dagger.Client) *InfraMapModule {
	return &InfraMapModule{
		client: client,
	}
}

// GenerateFromState generates an infrastructure diagram from a Terraform state file
func (m *InfraMapModule) GenerateFromState(ctx context.Context, stateFile string, format string) (string, error) {
	// Get the directory containing the state file
	workDir := m.client.Host().Directory(".")

	// Create container with InfraMap
	container := m.client.Container().
		From("cycloid/inframap:latest").
		WithDirectory("/opt", workDir).
		WithWorkdir("/opt")

	// Generate the diagram based on format
	var output string
	var err error

	switch format {
	case "png", "svg", "pdf":
		// Generate dot format first, then convert
		result := container.
			WithExec([]string{"/home/inframap/inframap", "generate", stateFile}, dagger.ContainerWithExecOpts{
				Expect: "ANY",
			})

		output, err = result.Stdout(ctx)
		if err != nil {
			stderr, _ := result.Stderr(ctx)
			return "", fmt.Errorf("failed to generate diagram: %w\nStderr: %s", err, stderr)
		}

	case "dot":
		// Generate raw dot format
		result := container.
			WithExec([]string{"/home/inframap/inframap", "generate", stateFile}, dagger.ContainerWithExecOpts{
				Expect: "ANY",
			})

		output, err = result.Stdout(ctx)
		if err != nil {
			stderr, _ := result.Stderr(ctx)
			return "", fmt.Errorf("failed to generate dot output: %w\nStderr: %s", err, stderr)
		}

	default:
		return "", fmt.Errorf("unsupported format: %s (supported: png, svg, pdf, dot)", format)
	}

	return output, nil
}

// GenerateFromHCL generates an infrastructure diagram from Terraform HCL files
func (m *InfraMapModule) GenerateFromHCL(ctx context.Context, directory string, format string) (string, error) {
	// Get the directory containing HCL files
	workDir := m.client.Host().Directory(directory)

	// Create container with InfraMap
	container := m.client.Container().
		From("cycloid/inframap:latest").
		WithDirectory("/opt", workDir).
		WithWorkdir("/opt")

	// Generate the diagram
	var output string
	var err error

	switch format {
	case "png", "svg", "pdf":
		// For HCL, we need to specify all .tf files
		result := container.
			WithExec([]string{"/home/inframap/inframap", "generate", "--hcl", "."}, dagger.ContainerWithExecOpts{
				Expect: "ANY",
			})

		output, err = result.Stdout(ctx)
		if err != nil {
			stderr, _ := result.Stderr(ctx)
			return "", fmt.Errorf("failed to generate diagram from HCL: %w\nStderr: %s", err, stderr)
		}

	case "dot":
		result := container.
			WithExec([]string{"/home/inframap/inframap", "generate", "--hcl", "."}, dagger.ContainerWithExecOpts{
				Expect: "ANY",
			})

		output, err = result.Stdout(ctx)
		if err != nil {
			stderr, _ := result.Stderr(ctx)
			return "", fmt.Errorf("failed to generate dot from HCL: %w\nStderr: %s", err, stderr)
		}

	default:
		return "", fmt.Errorf("unsupported format: %s (supported: png, svg, pdf, dot)", format)
	}

	return output, nil
}

// GenerateWithOptions generates a diagram with custom options
func (m *InfraMapModule) GenerateWithOptions(ctx context.Context, input string, options InfraMapOptions) (string, error) {
	workDir := m.client.Host().Directory(".")

	container := m.client.Container().
		From("cycloid/inframap:latest").
		WithDirectory("/opt", workDir).
		WithWorkdir("/opt")

	// Build command with options
	args := []string{"/home/inframap/inframap", "generate"}

	if options.Raw {
		args = append(args, "--raw")
	}

	if !options.Clean {
		args = append(args, "--clean=false")
	}

	if options.Provider != "" {
		args = append(args, "--provider", options.Provider)
	}

	args = append(args, input)

	// Add output format conversion if needed
	cmd := strings.Join(args, " ")
	if options.Format != "" && options.Format != "dot" {
		cmd = fmt.Sprintf("%s | dot -T%s", cmd, options.Format)
	}

	result := container.WithExec([]string{"sh", "-c", cmd}, dagger.ContainerWithExecOpts{
		Expect: "ANY",
	})

	output, err := result.Stdout(ctx)
	if err != nil {
		stderr, _ := result.Stderr(ctx)
		return "", fmt.Errorf("failed to generate diagram: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// PruneState removes unnecessary information from Terraform state
func (m *InfraMapModule) PruneState(ctx context.Context, stateFile string) (string, error) {
	workDir := m.client.Host().Directory(".")

	container := m.client.Container().
		From("cycloid/inframap:latest").
		WithDirectory("/opt", workDir).
		WithWorkdir("/opt")

	result := container.
		WithExec([]string{"/home/inframap/inframap", "prune", stateFile}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, err := result.Stdout(ctx)
	if err != nil {
		stderr, _ := result.Stderr(ctx)
		return "", fmt.Errorf("failed to prune state: %w\nStderr: %s", err, stderr)
	}

	return output, nil
}

// GetVersion returns the version of InfraMap
func (m *InfraMapModule) GetVersion(ctx context.Context) (string, error) {
	container := m.client.Container().
		From("cycloid/inframap:latest").
		WithExec([]string{"/home/inframap/inframap", "version"}, dagger.ContainerWithExecOpts{
			Expect: "ANY",
		})

	output, _ := container.Stdout(ctx)
	if output != "" {
		return output, nil
	}
	
	// InfraMap might output to stderr
	stderr, _ := container.Stderr(ctx)
	if stderr != "" {
		return stderr, nil
	}
	
	return "", fmt.Errorf("failed to get inframap version: no output received")
}

// InfraMapOptions contains options for diagram generation
type InfraMapOptions struct {
	// Raw shows all resources without InfraMap logic
	Raw bool
	// Clean removes unconnected nodes (default: true)
	Clean bool
	// Provider filters by specific provider (aws, google, azurerm, etc.)
	Provider string
	// Format output format (png, svg, pdf, dot)
	Format string
}